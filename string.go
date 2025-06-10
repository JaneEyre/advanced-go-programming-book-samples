package main

import (
	"fmt"
	//"reflect"
	//"unsafe"
)

/*
type StringHeader struct {
    Data uintptr
    Len  int
}

*/

// demonstrateStringSlicing shows how string slicing works in Go
// and how to inspect the underlying string header for length.
func demonstrateStringSlicing() {
	s := "hello, world"
	hello := s[:5] // Creates a new string "hello"
	world := s[7:] // Creates a new string "world"

	s1 := "hello, world"[:5] // String literal sliced directly, "hello"
	s2 := "hello, world"[7:] // String literal sliced directly, "world"

	fmt.Println("Original string 's':", s)
	fmt.Println("Slice 'hello':", hello)
	fmt.Println("Slice 'world':", world)
	fmt.Println("Literal slice 's1':", s1)
	fmt.Println("Literal slice 's2':", s2)

	/*
		fmt.Println("\n--- Inspecting StringHeader Lengths ---")
		// Using unsafe.Pointer and reflect.StringHeader to access internal string structure
		// Note: This is generally for understanding how Go strings work internally,
		// and is 'unsafe' as it bypasses Go's type safety.
		fmt.Println("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)     // 12
		fmt.Println("len(hello):", (*reflect.StringHeader)(unsafe.Pointer(&hello)).Len) // 5
		fmt.Println("len(world):", (*reflect.StringHeader)(unsafe.Pointer(&world)).Len) // 5
		fmt.Println("len(s1):", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len)   // 5
		fmt.Println("len(s2):", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)   // 5

		// Note: For actual string length, always use the built-in len() function:
		fmt.Println("\n--- Using built-in len() function ---")
		fmt.Println("len(s):", len(s))     // 12
		fmt.Println("len(hello):", len(hello)) // 5
		fmt.Println("len(world):", len(world)) // 5
		fmt.Println("len(s1):", len(s1))   // 5
		fmt.Println("len(s2):", len(s2))   // 5
	*/
	for i, c := range []byte("世界abc") {
		fmt.Println(i, c)
	}

	fmt.Printf("%#v\n", []rune("世界"))             // []int32{19990, 30028}
	fmt.Printf("%#v\n", string([]rune{'世', '界'})) // 世界

}

func main() {
	demonstrateStringSlicing()
}
