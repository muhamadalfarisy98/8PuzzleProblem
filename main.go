package main

import (
	"fmt"
	"math"
	"time"
)

var (
	// Rules - matrix direction
	DIRECTIONS = map[string][2]int{
		"U": {-1, 0},
		"D": {1, 0},
		"L": {0, -1},
		"R": {0, 1},
	}

	UP    = "UP"
	RIGHT = "RIGHT"
	DOWN  = "DOWN"
	LEFT  = "LEFT"

	ActionMoves = map[string]string{
		"U": UP,
		"R": RIGHT,
		"L": LEFT,
		"D": DOWN,
	}

	// State space
	START = [][]int{
		{1, 0, 2},
		{4, 5, 3},
		{7, 8, 6},
	}

	END = [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 0},
	}

	// unicode for draw puzzle in command prompt or terminal
	// source : https://pkg.go.dev/github.com/appgate-sdp-int/tabulate#section-readme (grid format)
	leftDownAngle  = '\u2514'
	rightDownAngle = '\u2518'
	rightUpAngle   = '\u2510'
	leftUpAngle    = '\u250C'

	middleJunction = '\u253C'
	topJunction    = '\u252C'
	bottomJunction = '\u2534'
	rightJunction  = '\u2524'
	leftJunction   = '\u251C'

	dash = '\u2500'

	// source : https://pkg.go.dev/github.com/pborman/ansi
	// Each escape sequence is represented by a string constant and a Sequence structure

	// text style
	boldText       = "\033[1m"
	boldTextYellow = "\033[1;33m"

	barStyle = "\033[1m\033[97m|\033[0m\033[0m"

	// line style
	firstLine  = "\033[1m\033[97m" + string(leftUpAngle) + string(dash) + string(dash) + string(dash) + string(topJunction) + string(dash) + string(dash) + string(dash) + string(topJunction) + string(dash) + string(dash) + string(dash) + string(rightUpAngle) + "\033[0m\033[0m"
	middleLine = "\033[1m\033[97m" + string(leftJunction) + string(dash) + string(dash) + string(dash) + string(middleJunction) + string(dash) + string(dash) + string(dash) + string(middleJunction) + string(dash) + string(dash) + string(dash) + string(rightJunction) + "\033[0m\033[0m"
	lastLine   = "\033[1m\033[97m" + string(leftDownAngle) + string(dash) + string(dash) + string(dash) + string(bottomJunction) + string(dash) + string(dash) + string(dash) + string(bottomJunction) + string(dash) + string(dash) + string(dash) + string(rightDownAngle) + "\033[0m\033[0m"
)

// Node - stores state space
type Node struct {
	currentNode  [][]int
	previousNode [][]int
	g            int
	h            int
	dir          string
}

// f - calculates heuristic function
func (n Node) f() int {
	return n.g + n.h
}

// getAdjNode - gets possible adjacent nodes
func (n Node) getAdjNode() []Node {
	var listNode []Node
	emptyPosRow, emptyPosCol := getPos(n.currentNode, 0)

	for key, direction := range DIRECTIONS {
		// define pointer of new replaced zero (empty) tile based on direction value
		newPosRow, newPosCol := emptyPosRow+direction[0], emptyPosCol+direction[1]
		if newPosRow >= 0 && newPosRow < len(n.currentNode) && newPosCol >= 0 && newPosCol < len(n.currentNode[0]) {
			newState := make([][]int, len(n.currentNode))
			for i := range n.currentNode {
				newState[i] = make([]int, len(n.currentNode[i]))
				copy(newState[i], n.currentNode[i])
			}
			newState[emptyPosRow][emptyPosCol] = n.currentNode[newPosRow][newPosCol]
			newState[newPosRow][newPosCol] = 0
			listNode = append(listNode, Node{currentNode: newState, previousNode: n.currentNode, g: n.g + 1, h: manhattanDistanceCost(newState), dir: key})
		}
	}
	return listNode
}

// getPos - get element's current position on 2D array
func getPos(currentState [][]int, element int) (int, int) {
	for row := 0; row < len(currentState); row++ {
		for col := 0; col < len(currentState[0]); col++ {
			if currentState[row][col] == element {
				return row, col
			}
		}
	}
	return -1, -1
}

// manhattanDistanceCost - calculates total distance between current state and goal state
func manhattanDistanceCost(currentState [][]int) int {
	cost := 0
	for row := 0; row < len(currentState); row++ {
		for col := 0; col < len(currentState[0]); col++ {
			posRow, posCol := getPos(END, currentState[row][col])
			cost += int(math.Abs(float64(row-posRow)) + math.Abs(float64(col-posCol)))
		}
	}
	return cost
}

// buildPath - build closed(visited node) set
func buildPath(closedSet map[string]Node) []map[string]interface{} {
	node := closedSet[fmt.Sprintf("%v", END)]
	branch := make([]map[string]interface{}, 0)

	for node.dir != "" {
		branch = append(branch, map[string]interface{}{"dir": node.dir, "node": node.currentNode})
		node = closedSet[fmt.Sprintf("%v", node.previousNode)]
	}
	branch = append(branch, map[string]interface{}{"dir": "", "node": node.currentNode})

	reverseBranch(branch)
	return branch
}

// reverseBranch - helps to reverse array
func reverseBranch(branches []map[string]interface{}) {
	for i := 0; i < len(branches)/2; i++ {
		j := len(branches) - i - 1
		branches[i], branches[j] = branches[j], branches[i]
	}
}

// getLowestCostNode - get minimum cost node path
func getLowestCostNode(openSet map[string]Node) Node {
	var testNode Node
	bestF := math.MaxInt64

	for _, node := range openSet {
		if testNode.currentNode == nil || node.f() < bestF {
			testNode = node
			bestF = node.f()
		}
	}
	return testNode
}

// printPuzzle - print array node
func printPuzzle(array [][]int) {
	for idx, row := range array {
		fmt.Print(barStyle)
		for _, val := range row {
			// blank tiles
			if val == 0 {
				fmt.Print("\033[103m   \033[0m\033[97m|\033[0m")
				continue
			}
			fmt.Printf(" %d \033[97m|\033[0m", val)
		}
		fmt.Println()
		if idx == 2 {
			fmt.Println(lastLine)
			continue
		}
		fmt.Println(middleLine)
	}
}

func main() {
	fmt.Println("--- Welcome to 8 Puzzle Problem Solver ---")
	fmt.Println()
	fmt.Println("Predefined start space :", START)
	fmt.Println("Predefined goal space :", END)
	fmt.Println()

	var menuInput int
	for {
		fmt.Println("Menu :")
		fmt.Println("1. Gunakan Predefined state space")
		fmt.Println("2. Input manual state space")
		fmt.Println("Masukkan no. pilihan :")
		fmt.Scanf("%d", &menuInput)

		if menuInput == 1 || menuInput == 2 {
			break
		}
		fmt.Println("no pilihan menu tidak tersedia, silahkan coba lagi!")
		fmt.Println()
	}

	arrStart := make([][]int, 3)
	arrEnd := make([][]int, 3)
	if menuInput == 2 {
		// input state
		fmt.Println("Enter Start Space 2-D Array:")
		for row := 0; row < 3; row++ {
			arrStart[row] = make([]int, 3)
			for col := 0; col < 3; col++ {
				fmt.Scanf("%d", &arrStart[row][col])
			}
		}
		START = arrStart
		fmt.Println("isi start arraynya : ", arrStart)

		// goal state
		var goalInput int
		for {
			fmt.Println("1. Gunakan Predefined goal state")
			fmt.Println("2. Input manual goal space")
			fmt.Println("Masukkan no. pilihan :")
			fmt.Scanf("%d", &goalInput)
			if menuInput == 1 || menuInput == 2 {
				break
			}
			fmt.Println("no pilihan menu tidak tersedia, silahkan coba lagi!")
			fmt.Println()
		}

		if goalInput == 2 {
			fmt.Println("Enter End Space 2-D Array:")
			for row := 0; row < 3; row++ {
				arrEnd[row] = make([]int, 3)
				for col := 0; col < 3; col++ {
					fmt.Scanf("%d", &arrEnd[row][col])
				}
			}
			END = arrEnd

			fmt.Println("isi end arraynya : ", END)
		}
	}

	startTime := time.Now()

	// main logic
	openSet := map[string]Node{fmt.Sprintf("%v", START): {currentNode: START, previousNode: START, g: 0, h: manhattanDistanceCost(START), dir: ""}}
	closedSet := make(map[string]Node)

	for {
		var testNode Node
		var ok bool

		testNode = getLowestCostNode(openSet)
		closedSet[fmt.Sprintf("%v", testNode.currentNode)] = testNode

		if fmt.Sprintf("%v", testNode.currentNode) == fmt.Sprintf("%v", END) {
			break
		}

		adjNode := testNode.getAdjNode()
		for _, node := range adjNode {
			nodeStr := fmt.Sprintf("%v", node.currentNode)

			// check duplicate or existing node
			if _, ok = closedSet[nodeStr]; ok || (ok && openSet[nodeStr].f() < node.f()) {
				continue
			}
			openSet[nodeStr] = node
		}
		// removes visited node
		delete(openSet, fmt.Sprintf("%v", testNode.currentNode))
	}

	// setup used node
	br := buildPath(closedSet)

	fmt.Println()
	fmt.Println(boldText, "  | Start State |")
	for _, b := range br {
		if b["dir"].(string) != "" {
			letter := ""
			if val, ok := ActionMoves[b["dir"].(string)]; ok {
				letter = val
			}
			fmt.Printf(" Move => %s \n", (boldTextYellow + letter))
		}
		fmt.Println(firstLine)

		printPuzzle(b["node"].([][]int))
		fmt.Println()
	}

	fmt.Println(boldText, "-- Achieved Goal State --")
	fmt.Println("Total steps:", len(br)-1)
	fmt.Println("Time taken to solve :", time.Since(startTime))
	fmt.Println()
}
