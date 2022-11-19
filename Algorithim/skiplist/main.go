package main

import "fmt"

const maxLevel = 32
const pFactor = 0.25

type SkipListNode struct {
	val     int
	forward []*SkipListNode
}

func main() {
	fmt.Println("hello")
}
