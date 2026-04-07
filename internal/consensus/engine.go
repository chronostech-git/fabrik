package consensus

// Engine interface holds the functions used in both PoW
// and PoI.
type Engine interface {
	// RunPow run hawk PoW consensus mechanism
	RunPoW(
		difficultyTarget uint64,
		block *BlockView,
	)

	// RunPoI run falcon PoI consensus mechanism
	// RunPoI(...)
}
