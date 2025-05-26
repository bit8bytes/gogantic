# Publish a module

This document contains the steps to publish a module.

First of all, run `go mod tidy` and `go test ./...`.

After this, follow the follwing steps:

1. `git add .` to add all new changes.
2. Select the version from GitHub if it is a minor change e.g. v0.0.1 will be v0.0.2
3. `git commit -m "<module>: <changes> for <version>"`
4. `git tag <version>`
5. `git push origin <version>`
6. `git push`
