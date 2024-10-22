#!/bin/bash

go clean -testcache
for D in $(ls -1d sd*); do
  go test "github.com/gaorx/stardust6/$D"
  if [ "$D" = "sdfile" ]; then
    go test "github.com/gaorx/stardust6/$D/sdfiletype"
    go test "github.com/gaorx/stardust6/$D/sdhttpfile"
  elif [ "$D" = "sdcache" ]; then
    go test "github.com/gaorx/stardust6/$D/sdcacheristretto"
  elif [ "$D" = "sdcodegen" ]; then
    go test "github.com/gaorx/stardust6/$D/sdgengo"
  elif [ "$D" = "sderr" ]; then
    go test "github.com/gaorx/stardust6/$D/sdnotfounderr"
  fi
done

