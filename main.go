package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jasonfantl/RPGEconomy/economy"
	"github.com/jasonfantl/RPGEconomy/gui"
)

type Game struct {
}

func (g *Game) Update() error {

	gui.Move()

	economy.Update()

	for _, p := range inpututil.AppendPressedKeys(make([]ebiten.Key, 1)) {
		if p == ebiten.KeyU {
			economy.DoStuff(1)
		}
		if p == ebiten.KeyD {
			economy.DoStuff(-1)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gui.Display(screen)
	economy.DrawGraphs(screen)
	// gui.DisplaySprite(1, 1, 2, screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gui.Dimensions()
}

func main() {
	rand.Seed(time.Now().Unix())
	wx, wy := gui.Dimensions()
	ebiten.SetWindowSize(wx, wy)
	ebiten.SetWindowTitle("Economy Sim")
	game := &Game{}

	gui.Initialize()
	economy.Initialize()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
