package simgo

import (
	"fmt"
)

func recoverer(maxPanics, id int, f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HERE", id)
			fmt.Println(err)
			if maxPanics == 0 {
				panic("SimGo exceeded max tries. Exiting...")
			} else {
				go recoverer(maxPanics-1, id, f)
			}
		}
	}()
	f()
}
