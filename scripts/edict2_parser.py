#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Author: Herman Schaaf
# Date: 2013

# Some parts lifted from Paul Goins's script for Edict 1, 2009

# EDICT2 FORMAT:
#
#    KANJI-1;KANJI-2 [KANA-1;KANA-2] /(general information) (see xxxx) gloss/gloss/.../
#    垜;安土;堋 [あずち] /(n) mound on which targets are placed (in archery)/firing mound/EntL2542010/

import os, re, gzip, gettext, pprint

# Part of speech codes
valid_pos_codes = list(set((
    "adj-i", "adj-na", "adj-no", "adj-pn", "adj-t", "adj-f", "adj",
    "adv", "adv-to", "aux", "aux-v", "aux-adj", "conj", "ctr", "exp",
    "int", "iv", "n", "n-adv", "n-suf", "n-pref", "n-t", "num", "pn",
    "pref", "prt", "suf", "v1", "v2a-s", "v4h", "v4r", "v5", "v5aru",
    "v5b", "v5g", "v5k", "v5k-s", "v5m", "v5n", "v5r", "v5r-i", "v5s",
    "v5t", "v5u", "v5u-s", "v5uru", "v5z", "vz", "vi", "vk", "vn",
    "vr", "vs", "vs-s", "vs-i", "vt",
    )))

# Field of application codes
valid_foa_codes = list(set((
    "Buddh", "MA", "comp", "food", "geom", "ling", "math", "mil",
    "physics", "chem", "biol"
    )))

# Miscellaneous marking codes
valid_misc_codes = list(set((
    "X", "abbr", "arch", "ateji", "chn", "col", "derog", "eK", "ek",
    "fam", "fem", "gikun", "hon", "hum", "iK", "id", "ik", "io",
    "m-sl", "male", "male-sl", "oK", "obs", "obsc", "ok", "on-mim",
    "poet", "pol", "rare", "sens", "sl", "uK", "uk", "vulg", "P"
    )))

# Dialect codes
valid_dialect_codes = list(set((
    "kyb", "osb", "ksb", "ktb", "tsb", "thb", "tsug", "kyu", "rkb",
    "nab"
    )))

re_kana = re.compile(r'\[(.*)\]')
re_tags = re.compile(r'\(((?:%s|[,]+)+)\:?\)' % '|'.join(valid_pos_codes + valid_misc_codes + valid_dialect_codes))
re_field_tags = re.compile(r'\{(%s)\}' % '|'.join(valid_foa_codes))
re_number_tag = re.compile(r'\((\d+)\)')
re_any_tag = re.compile(r'\(([^)]*)\)')
re_related_tag = re.compile(r'\(See ([^)]*)\)', re.IGNORECASE)

class EdictGloss(object):

    def __init__(self, english, tags, field, related=[], common=False):
        self.english = english
        self.tags = tags
        self.field = field
        self.related= related

        self.common = common

        if 'P' in self.tags:
            self.common = True

    def to_dict(self):
        d = {}
        for f in ['english', 'tags', 'field', 'related', 'common']:
            d[f] = getattr(self, f)
        return d

class EdictEntry(object):

    def __init__(self, glosses, japanese, furigana, tags=set(), kanji_tags=set(), kana_tags=set(), ent_seq=None, has_audio=False):
        self.glosses = glosses
        self.japanese = japanese
        self.furigana = furigana

        self.tags = tags
        self.kanji_tags = kanji_tags
        self.kana_tags = kana_tags

        self.ent_seq = ent_seq
        self.has_audio = has_audio

    def to_dict(self):
        pos = filter(lambda t: t in valid_pos_codes, self.tags)
        fields = filter(lambda t: t in valid_foa_codes, self.tags)
        tags = filter(lambda t: t in valid_misc_codes, self.tags)
        dialects = filter(lambda t: t in valid_dialect_codes, self.tags)
        furigana = self.japanese if not self.furigana else self.furigana
        d = {
            'japanese': self.japanese,
            'furigana': self.furigana,
            'glosses': [g.to_dict() for g in self.glosses],
            'pos': pos, # parts of speech
            'fields': fields, # fields of interest
            'tags': tags, # general tags
            'dialects': dialects, # dialects
            'kana_tags': self.kana_tags,
            'kanji_tags': self.kanji_tags,
            'common': 'P' in self.tags,
            'ent_seq': self.ent_seq,
            'has_audio': self.has_audio
        }
        return d

class Parser(object):

    def __init__(self, filename=None, encoding="EUC-JP"):
        if not filename is None and not os.path.exists(filename):
            raise Exception("Dictionary file does not exist.")
        self.filename = filename
        self.encoding = encoding
        self.cache = {}

    def extract_tags(self, word, expression=re_tags):
        t = expression.search(word)
        tags = []
        if t:
            groups = t.groups()
            tags = []
            for group in groups:
                tags += group.split(',')
            tags = tuple(tags)
            word = expression.sub('', word)
        return word, tags

    def extract_related(self, gloss):
        return self.extract_tags(word=gloss, expression=re_related_tag)

    def extract_fields(self, gloss):
        return self.extract_tags(word=gloss, expression=re_field_tags)

    def get_entries(self, raw_entry):
        #print raw_entry
        raw_words = raw_entry.split(' ')

        raw_kanji = raw_words[0].split(';')
        kana_match = re_kana.match(raw_words[1])
        if kana_match:
            raw_kana = kana_match.groups(0)[0].split(';')
        else:
            raw_kana = raw_kanji

        kanji_tagged = [self.extract_tags(k) for k in raw_kanji]
        kana_tagged = [self.extract_tags(k) for k in raw_kana]

        raw_english = raw_entry.split('/')[1:-2]

        english, main_tags = self.extract_tags(raw_english[0])
        english = [english] + raw_english[1:]

        if english[-1] == '(P)':
            main_tags = tuple(set(list(main_tags) + ['P']))
            english = english[:-1]

        # join numbered entries:
        joined_english = []
        has_numbers = False
        for e in english:
            clean, number = self.extract_tags(e, re_number_tag)
            clean = clean.strip()
            if number:
                has_numbers = True
                joined_english.append(clean)
            elif has_numbers:
                joined_english[-1] += '/' + clean
            else:
                joined_english.append(clean)

        english = joined_english

        glosses = []
        for gloss in english:
            clean_gloss, related_words = self.extract_related(gloss)
            clean_gloss, tags = self.extract_tags(clean_gloss)
            clean_gloss, fields = self.extract_fields(clean_gloss)

            if related_words:
                related_words = related_words[0].split(',')
            else:
                related_words = []

            field = fields[0] if fields else None
            cg = clean_gloss.strip()
            if cg:
                glosses.append(EdictGloss(english=cg, tags=tags, field=field, related=related_words))

        ent_seq = raw_entry.split('/')[-2]

        # entL sequences that end in X have audio clips
        has_audio = ent_seq[-1] == 'X'

        # throw away the entL and X part, keeping only the id
        ent_seq = ent_seq[4:]
        if has_audio:
            ent_seq = ent_seq[:-1]

        entries = []
        # create EdictEntry objects:
        for kana, ktag in kana_tagged:
            # special case for kana like this: おくび(噯,噯気);あいき(噯気,噫気,噯木)
            kana, matching_kanji = self.extract_tags(kana, expression=re_any_tag)
            if matching_kanji:
                matching_kanji = matching_kanji[0].split(',')

            for kanji, jtag in kanji_tagged:
                kwargs = {
                    "glosses": glosses,
                    "japanese": kanji,
                    "furigana": kana,
                    "tags": main_tags,
                    "kanji_tags": jtag,
                    "kana_tags": ktag,
                    "ent_seq": ent_seq,
                    "has_audio": has_audio,
                }
                if not matching_kanji:
                    entries.append(EdictEntry(**kwargs))
                else:
                    if kanji in matching_kanji:
                        entries.append(EdictEntry(**kwargs))
        return entries

    def parse(self):
        # Read from file
        if len(self.filename) >= 3 and self.filename[-3:] == ".gz":
            f = gzip.open(self.filename)
        else:
            f = open(self.filename, "rb")
        fdata = f.read()
        f.close()
        fdata = fdata.decode(self.encoding)

        lines = fdata.splitlines()
        lines = [line for line in lines if line and (line[0] != u"#")]

        data = {}
        entries = []
        for line in lines:
            entries += self.get_entries(line)

        for e in entries:
            yield e.to_dict()

if __name__ == '__main__':
    parser = Parser('../data/edict2')
    for e in parser.parse():
        print e
