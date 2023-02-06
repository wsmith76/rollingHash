package targetHandler

import (
	"rollingHash/hashing"
	"rollingHash/sourceHandler"
	"sync"
)

type TargetData struct {
	Data            []byte
	ChunkSize       int
	LastChunkSize   int
	Base13HashLarge hashing.Hashes
	Base13HashSmall hashing.Hashes
	Base53HashLarge hashing.Hashes
	Base53HashSmall hashing.Hashes
}

type Diff struct {
	ChunkSize     int
	LastChunkSize int
	Matches       map[uint32]uint32
	NewChunks     map[uint32][]byte
}

func New_Diff() Diff {
	var df Diff
	df.Matches = make(map[uint32]uint32)
	df.NewChunks = make(map[uint32][]byte)
	return df
}

func New_TargetData(dataIn []byte, chunkSize int, lastChunkSize int) TargetData {
	var td TargetData
	td.Data = append(dataIn, byte(0))
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		td.Base13HashLarge = hashing.GetRollingHash(dataIn, chunkSize, 13)
	}()
	go func() {
		defer wg.Done()
		td.Base13HashSmall = hashing.GetRollingHash(dataIn, lastChunkSize, 13)
	}()
	go func() {
		defer wg.Done()
		td.Base53HashLarge = hashing.GetRollingHash(dataIn, chunkSize, 53)
	}()
	go func() {
		defer wg.Done()
		td.Base53HashSmall = hashing.GetRollingHash(dataIn, lastChunkSize, 53)
	}()
	wg.Wait()
	return td
}

func (td *TargetData) DiffWithoutNewData(sourceData sourceHandler.SourcePacket) Diff {
	returnVal := New_Diff()
	base13Matches := make(map[uint32]uint32)
	base53Matches := make(map[uint32]uint32)
	returnVal.ChunkSize = sourceData.ChunkSize
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		base13Matches = findHashInSource(sourceData.Base13Hash, sourceData.Size, sourceData.ChunkSize, td.Base13HashSmall, td.Base13HashLarge)
	}()
	go func() {
		defer wg.Done()
		base53Matches = findHashInSource(sourceData.Base53Hash, sourceData.Size, sourceData.ChunkSize, td.Base53HashSmall, td.Base53HashLarge)
	}()
	wg.Wait()
	for k, v := range base13Matches {
		if base53Matches[k] == v {
			returnVal.Matches[k] = v
		}
	}
	for k, _ := range sourceData.Base13Hash.ByPosition {
		if _, found := returnVal.Matches[k]; !found {
			returnVal.NewChunks[k] = []byte{}
		}
	}
	return returnVal
}

func findHashInSource(sourceDataHash hashing.Hashes, sourceDataSize uint32, sourceDataChunkSize int, targetDataHashSmall hashing.Hashes, targetDataHashLarge hashing.Hashes) map[uint32]uint32 {
	returnVal := make(map[uint32]uint32)
	for sourcePosition, hash := range sourceDataHash.ByPosition {
		useSmallHash := sourcePosition > sourceDataSize-uint32(sourceDataChunkSize)
		var targetPosition uint32
		var ok bool
		if useSmallHash {
			targetPosition, ok = targetDataHashSmall.ByValue[hash]
		} else {
			targetPosition, ok = targetDataHashLarge.ByValue[hash]
		}
		if ok {
			returnVal[sourcePosition] = targetPosition
		}
	}
	return returnVal
}
