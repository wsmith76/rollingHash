package sourceHandler

import (
	"rollingHash/hashing"
)

type SourceData struct {
	data          []byte
	Size          uint32
	ChunkSize     int
	LastChunkSize int
	Base13Hash    hashing.Hashes
	Base53Hash    hashing.Hashes
}

type SourcePacket struct {
	Size          uint32
	ChunkSize     int
	LastChunkSize int
	Base13Hash    hashing.Hashes
	Base53Hash    hashing.Hashes
}

func New_SourceData(sourceBytes []byte) SourceData {
	var sd SourceData
	sd.data = sourceBytes
	sd.Size = uint32(len(sourceBytes))
	if len(sourceBytes) >= 1024 {
		sd.ChunkSize = 256
	} else {
		sd.ChunkSize = len(sourceBytes) / 4
	}
	sd.LastChunkSize = len(sourceBytes) % sd.ChunkSize
	if sd.LastChunkSize == 0 {
		sd.LastChunkSize = sd.ChunkSize
	}
	sd.Base13Hash = hashing.GetFixedHash(sourceBytes, sd.ChunkSize, 13)
	sd.Base53Hash = hashing.GetFixedHash(sourceBytes, sd.ChunkSize, 53)
	return sd
}

func (sd *SourceData) ToPacket() SourcePacket {
	var sp SourcePacket
	sp.Size = sd.Size
	sp.ChunkSize = sd.ChunkSize
	sp.LastChunkSize = sd.LastChunkSize
	sp.Base13Hash = sd.Base13Hash
	sp.Base53Hash = sd.Base53Hash
	return sp
}

func (sd *SourceData) NewChunks(request map[uint32][]byte) map[uint32][]byte {
	returnVal := make(map[uint32][]byte)
	for k, _ := range request {
		if k > sd.Size-uint32(sd.ChunkSize) {
			returnVal[k] = sd.data[k:]
		} else {
			returnVal[k] = sd.data[k : k+uint32(sd.ChunkSize)]
		}
	}
	return returnVal
}
