nihongo.io
=========

Open source Japanese Dictionary written in Go

### How to run:
1. `git clone https://github.com/gojp/nihongo.git`
2. Install dependencies: `go get ./...`
3. Populate the database (replace `MONGO_URI` and `PATH_TO_EDICT2` variables accordingly): `python populate.py`
4. Export `MONGOHQ_URL`: `export MONGOHQ_URL="mongodb://..."`
5. Run the app: `revel run github.com/gojp/nihongo
