package economy

type Good string

const (
	MONEY = "Money"
	FOOD  = "Food"
	WOOD  = "Wood"
	HOUSE = "House"
)

var goods = []Good{MONEY, FOOD, WOOD, HOUSE}
