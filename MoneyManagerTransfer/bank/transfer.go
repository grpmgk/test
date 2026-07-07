package bank

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Transfer(from, to *Account, mon int, wg *sync.WaitGroup) {
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
