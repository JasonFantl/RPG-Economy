package economy

import (
	"fmt"
)

var actors map[*Actor]bool

func Initialize() {
	actors = make(map[*Actor]bool)
	for i := 0; i < 100; i++ {
		actors[NewActor()] = true
	}
}

func Update() {

	sellers := make(map[*Actor]bool)
	buyers := make(map[*Actor]bool)

	for actor := range actors {
		if actor.isBuyer(ROCKET) && actor.canBuy(1) {
			buyers[actor] = true
		} else if actor.isSeller(ROCKET) && actor.canSell(ROCKET) {
			sellers[actor] = true
		}
	}
	fmt.Printf("BUYERS: %d, SELLERS: %d, NON: %d\n", len(buyers), len(sellers), len(actors)-(len(buyers)+len(sellers)))

	for actor := range actors {
		actor.Update()
	}
}

func DoStuff(d float64) {
	for actor := range actors {
		actor.assets[ROCKET] += int(d)
		// actor.personalValues[ROCKET] += d

	}
}
