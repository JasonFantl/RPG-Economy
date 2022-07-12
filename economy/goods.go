package economy

type Good string

const (
	MONEY  = "Money"
	FOOD   = "Food"
	WOOD   = "Wood"
	ROCKET = "Rocket"
)

var goods = []Good{MONEY, FOOD, WOOD, ROCKET}
