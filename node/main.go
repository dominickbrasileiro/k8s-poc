package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
)

func getData(w http.ResponseWriter) {
	files, err := os.ReadDir("./data/")

	if err != nil {
		panic(err)
	}

	ch := make(chan string, len(files))
	wg := sync.WaitGroup{}

	for _, fileInfo := range files {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			content, err := os.ReadFile("./data/" + name)

			if err != nil {
				panic(err)
			}

			ch <- string(content)
		}(fileInfo.Name())
	}

	wg.Wait()
	close(ch)

	contents := make([]string, 0)

	for c := range ch {
		contents = append(contents, c)
	}

	b, err := json.Marshal(contents)

	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func saveData(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()

	file, err := os.Create("./data/" + id + ".dat")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.ReadFrom(r.Body)
}

func main() {
	os.Mkdir("./data", 0644)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getData(w)
		case http.MethodPost:
			saveData(w, r)
		case http.MethodDelete:
			os.Exit(1)
		default:
			w.WriteHeader(405)
			w.Write([]byte("Method Not Allowed"))
		}
	})

	fmt.Println("Listening http://0.0.0.0:32000")

	http.ListenAndServe(":32000", nil)
}
