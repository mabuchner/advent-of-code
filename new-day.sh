#!/usr/bin/env bash
cargo new "day-$1" &&\
rm "./day-$1/src/main.rs" &&\
mkdir "./day-$1/src/bin" &&\
cp "./day-template/src/bin/part1.rs" "./day-$1/src/bin/." &&\
touch "./day-$1/src/bin/input.txt"
cp "./day-template/.gitignore" "./day-$1/."
