package blockchain

// all blockchain related constants go here
type CHAIN_TYPE int

// different chains for different purposes
const (
	MAIN_CHAIN CHAIN_TYPE = iota
	TEST_CHAIN
)

const TIMESTAMP_SLICE_SIZE = 16 // in case we use fixed size slices to represent timestamps instead of strings to save space
