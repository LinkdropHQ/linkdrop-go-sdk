#!/bin/bash

cd ./client || exit

if [ ! -f "package-lock.json" ]; then
  npm install
fi

# Generating link key on the client
link_key=$(node main.js)
echo "Client-side link key: $link_key"
cd ../server || exit

go run main.go "$link_key"

