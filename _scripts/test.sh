#!/bin/bash

PKG=github.com/gaorx/stardust6

go test $PKG/sderr
go test $PKG/sderr/sdnotfounderr
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
go test $PKG/sdcompress/sdgzip
go test $PKG/sdcompress/sdzip
go test $PKG/sdcompress/sdlz4
go test $PKG/sdhash/sdmd5
go test $PKG/sdhash/sdsha256
go test $PKG/sdhash/sdsha512
go test $PKG/sdjwt
go test $PKG/sdload
go test $PKG/sdparse
go test $PKG/sdtime
go test $PKG/sdreflect
go test $PKG/sdtemplate
go test $PKG/sdregexp
go test $PKG/sdstrings
go test $PKG/sdlocal
go test $PKG/sdsnowflake
go test $PKG/sdsecurity/sdauthn
go test $PKG/sdwebapp
go test $PKG/sdredis
go test $PKG/sdcache
go test $PKG/sdcache/sdcacheredis
go test $PKG/sdcache/sdcacheristretto
go test $PKG/sdcache/internal/doubletest
go test $PKG/sdsql
go test $PKG/sdbun/sdbunsqlite
go test $PKG/sdbun/sdbunmysql
