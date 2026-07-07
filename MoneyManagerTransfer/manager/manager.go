package manager

import (
	"fmt"
	"math/rand"
	"mnogo/bank"
	"sync"
	"time"
)

type Manager struct {
	accounts []*bank.Account
	find     map[string]*bank.Account
}

func NewManager(accounts []*bank.Account) *Manager {
	find := make(map[string]*bank.Account)
	for _, acc := range accounts {
		find[acc.ID()] = acc
	}
	return &Manager{
		accounts: accounts,
		find:     find,
	}
}
func Transfer(from, to *bank.Account, mon int, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))

	if from.ID() < to.ID() {
		from.Lock()
		to.Lock()
	} else {
		to.Lock()
		from.Lock()
	}
	defer from.Unlock()
	defer to.Unlock()

	if from.GetMoney() < mon {
		fmt.Printf("%s -> %s: не хватает денег (есть %d, нужно %d)\n",
			from.ID(), to.ID(), from.GetMoney(), mon)
		return
	}

	from.SubMoney(mon)
	to.AddMoney(mon)
	fmt.Printf("%s -> %s: %d\n", from.ID(), to.ID(), mon)
}

func PrintAll(accounts []*bank.Account) {
	total := 0
	for _, acc := range accounts {
		acc.RLock()
		fmt.Printf("%s: %d\n", acc.ID(), acc.GetMoney())
		total += acc.GetMoney()
		acc.RUnlock()
	}
	fmt.Printf("Общий баланс: %d\n", total)
}
