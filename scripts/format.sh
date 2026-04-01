#!/usr/bin/env bash

set -e

echo "==> Formatting Go source files"
go fmt ./...

echo "==> Running goimports (if installed)"
if command -v goimports >/dev/null 2>&1; then
    goimports -w .
else
    echo "goimports not found. Install it with:"
    echo "go install golang.org/x/tools/cmd/goimports@latest"
fi

echo "==> Tidying Go modules"
go mod tidy

echo "==> Done formatting Fabrik project"
