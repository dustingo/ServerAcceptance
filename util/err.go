package util

import "fmt"

func Eprint(e error) bool {
	if e != nil {
		fmt.Println(e)
		return true
	}
	return false
}
