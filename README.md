# Note Server

![Linter](https://github.com/sksmith/note-server/actions/workflows/lint.yml/badge.svg) ![Security](https://github.com/sksmith/note-server/actions/workflows/sec.yml/badge.svg) ![Test](https://github.com/sksmith/note-server/actions/workflows/test.yml/badge.svg)

A basic application that saves notes to and retrieves notes from an aws s3 storage.

## Running the Application Locally

This project comes with everything you need except an s3 storage. You'll need to set
that up before running the project. Once that's done just execute:

```shell
make run
```

If you want to create a deployable executable and run it:

```shell
make build
./bin/note-server
```
