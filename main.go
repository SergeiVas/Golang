// SearchGo project main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {

	var (
		count     = 0
		gorytines = 0
		k         = 5
		wg        sync.WaitGroup
	)
	ch := make(chan string) // Канал, в котором хранятся url
	go readFile(ch)
	for num := range ch {
		if gorytines < k {
			wg.Add(1)
			go countStrings(num, &count, &gorytines, &wg)
		} else { // Если кол - во горутин больше заданного ждем пока они выполнятся и запускаем снова
			wg.Wait()
			wg.Add(1)
			go countStrings(num, &count, &gorytines, &wg)
		}
	}
	wg.Wait() // Ждем выполнения всех горутин

	fmt.Print("Total: ")
	fmt.Println(count)
}

//Переходим по url'у и считаем кол - во искомых строк
func countStrings(url string, count *int, goryt *int, wg *sync.WaitGroup) {
	*goryt++
	response, _ := http.Get(url)
	defer response.Body.Close()
	s, _ := ioutil.ReadAll(response.Body)
	fmt.Print("Count for " + string(url) + ": ")
	fmt.Println(strings.Count(string(s), "Go"))
	*count += strings.Count(string(s), "Go")
	*goryt--
	wg.Done()
}

func initializationUrl(line string, num int, c chan string) {
	s := ""
	// Считываем посимвольно строку пока не встретим "/n" и отправляем в канал
	for i := num; i < len(line); i++ {
		if line[i] != 92 {
			s += string(line[i])
		} else {
			c <- string(s)
			s = ""
			i++
		}
	}
	// В конце "/n" не встречается, поэтому просто отправляем url в канал
	if len(s) > 5 {
		c <- string(s)
		close(c)
	}
}

func readFile(c chan string) {
	file, _ := os.Open("input.txt")
	f := bufio.NewReader(file)
	check := false
	// Считываем строки с файла и отправляем их на разбор url'ов
	for {
		line, err := f.ReadString('\n')
		line = strings.Replace(line, "\r\n'", "", -1)
		line = strings.Replace(line, "'", "", -1)
		if check == false {
			initializationUrl(line, 10, c) // если первая строка, то сразу игнорируем "$ echo -e "
			check = true
		} else {
			initializationUrl(line, 0, c)

		}
		if err != nil {
			break
		}

	}
}
