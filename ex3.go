package main

import (
	"container/list"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
	"math"
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

// Ex3.1
const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)
var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
func svgGreyLayout() {
	f, err := os.Create("file_31.svg")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, _ = fmt.Fprintln(f, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='", width, "' height='", height, "' >")
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)

			// all the conditions we need to ignore - the params must be init with some value
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			_, _ = fmt.Fprintln(f, "<polygon points='", ax, ",", ay, ",", bx, ",", by, ",", cx, ",", cy, ",", dx, ",", dy, "'/>")
		}
	}
	_, _ = fmt.Fprintln(f, "</svg>")

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}
func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

// Ex3.2 - prints a SVG rendering of a saddle.
type zFunc func(x, y float64) float64
func saddle(x, y float64) float64 {
	a := 25.0
	b := 17.0
	a2 := a * a
	b2 := b * b
	return (y*y/a2 - x*x/b2)
}
func cornerSaddle(i, j int , f zFunc) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}
func svg32(f zFunc) {
	w, err := os.Create("file_32.svg")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, _ = fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := cornerSaddle(i+1, j, f)
			bx, by := cornerSaddle(i, j, f)
			cx, cy := cornerSaddle(i, j+1, f)
			dx, dy := cornerSaddle(i+1, j+1, f)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			_, _ = fmt.Fprintf(w, "<polygon style='stroke: %s; fill: #222222' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				"#666666", ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	_, _ = fmt.Fprintln(w, "</svg>")

	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Ex3.3 - prints a SVG in color
func svgColor() {
	w, err := os.Create("file_33.svg")
	if err != nil {
		fmt.Println(err)
		return
	}

	zmin, zmax := minmax()
	_, _ = fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			_, _ = fmt.Fprintf(w, "<polygon style='stroke: %s; fill: #222222' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				color(i, j, zmin, zmax), ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	_, _ = fmt.Fprintln(w, "</svg>")
	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
// minmax returns the min and max values for z given the min/max values of x and y and assuming a square domain.
func minmax() (min float64, max float64) {
	min = math.NaN()
	max = math.NaN()
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			for xoff := 0; xoff <= 1; xoff++ {
				for yoff := 0; yoff <= 1; yoff++ {
					x := xyrange * (float64(i+xoff)/cells - 0.5)
					y := xyrange * (float64(j+yoff)/cells - 0.5)
					z := f(x, y)
					if math.IsNaN(min) || z < min {
						min = z
					}
					if math.IsNaN(max) || z > max {
						max = z
					}
				}
			}
		}
	}
	return
}
func color(i, j int, zmin, zmax float64) string {
	min := math.NaN()
	max := math.NaN()
	for xoff := 0; xoff <= 1; xoff++ {
		for yoff := 0; yoff <= 1; yoff++ {
			x := xyrange * (float64(i+xoff)/cells - 0.5)
			y := xyrange * (float64(j+yoff)/cells - 0.5)
			z := f(x, y)
			if math.IsNaN(min) || z < min {
				min = z
			}
			if math.IsNaN(max) || z > max {
				max = z
			}
		}
	}

	color := ""
	if math.Abs(max) > math.Abs(min) {
		red := math.Exp(math.Abs(max)) / math.Exp(math.Abs(zmax)) * 255
		if red > 255 {
			red = 255
		}
		color = fmt.Sprintf("#%02x0000", int(red))
	} else {
		blue := math.Exp(math.Abs(min)) / math.Exp(math.Abs(zmin)) * 255
		if blue > 255 {
			blue = 255
		}
		color = fmt.Sprintf("#0000%02x", int(blue))
	}
	return color
}

// Ex3.4
func svg(w io.Writer) {
	zmin, zmax := minmax()
	_, _ = fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			_, _ = fmt.Fprintf(w, "<polygon style='stroke: %s; fill: #222222' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				color(i, j, zmin, zmax), ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	_, _ = fmt.Fprintln(w, "</svg>")
}


func main() {

	fmt.Println("Ex4.3")
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
	fmt.Println("the function is: ", isValid(topoSort_result)) //Valid topolgical order

	fmt.Println("Ex5.11")
	printCycleTopologicalSort(cyclePrereqs) /*
		cycle: a -> b -> a
		cycle: data structures -> discrete math -> intro to programming -> data structures
		cycle: calculus -> linear algebra -> calculus
		1:	b
		2:	a
		3:	intro to programming
		4:	discrete math
		5:	data structures
		6:	algorithms
		7:	linear algebra
		8:	calculus
		9:	formal languages
		10:	computer organization
		11:	compilers
		12:	databases
		13:	operating systems
		14:	networks
		15:	programming languages
	*/

	fmt.Println("Ex5.12")
	callOutline([]string{"http://gopl.io"})
	/*
		<html>
		  <head>
		    <meta>
		    </meta>
		    <title>
		    </title>
		    <script>
		    </script>
		    <link>
		    </link>
		    <style>
		    </style>
		  </head>
		  <body>
		    <table>
		      <tbody>
		        <tr>
		          <td>
		            <a>
		              <img>
		              </img>
		            </a>
		            <br>
		            </br>
		            <div>
		              <a>
		                <img>
		                </img>
		              </a>
		              <a>
		                <img>
		                </img>
		              </a>
		              <a>
		                <img>
		                </img>
		              </a>
		            </div>
		            <br>
		            </br>
		          </td>
		          <td>
		            <h1>
		            </h1>
		            <p>
		              <br>
		              </br>
		              <br>
		              </br>
		              <br>
		              </br>
		              <tt>
		              </tt>
		              <tt>
		              </tt>
		              <tt>
		              </tt>
		            </p>
		            <div>
		              <table>
		                <tbody>
		                  <tr>
		                    <td>
		                      <h1>
		                        <a>
		                        </a>
		                      </h1>
		                      <h1>
		                        <a>
		                        </a>
		                      </h1>
		                      <h1>
		                        <a>
		                        </a>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                    </td>
		                    <td>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                      </h1>
		                      <h1>
		                        <a>
		                        </a>
		                      </h1>
		                    </td>
		                  </tr>
		                  <tr>
		                    <td>
		                      <h1>
		                        <a>
		                        </a>
		                        <a>
		                        </a>
		                        <a>
		                        </a>
		                        <a>
		                        </a>
		                      </h1>
		                    </td>
		                  </tr>
		                </tbody>
		              </table>
		            </div>
		            <p>
		              <a>
		                <code>
		                </code>
		              </a>
		              <a>
		                <code>
		                </code>
		              </a>
		              <a>
		                <code>
		                </code>
		              </a>
		              <a>
		                <code>
		                </code>
		              </a>
		            </p>
		            <p>
		              <a>
		              </a>
		              <a>
		              </a>
		            </p>
		          </td>
		        </tr>
		      </tbody>
		    </table>
		  </body>
		</html>
	*/

	fmt.Println("Ex5.13")
	callCrawler([]string{"https://golang.org"})
	/*
		Some of the output:

		https://golang.org
		https://support.eji.org/give/153413/#!/donation/checkout
		https://golang.org/
		https://golang.org/doc/
		https://golang.org/pkg/
		https://golang.org/project/
		https://golang.org/help/
		https://golang.org/blog/
		https://play.golang.org/
		https://golang.org/dl/
		https://tour.golang.org/
		https://blog.golang.org/
		https://golang.org/doc/copyright.html
		https://golang.org/doc/tos.html
		http://www.google.com/intl/en/policies/privacy/
		http://golang.org/issues/new?title=x/website:
		https://google.com
		https://golang.org/doc/install
		https://golang.org/doc/code.html
		https://golang.org/cmd/go/
		https://golang.org/doc/editors.html
		https://golang.org/doc/effective_go.html
		https://golang.org/doc/diagnostics.html
		https://golang.org/doc/faq
		https://golang.org/wiki
		golang.org
		├── blog
		├── cmd
		│   └── go
		├── dl
		├── doc
		│   ├── code.html
		│   ├── copyright.html
		│   ├── diagnostics.html
		│   ├── editors.html
		│   ├── effective_go.html
		│   ├── faq
		│   ├── install
		│   └── tos.html
		├── help
		├── issues
		├── pkg
		├── project
		└── wiki
	*/

	// Fill in path to a directory on your computer below scanFiles([]string{path})
	fmt.Println("Ex5.14")
	/*
		input: "C:\\Users\\yonis\\Desktop\\ComputerScience\\Test"
		output:
		Test
		a
		b
		c
		d
		e
		f
		test.txt
		g
		anotherTest.txt
	*/

	fmt.Println("Ex3.1")
	svgGreyLayout() // result in file_31.svg

	fmt.Println("Ex3.2")
	var f zFunc
	f = saddle
	svg32(f) // result in file_32.svg

	fmt.Println("Ex3.3")
	svgColor() // result in file_33.svg

	fmt.Println("Ex3.4")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		svg(w) // result in file_34.svg
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
