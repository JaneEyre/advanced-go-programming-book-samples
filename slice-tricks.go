/*
切片内存技巧
在本节开头的数组部分我们提到过有类似 [0]int 的空数组，空数组一般很少用到。但是对于切片来说，len 为 0 但是 cap 容量不为 0 的切片则是非常有用的特性。当然，如果 len 和 cap 都为 0 的话，则变成一个真正的空切片，虽然它并不是一个 nil 值的切片。在判断一个切片是否为空时，一般通过 len 获取切片的长度来判断，一般很少将切片和 nil 值做直接的比较。
比如下面的 TrimSpace 函数用于删除 []byte 中的空格。函数实现利用了 0 长切片的特性，实现高效而且简洁。
*/
func TrimSpace(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

// 其实类似的根据过滤条件原地删除切片元素的算法都可以采用类似的方式处理（因为是删除操作不会出现内存不足的情形）：
func Filter(s []byte, fn func(x byte) bool) []byte {
	b := s[:0]
	for _, x := range s {
		if !fn(x) {
			b = append(b, x)
		}
	}
	return b
}


/* 切片高效操作的要点是要降低内存分配的次数，尽量保证 append 操作不会超出 cap 的容量，降低触发内存分配的次数和每次分配内存大小。