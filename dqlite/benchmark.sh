export CGO_LDFLAGS_ALLOW="-Wl,-z,now"
go build -tags libsqlite3
./seguro-dqlite --db 127.0.0.1:9000 &
./seguro-dqlite --db 127.0.0.1:9001 --join 127.0.0.1:9000 &
./seguro-dqlite --db 127.0.0.1:9002 --join 127.0.0.1:9000 &