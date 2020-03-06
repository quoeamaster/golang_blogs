#!/usr/bin/env bash

# generate the assets first (refresh)
go run generateAssets.go

# after the above step, should have the udpated "assets.go" created
# rename the generateAssets.go to "disable"
mv generateAssets.go generateAssets.go.disable

# build it
go build -o tasksApp

# revert the rename
mv generateAssets.go.disable generateAssets.go

# run
./tasksApp
