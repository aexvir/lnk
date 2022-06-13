# lnk

Lnk is a basic url shortener that can be managed via rest api calls

## development

all tasks are managed via [mage](https://magefile.org/), to install it simply run

```shell
go install github.com/magefile/mage@latest
```

then on the root of the project, where the `magefile.go` file is located, you can
run `mage` to discover all the available tasks and their description

### without mage

if mage is not available, or you don't want to install third party dependencies, you can use go run to spin up the server

```shell
go run .
```

## deployment

no deployment strategy is provided at the moment

## api

the lnk api is built with grpc with a rest layer on top of it thanks to [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

### rest

for the rest api, once the service is running just visit `localhost:8000/api/docs` and browse the openapi schema

### grpc

the grpc server is running on the port 9000, and has reflection enabled
using a client like [evans](https://github.com/ktr0731/evans) you can open a repl where all rpc calls should be available

```shell
evans repl -r --host localhost --port 9000
```

## databases

### memory

the only database implemented currently is in-memory; this database is only intended for local development and testing, and it's not recommended for any serious use case

the database is feature complete, but it's process local, so horizontally scaling this service is not possible, as each process will have its own database
