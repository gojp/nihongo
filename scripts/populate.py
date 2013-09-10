from edict_parser import EdictEntry
from pymongo import Connection
import romkan

MONGO_URI = 'localhost'
c = Connection(MONGO_URI)

mongo_db = c['greenbook']
collection = mongo_db['edict']

PATH_TO_EDICT2 = '../data/edict'

inserts = []
with open(PATH_TO_EDICT2) as f:
    read_data = f.readlines()
    for line in [l.decode('EUC-JP') for l in read_data]:
    	d = EdictEntry(line).to_dict()
    	d['romaji'] = romkan.to_roma(d['furigana'])
        inserts.append(d)

collection.insert(inserts)
