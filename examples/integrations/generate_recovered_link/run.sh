#!/bin/bash

cd ./client || exit

if [ ! -f "package-lock.json" ]; then
  npm install
fi

node main.js