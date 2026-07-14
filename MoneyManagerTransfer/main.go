package main

import (
	"fmt"
	"math/rand"
	"mnogo/bank"
	"sync"
	"time"
)

func main() {
	// manager := bank.NewManager()
	// accA := manager.CreateAccount("A", 1000)
	// accB := manager.CreateAccount("B", 500)
	// accC := manager.CreateAccount("C", 300)

	manager := bank.NewManager("A", 1000)
	accB := manager.CreateAccount("B", 500)
	accC := manager.CreateAccount("C", 300)

	fmt.Println("начало")
	manager.PrintAll()

	stopPrint := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Println("\n середина")
				manager.PrintAll()
			case <-stopPrint:
				return
			}
		}
	}()

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	names := []string{"A", "B", "C"}

	for i := 0; i < 1000; i++ {
		fromName := names[rand.Intn(3)]
		toName := names[rand.Intn(3)]
		for fromName == toName {
			toName = names[rand.Intn(3)]
		}

		var from, to *bank.Account
		switch fromName {
		case "A":
			from = manager.Account
		case "B":
			from = accB
		case "C":
			from = accC
		}
		switch toName {
		case "A":
			to = manager.Account
		case "B":
			to = accB
		case "C":
			to = accC
		}

		amount := rand.Intn(500) + 1

		wg.Add(1)
		go bank.Transfer(from, to, amount, &wg)

		time.Sleep(5 * time.Millisecond)
	}
	wg.Wait()
	close(stopPrint)

	fmt.Println("\n итог")
	manager.PrintAll()
}
