package main

import (
	"fmt"
	"math/rand"
	"mnogo/bank"
	"mnogo/manager"
	"sync"
	"time"
)

func main() {
	manager := manager.NewManager([]*bank.Account{
		bank.NewAccount("A", 1000),
		bank.NewAccount("B", 500),
		bank.NewAccount("C", 300),
	})

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

		amount := rand.Intn(500) + 1

		from := manager.GetAccount(fromName)
		to := manager.GetAccount(toName)

		wg.Add(1)
		go bank.Transfer(from, to, amount, &wg)

		time.Sleep(5 * time.Millisecond)
	}

	wg.Wait()
	close(stopPrint)

	fmt.Println("\n итог")
	manager.PrintAll()
}
