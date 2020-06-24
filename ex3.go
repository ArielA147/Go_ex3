package main

import (
	"container/list"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Ex4.3 : Reverse an array using pointer and without slice.
func reverse(s *[5]int) {
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

var cyclePrereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"linear algebra":        {"calculus"},
	"intro to programming":  {"data structures"},
	"a":                     {"b"},
	"b":                     {"a"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

var regularPrereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func cycleTopoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	current := list.New()
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			current.PushFront(item)
			checkCycle(item, current)
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
			current.Init()
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

func checkCycle(item string, current *list.List) {
	cycle := "cycle: "
	if current.Len() < 2 {
		return
	}
	for e := current.Front().Next(); e != nil; e = e.Next() {
		if e.Value.(string) == item {
			for current.Back().Value.(string) != item {
				current.Remove(current.Back())
			}
			for a := current.Back(); a.Prev() != nil; a = a.Prev() {
				cycle += a.Value.(string) + " -> "
			}
			fmt.Println(cycle + current.Back().Value.(string))
			current.Init()
			return
		}
	}
}

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string]map[string]bool{
	"algorithms": {"data structures": true},
	"calculus":   {"linear algebra": true},
	"compilers": {
		"data structures":       true,
		"formal languages":      true,
		"computer organization": true,
	},
	"data structures":       {"discrete math": true},
	"databases":             {"data structures": true},
	"discrete math":         {"intro to programming": true},
	"formal languages":      {"discrete math": true},
	"networks":              {"operating systems": true},
	"operating systems":     {"data structures": true, "computer organization": true},
	"programming languages": {"data structures": true, "computer organization": true},
}

func topoSort(m map[string]map[string]bool) map[int]string {
	rank := 1
	mapOrder := make(map[int]string)
	seen := make(map[string]bool)
	var visitAll func(items map[string]bool)
	visitAll = func(items map[string]bool) {
		for item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				mapOrder[rank] = item
				rank++
			}
		}
	}
	for k := range m {
		visitAll(m[k])
		if !seen[k] {
			mapOrder[rank] = k
			seen[k] = true
			rank++
		}
	}
	return mapOrder
}

func printInOrder(rankToCourse map[int]string) {
	var len = len(rankToCourse)
	for i := 1; i < len+1; i++ {
		fmt.Printf("%d:\t%s\n", i, rankToCourse[i])
	}
}

func isValid(result map[int]string) interface{} {
	var getRank func(subject string) int
	getRank = func(subject string) int {
		for k, v := range result {
			if v == subject {
				return k
			}
		}
		return -1
	}
	for course, prerequesites := range prereqs {
		for subject := range prerequesites {
			if getRank(subject) > getRank(course) {
				return false
			}
		}
	}
	return true
}

func printCycleTopologicalSort(strings map[string][]string) {
	for i, course := range cycleTopoSort(cyclePrereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func outline2(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	var depth int

	startElement := func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
			depth++
		}
	}

	endElement := func(n *html.Node) {
		if n.Type == html.ElementNode {
			depth--
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}

	forEachNode(doc, startElement, endElement)
	return nil
}

// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func callOutline(args []string) {
	for _, url := range args {
		outline2(url)
	}
}

// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string, domain string) []string, worklist []string, domain string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item, domain)...)
			}
		}
	}
}

// Crawl the web breadth-first,
// starting from the command-line arguments.
func crawl(url string, domain string) []string {
	fmt.Println(url)
	if getDomain(url) == domain {
		saveLink(url)
	}
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func saveLink(url string) {
	path := strings.Split(url, "/")[2:]
	traverse(path)
}

func traverse(path []string) {
	exists, _ := exists(path[0])
	if !exists {
		os.Mkdir(path[0], os.ModeDir)
	}
	if len(path) > 1 {
		path[1] = path[0] + "/" + path[1]
		traverse(path[1:])
	}
	return
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getDomain(str string) string {
	domain := strings.Split(str, "/")[2]
	return domain
}

func callCrawler(args []string) {
	for _, arg := range args {
		breadthFirst(crawl, []string{arg}, getDomain(arg))
	}
}

func crawlFiles(fileName string, path string) []string {
	if fileName != "" {
		fmt.Println(getFileName(fileName))
	} else {
		fmt.Println(getFileName(path))
		fileName = path
	}

	fi, err := os.Stat(fileName)
	if err != nil {
		log.Print(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(fileName)
		if err != nil {
			log.Print(err)
		}
		filesAsString := []string{}
		for _, file := range files {
			filesAsString = append(filesAsString, filepath.Join(fileName, file.Name()))
		}
		return filesAsString
	}
	return []string{}
}

func scanFiles(paths []string) {
	for _, path := range paths {
		breadthFirst(crawlFiles, []string{""}, path)
	}
}

func getFileName(path string) string {
	parts := strings.Split(path, "\\")
	return parts[len(parts)-1]
}

func main() {

	fmt.Println("Ex4.4")
	s2 := [5]int{1, 2, 3, 4, 5}
	fmt.Println("the list before is: ", s2)
	reverse(&s2)
	fmt.Println("the list after is: ", s2) // the list after is :  [5 4 3 2 1]

	fmt.Println("Ex4.4")
	s := []int{1, 2, 3, 4, 5}
	fmt.Println("the list before is: ", s)
	rounds := 3
	rotate(s, rounds)
	fmt.Println("the list after", rounds, "rotaions:", s) // the list after 3 rotaions: [4 5 1 2 3]

	fmt.Println("Ex4.5")
	intSlice := []string{"1", "5", "5", "1", "1", "1", "3", "6", "9", "9", "4", "2", "6", "9", "6", "9", "6", "9", "3", "1", "5"}
	fmt.Println("ths slice before: ", intSlice)
	uniqueSlice := unique(intSlice)
	fmt.Println("ths slice after: ", uniqueSlice) // ths slice after :  [1 5 3 6 9 4 2]

	fmt.Println("Ex4.6")
	squashSpace_result := string(squashSpace([]byte("R I c \n k  A n D   M o   R t I \n \n \n y")))
	squashSpace_wanted := "R I c k A n D M o R t I y"
	fmt.Println("the function is: ", squashSpace_result == squashSpace_wanted) // True

	fmt.Println("Ex4.7")
	ReverseRune_result := string(ReverseRune([]byte("ArielAndYoni")))
	ReverseRune_wanted := "inoYdnAleirA"
	fmt.Println("the function is: ", ReverseRune_result == ReverseRune_wanted) // True

	fmt.Println("Ex5.10")
	topoSort_result := topoSort(prereqs)
	printInOrder(topoSort_result)
	fmt.Println("the function is: ", isValid(topoSort_result))

	fmt.Println("Ex5.11")
	printCycleTopologicalSort(cyclePrereqs)

	fmt.Println("Ex5.12")
	callOutline([]string{"http://gopl.io"})

	fmt.Println("Ex5.13")
	callCrawler([]string{"https://golang.org"})

	fmt.Println("Ex5.14")
	//Fill in path to a directory on your computer below
	//scanFiles([]string{path})

}
