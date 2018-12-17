package main

import (
	"fmt"
	"strconv"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

type node struct {
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

	return data, newNode
}

func parseTree(data []int) (tree *node) {
	_, tree = parseNode(data)
	return
}

func sumMetadata(tree *node) int {
	if tree == nil {
		return 0
	}

	sum := 0
	for _, val := range tree.metadata {
		sum += val
	}

	for _, child := range tree.children {
		sum += sumMetadata(child)
	}

	return sum
}

func main() {
	line := utils.ReadLinesOrDie("input.txt")[0]
	nums := stringToNums(line)

	tree := parseTree(nums)

	metadataSum := sumMetadata(tree)

	fmt.Printf("metadataSum: %v\n", metadataSum)
}
