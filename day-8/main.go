package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	numChildren int
	numMetadata int
	children    []*Node
	metadata    []int
}

func NewNode() *Node {
	var node Node
	node.children = []*Node{}
	node.metadata = []int{}
	return &node
}

func (n *Node) AddChild(o *Node) {
	n.children = append(n.children, o)
}

func (n *Node) AddMetadata(i int) {
	n.metadata = append(n.metadata, i)
}

func (n *Node) SumMetadata() int {
	sum := 0
	for _, value := range n.metadata {
		sum += value
	}
	for _, child := range n.children {
		sum += child.SumMetadata()
	}
	return sum
}

func (n *Node) GetValue() int {
	val := 0
	if len(n.children) > 0 {
		for _, index := range n.metadata {
			if index <= len(n.children) {
				val += n.children[index-1].GetValue()
			}
		}

	} else {
		for _, value := range n.metadata {
			val += value
		}
	}
	return val
}

func BuildTree(nums []int) *Node {
	_, root := buildTree(nums, 0)
	return root
}

func buildTree(nums []int, offset int) (int, *Node) {
	if offset == len(nums) {
		return 0, nil
	}

	numChildren := nums[offset]
	numMetadata := nums[offset+1]
	offset += 2

	root := NewNode()
	for i := 0; i < numChildren; i++ {
		n, node := buildTree(nums, offset)
		root.AddChild(node)
		offset = n
	}

	for i := 0; i < numMetadata; i++ {
		root.AddMetadata(nums[offset])
		offset += 1
	}

	return offset, root
}

func MemoryManeuver(input []int) int {
	root := BuildTree(input)
	return root.SumMetadata()
}

func MemoryManeuverPart2(input []int) int {
	root := BuildTree(input)
	return root.GetValue()
}

func main() {
	if len(os.Args) != 2 {
		message := fmt.Sprintf("Usage: %s <input-file>", os.Args[0])
		log.Fatal(message)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var input string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = scanner.Text()
	}

	nums := []int{}
	parts := strings.Split(input, " ")

	for _, part := range parts {
		num, err := strconv.ParseInt(part, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, int(num))
	}

	part1 := MemoryManeuver(nums)
	fmt.Println(part1)

	part2 := MemoryManeuverPart2(nums)
	fmt.Println(part2)
}
