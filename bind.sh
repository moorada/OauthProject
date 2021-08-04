#!/bin/bash

go-bindata -o views/view.go -ignore=view.go views/...

sed -i 's/package main/package views/' views/view.go
