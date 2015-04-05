#!/bin/bash

wget http://ftp.monash.edu.au/pub/nihongo/edict2.gz -O edict2.gz
gunzip edict2.gz
iconv -f EUC-JP -t UTF-8 < edict2 > edict2.utf-8
rm edict2
python edict2_parser.py > edict2.json
gzip edict2.json
