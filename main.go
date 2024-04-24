package main

import (
	"fmt"
	"math"
	"time"
)

var (
	// Rules - matrix direction
	DIRECTIONS = map[string][2]int{"U": {-1, 0}, "D": {1, 0}, "L": {0, -1}, "R": {0, 1}}

	// State space
	START = [][]int{{1, 0, 2}, {4, 5, 3}, {7, 8, 6}}
	END   = [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 0}}

	// unicode for draw puzzle in command promt or terminal
	leftDownAngle  = '\u2514'
	rightDownAngle = '\u2518'
	rightUpAngle   = '\u2510'
	leftUpAngle    = '\u250C'

	middleJunction = '\u253C'
	topJunction    = '\u252C'
	bottomJunction = '\u2534'
	rightJunction  = '\u2524'
	leftJunction   = '\u251C'

	// source : https://pkg.go.dev/github.com/pborman/ansi
	// Each escape sequence is represented by a string constant and a Sequence structure
	boldText = "\033[1m"
	barStyle = "\033[1m\033[97m|\033[0m\033[0m"
	dash     = '\u2500'

	firstLine  = "\033[1m\033[97m" + string(leftUpAngle) + string(dash) + string(dash) + string(dash) + string(topJunction) + string(dash) + string(dash) + string(dash) + string(topJunction) + string(dash) + string(dash) + string(dash) + string(rightUpAngle) + "\033[0m\033[0m"
	middleLine = "\033[1m\033[97m" + string(leftJunction) + string(dash) + string(dash) + string(dash) + string(middleJunction) + string(dash) + string(dash) + string(dash) + string(middleJunction) + string(dash) + string(dash) + string(dash) + string(rightJunction) + "\033[0m\033[0m"
	lastLine   = "\033[1m\033[97m" + string(leftDownAngle) + string(dash) + string(dash) + string(dash) + string(bottomJunction) + string(dash) + string(dash) + string(dash) + string(bottomJunction) + string(dash) + string(dash) + string(dash) + string(rightDownAngle) + "\033[0m\033[0m"
)

type Node struct {
	currentNode  [][]int
	previousNode [][]int
	g            int
	h            int
	dir          string
}

func (n Node) f() int {
	return n.g + n.h
}

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

func euclideanCost(currentState [][]int) int {
	cost := 0
	for row := 0; row < len(currentState); row++ {
		for col := 0; col < len(currentState[0]); col++ {
			posRow, posCol := getPos(END, currentState[row][col])
			cost += int(math.Abs(float64(row-posRow)) + math.Abs(float64(col-posCol)))
		}
	}
	return cost
}

func getAdjNode(node Node) []Node {
	var listNode []Node
	emptyPosRow, emptyPosCol := getPos(node.currentNode, 0)

	for key, direction := range DIRECTIONS {
		newPosRow, newPosCol := emptyPosRow+direction[0], emptyPosCol+direction[1]
		if newPosRow >= 0 && newPosRow < len(node.currentNode) && newPosCol >= 0 && newPosCol < len(node.currentNode[0]) {
			newState := make([][]int, len(node.currentNode))
			for i := range node.currentNode {
				newState[i] = make([]int, len(node.currentNode[i]))
				copy(newState[i], node.currentNode[i])
			}
			newState[emptyPosRow][emptyPosCol] = node.currentNode[newPosRow][newPosCol]
			newState[newPosRow][newPosCol] = 0
			listNode = append(listNode, Node{currentNode: newState, previousNode: node.currentNode, g: node.g + 1, h: euclideanCost(newState), dir: key})
		}
	}
	return listNode
}

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

func reverseBranch(branches []map[string]interface{}) {
	for i := 0; i < len(branches)/2; i++ {
		j := len(branches) - i - 1
		branches[i], branches[j] = branches[j], branches[i]
	}
}

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

func printPuzzle(array [][]int) {
	for idx, row := range array {
		fmt.Print(barStyle)
		for _, val := range row {
			//  kalau dia blank space
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
	openSet := map[string]Node{fmt.Sprintf("%v", START): {currentNode: START, previousNode: START, g: 0, h: euclideanCost(START), dir: ""}}
	closedSet := make(map[string]Node)

	for {
		var testNode Node
		var ok bool

		testNode = getLowestCostNode(openSet)
		closedSet[fmt.Sprintf("%v", testNode.currentNode)] = testNode

		if fmt.Sprintf("%v", testNode.currentNode) == fmt.Sprintf("%v", END) {
			break
		}

		adjNode := getAdjNode(testNode)
		for _, node := range adjNode {
			nodeStr := fmt.Sprintf("%v", node.currentNode)
			if _, ok = closedSet[nodeStr]; ok || (ok && openSet[nodeStr].f() < node.f()) {
				continue
			}
			openSet[nodeStr] = node
		}
		delete(openSet, fmt.Sprintf("%v", testNode.currentNode))
	}

	// setup used node
	br := buildPath(closedSet)

	fmt.Println()
	fmt.Println("  | Start State |")
	for _, b := range br {
		if b["dir"].(string) != "" {
			letter := ""
			switch b["dir"].(string) {
			case "U":
				letter = "UP"
			case "R":
				letter = "RIGHT"
			case "L":
				letter = "LEFT"
			case "D":
				letter = "DOWN"
			}
			fmt.Printf(" -- | Move = %s | --\n", letter)
		}
		fmt.Println(firstLine)
		// di casting karna dia interface
		printPuzzle(b["node"].([][]int))
		fmt.Println()
	}
	fmt.Println("-- Achieve Goal State --")
	fmt.Println("Total steps:", len(br)-1)
	fmt.Println("Time taken to solve :", time.Since(startTime))
}
