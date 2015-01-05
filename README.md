[![GoDoc](https://godoc.org/github.com/cihangir/gene?status.svg)](https://godoc.org/github.com/cihangir/gene)
[![Build Status](https://travis-ci.org/cihangir/gene.svg)](https://travis-ci.org/cihangir/gene)

gene
====

Tired of bootstrapping?

```
go install github.com/cihangir/gene/cmd/gene && ./bin/gene -schema ./command.json -target ./test/
```

For now
```
test
├── app
├── models
│   └── command.go
├── tests
│   └── testfuncs.go
└── workers
    └── command
        ├── clients
        │   └── command.go
        ├── cmd
        │   └── command
        │       └── main.go
        ├── commandapi
        │   └── command.go
        ├── errors
        │   └── command.go
        └── tests
            ├── command_test.go
            └── common_test.go

11 directories, 8 files
```
