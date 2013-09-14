from edict_parser import EdictEntry
from pymongo import Connection
import romkan
import simplejson as json
import subprocess
import sys

if len(sys.argv) <= 1:
    print 'usage: populate.py [es/mongo]'
    sys.exit(1)

es_or_mongo = sys.argv[1]

mongo = es_or_mongo == 'mongo'

if mongo:
    MONGO_URI = 'localhost'
    c = Connection(MONGO_URI)
    mongo_db = c['greenbook']
    collection = mongo_db['edict']
    inserts = []
else:
    ELASTICSEARCH_URI = 'http://localhost:9200'

PATH_TO_EDICT2 = '/Users/shawn/Downloads/edict2'

with open(PATH_TO_EDICT2) as f:
    read_data = f.readlines()
    for i, line in enumerate([l.decode('EUC-JP') for l in read_data]):
        d = EdictEntry(line).to_dict()
        d['romaji'] = romkan.to_roma(d['furigana'])
        if 'unparsed' in d:
            del(d['unparsed'])
        if mongo == 'mongo':
            inserts.append(d)
        else:
            subprocess.Popen(['curl', '-XPUT', '%s/edict/entry/%s' % (ELASTICSEARCH_URI, i), '-d', json.dumps(d)])

if mongo:
    collection.insert(inserts)
