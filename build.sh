#!/bin/bash

# CGO_ENABLED flag is needed to disable dynamic linking,
# otherwise the built binary doesn't run on alpine linux.
CGO_ENABLED=0 go build -o build/vault-snapshot -v ./...
