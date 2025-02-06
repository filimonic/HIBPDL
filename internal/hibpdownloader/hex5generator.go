package hibpdownloader

import (
	"fmt"
)

func hex5generator() (uint64, chan string) {
	const total = 0x100_000 // 0xFFFFF + 1
	ch := make(chan string)

	go func(ch chan string) {
		for i := range total {
			ch <- fmt.Sprintf("%05X", i)
		}
		close(ch)
	}(ch)

	return total, ch
}
