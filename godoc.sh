#!/bin/sh
godoc -html ./server > ./document/server/server.html
godoc -html ./server/config > ./document/server/config.html
godoc -html ./server/controller > ./document/server/controller.html
godoc -html ./server/model > ./document/server/model.html
godoc -html ./server/view > ./document/server/view.html
godoc -html ./server/lib > ./document/server/lib.html
