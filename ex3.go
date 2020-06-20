package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Ex4.3 : Reverse an array using pointer and without slice.
func reverse(s *[60]int) {
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
func squashSpace(bytes []byte) []byte {
	out := bytes[:0]
	var last rune

	for i := 0; i < len(bytes); {
		r, rune_size := utf8.DecodeRune(bytes[i:]) // rune and its size

		// check if the rune is a space character in Unicode
		if !unicode.IsSpace(r) {
			out = append(out, bytes[i:i+rune_size]...) // adding the bytes we want which not containing space
		} else if unicode.IsSpace(r) && !unicode.IsSpace(last) { // if found space but there are non space in the end - add the space
			out = append(out, ' ')
		}
		last = r       // the remaining rune
		i += rune_size // go to the next rune
	}
	return out
}

// Ex4.7 : Reverse the characters of a []byte slice that represents a UTF-8-encoded string, in place.
func rev(in []byte) {
	s := len(in)
	for i := 0; i < len(in)/2; i++ {
		in[i], in[s-1-i] = in[s-1-i], in[i]
	}
}

func ReverseRune(in []byte) []byte {
	for i := 0; i < len(in); {
		_, s := utf8.DecodeRune(in[i:]) // decoding
		rev(in[i : i+s])
		i += s
	}
	rev(in)
	return in

	/*
		buf := make([]byte, 0, len(in)) // a slice of length 0 and capacity len(in) that is backed by this underlying array.
		i := len(in)

		for i > 0 {
			_, s := utf8.DecodeLastRune(in[:i]) // decoding
			buf = append(buf, in[i-s:i]...)     // appending the needed word
			i -= s
		}
		copy(in, buf)

	*/
}

func main() {

	fmt.Println("Ex4.4")
	s := []int{1, 2, 3, 4, 5}
	fmt.Println("the list before is : ", s)
	rounds := 3
	rotate(s, rounds)
	fmt.Println("the list after", rounds, "rotaions:", s)

	fmt.Println("Ex4.5")
	intSlice := []string{"1", "5", "5", "1", "1","1",  "3", "6", "9", "9", "4", "2", "3", "1", "5"}
	fmt.Println("ths slice before : ", intSlice)
	uniqueSlice := unique(intSlice)
	fmt.Println("ths slice after : ", uniqueSlice)

	fmt.Println("Ex4.6")
	squashSpace_result := string(squashSpace([]byte("R I c \n k  A n D   M o   R t I \n \n \n y")))
	squashSpace_wanted := "R I c k A n D M o R t I y"
	fmt.Println("the function is :", squashSpace_result == squashSpace_wanted)

	fmt.Println("Ex4.7")
	ReverseRune_result := string(ReverseRune([]byte("Räksmörgås")))
	ReverseRune_wanted := "sågrömskäR"
	fmt.Println("the function is :", ReverseRune_result == ReverseRune_wanted)
}
