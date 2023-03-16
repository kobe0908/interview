package in_turn_print

import (
	"fmt"
	"sync"
)

func InturnPrint() {
	letter := make(chan bool)
	num := make(chan bool)

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for {
			select {
			case <-letter:
				fmt.Println("a")
				num <- true
			default:
				break
			}
		}
	}()

	wait.Add(1)
	go func() {
		for {
			select {
			case <-num:
				fmt.Println("1")
				letter <- true
			default:
				break
			}
		}
	}()
	num <- true
	wait.Wait()
}
