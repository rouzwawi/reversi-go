package main

import (
	"fmt"
	"runtime"
	// "github.com/rouzwawi/newmath"
	// "github.com/edsrzf/mmap-go"
)

const BOARD_SIZE = 8

const (
	EMPTY = iota
	P1
	P2
)

const DIRS = 8
const (
	E = iota
	NE
	N
	NW
	W
	SW
	S
	SE
)

type Game struct {
	state []int
	lines [][][]int
}

func (g Game) print() {
	SYMBOLS := []int{' ', 'x', 'o'}

	fmt.Print(" ")
	for i := 0; i < BOARD_SIZE; i++ {
		fmt.Printf(" %d", i)
	}
	fmt.Println("")
	for j := 0; j < BOARD_SIZE; j++ {
		line := fmt.Sprintf("%d ", j)
		for i := 0; i < BOARD_SIZE; i++ {
			symbol := SYMBOLS[g.state[idx(i, j)]]
			line += fmt.Sprintf("%c ", symbol)
		}
		fmt.Println(line)
	}
}


func idx(i int, j int) int {
	return j*BOARD_SIZE + i
}

func crd(i int) (int, int) {
	return i % BOARD_SIZE, i / BOARD_SIZE
}

func bound(i int) bool {
	return i >= 0 && i < BOARD_SIZE
}

func nxt(d int, i int, j int) (int, int) {
	switch d {
	case E:
		return i + 1, j
	case NE:
		return i + 1, j - 1
	case N:
		return i, j - 1
	case NW:
		return i - 1, j - 1
	case W:
		return i - 1, j
	case SW:
		return i - 1, j + 1
	case S:
		return i, j + 1
	case SE:
		return i + 1, j + 1
	default:
		panic(new(runtime.Error))
	}
}

func line(d int, i int, j int) []int {
	return _line(d, i, j, nil)
}

func _line(d int, i int, j int, list []int) []int {
	if !bound(i) || !bound(j) {
		return list
	}

	list = append(list, idx(i, j))

	i, j = nxt(d, i, j)
	return _line(d, i, j, list)
}

func newGame() *Game {
	game := new(Game)
	game.state = make([]int, BOARD_SIZE*BOARD_SIZE)
	game.lines = make([][][]int, BOARD_SIZE*BOARD_SIZE)

	for i := range game.state {
		lines := make([][]int, DIRS)
		for d := E; d <= SE; d++ {
			i, j := crd(i)
			lines[d] = line(d, i, j)
		}

		game.lines[i] = lines
		game.state[i] = EMPTY
	}
	half := BOARD_SIZE / 2
	game.state[idx(half, half)] = P1
	game.state[idx(half-1, half-1)] = P1
	game.state[idx(half-1, half)] = P2
	game.state[idx(half, half-1)] = P2

	return game
}

func main() {
	game := newGame()
	game.print()
}
