


### APPROXIMATION 1:
The Output Value V. Input Good is not used directly, only its 
derivative is used. This barely less accurate representation fails when two identical people try to tell if they should exchange a good. This model will say they value the good at the same price, but the buyer should consider how much value they would gain, and the seller should consider how much they would lose. Instead they both consider the derivative exactly at the point they are at, essentially assuming constant returns even when thats not the case.

to fix we would use

``` go
func (actor *Actor) calcOutputValueVInputGood(good Good, n int) float64 {
	x := float64(n)
	d := float64(actor.desiredAssets[good])
	b := actor.basePersonalValues[good]

	return b * (x - x*x/(d*2)) // here we use positive diminishing returns followed by negative returns
}

func (actor *Actor) changeInValueFrom(good Good, n int) float64 {
	currentValue := actor.calcOutputValueVInputGood(good, actor.assets[good])
	newValue := actor.calcOutputValueVInputGood(good, actor.assets[good]+n)
	return newValue - currentValue
}
```

and use `changeInValueFrom(good, 1)` for buyers and `-changeInValueFrom(good, -1)` for sellers instead of `adjustedValue[good]`.
