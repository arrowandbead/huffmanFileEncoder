package main

type heapNode struct {
	lChild *heapNode
	rChild *heapNode
	parent *heapNode
}

// type Model struct {
//     Data []T
// }

// type Node struct {
// 	lChild *Node
// 	rChild *Node
// 	parent *Node
// 	idNum  int
// 	isLead bool
// 	val    byte
// 	weight int
// }

// Compares two nodes, returns true if Node a has weight LESS THAN Node b, false otherwise
func compareNodesAGTB(a Node, b Node) bool {

	return a.weight > b.weight
}

func compareNodesALTB(a Node, b Node) bool {

	return a.weight < b.weight
}

func compareNumsALTB(a int, b int) bool {

	return a < b
}

func compareNumsAGTB(a int, b int) bool {

	return a > b
}

func heapRec[K any](data *[]**K, nodeNum int, numThings int, compareFunc func(K, K) bool) {
	best := nodeNum
	lChildNum := 2*nodeNum + 1
	rChildNum := 2*nodeNum + 2

	if lChildNum < numThings && compareFunc(**(*data)[best], **(*data)[lChildNum]) {
		best = lChildNum
	}

	if rChildNum < numThings && compareFunc(**(*data)[best], **(*data)[rChildNum]) {
		best = rChildNum
	}
	// fmt.Println((*data)[largest])
	if best != nodeNum {
		(*data)[nodeNum], (*data)[best] = (*data)[best], (*data)[nodeNum]
		heapRec(data, best, numThings, compareFunc)
	}
}

func heapUp[K any](data *[]**K, nodeNum int, numThings int, compareFunc func(K, K) bool) {
	if len((*data)) == 1 {
		return
	}
	if nodeNum == 0 {
		return
	}
	if len((*data)) == 1 {
		return
	}
	var parentNodeNum int
	if nodeNum%2 == 0 {
		parentNodeNum = (nodeNum - 2) / 2
	} else {
		parentNodeNum = (nodeNum - 1) / 2
	}

	parentNodeValue := (*data)[parentNodeNum]
	nodeValue := (*data)[nodeNum]

	if compareFunc(**nodeValue, **parentNodeValue) {
		(*data)[nodeNum], (*data)[parentNodeNum] = (*data)[parentNodeNum], (*data)[nodeNum]
		heapUp(data, parentNodeNum, numThings, compareFunc)
	}
}

func MakeHeap[K any](data *[]**K, compareFunc func(K, K) bool) {

	for i := len((*data)) - 1; i > -1; i-- {
		heapRec(data, i, len((*data)), compareFunc)
	}

}

func InsertIntoHeap[K any](toInsert *K, heap *[]**K, compareFunc func(K, K) bool) {
	*heap = (append((*heap), &toInsert))

	heapUp(heap, len((*heap))-1, len((*heap)), compareFunc)
}

func GetTopElement[K any](heap *[]**K, compareFunc func(K, K) bool) (topElement **K) {
	topElement = (*heap)[0]
	if len((*heap)) == 1 {
		*heap = make([]**K, 0)
		return topElement
	}
	lastElementIndex := len((*heap)) - 1
	topElement = (*heap)[0]

	(*heap)[lastElementIndex], (*heap)[0] = (*heap)[0], (*heap)[lastElementIndex]

	*heap = (*heap)[:lastElementIndex]

	heapRec(heap, 0, lastElementIndex, compareFunc)
	return topElement
}

func intPointersToInts(intPointerSlice []*int) (intSlice []int) {
	for _, v := range intPointerSlice {
		intSlice = append(intSlice, *v)
	}
	return intSlice
}

// func main() {
// 	b := []int{1, 2, 3, 4, 200, 555, 5, 8, 12, 17, 100, 150, 54, 33}

// 	var intPointers []*int
// 	for _, thing := range b {
// 		var newInt int = thing
// 		intPointers = append(intPointers, &newInt)
// 	}
// 	// for H := range b {
// 	// 	fmt.Println(string(H))
// 	// }

// 	MakeHeap[int](&intPointers, compareNumsAGTB)
// 	num := 130
// 	num2 := 17
// 	num3 := 54
// 	num4 := 190
// 	fmt.Println(intPointersToInts(intPointers))
// 	// for i := len((intPointers)) - 1; i >= 0; i-- {
// 	// 	fmt.Println(*(GetTopElement[int](&intPointers, compareNumsAGTB)))
// 	// }
// 	InsertIntoHeap[int](&num, &intPointers, compareNumsALTB)
// 	InsertIntoHeap[int](&num2, &intPointers, compareNumsALTB)
// 	InsertIntoHeap[int](&num3, &intPointers, compareNumsALTB)
// 	InsertIntoHeap[int](&num4, &intPointers, compareNumsALTB)

// 	fmt.Println(intPointersToInts(intPointers))
// 	for i := len((intPointers)) - 1; i >= 0; i-- {
// 		fmt.Println(*(GetTopElement[int](&intPointers, compareNumsAGTB)))
// 	}
// 	// for i := len((b)) - 1; i >= 0; i-- {
// 	// 	fmt.Println(*(GetTopElement[int](&intPointers, compareNumsAGTB)))
// 	// }
// 	// fmt.Println("----------")
// 	// fmt.Println(b)
// 	// thing := GetTopElement[int](&intPointers, compareNumsALTB)
// 	// fmt.Println(*thing)
// 	// thing = GetTopElement[int](&intPointers, compareNumsALTB)
// 	// fmt.Println(*thing)
// 	// fmt.Println(b)
// 	// thing = GetTopElement[int](&intPointers, compareNumsALTB)
// 	// fmt.Println(*thing)
// 	// fmt.Println(b)
// 	// fmt.Println("----------")

// }
