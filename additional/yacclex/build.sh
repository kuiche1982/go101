#!/bin/sh

# create parser
for y in *.y
do
    goyacc -o ${y%.y}.go -p Calc $y
done

# build go
go build