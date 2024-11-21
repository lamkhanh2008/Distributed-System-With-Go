package main

//component interface
type IPizza interface {
	getPrice() int
}

//concrete component
type VeggieMania struct {
}

func (p *VeggieMania) getPrice() int {
	return 13
}

//concrete decorator
type TomatoTopping struct {
	pizza IPizza
}

func (c *TomatoTopping) getPrice() int {
	pizzaPrice := c.pizza.getPrice()
	return pizzaPrice + 7
}

type CheeseTopping struct {
	pizza IPizza
}

func (c *CheeseTopping) getPrice() int {
	pizzaPrice := c.pizza.getPrice()
	return pizzaPrice + 10
}
