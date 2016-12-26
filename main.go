// SearchGo project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func countStrings(search string, c chan string, count *int, countUrl *int, goryt *int) int {
	*goryt++
	url := <-c
	response, _ := http.Get(url)
	defer response.Body.Close()
	s, _ := ioutil.ReadAll(response.Body)
	fmt.Print("Count for " + string(url) + ": ")
	fmt.Println(strings.Count(string(s), string(search)))
	*count += strings.Count(string(s), string(search))
	*countUrl--
	*goryt--
	return strings.Count(string(s), string(search))
}

func stringsUrl(allUrl string, c chan string, countUrl *int) {

	s := ""
	for i := 0; i < len(allUrl); i++ {
		if allUrl[i] != 92 {
			s += string(allUrl[i])
		} else {
			c <- string(s)
			s = ""
			i++
			*countUrl++
		}
	}
	c <- string(s)
}

func readFile() string {
	bs, _ := ioutil.ReadFile("input.txt")
	str := string(bs)

	formatStr := ""
	for i := strings.Index(str, "'") + 1; i < len(str)-1; i++ {
		formatStr += string(str[i])
	}
	return formatStr
}

func main() {
	var (
		search    = "Go"
		count     = 0
		countUrl  = 1
		gorytines = 0
	)
	ch := make(chan string)
	s := readFile()

	go stringsUrl(string(s), ch, &countUrl)

	for countUrl > 0 {
		for i := gorytines; i < 5; i++ {
			go countStrings(search, ch, &count, &countUrl, &gorytines)
		}
	}

	fmt.Print("Total: ")
	fmt.Println(count)

}
