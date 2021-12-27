# Note Server

![Linter](https://github.com/sksmith/note-server/actions/workflows/lint.yml/badge.svg)
![Security](https://github.com/sksmith/note-server/actions/workflows/sec.yml/badge.svg)
![Test](https://github.com/sksmith/note-server/actions/workflows/test.yml/badge.svg)

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

For doing local development, you'll want linting, and security tooling. Run this to install them.

```shell
make tools
```

To run the linter, security check, and tests:

```shell
make check
```

To just run the tests alone:

```shell
make test
```

## Setting Up AWS

I have quite a bit in this project automated. Linting, security, and unit testing all execute as
soon as the master branch gets merged. A deployment to AWS gets kicked off and the application runs
out in the wild. However, I still needed to setup quite a bit. Here's the list:

- **Elastic Container Registry** - this is where the application images get stored
  - A new repository needed created here
- **Elastic Container Service** 
  - A cluster needed created
  - A task definition needed defined
  - A new service needed created using that task definition
  - When creating the service, added a load balancer
  - Assigned a certificate to the load balancer
- **LogWatch** - this handles watching the tasks for their logs and gathering them up
  - A new **loggroup** needed created
