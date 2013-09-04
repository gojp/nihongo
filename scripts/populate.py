from edict_parser import EdictEntry
from pymongo import Connection

MONGO_URI = ''
c = Connection(MONGO_URI)

mongo_db = c['greenbook']
collection = mongo_db['edict']

PATH_TO_EDICT2 = ''

inserts = []
with open(PATH_TO_EDICT2) as f:
    read_data = f.readlines()
    for line in [l.decode('EUC-JP') for l in read_data]:
        inserts.append(EdictEntry(line).__dict__)

collection.insert(inserts)
