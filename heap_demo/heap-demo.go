package main


type MyHeap []int64


func (m MyHeap) Len() int {
	return len(m)
}

func (m MyHeap) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m MyHeap) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]

}

func (m *MyHeap) Push(x any) {
	*m = append(*m, x.(int64))
}

func (m *MyHeap) Pop() any {
	old := *m
	n := len(old)
	t:= old[n-1]
	*m = old[0:n-1]
	return t
}

//func main() {
//	h := &MyHeap{}
//	nums := []int64{2,6,4,5,3}
//	for _,val := range nums {
//		heap.Push(h, val)
//	}
//	for  i:= 0; i< len(nums); i++ {
//		fmt.Println(heap.Pop(h).(int64))
//	}
//}

