package main

import (
	"fmt"
	"strings"
)

func main() {
	list := []string{"aaaa", "bbbb", "cccc"}
	fmt.Println(strings.Join(list, `","`))
}
