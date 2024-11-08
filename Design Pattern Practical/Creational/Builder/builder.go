package main

import "fmt"

type Builder interface {
	BuildPartA()
	BuildPartB()
	GetResult() Product
}
type Product struct {
	PartA string
	PartB string
}
type ConcreteBuilder1 struct {
	product Product
}
type Director struct {
	builder Builder
}

func (b *ConcreteBuilder1) BuildPartA() {
	b.product.PartA = "PartA1"
}

func (b *ConcreteBuilder1) BuildPartB() {
	b.product.PartB = "PartB1"
}
func (b *ConcreteBuilder1) GetResult() Product {
	return Product{
		PartA: b.product.PartA,
		PartB: b.product.PartB,
	}
}

func (d *Director) Construct() Product {
	d.builder.BuildPartA()
	d.builder.BuildPartB()
	return d.builder.GetResult()
}

func main() {
	builder := &ConcreteBuilder1{}
	director := Director{builder}
	product := director.Construct()
	fmt.Printf("Product: %+v\n", product)
}
