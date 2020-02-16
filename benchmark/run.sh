#!/bin/sh

go test -bench . -test.cpu 2,4,8,12 -test.benchmem
