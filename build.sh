#!/bin/bash

mkdir bungkus
go build -buildmode=plugin -ldflags="-s -w" -o bungkus/contact.so
cp -Rvf LICENSE CHANGELOG  module.toml schema bungkus
mv bungkus contact
tar zcvvf contact.tar.gz contact
rm -Rvf contact
