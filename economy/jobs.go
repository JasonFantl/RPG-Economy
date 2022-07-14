package economy

import "math/rand"

type Job string

const (
	WOODWORKER = "Woodworker"
	FARMER     = "Farmer"
)

var jobToGood = map[Job]Good{
	WOODWORKER: WOOD,
	FARMER:     FOOD,
}
var goodToJob = map[Good]Job{
	WOOD: WOODWORKER,
	FOOD: FARMER,
}

func (actor *Actor) runJob(job Job) {
	if rand.Float64() > 0.1 {
		return
	}

	y := 1.0 / (float64(actor.skills[job] + 1))
	if rand.Float64() > y {
		actor.markets[jobToGood[job]].ownedAssets++
	}
}
