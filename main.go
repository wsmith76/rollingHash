package main

import (
	"io/ioutil"
	"log"
	"rollingHash/sourceHandler"
	"rollingHash/targetHandler"
)

func getDiff(sourceData []byte, targetData []byte) targetHandler.Diff {
	source := sourceHandler.New_SourceData(sourceData)
	target := targetHandler.New_TargetData(targetData, source.ChunkSize, source.LastChunkSize)
	diff := target.DiffWithoutNewData(source.ToPacket())
	diff.NewChunks = source.NewChunks(diff.NewChunks)
	return diff
}

func main() {
	sourceBytes, err := ioutil.ReadFile("Test1.txt")
	if err != nil {
		log.Fatal(err)
	}
	targetBytes, err := ioutil.ReadFile("Test2.txt")
	if err != nil {
		log.Fatal(err)
	}
	output := getDiff(sourceBytes, targetBytes)
	println(output)
}
