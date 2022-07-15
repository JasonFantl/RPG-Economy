package economy

import (
	"fmt"
	"math"
	"math/rand"
)

type marketInfo struct {
	basePersonalValue     float64 // if they had no goods, how much would they be willing to pay for this good
	adjustedPersonalValue float64 // changes based on if they have a good

	expectedValue float64

	ownedAssets   int
	desiredAssets int

	failedTransactionAttempts int // how many times we have to fail to buy/sell
	failedTransactionThresh   int // how many times we have to fail to buy/sell before updating our expected value
	priceSignalReactivity     float64
}

func newMarket(own int, desiredBottom, desiredTop int, personalValueBottom float64) *marketInfo {
	return &marketInfo{
		basePersonalValue:       rand.Float64()*personalValueBottom + personalValueBottom,
		expectedValue:           0,
		ownedAssets:             own,
		desiredAssets:           rand.Intn(desiredTop-desiredBottom+1) + desiredBottom,
		failedTransactionThresh: 4,
		priceSignalReactivity:   1.0,
	}
}

func (actor *Actor) runMarket(good Good) {
	if rand.Float64() > 0.01 { // throw some randomness in (helps for non-sellers/buyers who need to interact with the market)
		actor.markets[good].expectedValue += rand.Float64() - 0.5
	}

	// so we don't have to re-calculate every time
	actor.calcAdjustedValue(good)

	// we assume buyers are the initiators, so we ignore sellers
	if actor.isBuyer(good) {
		willingBuyPrice := actor.willingBuyPrice(good)
		// look for a seller, simulates going from shop to shop
		var seller *Actor = nil
		sellingPrice := 0
		for otherActor := range actors { // rely on random iteration
			// do we want to update our expected value after each failed shop visit, or only after exhausting every shop? Doing only after every for now
			if !otherActor.isSeller(good) || !otherActor.canSell(good) { // must be a seller with goods to sell
				continue
			}
			sellingPrice = otherActor.willingSellPrice(good) // looking at the price tag

			if sellingPrice >= 0 {
				if willingBuyPrice < sellingPrice { // price is too high
					continue
				}
				if !actor.canBuy(sellingPrice) { // buyer doesn't have the money
					continue
				}
			} else { // negative sell price
				if willingBuyPrice < sellingPrice { // price is too high
					continue
				}
				if !otherActor.canBuy(-sellingPrice) { // seller doesn't have the money
					continue
				}
			}

			// made it past all the checks, this is someone we can buy from
			seller = otherActor
			break
		}

		if seller != nil { // found a seller we can buy from
			// update expected values and asset count
			transact(actor, seller, good, sellingPrice)
		} else {
			actor.markets[good].failedTransactionAttempts++
		}
	}

	// sellers need to update if they haven't made a sale in a while
	if actor.isSeller(good) && actor.canSell(good) {
		actor.markets[good].failedTransactionAttempts++
	}

	if actor.markets[good].failedTransactionAttempts > actor.markets[good].failedTransactionThresh { // haven't transacted in a while, update expected values
		actor.markets[good].failedTransactionAttempts = 0
		if actor.isBuyer(good) {
			actor.markets[good].expectedValue += actor.markets[good].priceSignalReactivity
		} else if actor.isSeller(good) {
			actor.markets[good].expectedValue -= actor.markets[good].priceSignalReactivity
		}
	}
}

// approximation, see README: APPROXIMATION 1
func (actor *Actor) calcAdjustedValue(good Good) {
	x := float64(actor.markets[good].ownedAssets)
	d := float64(actor.markets[good].desiredAssets)
	b := actor.markets[good].basePersonalValue

	if d < 1 {
		actor.markets[good].adjustedPersonalValue = 0
	} else {
		// linear
		// actor.markets[good].adjustedPersonalValue = b * (1 - x/d)
		// exponential
		// actor.markets[good].adjustedPersonalValue = b * (1 - x*x/(d*d))
		// atan
		actor.markets[good].adjustedPersonalValue = (b / 2) * (1 - 1.5*math.Atan(math.Pi*0.5*(x/d-0.5)))

	}
}

func transact(buyer, seller *Actor, good Good, cost int) {
	seller.markets[good].failedTransactionAttempts = 0
	buyer.markets[good].failedTransactionAttempts = 0

	seller.markets[good].ownedAssets--
	buyer.markets[good].ownedAssets++

	seller.markets[good].expectedValue += seller.markets[good].priceSignalReactivity
	buyer.markets[good].expectedValue -= buyer.markets[good].priceSignalReactivity

	seller.money += cost
	buyer.money -= cost

	if buyer.money < 0 || seller.money < 0 {
		fmt.Printf("buyer money: %d, good: %s, cost: %d\n", buyer.money, good, cost)
		fmt.Printf("seller money: %d, good: %s, cost: %d\n", seller.money, good, cost)
	}
}

func (actor *Actor) willingBuyPrice(good Good) int {
	return int(math.Floor(math.Min(actor.markets[good].expectedValue, actor.markets[good].adjustedPersonalValue)))
}

func (actor *Actor) willingSellPrice(good Good) int {
	return int(math.Ceil(math.Max(actor.markets[good].expectedValue, actor.markets[good].adjustedPersonalValue)))
}

func (actor *Actor) isBuyer(good Good) bool {
	return actor.markets[good].expectedValue < actor.markets[good].adjustedPersonalValue
}

func (actor *Actor) isSeller(good Good) bool {
	return actor.markets[good].expectedValue > actor.markets[good].adjustedPersonalValue
}

func (actor *Actor) canBuy(cost int) bool {
	return actor.money >= cost
}

func (actor *Actor) canSell(good Good) bool {
	return actor.markets[good].ownedAssets > 0
}
