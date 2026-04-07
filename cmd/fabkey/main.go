package main

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	KeyDir    string // Data directory where the keystore is located
	WhoAmI    bool   // Prints the users key address if specified
	SignData  string // Sign data with private key
	VerifySig bool   // Verify given signature
}

func printAddress(ks *keystore.FileStore) {
	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}
	addr := key.Address.String()
	fmt.Println("Your address:", addr)
}

func signDataWithKey(ks *keystore.FileStore, data string) *crypto.Signature {
	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}

	hashedData := sha256.Sum256([]byte(data))

	sig, err := key.Sign(hashedData)
	if err != nil {
		log.Panic(err)
	}

	return sig
}

func verifySignature(
	ks *keystore.FileStore,
	dataHash types.Hash,
	sig *crypto.Signature,
) bool {
	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}

	valid := key.Verify(dataHash, sig)
	return valid
}

func main() {
	arg.MustParse(&args)

	ks := keystore.NewFileStore(args.KeyDir)

	if args.WhoAmI {
		printAddress(ks)
	}

	sig := signDataWithKey(ks, args.SignData)
	fmt.Printf("Signed: %s\n", sig.Hex())

	if args.VerifySig {
		dataHash := sha256.Sum256([]byte(args.SignData))
		valid := verifySignature(ks, dataHash, sig)
		fmt.Printf("%t\n", valid)
	}
}
