package bank

import (
	"sync"
)

type Account struct {
	id    string
	money int
	mu    sync.RWMutex
	//*Manager
}

func NewAccount(id string, money int, m *Manager) *Account {
	return &Account{id: id, money: money}
}

func (a *Account) ID() string { return a.id }

func (a *Account) GetMoney() int {
	defer a.mu.RUnlock()
	return a.money
}

func (a *Account) AddMoney(amount int) {
	defer a.mu.Unlock()
	a.money += amount
}

func (a *Account) SubMoney(amount int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.money < amount {
		return false
	}
	a.money -= amount
	return true
}

func (a *Account) Lock()    { a.mu.Lock() }
func (a *Account) Unlock()  { a.mu.Unlock() }
func (a *Account) RLock()   { a.mu.RLock() }
func (a *Account) RUnlock() { a.mu.RUnlock() }
