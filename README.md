## In-Place Slice Techniques

### Exercise 4.3
Rewrite `reverse` to use an array pointer instead of a slice.

### Exercise 4.4
Write a version of `rotate` that operates in a single pass.

### Exercise 4.5
Write an in-place function to eliminate adjacent duplicates in a `[]string`
slice.

### Exercise 4.6
Write an in-place function that squashes each run of adjacent Unicode spaces
(see `unicode.IsSpace`) in a UTF-8-encoded `[]byte` slice into a single ASCII
space.

### Exercise 4.7
Modify `reverse` to reverse the characters of a `[]byte` slice that represents
a UTF-8-encoded string, in place. Can you do it without allocating new memory?
