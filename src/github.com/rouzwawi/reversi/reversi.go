package main

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
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
	state        []int
	lines        [][][]int
	player       int
	showControls bool
	msg          string
	msgTime      time.Time
	anim         func()
	clock        *Clock
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

func (g *Game) setMessage(msg string) {
	g.msg = msg
	g.msgTime = time.Now()
}

func (g *Game) play(i, j int) {
	if !g.canMove(i, j, g.player) {
		return
	}

	plays := make([][]int, 0)
	lines := g.lines[idx(i, j)]
	maxlen := 0
	for _, line := range lines {
		if g.canMoveLine(line, g.player) {
			plays = append(plays, line)
			if len(line) > maxlen {
				maxlen = len(line)
			}
		}
	}

	g.showControls = false
	g.state[idx(i, j)] = g.player

	done := make([]bool, len(plays))
	nrem := len(plays)
	for i := 1; i < maxlen; i++ {
		for k, play := range plays {
			if len(play) <= i || done[k] {
				continue
			}
			p := play[i]
			if g.state[p] == g.player {
				done[k] = true
				nrem--
				continue
			}
			g.state[p] = g.player
		}
		if nrem == 0 {
			break
		}
		if g.anim != nil {
			g.anim()
		}
	}

	g.player = other(g.player)
	if !g.anyMoves(g.player) {
		g.player = other(g.player)
	}

	g.showControls = true
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

type Clock struct {
	ticker   *time.Ticker
	Tick     bool
	Duration time.Duration
	TickFunc func()
}

func NewClock() *Clock {
	clock := &Clock{
		ticker: time.NewTicker(time.Millisecond * 500),
		Tick:   true,
	}
	t0 := time.Now()

	go func() {
		for t := range clock.ticker.C {
			clock.Tick = !clock.Tick
			clock.Duration = t.Sub(t0)
			if clock.TickFunc != nil {
				clock.TickFunc()
			}
		}
	}()

	return clock
}

func New() *Game {
	game := &Game{
		state:        make([]int, BOARD_SIZE*BOARD_SIZE),
		lines:        make([][][]int, BOARD_SIZE*BOARD_SIZE),
		player:       P1,
		showControls: true,
		clock:        NewClock(),
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

func printGame(game *Game, ci, cj int) {
	const header = 3
	const d = termbox.ColorDefault
	const b = termbox.ColorBlue
	const m = termbox.ColorMagenta
	const g = termbox.ColorGreen
	const y = termbox.ColorYellow
	const r = termbox.ColorRed

	tbprint := func(x, y int, fg termbox.Attribute, msg string) {
		for _, c := range msg {
			termbox.SetCell(x, y, c, fg, termbox.ColorDefault)
			x++
		}
	}

	WX, _ := termbox.Size()
	LEFT := WX/2 - 10
	SYMBOLS := []string{" ", "●", "●", "+", "+"} // ○
	COLORS := []termbox.Attribute{d, b, r, d, d}

	var score [2]int

	// board numbers
	for i := 0; i < BOARD_SIZE; i++ {
		n := fmt.Sprintf("%d", i)
		tbprint(LEFT+i*2+2, header-1, d, n)
		tbprint(LEFT+i*2+2, BOARD_SIZE+header, d, fmt.Sprintf("%d", i))
		tbprint(LEFT, i+header, d, n)
		tbprint(LEFT+BOARD_SIZE*2+2, i+header, d, n)
	}

	// game state
	for j := 0; j < BOARD_SIZE; j++ {
		for i := 0; i < BOARD_SIZE; i++ {
			state := game.state[idx(i, j)]
			cl := COLORS[state]
			if state != EMPTY {
				score[state-1]++
			}

			if game.showControls && game.canMove(i, j, game.player) {
				if ci == i && cj == j {
					cl = COLORS[game.player]
					state = game.player
				} else {
					state = 2 + game.player
					cl = g
				}
			}
			symbol := SYMBOLS[state]

			tbprint(LEFT+i*2+2, j+header, cl, symbol)
		}
	}

	// selector
	tbprint(LEFT+ci*2+1, cj+header, COLORS[game.player], "[")
	tbprint(LEFT+ci*2+3, cj+header, COLORS[game.player], "]")

	// header and score
	tbprint(LEFT+4, 0, d, fmt.Sprintf("_ %2d - %-2d _", score[0], score[1]))
	tbprint(LEFT+4, 0, COLORS[P1], SYMBOLS[P1])
	tbprint(LEFT+14, 0, COLORS[P2], SYMBOLS[P2])

	pp := (map[int]int{P1: 4, P2: 14})[game.player]
	tbprint(LEFT+pp, 0, COLORS[game.player]|termbox.AttrUnderline, SYMBOLS[game.player])

	// message
	if len(game.msg) == 0 {
		for i := 0; i < len(game.state); i += 2 {
			if i == len(game.state)/2 {
				game.msg += "-"
			}
			game.msg += fmt.Sprintf("%x", game.state[i]<<2|game.state[i+1])
		}
	}

	tbprint(LEFT+9-len(game.msg)/2, BOARD_SIZE+header+2, b, game.msg)
	if time.Now().Sub(game.msgTime).Seconds() > 1 {
		game.msg = ""
	}

	mins := int(game.clock.Duration.Minutes())
	secs := int(game.clock.Duration.Seconds()) % 60
	deli := ":"
	if !game.clock.Tick {
		deli = " "
	}
	time := fmt.Sprintf("%02d%s%02d", mins, deli, secs)
	tbprint(LEFT+9-len(time)/2, BOARD_SIZE+header+3, b, time)
}

func main() {
	const coldef = termbox.ColorDefault
	var curev termbox.Event
	i, j := 0, 0

	game := New()

	var refresh = func() {
		termbox.Clear(coldef, coldef)
		printGame(game, i, j)
		termbox.Flush()
	}
	var animFunc = func() {
		refresh()
		time.Sleep(100 * time.Millisecond)
	}
	game.anim = animFunc
	game.clock.TickFunc = refresh

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	refresh()

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
				case 'a':
					if game.anim == nil {
						game.anim = animFunc
						game.setMessage("animation on")
					} else {
						game.anim = nil
						game.setMessage("animation off")
					}
				case 'n':
					newGame := New()
					newGame.anim = game.anim
					game = newGame
				}

				switch curev.Key {
				case termbox.KeyArrowRight:
					i, j = nxtBound(E, i, j)
				case termbox.KeyArrowLeft:
					i, j = nxtBound(W, i, j)
				case termbox.KeyArrowDown:
					i, j = nxtBound(S, i, j)
				case termbox.KeyArrowUp: // up
					i, j = nxtBound(N, i, j)
				case termbox.KeyEnter:
					game.play(i, j)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		refresh()
	}
}
