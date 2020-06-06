package main

import "fmt"

// Ex4.3 : Reverse an array using pointer and without slice.
func reverse(s *[100]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Ex4.4 : Rotate function in a single pass
// rotation - the number of rotations we want for the slice
func rotate(s []int, rotation int) {
	for i := 0; i < rotation; i = i + 1 {
		first := s[0]
		copy(s, s[1:])
		s[len(s)-1] = first
	}
}

// Ex4.5 : In-place function to eliminate adjacent duplicates in a []string slice
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Ex4.6 : In-place function hat squashes each run of adjacent Unicode spaces in a UTF-8-encoded []byte slice into a single ASCII space.

// Ex4.7 : Reverse the characters of a []byte slice that represents a UTF-8-encoded string, in place.

func main() {

	fmt.Println("Ex4.4")
	s := []int{1, 2, 3, 4, 5}
	fmt.Println("the list before is : ", s)
	rounds := 3
	rotate(s, rounds)
	fmt.Println("the list after", rounds, "rotaions:", s)

	fmt.Println("Ex4.5")
	intSlice := []string{"1", "5", "3", "6", "9", "9", "4", "2", "3", "1", "5"}
	fmt.Println("ths slice before : ", intSlice)
	uniqueSlice := unique(intSlice)
	fmt.Println("ths slice after : ", uniqueSlice)
}
