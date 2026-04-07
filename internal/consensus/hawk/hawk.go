package hawk

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/chronostech-git/fabrik/internal/consensus"
	"github.com/chronostech-git/fabrik/internal/types"
)

var (
	errMinerAlreadySet = errors.New("miner already set")
	errAlreadyMined    = errors.New("already mined")
)

const (
	MaxFutureBlockTime = 25
)

type Hawk struct {
	Difficulty uint64
	Workers    int
}

// NewPoW returns &Hawk{} with the number of workers based on
// the number of CPU cores usable by the current process (RunPoW)
func NewPoW() *Hawk {
	return &Hawk{
		Workers: runtime.NumCPU(),
	}
}

// SetDifficulty set's the difficulty target for the next block.
// This is used after CalcPoWDifficulty.
func (h *Hawk) SetDifficulty(difficulty uint64) {
	h.Difficulty = difficulty
}

// CalcPoWDifficulty Calculate the next pow difficulty using the previous difficulty params
func CalcPoWDifficulty(
	prevDifficulty uint64,
	prevTimestamp int64,
	currentTimestamp int64,
	targetBlockTime int64,
) uint64 {
	actualTime := time.Duration(currentTimestamp-prevTimestamp) * time.Second

	if actualTime <= 0 {
		return prevDifficulty
	}

	target := targetBlockTime

	ratio := float64(target) / float64(actualTime)

	if ratio > 4 {
		ratio = 4
	}
	if ratio < 0.25 {
		ratio = 0.25
	}

	newDifficulty := uint64(float64(prevDifficulty) * ratio)

	if newDifficulty < 1 {
		newDifficulty = 1
	}

	return newDifficulty
}

// RunPoW Runs the Proof-of-Work consensus (heavily simplified for now)
func (h *Hawk) RunPoW(
	difficultyTarget uint64,
	block *consensus.BlockView,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	target := strings.Repeat("0", int(difficultyTarget))

	var wg sync.WaitGroup
	var resultNonce uint64
	var resultHash types.Hash

	for i := 0; i < h.Workers; i++ {
		wg.Add(1)

		go func(start uint64) {
			defer wg.Done()

			nonce := start

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				block.Nonce = nonce

				// Calculate the hash of a BlockView block
				hash := block.CalcHawkHash()

				// Print hash until the correct prefix hash is calculated
				fmt.Printf("\rNonce: %d | Hash: %s\n", nonce, hash.String())

				// If the hash is calculated (found), stop the RunPoW function loop
				// and set the resultHash to the found hash.
				if strings.HasPrefix(hash.String(), target) {
					resultNonce = nonce
					resultHash = hash
					cancel()
					return
				}

				nonce += uint64(h.Workers)
			}
		}(uint64(i))
	}

	wg.Wait()

	// Finally, when the hash is found and the loop stops
	// set the blocks nonce to the result nonce, and set the hash
	// to the found hash.
	block.Nonce = resultNonce
	block.Hash = resultHash
}
