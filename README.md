[![GoDoc](https://godoc.org/github.com/cihangir/gene?status.svg)](https://godoc.org/github.com/cihangir/gene)
[![Build Status](https://travis-ci.org/cihangir/gene.svg)](https://travis-ci.org/cihangir/gene)

# gene

Tired of bootstrapping?

Json-schema based code generation with Go (golang).

## Why Code Generation?

Whenever a bootstrap is required for a project we are hustling with creating the
required folder, files, configs, api function, endpoints, clients, tests etc...

This package aims to ease that pain

## Features

#### Models
* Creating Models from json-schema definitions
* Creating Validations for Models from json-schema definitions
* Creating Constants for Model properties from json-schema definitions
* Creating JSON Tags for Model properties from json-schema definitions
* Adding golint-ed Documentations to the Models, Functions and Exported Variables
* Creating Constructor Functions for Models with their default values from json-schema definitions

#### SQL
* Creating Insert, Update, Delete, Select sql.DB.* compatible plain SQL statements without any reflection

#### Tests
* Providing simple Assert, Ok, Equals test functions for the app

#### Workers

#### API
* Creating rpc api endpoints for Create, Update, Select, Delete operations for every definition in json-schema

#### Client
* Creating Client code for communication with your endpoints

#### CMD
* Creating basic cli for the worker

#### Errors
* Creating idiomatic Go Errors for the api, validations etc.

#### Tests
* Creating tests for the generated api endpoints

## Install

Package itself is not go gettable, get the cli for generating your app
```
go install github.com/cihangir/gene/cmd/gene
```

## Usage

After having gene executable in your path
Pass schema flag for your base json-schema, and target as the existing path for your app

```
gene -schema ./command.json -target ./test/

```

For now, it is generating the following folder/file structure
```
test
├── app
├── models
│   ├── account.go
│   ├── account_statements.go
│   ├── config.go
│   ├── config_statements.go
│   ├── profile.go
│   └── profile_statements.go
├── tests
│   └── testfuncs.go
└── workers
    └── account
        ├── accountapi
        │   ├── account.go
        │   ├── config.go
        │   └── profile.go
        ├── clients
        │   ├── account.go
        │   └── profile.go
        ├── cmd
        │   └── account
        │       └── main.go
        ├── errors
        │   └── account.go
        └── tests
            ├── account_test.go
            ├── common_test.go
            ├── config_test.go
            └── profile_test.go

11 directories, 18 files
```
