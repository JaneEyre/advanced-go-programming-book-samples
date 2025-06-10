

/*
   reflect.SliceHeader

   	type SliceHeader struct {
   	    Data uintptr
   	    Len  int
   	    Cap  int
   	}
*/
var (
	a []int               // nil 切片, 和 nil 相等, 一般用来表示一个不存在的切片
	b = []int{}           // 空切片, 和 nil 不相等, 一般用来表示一个空的集合
	c = []int{1, 2, 3}    // 有 3 个元素的切片, len 和 cap 都为 3
	d = c[:2]             // 有 2 个元素的切片, len 为 2, cap 为 3
	e = c[0:2:cap(c)]     // 有 2 个元素的切片, len 为 2, cap 为 3
	f = c[:0]             // 有 0 个元素的切片, len 为 0, cap 为 3
	g = make([]int, 3)    // 有 3 个元素的切片, len 和 cap 都为 3
	h = make([]int, 2, 3) // 有 2 个元素的切片, len 为 2, cap 为 3
	i = make([]int, 0, 3) // 有 0 个元素的切片, len 为 0, cap 为 3
)

for i := range a {
	fmt.Printf("a[%d]: %d\n", i, a[i])
}
for i, v := range b {
	fmt.Printf("b[%d]: %d\n", i, v)
}
for i := 0; i < len(c); i++ {
	fmt.Printf("c[%d]: %d\n", i, c[i])
}
