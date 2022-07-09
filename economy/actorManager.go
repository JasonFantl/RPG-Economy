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
		actor.cooldown--
		if actor.expectedValues[ROCKET] < actor.personalValues[ROCKET] {
			buyers[actor] = true
		} else {
			sellers[actor] = true
		}
	}
	fmt.Printf("BUYERS: %d, SELLERS: %d\n", len(buyers), len(sellers))

	for actor := range actors {
		actor.Update()
	}
}

func ChangePVals(d float64) {
	for actor := range actors {
		actor.personalValues[ROCKET] += d
	}
}

func isBuyer(actor *Actor, good Good) bool {
	return (actor.expectedValues[good] < actor.personalValues[good]) && actor.cooldown <= 0
}

func isSeller(actor *Actor, good Good) bool {
	return actor.expectedValues[good] > actor.personalValues[good] && actor.cooldown <= 0
}
