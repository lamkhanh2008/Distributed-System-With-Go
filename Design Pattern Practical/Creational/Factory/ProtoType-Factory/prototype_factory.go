package main

import "fmt"

type Employee struct {
	Name, Position string
	AnnualIncome   int
}

const (
	Developer = iota
	Manager
)

func NewEmployyee(role int) *Employee {
	switch role {
	case Developer:
		return &Employee{"", "Developer", 6000}
	case Manager:
		return &Employee{"", "Manager", 10000}
	default:
		panic("unsupport type")
	}
}

func main() {
	fmt.Println(Manager, Developer)
	m := NewEmployyee(Manager)
	m.Name = "Sam"
	fmt.Println(m)
}
