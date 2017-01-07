// SearchGo project main.go
package main

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	//	"os"
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
	go readData(ch)
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
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("You can not go to url")
	} else {
		s, _ := ioutil.ReadAll(response.Body)
		fmt.Print("Count for " + string(url) + ": ")
		fmt.Println(strings.Count(string(s), "Go"))
		*count += strings.Count(string(s), "Go")
	}
	*goryt--
	wg.Done()
}

func readData(c chan string) {
	AllUrl := ""
	fmt.Scanln(&AllUrl)
	s := ""
	// Считываем посимвольно строку пока не встретим "/n" и отправляем в канал
	for i := 0; i < len(AllUrl); i++ {
		if AllUrl[i] != 92 {
			s += string(AllUrl[i])
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
