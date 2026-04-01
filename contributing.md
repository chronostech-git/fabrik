# Contributing to Fabrik

Thank you for your interest in contributing to **Fabrik**. Contributions help improve the project and make the virtual machine, tooling, and blockchain components more robust.

This document outlines how to propose changes, report issues, and submit code.

---

# Project Overview

Fabrik is a blockchain runtime environment written in Go that includes:

* **FVM** — a stack-based virtual machine
* **Assembler / compiler** for `.fab` smart contract files
* **CLI tools** for executing contracts
* **State management** for accounts, storage, and memory

Contributions may include:

* Bug fixes
* Performance improvements
* VM opcode implementations
* Compiler improvements
* Documentation
* Tests
* CLI tooling

---

# Getting Started

## 1. Fork the repository

Create your own fork and clone it locally:

```bash
git clone https://github.com/YOUR_USERNAME/fabrik.git
cd fabrik
```

Add the upstream repository:

```bash
git remote add upstream https://github.com/chronostech-git/fabrik.git
```

---

## 2. Install dependencies

Fabrik is written in **Go**, so ensure Go is installed:

```bash
go version
```

Install dependencies:

```bash
go mod tidy
```

---

## 3. Build the CLI

```bash
go build -o cli/fvm ./cli/fvm
```

Example usage:

```bash
cli/fvm --file contracts/deposit.fab --debug
```

or

```bash
cli/fvm --run <hex_bytecode> --debug
```

---

# Development Guidelines

## Code Style

Follow standard Go conventions:

* Use `gofmt` or `go fmt` before committing.
* Prefer clear naming over short names.
* Keep functions small and focused.
* Avoid unnecessary abstractions.

Format your code:

```bash
go fmt ./...
```

---

## Project Structure

```
fabrik/
│
├── cli/                CLI applications
├── internal/
│   ├── fvm/            Virtual machine implementation
│   ├── state/          Account and storage state
│   └── ...
│
├── contracts/          Example smart contracts (.fab)
└── docs/               Documentation
```

Key areas for contributions:

| Area             | Description                                 |
| ---------------- | ------------------------------------------- |
| `internal/fvm`   | VM opcodes, interpreter loop, stack, memory |
| `internal/state` | Account/state management                    |
| `cli`            | Command line tools                          |
| `contracts`      | Example smart contracts                     |
| `docs`           | Documentation                               |

---

# Reporting Issues

If you discover a bug or unexpected behavior:

1. Check existing issues first.
2. Open a new issue with:

   * Description of the problem
   * Steps to reproduce
   * Expected vs actual behavior
   * Example contract or bytecode
   * Environment details (OS, Go version)

Example:

```
Bug: MSTORE writes incorrect memory value

Steps:
1. Run contract X
2. Execute instruction Y

Expected:
Memory slot 0 contains 42

Actual:
Memory slot 0 contains 0
```

---

# Submitting Changes

## 1. Create a feature branch

```bash
git checkout -b feature/my-feature
```

Examples:

```
feature/new-opcode
fix/memory-overflow
improve/disassembler
```

---

## 2. Commit changes

Write clear commit messages:

```
fvm: implement SHA256 opcode

Adds SHA256 hashing opcode to the VM dispatch table
and corresponding stack handling.
```

---

## 3. Push your branch

```bash
git push origin feature/my-feature
```

---

## 4. Open a Pull Request

Submit a pull request against the `main` branch.

Include:

* Description of changes
* Reason for the change
* Related issue (if applicable)

---

# Testing

Before submitting a PR, verify:

```
go build ./...
```

Run example contracts:

```
cli/fvm --file contracts/deposit.fab --debug
```

Ensure:

* VM executes correctly
* Gas accounting works
* No unexpected stack errors occur

If your change affects the VM, include **example bytecode or contracts** demonstrating the change.

---

# Smart Contract Contributions

When contributing `.fab` contracts:

* Keep contracts **small and readable**
* Include comments explaining behavior
* Provide the compiled bytecode if relevant
* Ensure contracts run with the CLI

Example:

```
; simple addition contract

PUSH 10
PUSH 20
ADD
STOP
```

---

# Areas That Need Contributions

Contributions are especially welcome in:

* VM opcode coverage
* Contract debugging tools
* Gas accounting improvements
* Contract ABI / function dispatch
* Better CLI developer tooling
* Documentation

---

# Security

If you discover a **security vulnerability**, please do **not** open a public issue immediately. Instead contact the maintainers privately so the issue can be addressed responsibly.

---

# License

By contributing to Fabrik, you agree that your contributions will be licensed under the same license used by the project.

---

Thank you for helping improve Fabrik.
