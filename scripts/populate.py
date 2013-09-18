from edict2_parser import Parser
from pymongo import Connection
from pyes import ES
import romkan
import json
import subprocess
import sys

if len(sys.argv) <= 1:
    print 'usage: populate.py [es/mongo]'
    sys.exit(1)

es_or_mongo = sys.argv[1]

mongo = es_or_mongo == 'mongo'

inserts = []
if mongo:
    MONGO_URI = 'localhost'
    c = Connection(MONGO_URI)
    mongo_db = c['greenbook']
    collection = mongo_db['edict']
else:
    ELASTICSEARCH_URI = 'http://127.0.0.1:9200'
    conn = ES(ELASTICSEARCH_URI)
    conn.delete_index_if_exists("edict")
    conn.create_index("edict")
    conn.put_mapping("entry")

PATH_TO_EDICT2 = '../data/edict2'

parser = Parser(PATH_TO_EDICT2)
for e in parser.parse():
    e['romaji'] = romkan.to_roma(e['furigana'])
    inserts.append(e)

if mongo:
    collection.insert(inserts)
else:
    with open('output.json', 'w') as f:
        for i, d in enumerate(inserts):
            # When bulk inserting you need lines like this
            # before each entry
            # http://www.elasticsearch.org/guide/reference/api/bulk/
            index_line = '{"index":{"_index":"edict","_type": "entry", "_id": "%s"}}' % str(i + 1)
            f.write(index_line + '\n')
            f.write(json.dumps(d) + '\n')
        print "INSERTED %d ENTRIES!" % i

subprocess.Popen(['curl', '-s', '-XPOST', '%s/_bulk' % ELASTICSEARCH_URI, '--data-binary', '@output.json'])
