package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

const numThreads = 4

func HuffmanEncode(bytes []byte) (numBytes int, bitStringsToBytes map[string]byte, outBytes []byte) {
	var byteToCountMap = countBytes(bytes)

	backUpLeafSlice := make([]**Node, 0)
	nodeSlice := make([]**Node, 0)

	for k, v := range *byteToCountMap {
		newNode := Node{isLeaf: true, val: k, weight: v}
		newNodePointer := &newNode
		newNodePointerPointer := &newNodePointer
		nodeSlice = append(nodeSlice, newNodePointerPointer)
		backUpLeafSlice = append(backUpLeafSlice, newNodePointerPointer)
	}

	MakeHeap[Node](&nodeSlice, compareNodesAGTB)

	var root **Node

	for {
		if len(nodeSlice) == 1 {
			root = nodeSlice[0]
			break
		}

		var low **Node = GetTopElement[Node](&nodeSlice, compareNodesAGTB)
		var high **Node = GetTopElement[Node](&nodeSlice, compareNodesAGTB)

		newNode := Node{isLeaf: false, weight: (*low).weight + (*high).weight, lChild: *low, rChild: *high}
		(*low).parent = &newNode
		(*high).parent = &newNode

		InsertIntoHeap[Node](&newNode, &nodeSlice, compareNodesALTB)

	}
	return compressText(*root, backUpLeafSlice, bytes)
}

func compressText(treeRoot *Node, leaves []**Node, byteArray []byte) (numBytes int, bitStringsToBytes map[string]byte, outBytes []byte) {
	bytesToCountMap := make(map[byte]int)
	bytesToBitString := make(map[byte]string)

	for _, leaf := range leaves {
		var sb strings.Builder
		var parent = &((*leaf).parent)
		curr := leaf

		for {

			if (*parent).lChild == (*curr) {
				sb.WriteString("0")
			} else {

				sb.WriteString("1")
			}
			if (*parent) == treeRoot {
				break
			}
			curr = parent
			parent = &((*parent).parent)
		}
		stringRep := sb.String()
		bytesToBitString[(*leaf).val] = reverse(stringRep)
		bytesToCountMap[(*leaf).val] = (*leaf).weight
	}
	var compressedBuilder strings.Builder
	numBytes = len(byteArray)
	for _, b := range byteArray {

		compressedBuilder.WriteString(bytesToBitString[b])
	}
	compressedBits := compressedBuilder.String()

	var startBit = 0
	for {
		endBit := startBit + 7
		bits := ""
		if endBit > len(compressedBits) {
			bits = compressedBits[startBit:]
			for i := len(bits); i < 8; i++ {
				bits += "0"
			}
			intVersion, _ := strconv.ParseUint(bits, 2, 8)

			outBytes = append(outBytes, byte(intVersion))
			break
		}
		bits = compressedBits[startBit : endBit+1]
		intVersion, _ := strconv.ParseUint(bits, 2, 8)

		outBytes = append(outBytes, byte(intVersion))
		startBit += 8

	}
	bitStringsToBytes = reverseMap(bytesToBitString)

	var bitRep = make([]int, 8*len(outBytes))
	for i, b := range outBytes {

		for j := 0; j < 8; j++ {
			bitRep[i*8+j] = int(b >> uint(7-j) & 0x01)
		}

	}

	var bsSb strings.Builder
	for i := 0; i < len(bitRep); i++ {
		bsSb.WriteString(strconv.Itoa(bitRep[i]))
	}
	bitString := bsSb.String()

	decompressedBytes := make([]byte, 0)
	numBytesMade := 0

	currBits := ""
	for i := 0; i < len(bitString); i++ {
		if numBytes == numBytesMade {
			break
		}
		currBits += string(bitString[i])
		if val, ok := bitStringsToBytes[currBits]; ok {
			numBytesMade += 1
			currBits = ""
			decompressedBytes = append(decompressedBytes, val)
		}
	}
	return numBytes, bitStringsToBytes, outBytes
}

func DecompressText(numBytes int, bitStringsToBytes map[string]byte, outBytes []byte) (decompressedBytes []byte) {
	var bitRep = make([]int, 8*len(outBytes))
	for i, b := range outBytes {

		for j := 0; j < 8; j++ {
			bitRep[i*8+j] = int(b >> uint(7-j) & 0x01)
		}

	}

	var bsSb strings.Builder
	for i := 0; i < len(bitRep); i++ {
		bsSb.WriteString(strconv.Itoa(bitRep[i]))
	}
	bitString := bsSb.String()

	numBytesMade := 0

	currBits := ""
	for i := 0; i < len(bitString); i++ {
		if numBytes == numBytesMade {
			break
		}
		currBits += string(bitString[i])
		if val, ok := bitStringsToBytes[currBits]; ok {
			numBytesMade += 1
			currBits = ""
			decompressedBytes = append(decompressedBytes, val)
		}
	}
	return decompressedBytes
}

func nodesToWeightVector(nodes *[]**Node) (weights []int) {
	for _, node := range *nodes {
		weights = append(weights, (*node).weight)
	}
	return weights
}

func countReducer(byteChunk []byte) map[byte]int {

	byteCounts := make(map[byte]int)

	for _, b := range byteChunk {
		if _, ok := byteCounts[b]; ok {
			byteCounts[b] += 1
		} else {
			byteCounts[b] = 1
		}
	}

	return byteCounts
}
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func reverseMap(m map[byte]string) map[string]byte {
	n := make(map[string]byte, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

func countBytes(bytes []byte) *map[byte]int {

	chunkSize := int((len(bytes) + numThreads - 1) / numThreads)
	var divided [][]byte
	for i := 0; i < len(bytes); i += chunkSize {
		end := i + chunkSize
		if end > len(bytes) {
			end = len(bytes)
		}
		divided = append(divided, bytes[i:end])
	}

	var countMapReceiver chan map[byte]int = make(chan map[byte]int)

	var wg sync.WaitGroup

	var total = 0
	for _, chunk := range divided {
		total += len(chunk)
		wg.Add(1)
		go func(goChunk []byte) {

			defer wg.Done()
			countMap := countReducer(goChunk)
			countMapReceiver <- countMap
		}(chunk)
	}

	var finalMap map[byte]int = make(map[byte]int)

	var numDone = 0
	for {

		counts, good := <-countMapReceiver
		if !good {
			continue
		} else {
			numDone += 1
		}
		for k, v := range counts {
			if _, ok := finalMap[k]; ok {
				finalMap[k] += v
			} else {
				finalMap[k] = v
			}
		}
		if numDone == len(divided) {
			break
		}
	}
	wg.Wait()

	return &finalMap
}

func main() {
	data, _ := os.ReadFile("./longRandom.txt")
	numBytes, bitStringsToBytes, outBytes := HuffmanEncode(data)
	fmt.Println(len(data))
	fmt.Println(len(outBytes))
	decompressed := DecompressText(numBytes, bitStringsToBytes, outBytes)
	fmt.Println(len(string(decompressed)))
}
