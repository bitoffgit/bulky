#! /bin/bash
set -e

if ! go version &> /dev/null
then
    echo "Golang could not be found in you machine"
    exit
fi

go build -o bin/bulky main.go

BULKY_PATH=/usr/local/bin/bulky

if [ -f "$BULKY_PATH" ] ; then
    sudo rm "$BULKY_PATH"
fi

sudo ln -s "$(pwd)/bin/bulky" "$BULKY_PATH"
