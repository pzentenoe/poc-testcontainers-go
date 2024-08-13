# poc-testcontainers-go

**poc-testcontainers-go** is a proof-of-concept (POC) project designed to demonstrate how to use Testcontainers in Go to manage and test interactions with Docker containers. This project specifically showcases how to interact with an SFTP server running inside a Docker container, using tools like Ginkgo and Gomega for behavior-driven development (BDD) testing.

## Overview

This project includes:

* A Go-based implementation to start an SFTP server in a Docker container.
* Automated tests to upload and verify files using an SFTP client.
* Usage of Ginkgo and Gomega for BDD-style tests.
* Testcontainers to simplify integration testing with Docker.

## Dependencies
To run this project, you will need to install the following dependencies:

* Go: The Go programming language. Ensure you have Go installed and set up properly. You can download it from golang.org.
* Testcontainers-Go: A Go package for managing Docker containers in tests.

```bash
go get github.com/testcontainers/testcontainers-go
```
* Ginkgo: A BDD-style testing framework for Go.
```bash
go get github.com/onsi/ginkgo/v2 
```
* Gomega: An assertion library that integrates with Ginkgo.
```bash
go get github.com/onsi/gomega/... 
```

* SFTP: A Go package for SSH file transfer protocol (SFTP) interactions.
```bash
go get github.com/pkg/sftp 
```

* SSH: A Go package for working with SSH.
```bash 
go get golang.org/x/crypto/ssh
```

* Testify (optional): Provides additional testing utilities, though Ginkgo and Gomega are the primary tools used.
```bash
go get github.com/stretchr/testify 
```

## Installation
```bash
git clone https://github.com/pzentenoe/poc-testcontainers-go.git

cd poc-testcontainers-go
```

## Running Tests

```bash
ginkgo -v ./...
```

## Author
- **Pablo Zenteno** - _Full Stack Developer_ - [pzentenoe](https://github.com/pzentenoe)