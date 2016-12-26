// SearchGo project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Переходим по url'у и считаем кол - во искомых строк
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

// Разделяем url'ы и записываем их в канал
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

// Считываем url'ы
func readFile() string {
	bs, _ := ioutil.ReadFile("input.txt")
	str := string(bs)
	formatStr := ""
	//Начинаем читать с '  и заканчиваем ' чтобы получить только строки url'ов
	for i := strings.Index(str, "'") + 1; i < len(str)-1; i++ {
		formatStr += string(str[i])
	}
	return formatStr // Возвращаем строку из url'ов
}

func main() {
	var (
		search    = "Go"
		count     = 0
		countUrl  = 1 // Кол - во url'ов, которые необходимо будет обработать
		gorytines = 0
	)
	ch := make(chan string) // Канал, в котором хранятся url
	s := readFile()         // Записываем в строку необработанные url'ы

	go stringsUrl(string(s), ch, &countUrl) // Обрабатываем url'ы и записываем их количество

	for countUrl > 0 {
		for i := gorytines; i < 5; i++ {
			go countStrings(search, ch, &count, &countUrl, &gorytines)
		}
	}

	fmt.Print("Total: ")
	fmt.Println(count)

}
