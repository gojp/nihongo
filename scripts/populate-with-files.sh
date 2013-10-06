if [ $# -eq 0 ]
  then
    echo "output JSON from edict2_parser.py is required"
fi

curl -XDELETE 'http://localhost:9200/edict/'

split $1 -l 2000

for i in `ls`; do echo ''>>$i; done;

for i in `ls`; do curl -s -XPOST "localhost:9200/_bulk" --data-binary @$i; done;
