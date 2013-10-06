from edict2_parser import Parser
from pymongo import Connection
import rawes
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
    mongo_db = c['nihongo']
    collection = mongo_db['edict']
else:
    ELASTICSEARCH_URI = 'localhost:9200'
    es = rawes.Elastic(ELASTICSEARCH_URI)

    try:
        print "Dropping index if exist..."
        # drop existing index
        es.delete('edict')
    except:
        print "No pre-existing index found"

    print "Creating a new index..."
    mapping = {
        "mappings" : {
            "entry" : {
                "properties" : {
                    "english": {"type": "string"},
                    "japanese" : {"type" : "string"},
                    "furigana" : {"type" : "string"},
                    "romaji": {"type": "string"},
                    "glosses": {
                        "properties": {
                            "english": {"type" : "string"},
                            "tags": {"type" : "string"},
                            "field": {"type" : "string"},
                            "related": {"type" : "string"}
                        },
                        "index name": "gloss",
                        "index": "no"
                    },
                    "common" : {"type" : "boolean"},
                    "ent_seq" : {"type" : "string"}
                },
                "_boost" : {"name" : "common_boost", "null_value" : 1.0},
            }
        }
    }
    es.post('edict', data=mapping)

PATH_TO_EDICT2 = '../data/edict2'


def get_inserts(max_chunk=10000):
    inserts = []
    parser = Parser(PATH_TO_EDICT2)
    i = 0
    for e in parser.parse():
        i += 1
        e['english'] = [g['english'] for g in e['glosses']]
        e['romaji'] = romkan.to_roma(e['furigana'])
        e['common_boost'] = 2.0 if e['common'] == True else 1.0
        inserts.append(e)
        if i % max_chunk == 0:
            yield inserts
            inserts = []
    yield inserts

counter = 0

if mongo:
    collection.drop()
    for insert_chunk in get_inserts():
        bulk_list = []
        for d in insert_chunk:
            counter += 1
            bulk_list.append(d)
        collection.insert(bulk_list)
        print("Inserted {}".format(counter))
else:
    f = open('output.json', 'w')
    inserts = get_inserts()
    inserted_count = 0


    bulk_list = []
    for insert_chunk in get_inserts():
        for d in insert_chunk:
            counter += 1
            # When bulk inserting you need lines like this
            # before each entry
            # http://www.elasticsearch.org/guide/reference/api/bulk/

            bulk_list.append({
                "index":
                    {"_index":"edict","_type": "entry", "_id": str(counter + 1)}
                })
            bulk_list.append(d)

        if bulk_list:
            bulk_body = '\n'.join(map(json.dumps, bulk_list))+'\n'
            es.post('edict/entry/_bulk', data=bulk_body)
            bulk_list = []

            print "INSERTED %d WORDS!" % counter

            # Uncomment this to only insert a certain amount of words:
            # ----
            # if counter > 100000:
            #     return False

        

