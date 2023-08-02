package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type User struct {
	name     string
	age      int
	email    string
	bookList []string
}

func main() {
	fmt.Println("vim-go")

	userCHN := make(chan User)

	user := getUsers(userCHN)

	fmt.Println("USER DATA: ")
	fmt.Println("name", user.name)
	fmt.Println("age", user.age)
	fmt.Println("email", user.email)
	fmt.Println("books", user.bookList)

}

func getUsers(userCH chan User) User {
	fmt.Println("FETCHING DATA")
	nameCHN := make(chan string)
	ageCHN := make(chan int)
	booksCHN := make(chan string)

	var bookList []string

	go getName(nameCHN)
	go getAge(ageCHN)
	go getBooks(booksCHN)

	name := <-nameCHN
	age := <-ageCHN

	for book := range booksCHN {
		bookList = append(bookList, book)
	}

	user := User{
		name:     name,
		age:      age,
		email:    name + "@gmail.com",
		bookList: bookList,
	}
	return user
}

func getName(namesCHN chan string) {
	counter := seededRand.Intn(4) + 1
	fmt.Println("fetching user name...")
	time.Sleep(2 * time.Second)

	names := Names()
	name := names[seededRand.Intn(47)+1]
	namesCHN <- string(name)
	fmt.Println("name fetched in: ", counter)

	close(namesCHN)
}

func getAge(ageCHN chan int) {
	counter := seededRand.Intn(4) + 1

	fmt.Println("fetching user age...")
	time.Sleep(time.Duration(counter))

	rand.Seed(time.Now().UnixNano())

	age := rand.Intn(98) + 1

	ageCHN <- age
	fmt.Println("age fetched in: ", counter)

	close(ageCHN)
}

func getBooks(bookCHN chan string) {

	fmt.Println("fetching user books...")
	nBooks := seededRand.Intn(9) + 1

	var wg sync.WaitGroup

	for i := 0; i < nBooks; i++ {
		counter := seededRand.Intn(4) + 1
		wg.Add(1)
		go getBook((seededRand.Intn(9) + 1), bookCHN, counter, &wg)
	}

	wg.Wait()
	fmt.Println(fmt.Sprintf("found %d books", nBooks))

	close(bookCHN)
}

func getBook(bookId int, bookCHN chan string, timeto int, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(timeto * int(time.Second)))
	books := Books()

	bookCHN <- books[bookId]
	fmt.Println("book fetched in", timeto)
}
