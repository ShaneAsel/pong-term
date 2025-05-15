package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

var screen tcell.Screen

const paddleHeight = 9
const paddleSpeed = 3
const cpuSpeed = 1
const cpuFrameDelay = 2
const winningScore = 3

type Paddle struct {
	y int
}

type Ball struct {
	x, y   int
	vx, vy int
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	screen = s
	defer screen.Fini()
	screen.Init()
	screen.Clear()

	w, h := screen.Size()
	gameOver := false
	winner := ""
	player := Paddle{y: h / 2}
	opponent := Paddle{y: h / 2}
	ball := Ball{x: w / 2, y: h / 2, vx: 1, vy: 1}
	playerScore := 0
	opponentScore := 0
	cpuFrameCounter := 0

	ticker := time.NewTicker(50 * time.Millisecond)
	go listenInput(&player, &ball, &opponent, &playerScore, &opponentScore, &gameOver, &winner, w, h)

	for range ticker.C {
		if !gameOver {
			update(&ball, &player, &opponent, w, h, &playerScore, &opponentScore, &cpuFrameCounter)

			if playerScore >= winningScore {
				gameOver = true
				winner = "You Win!"
			} else if opponentScore >= winningScore {
				gameOver = true
				winner = "CPU Wins!"
			}
		}

		draw(&ball, &player, &opponent, w, h, playerScore, opponentScore, gameOver, winner)
	}
}

func listenInput(p *Paddle, b *Ball, o *Paddle, playerScore, opponentScore *int, gameOver *bool, winner *string, w, h int) {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if p.y-paddleHeight/2 > 0 {
					p.y -= paddleSpeed
				}
			case tcell.KeyDown:
				if p.y+paddleHeight/2 < h-1 {
					p.y += paddleSpeed
				}
			case tcell.KeyEscape, tcell.KeyCtrlC:
				screen.Fini()
				log.Fatal("Game exited")
			case tcell.KeyEnter:
				if *gameOver {
					*playerScore = 0
					*opponentScore = 0
					*gameOver = false
					*winner = ""
					p.y = h / 2
					o.y = h / 2
					b.x, b.y = w/2, h/2
					b.vx = 1
					b.vy = 1
				}
			}
		}
	}
}

func update(b *Ball, p *Paddle, o *Paddle, w, h int, pScore *int, oScore *int, cpuFrameCounter *int) {
	if b.y <= 0 || b.y >= h-1 {
		b.vy *= -1
	}

	if b.x == 2 && b.y >= p.y-paddleHeight/2 && b.y <= p.y+paddleHeight/2 {
		b.vx *= -1
	}

	if b.x == w-3 && b.y >= o.y-paddleHeight/2 && b.y <= o.y+paddleHeight/2 {
		b.vx *= -1
	}

	if b.x <= 0 {
		*oScore++
		b.x, b.y = w/2, h/2
		b.vx = 1
	}

	if b.x >= w-1 {
		*pScore++
		b.x, b.y = w/2, h/2
		b.vx = -1
	}

	b.x += b.vx
	b.y += b.vy

	if *cpuFrameCounter%cpuFrameDelay == 0 {
		if b.y > o.y && o.y+paddleHeight/2 < h-1 {
			o.y += cpuSpeed
		} else if b.y < o.y && o.y-paddleHeight/2 > 0 {
			o.y -= cpuSpeed
		}
	}
}

func draw(b *Ball, p *Paddle, o *Paddle, w, h, pScore, oScore int, gameOver bool, winner string) {
	screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	scoreText := fmt.Sprintf("You: %d  CPU: %d", pScore, oScore)
	for i, ch := range scoreText {
		screen.SetContent(w/2-len(scoreText)/2+i, 0, ch, nil, style)
	}

	if !gameOver {
		screen.SetContent(b.x, b.y, '●', nil, style)

		for i := -paddleHeight / 2; i <= paddleHeight/2; i++ {
			screen.SetContent(1, p.y+i, '|', nil, style)
		}

		for i := -paddleHeight / 2; i <= paddleHeight/2; i++ {
			screen.SetContent(w-2, o.y+i, '|', nil, style)
		}
	} else {
		lines := []string{
			"GAME OVER - " + winner,
			"",
			"ESC to quit, ENTER to play again",
		}

		startY := h/2 - len(lines)/2
		for i, line := range lines {
			for j, ch := range line {
				screen.SetContent(w/2-len(line)/2+j, startY+i, ch, nil, style)
			}
		}
	}

	screen.Show()
}
