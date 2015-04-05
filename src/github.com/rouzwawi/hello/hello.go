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
	player int
}

func (g Game) canMoveLine(line []int, player int) bool {
	if g.state[line[0]] != EMPTY {
		return false
	}

	other := other(player)
	seen_other := false
	for _, p := range line[1:] {
		switch g.state[p] {
		case EMPTY:
			return false
		case other:
			seen_other = true
		case player:
			return seen_other
		}
	}

	return false
}

func (g Game) canMove(i int, j int, player int) bool {
	lines := g.lines[idx(i, j)]
	for _, line := range lines {
		if g.canMoveLine(line, player) {
			return true
		}
	}

	return false
}

func (g Game) anyMoves(player int) bool {
	for j := 0; j < BOARD_SIZE; j++ {
		for i := 0; i < BOARD_SIZE; i++ {
			if g.canMove(i, j, player) {
				return true
			}
		}
	}
	return false
}

func (g *Game) play(i int, j int) {
	if !g.canMove(i, j, g.player) {
		return
	}

	plays := make([][]int, 0)
	lines := g.lines[idx(i, j)]
	for _, line := range lines {
		if g.canMoveLine(line, g.player) {
			plays = append(plays, line)
		}
	}

	for _, play := range plays {
		for i, p := range play {
			if i > 0 && g.state[p] == g.player {
				break
			}
			g.state[p] = g.player
		}
	}

	g.player = other(g.player)
	if !g.anyMoves(g.player) {
		g.player = other(g.player)
	}
}

func (g Game) print() {
	SYMBOLS := []int{' ', '●', '○', '+', '+'} // ◆◇

	fmt.Print(" ")
	for i := 0; i < BOARD_SIZE; i++ {
		fmt.Printf(" %d", i)
	}
	fmt.Println("")
	for j := 0; j < BOARD_SIZE; j++ {
		line := fmt.Sprintf("%d ", j)
		for i := 0; i < BOARD_SIZE; i++ {
			state := g.state[idx(i, j)]
			if g.canMove(i, j, g.player) {
				state = 2 + g.player
			}
			symbol := SYMBOLS[state]
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

func other(player int) int {
	if player == P1 {
		return P2
	} else {
		return P1
	}
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
	game.player = P1

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
	game.play(4,2)
	game.print()
	game.play(3,2)
	game.print()
}
