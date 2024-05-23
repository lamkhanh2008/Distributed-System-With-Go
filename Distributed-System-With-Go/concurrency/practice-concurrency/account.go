package main

import (
	"errors"
	"sync"
)

type Account struct {
	Id      string
	Balance int
	sync.Mutex
}

func (acc *Account) WithDraw(mount int) error {
	acc.Lock()
	defer acc.Unlock()
	if acc.Balance+mount < 0 {
		return errors.New("Tai khoan k du")
	}

	return nil
}

func Tranfer(from, to *Account, amt int) error {
	from.Lock()
	defer from.Unlock()
	to.Lock()
	defer to.Unlock()
	if from.Balance < amt {
		return errors.New("Tai khoan k du")
	}
	from.Balance -= amt
	to.Balance += amt
	return nil
}

// func main() {
// 	acc1 := Account{Id: "1", Balance: 10}
// 	acc2 := Account{Id: "2", Balance: 20}
// 	go Tranfer(&acc1, &acc2, 3)
// 	go Tranfer(&acc2, &acc1, 2)
// }
