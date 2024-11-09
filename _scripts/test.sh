#!/bin/bash

PKG=github.com/gaorx/stardust6

go test $PKG/sderr
go test $PKG/sdslog
go test $PKG/sdbytes
go test $PKG/sdrand
go test $PKG/sdreflect
go test $PKG/sdtextenc
go test $PKG/sdfile
go test $PKG/sdfile/sdfiletype
go test $PKG/sdfile/sdhttpfile
go test $PKG/sdcodegen
go test $PKG/sdcodegen/sdgogen
go test $PKG/sdjson
go test $PKG/sdresty