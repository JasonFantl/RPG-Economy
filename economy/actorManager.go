package economy

var actors []*Actor

func Initialize() {
	actors = make([]*Actor, 10)
	for i := 0; i < len(actors); i++ {
		actors[i] = NewActor()
	}
}

func Update() {

	for _, actor := range actors {
		actor.update()
	}
}

func makeExchange(buyer, seller *Actor, good Good, amount, cost int) {
	seller.assets[MONEY] += cost
	seller.assets[good] += amount
	buyer.assets[MONEY] -= cost
	buyer.assets[good] += amount
}
