#!/bin/bash

if [ $# -eq 0 ]
  then
      echo "output JSON from populate.py is required (./populate.py print)"
      exit 1
fi

echo "Deleting current edict index..."
curl -XDELETE 'http://localhost:9200/edict/'

if [ -d "json" ]; then
  rm json/*
else
  mkdir json
fi

echo "Splitting JSON output into 2000 files..."
cd json && split ../$1 -l 2000

# Append a newline to each file
for i in `ls`; do echo ''>>$i; done;

echo "Bulk inserting into ES..."
for i in `ls`; do curl -s -XPOST "localhost:9200/_bulk" --data-binary @$i; done;

echo "Done!\n"
