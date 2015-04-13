package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"reversi"
	"time"
)

const BOARD_SIZE = reversi.BOARD_SIZE

func printGame(game *reversi.Game, ci, cj int) {
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
			state := game.State()[reversi.Idx(i, j)]
			cl := COLORS[state]
			if state != reversi.EMPTY {
				score[state-1]++
			}

			if game.ShowControls && game.CanMove(i, j, game.Player) {
				if ci == i && cj == j {
					cl = COLORS[game.Player]
					state = game.Player
				} else {
					state = 2 + game.Player
					cl = g
				}
			}
			symbol := SYMBOLS[state]

			tbprint(LEFT+i*2+2, j+header, cl, symbol)
		}
	}

	// selector
	tbprint(LEFT+ci*2+1, cj+header, COLORS[game.Player], "[")
	tbprint(LEFT+ci*2+3, cj+header, COLORS[game.Player], "]")

	// header and score
	tbprint(LEFT+4, 0, d, fmt.Sprintf("_ %2d - %-2d _", score[0], score[1]))
	tbprint(LEFT+4, 0, COLORS[reversi.P1], SYMBOLS[reversi.P1])
	tbprint(LEFT+14, 0, COLORS[reversi.P2], SYMBOLS[reversi.P2])

	pp := (map[int]int{reversi.P1: 4, reversi.P2: 14})[game.Player]
	tbprint(LEFT+pp, 0, COLORS[game.Player]|termbox.AttrUnderline, SYMBOLS[game.Player])

	// message
	msg := game.Message()
	state := game.State()
	if len(msg) == 0 {
		for i := 0; i < len(state); i += 2 {
			if i == len(state)/2 {
				msg += "-"
			}
			msg += fmt.Sprintf("%x", state[i]<<2|state[i+1])
		}
	}

	tbprint(LEFT+9-len(msg)/2, BOARD_SIZE+header+2, b, msg)

	mins := int(game.Clock.Duration.Minutes())
	secs := int(game.Clock.Duration.Seconds()) % 60
	deli := ":"
	if !game.Clock.Tick {
		deli = " "
	}
	time := fmt.Sprintf("%02d%s%02d", mins, deli, secs)
	tbprint(LEFT+9-len(time)/2, BOARD_SIZE+header+3, b, time)
}

func main() {
	const coldef = termbox.ColorDefault
	var curev termbox.Event
	i, j := 0, 0

	game := reversi.New()

	var refresh = func() {
		termbox.Clear(coldef, coldef)
		printGame(game, i, j)
		termbox.Flush()
	}
	var animFunc = func() {
		refresh()
		time.Sleep(100 * time.Millisecond)
	}
	game.Anim = animFunc
	game.Draw = refresh
	game.Clock.TickFunc = refresh

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
					if game.Anim == nil {
						game.Anim = animFunc
						game.SetMessage("animation on")
					} else {
						game.Anim = nil
						game.SetMessage("animation off")
					}
				case 'n':
					newGame := reversi.New()
					newGame.Anim = game.Anim
					game.Draw = refresh
					game.Clock.TickFunc = refresh
					game = newGame
				}

				switch curev.Key {
				case termbox.KeyArrowRight:
					i, j = reversi.NextBound(reversi.E, i, j)
				case termbox.KeyArrowLeft:
					i, j = reversi.NextBound(reversi.W, i, j)
				case termbox.KeyArrowDown:
					i, j = reversi.NextBound(reversi.S, i, j)
				case termbox.KeyArrowUp: // up
					i, j = reversi.NextBound(reversi.N, i, j)
				case termbox.KeyEnter:
					game.Play(i, j)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		refresh()
	}
}
