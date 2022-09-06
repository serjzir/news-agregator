package repeatable

import (
	"fmt"
	"time"
)

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		fmt.Println(attempts)
		fmt.Println("Connect")
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return err
}
