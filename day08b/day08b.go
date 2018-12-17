package main

import (
	"fmt"
	"strconv"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

type node struct {
	value    int
	metadata []int
	children []*node
}

func stringToNums(s string) []int {
	numStrs := strings.Split(s, " ")
	nums := make([]int, len(numStrs))
	for i, numStr := range numStrs {
		nums[i], _ = strconv.Atoi(numStr)
	}
	return nums
}

// This assumes child values are already calculated.
func calcNodeValue(n *node) {
	n.value = 0

	if len(n.children) == 0 {
		for _, val := range n.metadata {
			n.value += val
		}
	} else {
		for _, childNum := range n.metadata {
			childIdx := childNum - 1
			if childIdx >= 0 && childIdx < len(n.children) {
				n.value += n.children[childIdx].value
			}
		}
	}
}

func parseNode(data []int) ([]int, *node) {
	if len(data) == 0 {
		return nil, nil
	}

	childCount := data[0]
	metadataCount := data[1]
	data = data[2:]

	newNode := new(node)
	newNode.children = make([]*node, childCount)
	newNode.metadata = make([]int, metadataCount)

	for i := 0; i < childCount; i++ {
		data, newNode.children[i] = parseNode(data)
	}

	for i := 0; i < metadataCount; i++ {
		newNode.metadata[i] = data[0]
		data = data[1:]
	}

	calcNodeValue(newNode)

	return data, newNode
}

func parseTree(data []int) (tree *node) {
	_, tree = parseNode(data)
	return
}

func main() {
	line := utils.ReadLinesOrDie("input.txt")[0]
	nums := stringToNums(line)

	tree := parseTree(nums)

	fmt.Printf("tree value: %v\n", tree.value)
}
