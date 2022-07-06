package gui

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Image    *ebiten.Image
	Options  ebiten.DrawImageOptions
	Priority int
}

type SpriteID byte

const (
	GRASS_FLOOR SpriteID = iota
	DIRT_FLOOR
	PERSON
)

var sprites map[SpriteID]*ebiten.Image

const (
	tileSize = 8 // pixel size in sprite sheet
)

func LoadSprites() {
	sprites = make(map[SpriteID]*ebiten.Image, 0)

	tileSheet, _, err := ebitenutil.NewImageFromFile("data/spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}

	extractImage := func(x, y int) *ebiten.Image {
		x, y = x*tileSize, y*tileSize
		return tileSheet.SubImage(image.Rect(x, y, x+tileSize, y+tileSize)).(*ebiten.Image)
	}

	sprites[GRASS_FLOOR] = extractImage(0, 0)
	sprites[DIRT_FLOOR] = extractImage(1, 0)

	sprites[PERSON] = extractImage(0, 4)

}
