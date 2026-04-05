# Fabrik

> An experimental blockchain runtime and virtual machine written in Go.

**Status:** Early-stage / experimental — interfaces and commands may change frequently.
For a running log of changes and upcoming work, see [`todos.txt`](./todos.txt).

---

## What is Fabrik?

Fabrik is a research and experimentation platform for understanding the fundamental architecture behind programmable blockchains. It is built from scratch in Go and includes:

- **FVM** — a stack-based virtual machine with its own bytecode, opcodes, and gas metering
- **Fabnet** — a custom peer-to-peer networking layer for node communication
- **Chain** — genesis block creation, coinbase transactions, and persistent chain storage
- **Account** — external and contract account modeling with optional staking
- **Wallet** — key generation and keystore management

---

## Repository Layout

```
cmd/
  chain/     — blockchain execution and genesis CLI
  account/   — account creation and staking CLI
  fvm/       — virtual machine runner CLI
  fabnet/    — P2P server and peer connection CLI
  wallet/    — (in progress) wallet/transaction CLI

internal/
  accounts/        — account types (external, contract)
  blockchain/      — block, chain, genesis, consensus, transactions, wallet
  crypto/          — key generation, signing, address derivation
  fvm/             — compiler, executor, opcodes, gas, stack
  p2p/             — server, peer discovery, dialer, message protocol
  serialize/rlp/   — RLP encoding and decoding
  state/           — account state and chain state management
  storage/         — LevelDB and keystore backends
  types/           — shared primitive types (Address, Amount, etc.)

contracts/
  deposit.fab      — example FVM contract

scripts/
  clean.sh         — remove build artifacts
  format.sh        — format Go source
```

---

## Requirements

| Requirement | Version |
|-------------|---------|
| Go          | 1.25.8+ (as declared in `go.mod`) |
| make        | any     |

### Primary Dependencies

| Package | Version |
|---------|---------|
| `github.com/alexflint/go-arg` | v1.6.1 |
| `github.com/holiman/uint256`  | v1.3.2 |
| `github.com/syndtr/goleveldb` | v1.0.0 |
| `golang.org/x/crypto`         | v0.49.0 |

---

## Build

Install dependencies and build all CLI binaries:

```bash
make
```

This outputs binaries to `./cli/`:

```
cli/
  chain
  account
  fvm
  fabnet
```

**Individual targets:**

```bash
make chain     # build chain binary only
make account   # build account binary only
make fvm       # build fvm binary only
make fabnet    # build fabnet binary only
make deps      # run go mod tidy + download
make clean     # remove ./cli directory
```

---

## CLI Usage

### `chain` — Create and Persist a Genesis Block

Initializes a new chain in the given data directory. Creates a coinbase transaction, signs it with the keystore key, and writes the genesis block to LevelDB.

```bash
./cli/chain <datadir> --gaslimit <gaslimit> [--debug]
```

| Flag          | Required | Description                                 |
|---------------|----------|---------------------------------------------|
| `<datadir>`   | yes      | Directory for keystore, genesis, and LevelDB |
| `--gaslimit`  | yes      | Gas limit for the genesis block             |
| `--debug`     | no       | Print the full chain state after creation   |

**Example:**

```bash
./cli/chain data --gaslimit 1000
```

**Output:**

```
2026/04/05 01:39:52 Coinbase transaction signed and verified
        signer: 0x013bbe45ba4506d5886601d47db697fb823edaf2
        sig:    0x4c2f44dcd24eeb991af6f379c02eae0c84c2e9532fa71a08a5efc91e523b5e272a19...
        hash:   0xac7327f8b30d32e223276f853899fa6df7b2b84eea92c8b45e96f7082e937ac7
        valid:  true

2026/04/05 01:39:52 Genesis block created at 1775367592 unix time with hash 0x5adb616b34ff5d0c...
```

---

### `account` — Create an Account Record

Creates a new account in the given data directory. Supports external (EOA) and contract accounts. Optionally generates a new wallet and runs the staking deposit contract via FVM.

```bash
./cli/account <datadir> --type <external|contract> [--wallet] [--stake <amount>] [--gas <limit>] [--debug]
```

| Flag        | Required | Description                                              |
|-------------|----------|----------------------------------------------------------|
| `<datadir>` | yes      | Data directory for keystore and account state            |
| `--type`    | yes      | Account type: `external` or `contract`                   |
| `--wallet`  | no       | Generate a new wallet keypair in the data directory      |
| `--stake`   | no       | Stake amount (default: `0`)                              |
| `--gas`     | no       | Gas limit for the staking deposit contract execution     |
| `--debug`   | no       | Print FVM debug info during deposit contract execution   |

**Examples:**

```bash
# Create a new external account with a fresh wallet
./cli/account data --type external --wallet

# Create a contract account with a stake
./cli/account data --type contract --stake 32 --gas 100000

# Debug mode: trace FVM execution during staking
./cli/account data --type external --stake 32 --gas 100000 --debug
```

---

### `fvm` — Run the Fabrik Virtual Machine

Executes a `.fab` smart contract file (or raw bytecode) using the FVM. Supports gas metering, debug tracing, and a caller address for contract context.

```bash
./cli/fvm [--file <contract.fab>] [--run <hex-bytecode>] [--caller <address>] [--gas <limit>] [--debug]
```

| Flag       | Required | Default  | Description                              |
|------------|----------|----------|------------------------------------------|
| `--file`   | no*      | —        | Path to a `.fab` source file             |
| `--run`    | no*      | —        | Raw bytecode as a hex string             |
| `--caller` | no       | —        | Sender address for contract context      |
| `--gas`    | no       | `100000` | Gas limit for execution                  |
| `--debug`  | no       | —        | Print step-by-step VM execution trace    |

> \* Provide either `--file` or `--run`.

**Example:**

```bash
./cli/fvm --file contracts/deposit.fab --gas 100000 --debug
```

**Example `.fab` contract** (`contracts/deposit.fab`):

```asm
# deposit.fab — Simple deposit contract for Fabrik VM

START:
    CALLER       # Push sender address onto the stack
    CALLVALUE    # Push the call value onto the stack
    SHA256       # Compute storage slot (hash of address)
    SSTORE       # storage[slot] = value
    SLOAD        # Reload and return the stored value
    STOP
```

---

### `fabnet` — P2P Network Node

Starts a Fabnet node that listens for peer connections. Peers can be dialed at startup using `--connect`. Nodes exchange messages using the Fabnet protocol (e.g., `PING`/`PONG`).

```bash
./cli/fabnet <datadir> --port <port> [--ipaddr <host>] [--connect <peer-addr>] ...
```

| Flag        | Required | Default       | Description                                           |
|-------------|----------|---------------|-------------------------------------------------------|
| `<datadir>` | yes      | —             | Data directory for peer storage                       |
| `--port`    | yes      | —             | Port to bind the server                               |
| `--ipaddr`  | no       | `127.0.0.1`   | Host address to bind the server                       |
| `--connect` | no       | —             | Peer address(es) to connect to (repeatable)           |

**Starting a server (Terminal A):**

```bash
./cli/fabnet data --ipaddr 127.0.0.1 --port 8000
```

```
2026/04/02 16:27:33 [FABNET] Server started on 127.0.0.1:8000
```

**Connecting as a peer (Terminal B):**

```bash
./cli/fabnet data2 --ipaddr 127.0.0.1 --port 8001 --connect 127.0.0.1:8000
```

```
2026/04/02 16:27:39 [FABNET] Peer connected with address 127.0.0.1:33848
```

**Testing the connection:**

Type `PING` in Terminal B:

```
> PING
[From 127.0.0.1:8000] PONG
```

---

## Development Notes

- All data is **datadir-centric** — each node stores its `keystore`, `genesis`, and LevelDB chain in its own directory.
- Module path: `github.com/chronostech-git/fabrik`
- `scripts/clean.sh` is available as an alternative to `make clean`.
- `scripts/format.sh` runs `gofmt` across the source tree.
- Consensus engine, block header gas metrics, and `ChainIterator` were significantly updated on 04/05/2026.

---

## Contributing

Fabrik is open source and welcomes contributions. Please review [`contributing.md`](./contributing.md) before submitting pull requests. Since the project is in an early experimental phase, contributors are encouraged to explore the architecture and help improve core components.

---

*Built by [ChronosTech](https://github.com/chronostech-git)*
