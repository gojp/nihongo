#!/bin/bash

mkdir -p data
wget http://ftp.monash.edu.au/pub/nihongo/edict2.gz -O data/edict2.gz
gunzip data/edict2.gz
