// Session 14 — the sort package, including sorting structs and map keys.
// Run:  go run examples/session14/sortpkg/sortpkg.go
package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// Sort built-in slices in place.
	nums := []int{5, 2, 8, 1, 9}
	sort.Ints(nums)
	fmt.Println("sorted ints:", nums)

	words := []string{"banana", "apple", "cherry"}
	sort.Strings(words)
	fmt.Println("sorted strings:", words)

	// Sort a slice of structs with a custom comparison (sort.Slice).
	people := []Person{
		{"Sara", 30},
		{"Ali", 25},
		{"Sobhan", 28},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age // ascending by age
	})
	fmt.Println("by age:", people)

	// The idiom for stable, ordered output from a (randomly-ordered) map:
	scores := map[string]int{"math": 90, "art": 75, "bio": 82}
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s=%d\n", k, scores[k])
	}
}
