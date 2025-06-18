#!/usr/bin/env bash
cargo new day-08 &&\
rm ./day-08/src/main.rs &&\
mkdir ./day-08/src/bin &&\
cp ./day-template/src/bin/part1.rs ./day-08/src/bin/. &&\
cp ./day-template/.gitignore ./day-08/.
