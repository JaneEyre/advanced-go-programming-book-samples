package main

import "fmt"

func defernewvar() {
    for i := 0; i < 3; i++ {
        i := i // 定义一个循环体内局部变量 i
        defer func(){ println(i) } ()
    }
    fmt.Println("--- End of New var ---")
}

func deferinstcalli() {
    for i := 0; i < 3; i++ {
        // 通过函数传入 i
        // defer 语句会马上对调用参数求值
	// !!! This entire construct func(i int){ println(i) } is a function literal. It's a Anonymous function.
	// the " (i) " follows "defer func(i int){ println(i) }" , is calling this function
        defer func(i int){ println(i) } (i)
    }
    fmt.Println("--- End of instcall ---")
}

func main() {
	defernewvar()
	deferinstcalli()
}
