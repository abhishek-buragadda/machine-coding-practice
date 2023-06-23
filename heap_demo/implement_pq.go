package main

import (
	"container/heap"
	"fmt"
)

/*

Implement priority queue,
 - datatypes to support -> generic
 - min/max queue? - Min
 - order based on ? - Status which is enum
 - methods push(), pop(),
 */



type CustomElement struct {
	Value    interface{}
	Priority int64
}


type PriorityQueue []CustomElement

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *PriorityQueue) Push(x interface{}) {
	*p = append(*p, x.(CustomElement))
}

func (p *PriorityQueue) Pop() interface{} {
	old := *p
	n := len(old)
	ele := old[n-1]
	*p = old[0:n-1]
	return ele
}

func main() {
	pq := &PriorityQueue{}

	nums := [][]int64{{1,2},{2,3},{3,5},{4,1}}

	for _, value:= range nums {
		heap.Push(pq, CustomElement{
			Value:    value[0],
			Priority: value[1],
		})
	}

	for i:=0; i< len(nums); i++{
		fmt.Println( heap.Pop(pq))

	}

}
