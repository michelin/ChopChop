#!/bin/bash

go build -o ../bin/chopchop ../cmd/main.go
robot --pythonpath libraries/ --outputdir out/ tests/
rm -r ../bin
