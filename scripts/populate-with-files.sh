#!/bin/bash

if [ $# -eq 0 ]
  then
      echo "output JSON from populate.py is required (./populate.py print > output.json then re-run this script with output.json as the argument)"
      exit 1
fi

printf "Deleting current edict index...\n"
curl -XDELETE 'http://localhost:9200/edict/'

if [ -d "json" ]; then
  rm json/*
else
  mkdir json
fi

printf "\nSplitting JSON output into 2000 files...\n"
cd json && split ../$1 -l 2000

# Append a newline to each file
for i in `ls`; do echo ''>>$i; done;

printf "Bulk inserting into ES...\n"
for i in `ls`; do curl -s -XPOST "localhost:9200/_bulk" --data-binary @$i; done;

printf "Done!\n"
