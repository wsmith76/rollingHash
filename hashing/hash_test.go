package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New_Hashes(t *testing.T) {
	h := New_Hashes()
	assert.NotNil(t, h)
	assert.NotPanics(t, func() { h.ByPosition[uint32(3)] = uint64(43) })
	assert.NotPanics(t, func() { h.ByValue[uint64(3)] = uint32(43) })
}

func Test_UInt64Pow(t *testing.T) {
	assert.Equal(t, uInt64Pow(uint64(2), uint64(3)), uint64(8))
	assert.Equal(t, uInt64Pow(uint64(7), uint64(0)), uint64(1))
	assert.Equal(t, uInt64Pow(uint64(7), uint64(1)), uint64(7))
}

func Test_GetHash(t *testing.T) {
	testString := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	testBytes := []byte(testString)
	assert.Equal(t, uint64(2013743), getHash(testBytes[0:5], 13))
	assert.Equal(t, uint64(125257871), getHash(testBytes[0:5], 37))
	assert.Equal(t, uint64(2044684), getHash(testBytes[1:6], 13))
	assert.Equal(t, uint64(127184092), getHash(testBytes[1:6], 37))
	assert.Equal(t, uint64(2075625), getHash(testBytes[2:7], 13))
	assert.Equal(t, uint64(129110313), getHash(testBytes[2:7], 37))
}

func Test_GetRollingHash(t *testing.T) {
	testString := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	testBytes := []byte(testString)
	rollingHash13 := GetRollingHash(testBytes, 5, 13)
	rollingHash37 := GetRollingHash(testBytes, 5, 37)
	assert.Equal(t, uint64(2044684), rollingHash13.ByPosition[1])
	assert.Equal(t, uint64(127184092), rollingHash37.ByPosition[1])
	assert.Equal(t, len(rollingHash13.ByValue), len(rollingHash13.ByPosition))
	assert.Equal(t, len(rollingHash37.ByValue), len(rollingHash37.ByPosition))
}

func Test_divideIntoChunks(t *testing.T) {
	testString := "Here comes the candle to light you to bed.  Here comes the chopper to chop off your head."
	testBytes := []byte(testString)
	assert.Equal(t, len(divideIntoChunks(testBytes, 5)), 18)
	assert.Equal(t, len(divideIntoChunks(testBytes, 5)[16]), 5)
	assert.Equal(t, len(divideIntoChunks(testBytes, 5)[17]), 4)
	assert.Equal(t, len(divideIntoChunks(testBytes, 256)), 1)
	assert.Equal(t, divideIntoChunks(testBytes, 256)[0], testBytes)
}

func Test_GetFixedHash(t *testing.T) {
	testString := "Here comes the candle to light you to bed.  Here comes the chopper to chop off your head."
	testBytes := []byte(testString)
	fixedHash := GetFixedHash(testBytes, 10, 13)
	assert.Equal(t, 9, len(fixedHash.ByPosition))
	assert.Equal(t, len(fixedHash.ByPosition), len(fixedHash.ByValue))
}

func Test_Integration_RollingHash_FixedHash(t *testing.T) {
	testString := "Here comes the candle to light you to bed.  Here comes the chopper to chop off your head."
	testBytes := []byte(testString)
	fixedHash := GetFixedHash(testBytes, 10, 13)
	rollingHash := GetRollingHash(testBytes, 10, 13)
	assert.Equal(t, fixedHash.ByPosition[10], rollingHash.ByPosition[10])
	assert.Equal(t, fixedHash.ByPosition[20], rollingHash.ByPosition[20])
}
