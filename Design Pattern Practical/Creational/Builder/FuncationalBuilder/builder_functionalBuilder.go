package main

import "fmt"

type Person struct {
	name, position string
}

type personMod func(*Person)
type PersonBuilder struct {
	actions []personMod
}

func (p *PersonBuilder) Caller(name string) *PersonBuilder {
	p.actions = append(p.actions, func(per *Person) {
		per.name = name
	})
	return p
}

func (p *PersonBuilder) WorkAsA(position string) *PersonBuilder {
	p.actions = append(p.actions, func(per *Person) {
		per.position = position
	})
	return p
}

func (p *PersonBuilder) Build() *Person {
	per := Person{}
	for _, action := range p.actions {
		action(&per)
	}
	return &per
}
func main() {
	b := PersonBuilder{}
	p := b.Caller("Dmitri").WorkAsA("dev").Build()
	fmt.Println(*p)
}
