#!/bin/sh

for REPO in $@
do

    echo ${REPO}
    go run cmd/update-alt-label/main.go ${REPO}
    cp export-alt.py ${REPO}/
    cd ${REPO}
    git status --porcelain --untracked-files=all | egrep '.geojson' | awk '{ print $2 }' > new.txt
    python export-alt.py -r . -f new.txt
    rm new.txt
    rm export-alt.py
done
