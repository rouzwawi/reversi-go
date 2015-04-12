package reversi

import (
	"errors"
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
	Player       int
	ShowControls bool
	msg          string
	msgTimer     *time.Timer
	Anim         func()
	Draw         func()
	Clock        *Clock
}

func New() *Game {
	game := &Game{
		state:        make([]int, BOARD_SIZE*BOARD_SIZE),
		lines:        make([][][]int, BOARD_SIZE*BOARD_SIZE),
		Player:       P1,
		ShowControls: true,
		Clock:        NewClock(),
	}

	for i := range game.state {
		lines := make([][]int, DIRS)
		for d := E; d <= SE; d++ {
			i, j := Crd(i)
			lines[d] = line(d, i, j)
		}

		game.lines[i] = lines
		game.state[i] = EMPTY
	}
	half := BOARD_SIZE / 2
	game.state[Idx(half, half)] = P1
	game.state[Idx(half-1, half-1)] = P1
	game.state[Idx(half-1, half)] = P2
	game.state[Idx(half, half-1)] = P2

	return game
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

func (g *Game) Play(i, j int) {
	if !g.CanMove(i, j, g.Player) {
		return
	}

	plays := make([][]int, 0)
	lines := g.lines[Idx(i, j)]
	maxlen := 0
	for _, line := range lines {
		if g.canMoveLine(line, g.Player) {
			plays = append(plays, line)
			if len(line) > maxlen {
				maxlen = len(line)
			}
		}
	}

	g.ShowControls = false
	g.state[Idx(i, j)] = g.Player

	done := make([]bool, len(plays))
	nrem := len(plays)
	for i := 1; i < maxlen; i++ {
		for k, play := range plays {
			if len(play) <= i || done[k] {
				continue
			}
			p := play[i]
			if g.state[p] == g.Player {
				done[k] = true
				nrem--
				continue
			}
			g.state[p] = g.Player
		}
		if nrem == 0 {
			break
		}
		if g.Anim != nil {
			g.Anim()
		}
	}

	g.Player = other(g.Player)
	if !g.anyMoves(g.Player) {
		g.Player = other(g.Player)
	}

	g.ShowControls = true
}

func (g *Game) State() []int {
	s := make([]int, len(g.state))
	copy(s, g.state)
	return s
}

func (g *Game) Message() string {
	return g.msg
}

func (g *Game) CanMove(i, j, player int) bool {
	lines := g.lines[Idx(i, j)]
	for _, line := range lines {
		if g.canMoveLine(line, player) {
			return true
		}
	}

	return false
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

func (g *Game) anyMoves(player int) bool {
	for j := 0; j < BOARD_SIZE; j++ {
		for i := 0; i < BOARD_SIZE; i++ {
			if g.CanMove(i, j, player) {
				return true
			}
		}
	}
	return false
}

func (g *Game) triggerRefresh() {
	if g.Draw != nil {
		g.Draw()
	}
}

func (g *Game) SetMessage(msg string) {
	g.msg = msg
	g.triggerRefresh()

	if g.msgTimer != nil {
		g.msgTimer.Reset(time.Second)
	} else {
		g.msgTimer = time.NewTimer(time.Second)
		go func() {
			for t := range g.msgTimer.C {
				var _ = t
				g.msg = ""
				g.triggerRefresh()
			}
		}()
	}
}

func Idx(i, j int) int {
	return j*BOARD_SIZE + i
}

func Crd(i int) (int, int) {
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

func NextBound(d, i, j int) (int, int) {
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

	list = append(list, Idx(i, j))

	i, j = nxt(d, i, j)
	return _line(d, i, j, list)
}
