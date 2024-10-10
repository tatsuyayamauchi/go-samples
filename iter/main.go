package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	fmt.Println("run SliceSample")
	SliceSample()

	fmt.Println("")

	fmt.Println("run MapSample")
	MapSample()
}

// Result
/*
run SliceSample
[a b c]
a
b
c
[a b c]
0 a
1 b
2 c
2 c
1 b
0 a
[a b c a b c]
[a b]
[c a]
[b c]

run MapSample
map[a:1 b:2 c:3]
c 3
a 1
b 2
map[a:1 b:2 c:3]
a
b
c
1
2
3
map[a:1 b:2 c:3 d:4 e:5]
*/

func SliceSample() {
	hoge := []string{"a", "b", "c"}
	fmt.Println(hoge)

	foo := slices.Values(hoge)
	for v := range foo {
		fmt.Println(v)
	}
	hoge2 := slices.Collect(foo)
	fmt.Println(hoge2)

	foo2 := slices.All(hoge)
	for k, v := range foo2 {
		fmt.Println(k, v)
	}

	foo3 := slices.Backward(hoge)
	for k, v := range foo3 {
		fmt.Println(k, v)
	}

	fmt.Println(slices.AppendSeq(hoge, slices.Values(hoge)))

	for v := range slices.Chunk(slices.AppendSeq(hoge, slices.Values(hoge)), 2) {
		fmt.Println(v)
	}
}

func MapSample() {
	hoge := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	fmt.Println(hoge)

	foo := maps.All(hoge)
	for k, v := range foo {
		fmt.Println(k, v)
	}

	fmt.Println(maps.Collect(foo))

	for v := range maps.Keys(hoge) {
		fmt.Println(v)
	}

	for v := range maps.Values(hoge) {
		fmt.Println(v)
	}

	newVal := map[string]int{
		"d": 4,
		"e": 5,
	}

	maps.Insert(hoge, func(yield func(string, int) bool) {
		for k, v := range newVal {
			if !yield(k, v) {
				return
			}
		}
	})
	fmt.Println(hoge)
}
