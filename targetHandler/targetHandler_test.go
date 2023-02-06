package targetHandler

import (
	"github.com/stretchr/testify/assert"
	"rollingHash/hashing"
	"testing"
)

func TestNew_Diff(t *testing.T) {
	d := New_Diff()
	assert.NotNil(t, d)
	assert.NotPanics(t, func() { d.Matches[uint32(7)] = uint32(12) })
	assert.NotPanics(t, func() { d.NewChunks[uint32(7)] = []byte{} })
}

func TestNew_TargetData(t *testing.T) {
	testString := "ABCDEFGHIJKLMNOPQUSTUVWXYZ"
	testBytes := []byte(testString)
	td := New_TargetData(testBytes, 10, 5)
	assert.Equal(t, 27, len(td.Data))
	assert.Equal(t, 16, len(td.Base13HashLarge.ByPosition))
	assert.Equal(t, 21, len(td.Base13HashSmall.ByPosition))
}

func Test_findHashInSource(t *testing.T) {
	var sourceDataHash hashing.Hashes
	var targetDataHashSmall hashing.Hashes
	var targetDataHashLarge hashing.Hashes
	sourceDataHash.ByValue[uint64(12345)] = uint32(0)
	sourceDataHash.ByValue[uint64(65432)] = uint32(256)
	sourceDataHash.ByValue[uint64(666665)] = uint32(512)
	sourceDataHash.ByPosition[uint32(0)] = uint64(12345)
	sourceDataHash.ByPosition[uint32(256)] = uint64(65432)
	sourceDataHash.ByPosition[uint32(512)] = uint64(666665)
	sourceDataSize := uint32(600)
	sourceDataChunkSize := 256
	for i := uint32(0); i < uint32(600); i++ {
		targetDataHashSmall.ByPosition[i] = uint64(999999)
		targetDataHashSmall.ByValue[uint64(999999)] = i
		targetDataHashLarge.ByPosition[i] = uint64(999999)
		targetDataHashLarge.ByValue[uint64(999999)] = i
	}
	targetDataHashLarge.ByPosition[uint32(75)] = uint64(12345)
	targetDataHashLarge.ByValue[uint64(12345)] = uint32(75)
	targetDataHashSmall.ByPosition[uint32(270)] = uint64(666665)
	targetDataHashSmall.ByValue[uint64(666665)] = uint32(270)
	hashResult := findHashInSource(sourceDataHash, sourceDataSize, sourceDataChunkSize, targetDataHashSmall, targetDataHashLarge)
	assert.Equal(t, uint32(75), hashResult[0])
	assert.Equal(t, uint32(270), hashResult[512])
	assert.Equal(t, nil, hashResult[256])
}
