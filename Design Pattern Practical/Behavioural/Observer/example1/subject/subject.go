package subject

import "observer/example1/observer"

type Subject interface {
	Register(ob observer.Observer)
	Deregister(ob observer.Observer)
	NotifyAll()
}

type Producer struct {
	ObserverList []observer.Observer
	Name         string
	InStock      bool
}

func NewProducer(name string, instock bool) *Producer {
	return &Producer{
		ObserverList: []observer.Observer{},
		Name:         name,
		InStock:      instock,
	}
}
func (p *Producer) Register(ob observer.Observer) {
	p.ObserverList = append(p.ObserverList, ob)
}
func (p *Producer) Deregister(ob observer.Observer) {
	for i, instance := range p.ObserverList {
		if instance.GetID() == ob.GetID() {
			p.ObserverList = append(p.ObserverList[:i], p.ObserverList[i+1:]...)
		}
	}
}

func (p *Producer) NotifyAll() {
	for _, observer := range p.ObserverList {
		observer.Update(p.Name)
	}
}
