package main

import (
	"bytes"
	_ "embed"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	worldScale   = 1.45
	marioScale   = 0.60
	moveSpeed    = 4.0
	gravity      = 0.50
	jumpStrength = 11.0
	worldSpeed   = 0.65
	worldOffsetY = -145.0
	groundY      = 408.0
	worldLoopW   = 2550.0
)

type Game struct {
	backgroundImage *ebiten.Image
	marioImage      *ebiten.Image
	backgroundX     float64
	marioX          float64
	marioY          float64
	marioWidth      float64
	marioHeight     float64
	velocityY       float64
	isJumping       bool
	facingLeft      bool
}

//go:embed assets/background.png
var backgroundPNG []byte

//go:embed assets/mario.png
var marioPNG []byte

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.marioX -= moveSpeed
		g.facingLeft = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.marioX += moveSpeed
		g.facingLeft = false
	}

	if g.marioX < 0 {
		g.marioX = 0
	}
	maxX := float64(screenWidth) - g.marioWidth
	if g.marioX > maxX {
		g.marioX = maxX
	}

	if (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp)) && !g.isJumping {
		g.velocityY = -jumpStrength
		g.isJumping = true
	}

	g.velocityY += gravity
	g.marioY += g.velocityY

	groundMarioY := groundY - g.marioHeight
	if g.marioY >= groundMarioY {
		g.marioY = groundMarioY
		g.velocityY = 0
		g.isJumping = false
	}

	g.backgroundX += worldSpeed
	if g.backgroundX >= worldLoopW {
		g.backgroundX -= worldLoopW
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{135, 206, 235, 255})

	for x := -g.backgroundX; x < float64(screenWidth); x += worldLoopW {
		bgOp := &ebiten.DrawImageOptions{}
		bgOp.GeoM.Scale(worldScale, worldScale)
		bgOp.GeoM.Translate(x, worldOffsetY)
		screen.DrawImage(g.backgroundImage, bgOp)
	}

	marioOptions := &ebiten.DrawImageOptions{}
	if g.facingLeft {
		marioOptions.GeoM.Scale(-marioScale, marioScale)
		marioOptions.GeoM.Translate(g.marioX+g.marioWidth, g.marioY)
	} else {
		marioOptions.GeoM.Scale(marioScale, marioScale)
		marioOptions.GeoM.Translate(g.marioX, g.marioY)
	}
	screen.DrawImage(g.marioImage, marioOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}

	var err error

	game.backgroundImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(backgroundPNG))
	if err != nil {
		log.Fatal(err)
	}

	game.marioImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(marioPNG))
	if err != nil {
		log.Fatal(err)
	}

	game.marioWidth = float64(game.marioImage.Bounds().Dx()) * marioScale
	game.marioHeight = float64(game.marioImage.Bounds().Dy()) * marioScale
	game.marioX = 40
	game.marioY = groundY - game.marioHeight

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mario Prototype in Go")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
