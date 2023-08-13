package blockchain

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	EarliestBlockNumber = BlockNumber(0)
	InvalidBlockNumber  = BlockNumber(-1)
)

type BlockNumber int64

func (bn BlockNumber) ToInt64() int64 {
	return (int64)(bn)
}

func (bn BlockNumber) ToHex() string {
	hexStr := fmt.Sprintf("0x%x", bn.ToInt64())
	return hexStr
}

func (bn BlockNumber) Inc() BlockNumber {
	return bn + 1
}

// ---- builder methods ----

type BlockNumberBuilder struct {
	value *BlockNumber
}

func NewBlockNumberBuilder() *BlockNumberBuilder {
	return &BlockNumberBuilder{}
}

func (b *BlockNumberBuilder) FromInt64(i int64) *BlockNumberBuilder {
	bn := BlockNumber(i)
	b.value = &bn
	return b
}

// FromHexString parses a block number from a hex string.
func (b *BlockNumberBuilder) FromHexString(hexStr string) *BlockNumberBuilder {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Using ParseInt to convert the hex string to int64
	decimal, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		log.Println("Error parsing hex string to int64")
	}
	blockNumber := BlockNumber(decimal)
	b.value = &blockNumber
	return b
}

func (b *BlockNumberBuilder) Value() BlockNumber {
	if b.value == nil {
		return BlockNumber(0)
	}
	return *b.value
}

func (b *BlockNumberBuilder) Pointer() *BlockNumber {
	return b.value
}
