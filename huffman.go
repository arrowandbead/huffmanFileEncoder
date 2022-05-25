package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

const numThreads = 4

func HuffmanEncode(bytes []byte) (compressed []byte) {
	var byteToCountMap = countBytes(bytes)

	// type Node struct {
	// 	lChild *Node
	// 	rChild *Node
	// 	parent *Node
	// 	idNum  int
	// 	isLead bool
	// 	val    int
	// 	weight int
	// }
	fmt.Println(len(bytes))
	//make nodes
	backUpLeafSlice := make([]**Node, 0)
	nodeSlice := make([]**Node, 0)

	for k, v := range *byteToCountMap {
		// fmt.Println(int(k))
		newNode := Node{isLeaf: true, val: k, weight: v}
		newNodePointer := &newNode
		newNodePointerPointer := &newNodePointer
		nodeSlice = append(nodeSlice, newNodePointerPointer)
		backUpLeafSlice = append(backUpLeafSlice, newNodePointerPointer)
	}

	// backUpLeafSlice = backUpLeafSlice[:len(backUpLeafSlice)-1]
	// // fmt.Println(backUpLeafSlice)
	// nodeSlice = nodeSlice[:len(nodeSlice)-1]

	var root **Node

	for {
		if len(nodeSlice) == 1 {
			root = nodeSlice[0]
			break
		}
		// fmt.Println("------")
		// fmt.Println(nodesToWeightVector(&nodeSlice))
		// fmt.Println(len(nodeSlice))

		var low **Node = GetTopElement[Node](&nodeSlice, compareNodesAGTB)
		var high **Node = GetTopElement[Node](&nodeSlice, compareNodesAGTB)

		// fmt.Println(low.weight)
		// fmt.Println(high.weight)
		newNode := Node{isLeaf: false, weight: (*low).weight + (*high).weight, lChild: *low, rChild: *high}
		// fmt.Println(newNode.weight)
		(*low).parent = &newNode
		(*high).parent = &newNode
		// fmt.Println("-----------")
		// fmt.Println(nodeSlice)
		InsertIntoHeap[Node](&newNode, &nodeSlice, compareNodesALTB)

		// fmt.Println(nodesToWeightVector(&nodeSlice))
		// fmt.Println("------")

	}
	compressText(*root, backUpLeafSlice, bytes)

	return make([]byte, 2)
}

func compressText(treeRoot *Node, leaves []**Node, bytes []byte) {
	bytesToCountMap := make(map[byte]int)
	bytesToBitString := make(map[byte]string)
	for b := range bytes {
		fmt.Println(int(b))
	}
	// fmt.Println(leaves)
	// fmt.Println(leaves)
	// for _, leaf := range leaves {
	// 	fmt.Println(*leaf)
	// }
	// os.Exit(2)
	fmt.Println("-----------------")
	for _, leaf := range leaves {
		var sb strings.Builder
		var parent = (*leaf).parent
		fmt.Println(leaf)
		for {

			if parent.lChild == *leaf {
				sb.WriteString("0")
			} else {
				sb.WriteString("1")
			}
			if parent == treeRoot {
				break
			}
			parent = parent.parent
		}
		stringRep := sb.String()

		bytesToBitString[(*leaf).val] = stringRep
		bytesToCountMap[(*leaf).val] = (*leaf).weight

	}
	var compressedBuilder strings.Builder
	// var numBytes int = len(bytes)
	for _, b := range bytes {
		compressedBuilder.WriteString(bytesToBitString[b])
	}

	compressedBits := compressedBuilder.String()
	var startByte = 0
	var outBytes = make([]byte, 0)
	for {
		endByte := startByte + 7

		if endByte > len(compressedBits) {
			bits := compressedBits[startByte:]
			for i := len(bits); i < 9; i++ {
				bits += "0"
			}
			intVersion, _ := strconv.ParseInt(bits, 2, 8)
			outBytes = append(outBytes, byte(intVersion))
			break
		}
		bits := compressedBits[startByte : endByte+1]

		intVersion, _ := strconv.ParseInt(bits, 2, 8)

		outBytes = append(outBytes, byte(intVersion))
		startByte += 8

	}

	bitString := ""
	count := 0
	for _, b := range outBytes {
		// fmt.Println(int64(b))
		if uint64(b) == 127 {
			count += 1
		}
		thing := strconv.FormatUint(uint64(b), 2)

		// fmt.Println(int(b))
		bitString += thing
	}

	// nextBit = 0
	// for i := 0; i < numBytes; i++ {
	// 	currNode := treeRoot
	// 	numBitsUsed += 1
	// 	for {
	// 		if currNode.isLeaf {

	// 		}
	// 	}
	// }
}

// func bytesAsBinary(bytes []byte)(bitString string){
// 	for i := 0; i < len(bytes); i++ {
//         for j := 0; j < 8; j++ {
//             zeroOrOne := bytes[i] >> (7 - j) & 1
//             fmt.Printf("%c", '0'+zeroOrOne)
//         }
//         fmt.Print(" ")
//     }
// }
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
	data, _ := os.ReadFile("./longerlonger.txt")
	HuffmanEncode(data)
}
