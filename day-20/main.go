package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type TokenType int

const (
	Begin      TokenType = 0
	End        TokenType = 1
	Seq        TokenType = 2
	LeftParen  TokenType = 3
	RightParen TokenType = 4
	Pipe       TokenType = 5
)

type Token struct {
	kind   TokenType
	lexeme string
}

func tokenize(regex string) []Token {
	tokens := []Token{}

	c := 0
	next := func() rune {
		char := regex[c]
		c += 1
		return rune(char)
	}

	for c := next(); c != '$'; {
		switch c {
		case '^':
			token := Token{Begin, ""}
			tokens = append(tokens, token)
			c = next()
		case '(':
			token := Token{LeftParen, ""}
			tokens = append(tokens, token)
			c = next()
		case ')':
			token := Token{RightParen, ""}
			tokens = append(tokens, token)
			c = next()
		case '|':
			token := Token{Pipe, ""}
			tokens = append(tokens, token)
			c = next()
		default:
			var lexeme strings.Builder
			for c == 'N' || c == 'S' || c == 'E' || c == 'W' {
				lexeme.WriteRune(c)
				c = next()
			}
			token := Token{Seq, lexeme.String()}
			tokens = append(tokens, token)
		}
	}

	token := Token{End, ""}
	tokens = append(tokens, token)

	return tokens
}

type NodeType string

const (
	Root     NodeType = "Root"
	TopLevel NodeType = "TopLevel"
	Choice   NodeType = "Choice"
	Val      NodeType = "Val"
)

type Node struct {
	val      string
	kind     NodeType
	children []*Node
}

func NewNode(val string, kind NodeType) *Node {
	n := Node{val, kind, []*Node{}}
	return &n
}

func (n *Node) Print(s int) {
	for i := 0; i < s; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("[%s] %s\n", n.kind, n.val)
	for _, c := range n.children {
		c.Print(s + 2)
	}
}

func (n *Node) PrintToString(s int) string {
	var sb strings.Builder
	padding := ""
	for i := 0; i < s; i++ {
		padding += " "
	}
	sb.WriteString(padding)
	node := fmt.Sprintf("[%s] %s\n", n.kind, n.val)
	sb.WriteString(node)
	for _, c := range n.children {
		sb.WriteString(c.PrintToString(s + 2))
	}
	return sb.String()
}

func (n *Node) PrintToFile(name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(n.PrintToString(0))
	w.Flush()
}

func parse(tokens []Token) *Node {
	root := NewNode("", Root)

	c := 0
	next := func() Token {
		token := tokens[c]
		c += 1
		return token
	}

	for token := next(); token.kind != End; {
		if token.kind == Begin {
			token = next()
		}
		n, t := parseTopLevels(token, next)
		token = t
		root = n
	}

	return root
}

func parseTopLevels(token Token, next func() Token) (*Node, Token) {
	// fmt.Println("parseTopLevels", token)
	if token.kind == Seq {
		n := NewNode("", TopLevel)
		child, t := parseTopLevel(token, next)
		n.children = append(n.children, child)
		token = t
		for token.kind == Seq {
			child, t := parseTopLevel(token, next)
			n.children = append(n.children, child)
			token = t
		}
		return n, token
	}
	return nil, token
}

func parseTopLevel(token Token, next func() Token) (*Node, Token) {
	// fmt.Println("parseTopLevel", token)
	if token.kind == Seq {
		n := NewNode(token.lexeme, Val)
		token = next()
		if token.kind == LeftParen {
			token = next()
			child, t := parseOptions(token, next)
			token = t
			n.children = append(n.children, child)
		}
		return n, token
	}
	return nil, token
}

func parseOptions(token Token, next func() Token) (*Node, Token) {
	// fmt.Println("parseOptions", token)
	option := NewNode("", Choice)
	for token.kind != RightParen {
		n, t := parseOption(token, next)
		token = t
		option.children = append(option.children, n)
		for token.kind == Pipe {
			token = next()
			n, t := parseOption(token, next)
			token = t
			option.children = append(option.children, n)
		}
	}
	token = next()
	return option, token
}

func parseOption(token Token, next func() Token) (*Node, Token) {
	// fmt.Println("parseOption", token)
	if token.kind == Seq {
		n, t := parseTopLevels(token, next)
		token = t
		return n, token
	} else if token.kind == LeftParen {
		token = next()
		n, t := parseOptions(token, next)
		token = t
		return n, token
	} else {
		n := NewNode("", Val)
		return n, token
	}
}

func getShortestPath(root *Node) string {
	path := ""

	if len(root.children) == 0 {
		path = root.val
		return path
	}

	if root.kind == TopLevel {
		path = root.val
		for _, child := range root.children {
			path += getShortestPath(child)
		}
	}

	if root.kind != TopLevel {
		minSubPath := getShortestPath(root.children[0])
		for _, child := range root.children {
			subPath := getShortestPath(child)
			if len(subPath) < len(minSubPath) {
				minSubPath = subPath
			}
		}
		path = root.val + minSubPath
	}

	return path
}

func getLongestPath(root *Node) string {
	// fmt.Println(root)
	path := ""

	if len(root.children) == 0 {
		path = root.val
		return path
	}

	if root.kind == TopLevel {
		path = root.val
		for _, child := range root.children {
			path += getLongestPath(child)
		}
	}

	if root.kind != TopLevel {
		maxSubPath := getLongestPath(root.children[0])
		for _, child := range root.children {
			subPath := getLongestPath(child)
			if len(subPath) > len(maxSubPath) {
				maxSubPath = subPath
			}
		}
		path = root.val + maxSubPath
	}

	return path
}

func getPossiblePaths(root *Node) []string {
	paths := []string{}

	if len(root.children) == 0 {
		paths = append(paths, root.val)
		return paths
	}

	if root.kind == TopLevel {
		prevSubPaths := []string{""}
		for _, child := range root.children {
			subPaths := getPossiblePaths(child)
			newSubPaths := []string{}
			for _, prevSubPath := range prevSubPaths {
				for _, subPath := range subPaths {
					newSubPaths = append(newSubPaths, prevSubPath+subPath)
				}
			}
			prevSubPaths = newSubPaths
		}
		paths = prevSubPaths
	}

	if root.kind != TopLevel {
		for _, child := range root.children {
			subPaths := getPossiblePaths(child)
			for _, path := range subPaths {
				paths = append(paths, root.val+path)
			}
		}
	}

	return paths
}

func RegularMap(regex string) int {
	fmt.Println("Begin tokenize...")
	tokens := tokenize(regex)
	fmt.Println("End tokenize...")
	fmt.Println("Begin parse...")
	root := parse(tokens)
	fmt.Println("End parse...")

	fmt.Println("Printing to file...")
	root.PrintToFile("please")
	fmt.Println("End print...")

	fmt.Println("Begin looking for path...")
	path := root.val
	for i := 0; i < len(root.children); i++ {
		child := root.children[i]
		if i == len(root.children)-1 {
			subPath := getLongestPath(child)
			path += subPath
		} else {
			subPath := getShortestPath(child)
			path += subPath
		}
	}
	fmt.Println("End looking for path...")

	return len(path)
}

// Other approach

type Stack []Pos

func (s *Stack) Pop() Pos {
	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item
}

func (s *Stack) Push(item Pos) {
	*s = append(*s, item)
}

func (s *Stack) Top() Pos {
	return (*s)[len(*s)-1]
}

type Pos struct {
	x    int
	y    int
	kind rune
}

type Map struct {
	layout        [][]rune
	start         Point
	width, height int
}

func NewMap(directions []Pos) *Map {
	m := Map{}

	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32

	for _, d := range directions {
		if d.x < minX {
			minX = d.x
		}
		if d.x > maxX {
			maxX = d.x
		}
		if d.y < minY {
			minY = d.y
		}
		if d.y > maxY {
			maxY = d.y
		}
	}

	m.width = maxX - minX + 3
	m.height = maxY - minY + 3

	m.layout = make([][]rune, m.height)
	for row := 0; row < m.height; row++ {
		m.layout[row] = make([]rune, m.width)
		for col := 0; col < m.width; col++ {
			m.layout[row][col] = '?'
		}
	}

	for _, d := range directions {
		row := d.y - minY + 1
		col := d.x - minX + 1
		m.layout[row][col] = d.kind
		if d.kind == 'X' {
			m.start = Point{col, row}
		}
	}

	for row := 0; row < m.height; row++ {
		for col := 0; col < m.width; col++ {
			if m.layout[row][col] == '?' {
				m.layout[row][col] = '#'
			}
		}
	}

	return &m
}

func (m *Map) Print() {
	for _, line := range m.layout {
		for _, tile := range line {
			fmt.Printf("%c", tile)
		}
		fmt.Println()
	}
}

func (m *Map) PrintToFile(file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range m.layout {
		for _, tile := range line {
			w.WriteString(fmt.Sprintf("%c", tile))
		}
		w.WriteString("\n")
	}
	w.Flush()
}

type Point struct {
	x, y int
}

type Queue []Point

func (q *Queue) Push(p Point) {
	*q = append(*q, p)
}

func (q *Queue) Pop() Point {
	i := (*q)[len(*q)-1]
	(*q) = (*q)[:len(*q)-1]
	return i
}

func (q *Queue) Empty() bool {
	return len(*q) == 0
}

func (q *Queue) Contains(p Point) bool {
	for _, i := range *q {
		if i == p {
			return true
		}
	}
	return false
}

type Set []Point

func (s *Set) Push(p Point) {
	*s = append(*s, p)
}

func (s *Set) Contains(p Point) bool {
	for _, i := range *s {
		if i == p {
			return true
		}
	}
	return false
}

func (m *Map) GetNeighbors(from Point) []Point {
	neighbors := []Point{}

	if from.x-2 >= 0 && m.layout[from.y][from.x-1] != '#' && m.layout[from.y][from.x-2] == '.' {
		neighbors = append(neighbors, Point{from.x - 2, from.y})
	}
	if from.y-2 >= 0 && m.layout[from.y-1][from.x] != '#' && m.layout[from.y-2][from.x] == '.' {
		neighbors = append(neighbors, Point{from.x, from.y - 2})
	}
	if from.x+2 < m.width && m.layout[from.y][from.x+1] != '#' && m.layout[from.y][from.x+2] == '.' {
		neighbors = append(neighbors, Point{from.x + 2, from.y})
	}
	if from.y+2 < m.height && m.layout[from.y+1][from.x] != '#' && m.layout[from.y+2][from.x] == '.' {
		neighbors = append(neighbors, Point{from.x, from.y + 2})
	}

	return neighbors
}

func (m *Map) BFS() (int, int) {
	open := Queue{}
	closed := Set{}
	meta := map[Point]int{}

	meta[m.start] = 0
	open.Push(m.start)

	for !open.Empty() {
		tile := open.Pop()
		for _, n := range m.GetNeighbors(tile) {
			if closed.Contains(n) {
				continue
			}
			if !open.Contains(n) {
				meta[n] = meta[tile] + 1
				open.Push(n)
			}
		}
		closed.Push(tile)
	}

	max := 0
	count := 0
	for _, v := range meta {
		if v > max {
			max = v
		}
		if v >= 1000 {
			count++
		}
	}
	return max, count
}

func RegularMapV2(regex string) (int, int) {
	history := []Pos{}
	state := Stack{}
	pos := Pos{0, 0, 'X'}
	history = append(history, pos)
	for _, c := range regex {
		switch c {
		case '(':
			state.Push(pos)
		case ')':
			pos = state.Pop()
		case '|':
			pos = state.Top()
		case 'N':
			pos.y -= 1
			pos.kind = '-'
			history = append(history, pos)
			pos.y -= 1
			pos.kind = '.'
			history = append(history, pos)
		case 'S':
			pos.y += 1
			pos.kind = '-'
			history = append(history, pos)
			pos.y += 1
			pos.kind = '.'
			history = append(history, pos)
		case 'W':
			pos.x -= 1
			pos.kind = '|'
			history = append(history, pos)
			pos.x -= 1
			pos.kind = '.'
			history = append(history, pos)
		case 'E':
			pos.x += 1
			pos.kind = '|'
			history = append(history, pos)
			pos.x += 1
			pos.kind = '.'
			history = append(history, pos)
		}
	}

	m := NewMap(history)
	m.Print()

	return m.BFS()
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

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(RegularMapV2(scanner.Text()))
}
