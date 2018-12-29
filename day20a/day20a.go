package main

import (
	"fmt"
	"math"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

const (
	clearSpace = '.'
	horizDoor  = '-'
	vertDoor   = '|'
	wallSpace  = '#'
	startSpot  = 'X'
)

type mapCoord struct {
	x, y int
}

type regexNode struct {
	visited  bool
	path     string
	branches []*regexNode
	next     *regexNode
}

var numVisited int

func doPrintRegexTree(node *regexNode, depth int) {
	nodeNum := 1
	for node != nil {
		indent := make([]byte, depth)
		for i := range indent {
			indent[i] = ' '
		}
		pathStr := node.path
		if pathStr == "" {
			pathStr = "<empty>"
		}
		fmt.Printf("%s%d. %s\n", string(indent), nodeNum, pathStr)

		for _, branch := range node.branches {
			doPrintRegexTree(branch, depth+2)
		}

		node = node.next
		nodeNum++
	}
}

func printRegexTree(node *regexNode) {
	doPrintRegexTree(node, 0)
}

type tokenType int

const (
	startRe tokenType = iota
	endRe
	lparen
	rparen
	pipe
	path
)

func parseError(pos int, msg string) {
	panic(fmt.Sprintf("parse error at %d: %s", pos, msg))
}

type token struct {
	tokType tokenType
	text    string
	pos     int
}

type lexer struct {
	nextTokPos int
	nextTok    *token
	input      string
}

func (l *lexer) init(input string) {
	l.input = input
	l.pullNextToken()
}

func (l *lexer) pullNextToken() {
	l.nextTok = new(token)
	l.nextTok.pos = l.nextTokPos

	idx := strings.IndexAny(l.input, "^()|$")
	l.nextTokPos = idx + 1

	if idx == 0 {
		switch l.input[0] {
		case '^':
			l.nextTok.tokType = startRe
			break
		case '(':
			l.nextTok.tokType = lparen
			break
		case ')':
			l.nextTok.tokType = rparen
			break
		case '|':
			l.nextTok.tokType = pipe
			break
		case '$':
			l.nextTok.tokType = endRe
			break
		}

		l.nextTok.text = l.input[:1]
		l.input = l.input[1:]
		l.nextTokPos++
	} else if idx >= 1 {
		l.nextTok.tokType = path
		l.nextTok.text = l.input[:idx]
		l.input = l.input[idx:]
		l.nextTokPos = idx
	} else {
		l.nextTok.tokType = path
		l.nextTok.text = l.input
		l.input = ""
	}
}

func matchToken(expectedType tokenType, theLexer *lexer) *token {
	theTok := theLexer.nextTok
	if theTok.tokType == expectedType {
		theLexer.pullNextToken()
	} else {
		parseError(theLexer.nextTok.pos,
			fmt.Sprintf("unexpected token: \"%s\"", theLexer.nextTok.text))
	}

	return theTok
}

func parseBranchListTail(theLexer *lexer) []*regexNode {
	var listNodes []*regexNode

	switch theLexer.nextTok.tokType {
	case path, lparen:
		startNode := parseSpec(theLexer)
		var endNodes []*regexNode
		if theLexer.nextTok.tokType == pipe {
			matchToken(pipe, theLexer)
			endNodes = parseBranchListTail(theLexer)
		}
		listNodes = append([]*regexNode{startNode}, endNodes...)
	case rparen:
		// epsilon production
		listNodes = make([]*regexNode, 1)
		listNodes[0] = new(regexNode)
	default:
		parseError(theLexer.nextTok.pos,
			fmt.Sprintf("unknown token: \"%s\"", theLexer.nextTok.text))
	}

	return listNodes
}

func parseBranchList(theLexer *lexer) []*regexNode {
	var listNodes []*regexNode

	switch theLexer.nextTok.tokType {
	case path, lparen:
		startNode := parseSpec(theLexer)
		matchToken(pipe, theLexer)
		endNodes := parseBranchListTail(theLexer)
		listNodes = append([]*regexNode{startNode}, endNodes...)
	default:
		parseError(theLexer.nextTok.pos,
			fmt.Sprintf("unknown token: \"%s\"", theLexer.nextTok.text))
	}

	return listNodes
}

func parseComponent(theLexer *lexer) *regexNode {
	var compNode regexNode

	switch theLexer.nextTok.tokType {
	case path:
		theTok := matchToken(path, theLexer)
		compNode.path = theTok.text
	case lparen:
		matchToken(lparen, theLexer)
		compNode.branches = parseBranchList(theLexer)
		matchToken(rparen, theLexer)
	default:
		parseError(theLexer.nextTok.pos,
			fmt.Sprintf("unknown token: \"%s\"", theLexer.nextTok.text))
	}

	return &compNode
}

func parseComponentTail(theLexer *lexer) *regexNode {
	var compNode *regexNode

	switch theLexer.nextTok.tokType {
	case path, lparen:
		compNode = parseSpec(theLexer)
	default:
		// epsilon production
	}

	return compNode
}

func parseSpec(theLexer *lexer) *regexNode {
	var node *regexNode

	switch theLexer.nextTok.tokType {
	case path, lparen:
		node = parseComponent(theLexer)
		node.next = parseComponentTail(theLexer)
	default:
		parseError(theLexer.nextTok.pos,
			fmt.Sprintf("unknown token: \"%s\"", theLexer.nextTok.text))
	}

	return node
}

func parseBody(theLexer *lexer) *regexNode {
	var node *regexNode
	switch theLexer.nextTok.tokType {
	case startRe:
		matchToken(startRe, theLexer)
		node = parseSpec(theLexer)
		matchToken(endRe, theLexer)
	default:
		parseError(theLexer.nextTok.pos,
			fmt.Sprintf("unknown token: \"%s\"", theLexer.nextTok.text))
	}

	return node
}

func parseRegex(regexStr string) *regexNode {
	var theLexer lexer
	theLexer.init(regexStr)

	node := parseBody(&theLexer)

	return node
}

func reconstructRegexStrBody(
	node *regexNode, strBuilder *strings.Builder) {

	for node != nil {
		if len(node.branches) == 0 {
			strBuilder.WriteString(node.path)
		} else {
			strBuilder.WriteByte('(')
			for i, branch := range node.branches {
				if i != 0 {
					strBuilder.WriteByte('|')
				}
				reconstructRegexStrBody(branch, strBuilder)
			}
			strBuilder.WriteByte(')')
		}

		node = node.next
	}
}

func reconstructRegexStr(node *regexNode) string {
	var strBuilder strings.Builder

	strBuilder.WriteByte('^')
	reconstructRegexStrBody(node, &strBuilder)
	strBuilder.WriteByte('$')

	return strBuilder.String()
}

func doMove(move byte, pos mapCoord, board [][]byte) mapCoord {
	if move == 'N' {
		pos.y--
		board[pos.x][pos.y] = horizDoor
		pos.y--
		board[pos.x][pos.y] = clearSpace
	} else if move == 'E' {
		pos.x++
		board[pos.x][pos.y] = vertDoor
		pos.x++
		board[pos.x][pos.y] = clearSpace
	} else if move == 'S' {
		pos.y++
		board[pos.x][pos.y] = horizDoor
		pos.y++
		board[pos.x][pos.y] = clearSpace
	} else if move == 'W' {
		pos.x--
		board[pos.x][pos.y] = vertDoor
		pos.x--
		board[pos.x][pos.y] = clearSpace
	}

	return pos
}

func doTraceRoutes(
	node *regexNode, board [][]byte, pos mapCoord) []mapCoord {

	if node == nil {
		return []mapCoord{pos}
	}

	if !node.visited {
		node.visited = true
		numVisited++
	}

	for _, c := range node.path {
		pos = doMove(byte(c), pos, board)
	}

	endPoints := []mapCoord{pos}

	for _, branch := range node.branches {
		results := doTraceRoutes(branch, board, pos)
		endPoints = append(endPoints, results...)
	}

	endPointMap := map[mapCoord]bool{}
	for _, coord := range endPoints {
		endPointMap[coord] = true
	}

	endPoints = endPoints[:0]
	for k := range endPointMap {
		endPoints = append(endPoints, k)
	}

	finalPoints := []mapCoord{}

	for _, p := range endPoints {
		results := doTraceRoutes(node.next, board, p)
		finalPoints = append(finalPoints, results...)
	}

	return finalPoints
}

func traceRoutes(node *regexNode, board [][]byte, pos mapCoord) {
	numVisited = 0
	doTraceRoutes(node, board, pos)
}

func trimBoard(board [][]byte, pos mapCoord) ([][]byte, mapCoord) {
	cols := len(board)
	rows := len(board[0])

	minX := math.MaxInt64
	maxX := math.MinInt64
	minY := math.MaxInt64
	maxY := math.MinInt64

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			if board[x][y] != wallSpace {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	trimmedCols := maxX - minX + 1 + 2
	trimmedRows := maxY - minY + 1 + 2

	pos = mapCoord{
		pos.x - minX + 1,
		pos.y - minY + 1,
	}

	newBoard := make([][]byte, trimmedCols)
	for x := range newBoard {
		newBoard[x] = make([]byte, trimmedRows)
		for y := range newBoard[x] {
			newBoard[x][y] = board[x+minX-1][y+minY-1]
		}
	}

	fmt.Printf("big board range: %d, %d, %d, %d\n", minX, minY, maxX, maxY)
	fmt.Printf("big board size: %dx%d\n", maxX-minX+1, maxY-minY+1)

	newBoard[pos.x][pos.y] = startSpot

	return newBoard, pos
}

func printBoard(board [][]byte) {
	cols := len(board)
	rows := len(board[0])
	bufLen := (cols + 1) * rows
	byteBuf := make([]byte, bufLen)
	idx := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			byteBuf[idx] = board[x][y]
			idx++
		}
		byteBuf[idx] = '\n'
		idx++
	}

	fmt.Printf(string(byteBuf))
}

func countNodes(node *regexNode) int {
	if node == nil {
		return 0
	}

	cnt := 1
	for _, branch := range node.branches {
		cnt += countNodes(branch)
	}

	cnt += countNodes(node.next)

	return cnt
}

func buildBoard(cols, rows int) [][]byte {
	board := make([][]byte, cols)
	for x := range board {
		board[x] = make([]byte, rows)
		for y := range board[x] {
			board[x][y] = wallSpace
		}
	}

	return board
}

func floodForDist(board [][]byte, distances [][]int,
	pointsOfInterest []mapCoord, dist int) {

	if len(pointsOfInterest) <= 0 {
		return
	}

	newPointsOfInterest := []mapCoord{}

	canMoveTo := func(p mapCoord) bool {
		return board[p.x][p.y] != wallSpace && distances[p.x][p.y] == -1
	}

	checkSpace := func(point mapCoord, dx, dy int) {
		newPoint := point
		newPoint.x += dx
		newPoint.y += dy
		if canMoveTo(newPoint) {
			distances[newPoint.x][newPoint.y] = dist
			newPoint.x += dx
			newPoint.y += dy
			if canMoveTo(newPoint) {
				distances[newPoint.x][newPoint.y] = dist + 1
				newPointsOfInterest = append(newPointsOfInterest, newPoint)
			}
		}
	}

	for _, poi := range pointsOfInterest {
		checkSpace(poi, 0, -1)
		checkSpace(poi, 1, 0)
		checkSpace(poi, 0, 1)
		checkSpace(poi, -1, 0)
	}

	floodForDist(board, distances, newPointsOfInterest, dist+1)
}

func findMaxDist(board [][]byte, pos mapCoord) int {
	cols := len(board)
	rows := len(board[0])

	distances := make([][]int, cols)
	for x := range distances {
		distances[x] = make([]int, rows)
		for y := range distances[x] {
			distances[x][y] = -1
		}
	}

	distances[pos.x][pos.y] = 0
	pointsOfInterest := []mapCoord{pos}

	floodForDist(board, distances, pointsOfInterest, 0)

	maxDist := 0
	for x := range distances {
		for y := range distances[x] {
			if distances[x][y] > maxDist {
				maxDist = distances[x][y]
			}
		}
	}

	return maxDist
}

func getMaxDepth(s string) int {
	depth := 0
	maxDepth := 0
	for _, c := range s {
		if c == '(' {
			depth++
		}
		if c == ')' {
			depth--
		}
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	// lines := utils.ReadLinesOrDie("sample_input1.txt")

	line := lines[0]

	maxDepth := getMaxDepth(line)
	fmt.Printf("max depth: %d\n", maxDepth)

	theRegex := parseRegex(line)
	fmt.Println(theRegex)
	fmt.Println("done parsing regex")

	reconstructed := reconstructRegexStr(theRegex)
	// fmt.Printf("original: %s\n", line)
	// fmt.Printf("reconstr: %s\n", reconstructed)
	if line == reconstructed {
		fmt.Println("original regex matches reconstructed")
	} else {
		fmt.Println("original regex DIFFERS FROM reconstructed")
	}

	fmt.Println(len(line))
	maxSideLen := len(line) * 2
	board := buildBoard(maxSideLen, maxSideLen)

	numNodes := countNodes(theRegex)
	fmt.Printf("num nodes: %d\n", numNodes)

	pos := mapCoord{maxSideLen / 2, maxSideLen / 2}
	traceRoutes(theRegex, board, pos)
	board, pos = trimBoard(board, pos)
	fmt.Printf("start pos: %v\n", pos)
	fmt.Printf("num visited: %d\n", numVisited)

	maxDist := findMaxDist(board, pos)
	fmt.Printf("max doors traversed: %v\n", maxDist)

	// fmt.Println()
	// fmt.Println(line)
	// fmt.Println()
	// printRegexTree(theRegex)
	fmt.Println()
	printBoard(board)
}
