package main

import (
	"container/heap"
	"fmt"
)



type Elem struct {
	Name string
	votes int
}

type Elems []Elem
func (e Elems) Len() int {
	return len(e)
}

func (e Elems) Less(i, j int) bool {
	return e[i].votes >= e[j].votes
}

func (e Elems) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *Elems) Push(x any) {
	element := x.(Elem)
	*e = append(*e, element)

}

func (e *Elems) Pop() any {
	tmp := *e
	l := len(tmp)
	deletedEle := tmp[l-1]
	*e = tmp[:l-1]
	return deletedEle
}


func main(){

	input := [][]string{
		{
			"dan","alice", "bob","charlie",
		},
		{
			"dan","alice", "charlie", "bob",
		},
		{
			"charlie","dan","alice", "bob",
		},
	}
	fmt.Println(FindWinners(input))

}


func FindWinners(votes [][]string) []string {

	countMap := make(map[string]int)
	for _, vote := range votes {
		for i:=0; i< len(vote); i++ {
			candidate := vote[i]
			voteCount :=0
			if 3> i {
				voteCount = 3-i
			}
			if count, ok :=countMap[candidate]; !ok {
				countMap[candidate] = voteCount
			}else{
				countMap[candidate] = count+ voteCount
			}
		}
	}

	elements := make(Elems, len(countMap))
	i := 0
	for key, value := range countMap {
		elements[i].votes = value
		elements[i].Name = key
		i++
	}
	heap.Init(&elements)
	count := len(countMap)
	var res []string
	for i := 0; i < count; i++ {
		temp := heap.Pop(&elements).(Elem)
		res = append(res, temp.Name)
	}
	return res
}


