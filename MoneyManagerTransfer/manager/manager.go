package manager

import (
	"fmt"
	"mnogo/bank"
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

func (m *Manager) GetAccount(id string) *bank.Account {
	return m.find[id]
}

func (m *Manager) PrintAll() {
	total := 0
	for _, acc := range m.accounts {
		fmt.Printf("%s: %d\n", acc.ID(), acc.GetMoney())
		total += acc.GetMoney()
	}
	fmt.Printf("Общий баланс: %d\n", total)
}