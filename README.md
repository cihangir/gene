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
└── workers
    └── command
        ├── cmd
        │   └── command
        │       └── main.go
        ├── commandapi
        │   └── command.go
        ├── errors
        │   └── command.go
        └── tests

9 directories, 4 files
```