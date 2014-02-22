#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Author: Herman Schaaf
# Date: 2013

import unittest
import pprint
from edict2_parser import Parser

test_cases = {
    "ip": "ＩＰアドレス [アイピーアドレス] /(n) (See インターネットプロトコルアドレス) Internet Protocol address/IP address/EntL2159840X/",
    "biology": "ＭＨＣ [エムエッチシー] /(n) {biol} major histocompatibility complex/MHC/EntL2773110/",
    "seasonal": "季節的 [きせつてき] /(adj-na) seasonal/EntL1222860X/",
    "late_spring": "季春 [きしゅん] /(n) (1) late spring/(2) (obs) third month of the lunar calendar/EntL1639680X/",
    "quarterly": "季刊誌 [きかんし] /(n) (See 季刊雑誌) a quarterly (magazine)/EntL1757050X/",
    "mouth": "口 [くち] /(n) (1) mouth/(2) opening/hole/gap/orifice/(3) mouth (of a bottle)/spout/nozzle/mouthpiece/(4) gate/door/entrance/exit/(5) (See 口を利く・1) speaking/speech/talk (i.e. gossip)/(6) (See 口に合う) taste/palate/(7) mouth (to feed)/(8) (See 働き口) opening (i.e. vacancy)/available position/(9) (See 口がかかる) invitation/summons/(10) kind/sort/type/(11) opening (i.e. beginning)/(suf,ctr) (12) counter for mouthfuls, shares (of money), and swords/(P)/EntL1275640X/",
    "demon": "デーモン /(n) (1) demon/(2) {comp} daemon (in Unix, etc.)/(P)/EntL1081470X/"
}

pp = pprint.PrettyPrinter(indent=4)


class TestEdict2Parser(unittest.TestCase):

    def setUp(self):
        self.parser = Parser()

    def tearDown(self):
        pass

    def test_ip_case(self):
        entries = self.parser.get_entries(test_cases['ip'])
        entry = entries[0]

        # test main entry
        self.assertEqual(entry.japanese, "ＩＰアドレス")
        self.assertEqual(entry.furigana, "アイピーアドレス")
        self.assertEqual(entry.ent_seq, "2159840")
        self.assertEqual(entry.has_audio, True)
        self.assertEqual(entry.tags, ("n",))
        self.assertEqual(len(entry.glosses), 2)

        # test glosses
        self.assertEqual(entry.glosses[0].tags, [])
        self.assertEqual(entry.glosses[0].related, ["インターネットプロトコルアドレス"])
        self.assertEqual(entry.glosses[0].english, "Internet Protocol address")

        self.assertEqual(entry.glosses[1].tags, [])
        self.assertEqual(entry.glosses[1].english, "IP address")

    def test_field_case(self):
        entries = self.parser.get_entries(test_cases['biology'])
        entry = entries[0]

        self.assertEqual(entry.japanese, "ＭＨＣ")
        self.assertEqual(entry.furigana, "エムエッチシー")
        self.assertEqual(entry.tags, ("n",))
        self.assertEqual(len(entry.glosses), 2)

        self.assertEqual(entry.glosses[0].field, "biol")

    def test_seasonal_case(self):
        entries = self.parser.get_entries(test_cases['seasonal'])
        entry = entries[0]

        self.assertEqual(entry.japanese, "季節的")
        self.assertEqual(entry.furigana, "きせつてき")
        self.assertEqual(entry.tags, ("adj-na",))
        self.assertEqual(len(entry.glosses), 1)
        self.assertEqual(entry.ent_seq, "1222860")

    def test_quarterly_case(self):
        entries = self.parser.get_entries(test_cases['quarterly'])
        entry = entries[0]

        self.assertEqual(entry.japanese, "季刊誌")
        self.assertEqual(entry.furigana, "きかんし")
        self.assertEqual(entry.tags, ("n",))
        self.assertEqual(len(entry.glosses), 1)

        self.assertEqual(entry.glosses[0].related, ["季刊雑誌"])

    def test_computer_term_case(self):
        entries = self.parser.get_entries(test_cases['demon'])
        entry = entries[0]

        self.assertEqual(entry.japanese, "デーモン")
        self.assertEqual(entry.furigana, "デーモン")
        self.assertEqual(entry.tags, ("P", "n"))
        self.assertEqual(len(entry.glosses), 2)

        self.assertEqual(entry.glosses[1].field, "comp")

    def test_many_glosses_case(self):
        entries = self.parser.get_entries(test_cases['mouth'])
        entry = entries[0]

        self.assertEqual(entry.japanese, "口")
        self.assertEqual(entry.furigana, "くち")
        self.assertEqual(entry.tags, ("P", "n"))
        self.assertEqual(len(entry.glosses), 12)

        self.assertEqual(entry.glosses[0].english, "mouth")
        self.assertEqual(entry.glosses[1].english, "opening/hole/gap/orifice")
        self.assertEqual(entry.glosses[4].related, ["口を利く・1"])
        self.assertEqual(entry.glosses[5].related, ["口に合う"])
        self.assertEqual(entry.glosses[11].tags, ('suf', 'ctr'))
        self.assertEqual(entry.glosses[11].tags, ('suf', 'ctr'))
        self.assertEqual(entry.glosses[11].english, "counter for mouthfuls, shares (of money), and swords")

if __name__ == '__main__':
    unittest.main()
