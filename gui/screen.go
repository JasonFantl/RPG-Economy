package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1400
	screenHeight = 1000
)

var (
	zoom       = 1.0
	xOff, yOff = 0, 0
)

func Initialize() {
	LoadSprites()
	generateSpace(20)
}

func Dimensions() (int, int) {
	return screenWidth, screenHeight
}

func Display(screen *ebiten.Image) {

	// TODO : use screen state to restrict which tiles we render, many unnecessarily rendered

	for i := 0; i < len(layout); i++ {
		for j := 0; j < len(layout[0]); j++ {
			DisplaySprite(i, j, layout[i][j], screen)
		}
	}
}

func DisplaySprite(x, y int, spriteID SpriteID, screen *ebiten.Image) {
	sx, sy := float64(x*tileSize)*zoom-float64(xOff), float64(y*tileSize)*zoom-float64(yOff)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(zoom, zoom) // TODO : fix the zooming
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(sprites[spriteID], op)
}

func Move() {
	// cursor move
	moveSpeed := 5
	activateMoveSpace := 0.1

	x, y := ebiten.CursorPosition()
	fx, fy := float64(x)/screenWidth, float64(y)/screenHeight

	if fx >= 0.0 && fx <= 1.0 && fy >= 0 && fy <= 1.0 {
		if fx < activateMoveSpace {
			xOff -= moveSpeed
		} else if fx > 1.0-activateMoveSpace {
			xOff += moveSpeed
		}
		if fy < activateMoveSpace {
			yOff -= moveSpeed
		} else if fy > 1.0-activateMoveSpace {
			yOff += moveSpeed
		}
	}

	// zooming
	_, w := ebiten.Wheel()
	if w > 0 {
		zoom += 0.1
	} else if w < 0 {
		zoom -= 0.1
	}
}
