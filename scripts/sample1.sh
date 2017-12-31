#!/bin/bash -e

./bin/main -t samples/1 -k samples/mikan -s samples/1.synonyms \
    -v --match-length=0 --confirm=0 \
    -w 4 -m 30 --one --code --parallel=0 \
    --word-pattern=5 --only-keywords=1 \
    --all-word-length=1 --progress --kanji=1 $@

