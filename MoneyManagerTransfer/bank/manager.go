package bank

import (
	"fmt"
	"sync"
)

type Manager struct {
	*Account
	accounts []*Account
	find     map[string]*Account
}

func NewManager(id string, money int) *Manager {
	mainAcc := &Account{
		id:    id,
		money: money,
		mu:    sync.RWMutex{},
	}
	m := &Manager{
		Account:  mainAcc,
		accounts: []*Account{mainAcc},
		find:     map[string]*Account{id: mainAcc},
	}
	return m
}

func (m *Manager) CreateAccount(id string, money int) *Account {
	acc := NewAccount(id, money, m) // передаём себя
	m.accounts = append(m.accounts, acc)
	m.find[id] = acc
	return acc
}

// func (m *Manager) GetAccount(id string) *Account {
// 	return m.find[id]
// }

func (m *Manager) PrintAll() {
	total := 0
	for _, acc := range m.accounts {
		fmt.Printf("%s: %d\n", acc.ID(), acc.GetMoney())
		total += acc.GetMoney()
	}
	fmt.Printf("Общий баланс: %d\n", total)
}
