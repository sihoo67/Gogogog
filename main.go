package main

import (
    "image/color"
    "math/rand"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/driver/desktop"
    "fyne.io/fyne/v2/widget"
)

const (
    gridSize   = 20
    cellSize   = 20
    width      = 30
    height     = 30
    tickMillis = 150
)

type direction int

const (
    up direction = iota
    down
    left
    right
)

type point struct {
    x, y int
}

type snakeGame struct {
    head      point
    dir       direction
    body      []point
    food      point
    canvas    *fyne.Container
    gameOver  bool
    scoreText *widget.Label
}

func newGame(canvas *fyne.Container, scoreText *widget.Label) *snakeGame {
    rand.Seed(time.Now().UnixNano())
    g := &snakeGame{
        head:      point{x: width / 2, y: height / 2},
        dir:       right,
        body:      []point{},
        canvas:    canvas,
        gameOver:  false,
        scoreText: scoreText,
    }
    g.spawnFood()
    return g
}

func (g *snakeGame) spawnFood() {
    g.food = point{x: rand.Intn(width), y: rand.Intn(height)}
}

func (g *snakeGame) move() {
    if g.gameOver {
        return
    }

    // body 따라감
    g.body = append([]point{g.head}, g.body...)
    if g.head == g.food {
        g.spawnFood()
    } else {
        g.body = g.body[:len(g.body)-1]
    }

    switch g.dir {
    case up:
        g.head.y--
    case down:
        g.head.y++
    case left:
        g.head.x--
    case right:
        g.head.x++
    }

    // 충돌 검사
    if g.head.x < 0 || g.head.x >= width || g.head.y < 0 || g.head.y >= height {
        g.gameOver = true
    }
    for _, b := range g.body {
        if b == g.head {
            g.gameOver = true
        }
    }

    g.draw()
}

func (g *snakeGame) draw() {
    g.canvas.Objects = nil

    // 뱀
    headRect := canvas.NewRectangle(color.RGBA{0, 255, 0, 255})
    headRect.Move(fyne.NewPos(float32(g.head.x*cellSize), float32(g.head.y*cellSize)))
    headRect.Resize(fyne.NewSize(cellSize, cellSize))
    g.canvas.Add(headRect)

    for _, b := range g.body {
        bodyRect := canvas.NewRectangle(color.RGBA{0, 200, 0, 255})
        bodyRect.Move(fyne.NewPos(float32(b.x*cellSize), float32(b.y*cellSize)))
        bodyRect.Resize(fyne.NewSize(cellSize, cellSize))
        g.canvas.Add(bodyRect)
    }

    // 음식
    foodRect := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
    foodRect.Move(fyne.NewPos(float32(g.food.x*cellSize), float32(g.food.y*cellSize)))
    foodRect.Resize(fyne.NewSize(cellSize, cellSize))
    g.canvas.Add(foodRect)

    if g.gameOver {
        text := canvas.NewText("Game Over", color.White)
        text.TextSize = 36
        text.Move(fyne.NewPos(180, 250))
        g.canvas.Add(text)
    }

    g.scoreText.SetText("점수: " + string(rune(len(g.body))))
    g.canvas.Refresh()
}

func main() {
    a := app.New()
    w := a.NewWindow("🐍 SnakeFyne")
    w.Resize(fyne.NewSize(width*cellSize, height*cellSize+30))

    canvas := container.NewWithoutLayout()
    scoreText := widget.NewLabel("점수: 0")

    game := newGame(canvas, scoreText)

    content := container.NewVBox(
        scoreText,
        canvas,
    )

    if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
        deskCanvas.SetOnKeyDown(func(ev *fyne.KeyEvent) {
            switch ev.Name {
            case fyne.KeyUp:
                if game.dir != down {
                    game.dir = up
                }
            case fyne.KeyDown:
                if game.dir != up {
                    game.dir = down
                }
            case fyne.KeyLeft:
                if game.dir != right {
                    game.dir = left
                }
            case fyne.KeyRight:
                if game.dir != left {
                    game.dir = right
                }
            }
        })
    }

    go func() {
        ticker := time.NewTicker(tickMillis * time.Millisecond)
        for range ticker.C {
            game.move()
        }
    }()

    w.SetContent(content)
    w.ShowAndRun()
}
