# Fabrik

Fabrik is a minimal, experimental blockchain implementation written in Go. It includes a custom state system, transaction processing, and a simple virtual machine (FVM) for future smart contract execution.

---

## ⚙️ Features

* Custom blockchain implementation
* Genesis block generation
* Wallet + key management
* In-memory or persistent node execution
* Early-stage virtual machine (FVM)

---

## 📦 Build

Fabrik uses a simple `Makefile` to build all binaries.

```bash
make
```

This will generate the following binaries inside `fabrik-cli/`:

* `genesis`
* `wallet`
* `node`
* `fvm` (unfinished)

To clean build artifacts:

```bash
make clean
```

---

## 🚀 Usage

### 1. Create a Wallet

```bash
fabrik-cli/wallet --datadir <datadir-name>
```

* Generates a new keypair
* Stores it in `<datadir-name>/keystore`
* Prints the address

---

### 2. Create Genesis

```bash
fabrik-cli/genesis --datadir <datadir-name> --extra <extra_data>
```

* Uses your wallet key
* Creates a genesis file at:

  ```
  <datadir-name>/genesis/genesis.dat
  ```
* Allocates initial balance to the coinbase address

---

### 3. Run Node

```bash
fabrik-cli/node --datadir <datadir-name> --usememory
```

* Loads the genesis file
* Starts the chain
* Applies a test block

Flags:

* `--datadir` → required
* `--usememory` → optional (boolean)

  * `true` → in-memory storage
  * `false` → persistent storage (default behavior)

---

### 4. FVM (Experimental)

```bash
fabrik-cli/fvm
```

* Currently unfinished
* Intended for executing smart contract bytecode

---

## 📁 Project Structure

```
cmd/
  genesis/   → genesis creation CLI
  wallet/    → wallet + key management
  node/      → blockchain node
  fvm/       → virtual machine (WIP)

internal/
  blockchain/
  state/
  crypto/
  serialize/
  storage/
  types/
```

---

## 🧠 Design Notes

* Uses custom RLP-like encoding for serialization
* State is managed via a simple balance map
* Addresses are derived from ECDSA public keys
* VM is stack-based and under active development

---

## ⚠️ Status

This project is **experimental and not production-ready**.

Missing features include:

* Networking (P2P)
* Consensus (beyond basic PoW scaffolding)
* Full smart contract support
* Persistent state robustness
* Security hardening

---

## 🛠️ Development

To rebuild specific components:

```bash
make wallet
make genesis
make node
make fvm
```

---

## 📌 Notes

* Always generate a wallet before creating genesis
* Genesis must exist before running the node
* Data is stored relative to the provided `--datadir`

---

## 📜 License

MIT (or specify if different)

---

## 👨‍💻 Author

ChronosTech
