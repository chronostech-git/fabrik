# Fabrik

> An experimental blockchain runtime and virtual machine written in Go.

**Status:** Early-stage / experimental — interfaces and commands may change frequently.
For a running log of changes and upcoming work, see [`todos.txt`](./todos.txt).

---

## What is Fabrik?

Fabrik is a research and experimentation platform for understanding the fundamental architecture behind programmable blockchains. It is built from scratch in Go and includes:

- **Account model** -- Contract and external accounts 
- **Chain operations** -- Genesis block creation and chain initialization using cli/chain.exe
- **P2P networking** -- Peer to peer networking using fabnet 
- **Virtual machine** -- Virtual machine for smart contract execution and deployment
- **Dual consensus** -- Proof Of Work and Proof Of Importance consensus
- **Smart contracts** -- A stack-based language for smart contract (.fab) files

---

## Building Fabrik

1. Clone Fabrik repository

```bash
$ git clone github.com/chronostech-git/fabrik.git
```

2. Build using GNU make

```bash
# Build all binaries
$ make all 
```

3. Use CLI tools 

>Command line tools will be generated under the 'cli/' directory.

---

### Here is a list of binaries generated with make

* cli/account.exe - Create an account given data directory (--datadir), with the option of --withwallet 
* cli/chain.exe   - Chain operations and genesis creation
* cli/fabkey.exe  - Use your key to sign data, verify signatures and retrieve your wallet address using --whoami
* cli/fabnet.exe  - P2P network operations: start fabnet server, connect to peers, exchange data etc...
* cli/fvm.exe     - Execute programs using either raw hexidecimal bytecode or provide a .fab contract file

---

### Command Line Tools

### Generate an account with a wallet
```bash
$ ./cli/account --datadir ./data --withwallet
```

Output
```
2026/04/06 21:56:52 An external account with a balance of 0 was created
2026/04/06 21:56:52 A new fabkey has been created and was used to create your account
2026/04/06 21:56:52 Wallet address: 0xa360fd49f0d7288ee6979bfc5c36d13e9160f323
```

### Start the chain by creating a genesis block
```bash
$ ./cli/chain --datadir ./data --mechanism hawk-pow [--usememory]
```

Output
```
Nonce: 6 | Hash: 0x438d209238bf6875e33d9945641a57a2019dfe06ae9c4843bdada52abccdc722
Nonce: 7 | Hash: 0x933da91b017d7be1afd5fa511962eea3f46b6f077e24cf4717449e6dd1066ff9
Nonce: 0 | Hash: 0x438d209238bf6875e33d9945641a57a2019dfe06ae9c4843bdada52abccdc722
Nonce: 1 | Hash: 0x933da91b017d7be1afd5fa511962eea3f46b6f077e24cf4717449e6dd1066ff9
Nonce: 2 | Hash: 0x933da91b017d7be1afd5fa511962eea3f46b6f077e24cf4717449e6dd1066ff9
....

Genesis mining complete.
```

### Use your key to sign data and verify signature
```bash
./cli/fabkey --keydir ./data --signdata "Hello, world!" [--verifysig]
```

Output
```
Signed: 0xd3f13a1124f9...a07e9d7f2ff3a3ade71d
```

### Start a fabnet server and connect to a peer
```bash
# Terminal (A) - start fabnet server
$ ./cli/fabnet --datadir ./data --ipaddr 127.0.0.1 --host 8000

# Terminal (B) - create a peer and connect to server
$ ./cli/fabnet --datadir ./data --ipaddr 127.0.0.1 -host 8001 --connect 127.0.0.1:8000
```

Output (terminal A)
```
2026/04/07 08:15:34 [FABNET] Server started on 127.0.0.1:8000
2026/04/07 08:16:18 [FABNET] Peer connected with address 127.0.0.1:64905
```

> Any discovered peers will be saved to <datadir>/peers.json

### Use FVM (Fabrik Virtual Machine)

```bash
# Execute smart contract from file
$ ./cli/fvm --file example.fab [--debug]

# Execute smart contract from hex
$ ./cli/fvm --run <hex> [--debug]
```

---

## Contributing

Fabrik is open source and welcomes contributions. Please review [`contributing.md`](./contributing.md) before submitting pull requests. Since the project is in an early experimental phase, contributors are encouraged to explore the architecture and help improve core components.

---

*Built by [ChronosTech](https://github.com/chronostech-git)*
