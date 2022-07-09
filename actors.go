package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jasonfantl/RPGEconomy/economy"
	"github.com/jasonfantl/RPGEconomy/gui"
)

type Tuple struct {
	x, y int
}
type Actor struct {
	economy  economy.Actor
	position Tuple
}

func NewActor() *Actor {
	return &Actor{
		economy:  *economy.NewActor(),
		position: Tuple{3, 3},
	}
}

func (actor *Actor) Update() {
	actor.position.x += rand.Intn(2) - 1
	actor.position.y += rand.Intn(2) - 1

}

func (actor Actor) Display(screen *ebiten.Image) {
	gui.DisplaySprite(actor.position.x, actor.position.y, gui.PERSON, screen)
}
