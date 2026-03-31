# Fabrik

Fabrik is an experimental blockchain project written in Go. It includes foundational components for wallets, genesis generation, chain persistence, account modeling, and an early virtual machine (FVM).

> **Status:** early-stage / experimental. This repository is not production-ready.

## What this project currently includes

- Wallet generation and key storage (`cmd/wallet`)
- Genesis file creation (`cmd/genesis`)
- A chain runner that loads genesis and writes blocks (`cmd/chain`)
- Basic account creation flow for external and contract account types (`cmd/account`)
- A prototype VM pipeline (parse → compile → execute) via `cmd/fvm`
- Internal packages for blockchain, state, storage, serialization, cryptography, networking scaffolding, and types

## Repository layout

```text
cmd/
  account/   # account CLI
  chain/     # chain execution CLI (loads genesis, appends test block)
  fvm/       # virtual machine runner CLI
  genesis/   # genesis generation CLI
  node/      # placeholder CLI (currently no logic)
  wallet/    # wallet/key generation CLI

internal/
  accounts/
  blockchain/
  crypto/
  fvm/
  network/
  serialize/rlp/
  state/
  storage/
  types/

scripts/
  clean.sh
```

## Requirements

- Go **1.25.8** (as declared in `go.mod`)
- `make`

## Dependencies

Primary module dependencies:

- `github.com/alexflint/go-arg`
- `github.com/holiman/uint256`
- `github.com/syndtr/goleveldb`
- `golang.org/x/crypto`

## Build

Build all CLIs:

```bash
make
```

This creates binaries in `./cli`:

- `cli/genesis`
- `cli/wallet`
- `cli/node`
- `cli/fvm`
- `cli/account`

Other available make targets:

```bash
make genesis
make wallet
make node
make chain
make fvm
make account
make clean
```

## CLI usage

### 1) Create a wallet

```bash
./cli/wallet <datadir>
```

Behavior:
- Ensures `<datadir>` exists
- Creates/loads a key in `<datadir>/keystore`
- Prints public key information and a private-key warning

---

### 2) Create genesis

```bash
./cli/genesis <datadir> --data "optional extra data"
```

Behavior:
- Loads key material from the datadir keystore
- Builds a genesis transaction allocating initial funds to coinbase
- Writes genesis file to `<datadir>/genesis/genesis.dat`

---

### 3) Run chain

```bash
./cli/chain <datadir> [--use-memory] [--debug] [--dump]
```

Behavior:
- Loads genesis from `<datadir>`
- Initializes chain storage (currently opens LevelDB path `<datadir>/manifest`)
- Adds a test block containing a sample transaction
- Flushes cached blocks to disk
- Pretty-prints chain output when `--dump` is set

> Note: `--use-memory` is present, but current control flow still proceeds to initialize LevelDB.

---

### 4) Create an account record

```bash
./cli/account <datadir> --type <external|contract>
```

Behavior:
- Loads wallet key from datadir keystore
- Creates either an external or contract account with that address
- Adds it to in-memory account state

---

### 5) Run FVM (prototype)

```bash
./cli/fvm --file <contract.fab> [--gas 100000] [--debug]
```

Behavior:
- Parses `.fab` input
- Compiles to bytecode
- Executes bytecode on FVM with configured gas limit
- Optional debug output for stack and remaining gas

---

### 6) `node` command status

`cmd/node` currently contains an empty `main()` and does not perform node operations yet.

## Typical local flow

```bash
make
./cli/wallet --datadir ./data
./cli/genesis --datadir ./data --data "local dev genesis"
./cli/chain --datadir ./data --dump
./cli/account --datadir ./data --type external
```

## Known limitations

- No production networking/peer sync behavior exposed via CLI
- `node` command is currently a placeholder
- FVM is early-stage and intended for experimentation
- Limited hardening and validation for production security/performance

## Development notes

- Data is datadir-centric (`keystore`, `genesis`, and chain database paths)
- `scripts/clean.sh` exists in addition to `make clean`
- Project module path: `github.com/chronostech-git/fabrik`

