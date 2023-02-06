package hashing

import "sync"

type Hashes struct {
	ByValue    map[uint64]uint32
	ByPosition map[uint32]uint64
}

func New_Hashes() Hashes {
	var h Hashes
	h.ByValue = make(map[uint64]uint32)
	h.ByPosition = make(map[uint32]uint64)
	return h
}
func uInt64Pow(base uint64, exp uint64) uint64 {
	output := uint64(1)
	if exp == 0 {
		return 1
	}
	for i := uint64(1); i <= exp; i++ {
		output = output * base
	}
	return output
}
func getHash(bytesIn []byte, base uint64) uint64 {
	hashValue := uint64(0)
	for i := 0; i < len(bytesIn); i++ {
		hashValue += uint64(bytesIn[i]) * uInt64Pow(base, uint64(len(bytesIn)-i-1))
	}
	return hashValue
}
func GetRollingHash(bytesIn []byte, chunkSize int, base uint64) Hashes {
	returnValue := New_Hashes()
	initialHash := getHash(bytesIn[0:chunkSize], base)
	returnValue.ByPosition[0] = initialHash
	returnValue.ByValue[initialHash] = 0
	for i := 1; i < len(bytesIn)-chunkSize; i++ {
		nextHash := (initialHash-uint64(bytesIn[i-1])*uInt64Pow(base, uint64(chunkSize-1)))*base + (uint64(bytesIn[chunkSize+i-1]))
		returnValue.ByPosition[uint32(i)] = nextHash
		returnValue.ByValue[nextHash] = uint32(i)
		initialHash = nextHash
	}
	return returnValue
}
func fixedHash(i uint32, chunk []byte, base uint64, chunkSize int, wg *sync.WaitGroup, mutex *sync.Mutex, returnVal *Hashes) {
	defer wg.Done()
	hashVal := getHash(chunk, base)
	mutex.Lock()
	returnVal.ByPosition[i*uint32(chunkSize)] = hashVal
	returnVal.ByValue[hashVal] = i * uint32(chunkSize)
	mutex.Unlock()
}
func GetFixedHash(bytesIn []byte, chunkSize int, base uint64) Hashes {
	returnValue := New_Hashes()
	chunks := divideIntoChunks(bytesIn, chunkSize)
	wg := sync.WaitGroup{}
	wg.Add(len(chunks))
	var mutex = &sync.Mutex{}
	for i, chunk := range chunks {
		go fixedHash(uint32(i), chunk, base, chunkSize, &wg, mutex, &returnValue)
	}
	wg.Wait()
	return returnValue
}
func divideIntoChunks(fileData []byte, chunkSize int) [][]byte {
	position := 0
	var chunks [][]byte
	fileLength := len(fileData)
	for position < fileLength {
		var ubound int
		if position+chunkSize <= fileLength {
			ubound = chunkSize
		} else {
			ubound = len(fileData) - position
		}
		chunks = append(chunks, fileData[position:position+ubound])
		position += chunkSize
	}
	return chunks
}
