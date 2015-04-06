package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"errors"
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
	state  []int
	lines  [][][]int
	player int
}

func (g *Game) canMoveLine(line []int, player int) bool {
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

func (g *Game) canMove(i, j, player int) bool {
	lines := g.lines[idx(i, j)]
	for _, line := range lines {
		if g.canMoveLine(line, player) {
			return true
		}
	}

	return false
}

func (g *Game) anyMoves(player int) bool {
	for j := 0; j < BOARD_SIZE; j++ {
		for i := 0; i < BOARD_SIZE; i++ {
			if g.canMove(i, j, player) {
				return true
			}
		}
	}
	return false
}

func (g *Game) play(i, j int) {
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

func idx(i, j int) int {
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

func nxt(d, i, j int) (int, int) {
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
		panic(errors.New("unknown direction"))
	}
}

func nxtBound(d, i, j int) (int, int) {
	_i, _j := nxt(d, i, j)
	if bound(_i) && bound(_j) {
		return _i, _j
	} else {
		return i, j
	}
}

func line(d, i, j int) []int {
	return _line(d, i, j, nil)
}

func _line(d, i, j int, list []int) []int {
	if !bound(i) || !bound(j) {
		return list
	}

	list = append(list, idx(i, j))

	i, j = nxt(d, i, j)
	return _line(d, i, j, list)
}

func NewGame() *Game {
	game := &Game{
		state:  make([]int, BOARD_SIZE*BOARD_SIZE),
		lines:  make([][][]int, BOARD_SIZE*BOARD_SIZE),
		player: P1,
	}

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

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func mouse_button_num(k termbox.Key) int {
	switch k {
	case termbox.MouseLeft:
		return 0
	case termbox.MouseMiddle:
		return 1
	case termbox.MouseRight:
		return 2
	}
	return 0
}

func printGame(game *Game, ci, cj int) {
	const header = 3
	const c = termbox.ColorDefault
	const b = termbox.ColorCyan
	const g = termbox.ColorGreen
	const r = termbox.ColorRed
	var SYMBOLS = []string{" ", "●", "○", "+", "+"}
	var COLORS = []termbox.Attribute{c, b, r, c, c}

	var score [2]int

	// board numbers
	for i := 0; i < BOARD_SIZE; i++ {
		n := fmt.Sprintf("%d", i)
		tbprint(i*2+2, header-1, c, c, n)
		tbprint(i*2+2, BOARD_SIZE+header, c, c, fmt.Sprintf("%d", i))
		tbprint(0, i+header, c, c, n)
		tbprint(BOARD_SIZE*2+2, i+header, c, c, n)
	}

	// game state
	for j := 0; j < BOARD_SIZE; j++ {
		for i := 0; i < BOARD_SIZE; i++ {
			state := game.state[idx(i, j)]
			cl := COLORS[state]
			if state != EMPTY {
				score[state-1]++
			}

			if game.canMove(i, j, game.player) {
				if ci == i && cj == j {
					cl = COLORS[game.player]
					state = game.player
				} else {
					state = 2 + game.player
				}
			}
			symbol := SYMBOLS[state]

			tbprint(i*2+2, j + header, cl, c, symbol)
		}
	}

	// selector
	tbprint(ci*2+1, cj+header, COLORS[game.player], c, "[")
	tbprint(ci*2+3, cj+header, COLORS[game.player], c, "]")

	// header and score
	tbprint(4, 0, c, c, fmt.Sprintf("_ %2d - %-2d _", score[0], score[1]))
	tbprint(4, 0, COLORS[P1], c, SYMBOLS[P1])
	tbprint(14, 0, COLORS[P2], c, SYMBOLS[P2])

	pp := (map[int]int{P1:4, P2:14})[game.player]
	tbprint(pp, 0, COLORS[game.player] | termbox.AttrUnderline, c, SYMBOLS[game.player])
}

func main() {
	var curev termbox.Event
	i, j := 0, 0

	game := NewGame()

	const coldef = termbox.ColorDefault

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)

	termbox.Clear(coldef, coldef)
	printGame(game, i, j)
	termbox.Flush()

	data := make([]byte, 0, 64)
mainloop:
	for {
		if cap(data)-len(data) < 32 {
			newdata := make([]byte, len(data), len(data)+32)
			copy(newdata, data)
			data = newdata
		}
		beg := len(data)
		d := data[beg : beg+32]
		switch ev := termbox.PollRawEvent(d); ev.Type {
		case termbox.EventRaw:
			data = data[:beg+ev.N]
			curev = termbox.ParseEvent(data)
			if curev.N > 0 {
				copy(data, data[curev.N:])
				data = data[:len(data)-curev.N]
			}

			switch curev.Type {
			case termbox.EventKey:
				switch curev.Ch {
				case 'q':
					break mainloop
				}

				switch curev.Key {
				case 65514: // right
					i, j = nxtBound(E, i, j)
				case 65515: // left
					i, j = nxtBound(W, i, j)
				case 65516: // down
					i, j = nxtBound(S, i, j)
				case 65517: // up
					i, j = nxtBound(N, i, j)
				case 13: // enter
					game.play(i, j)
				}

				// tbprint(0, 2, coldef, coldef,
				// 	fmt.Sprintf("EventKey: k: %d, c: %c, mod: %d", curev.Key, curev.Ch, curev.Mod))
				// case termbox.EventMouse:
				// 	tbprint(0, 2, coldef, coldef,
				// 		fmt.Sprintf("EventMouse: x: %d, y: %d, b: %d", curev.MouseX, curev.MouseY, mouse_button_num(curev.Key)))
				// case termbox.EventNone:
				// 	tbprint(0, 2, coldef, coldef, "EventNone")
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		termbox.Clear(coldef, coldef)
		printGame(game, i, j)
		termbox.Flush()
	}
}
