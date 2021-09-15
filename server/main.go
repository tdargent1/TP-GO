package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type DataEntry struct {
	entry  string
	author string
}

func getDate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(404)
		fmt.Println("Something went bad")
		fmt.Fprintln(w, "Something went bad")
	} else {
		result := fmt.Sprintf("%02dh%02d", time.Now().Hour(), time.Now().Minute())
		fmt.Fprintln(w, result)
	}
}

func addEntry(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(404)
		fmt.Println("Something went bad")
		fmt.Fprintln(w, "Something went bad")
	} else {
		if err := req.ParseForm(); err != nil {
			w.WriteHeader(404)
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}

		var dataEntry DataEntry
		dataEntry.entry = req.FormValue("entry")
		dataEntry.author = req.FormValue("author")

		saveEntry(&dataEntry)
		result := fmt.Sprint(dataEntry.author, ":", dataEntry.entry)
		fmt.Fprintln(w, result)
	}
}

func saveEntry(dataEntry *DataEntry) {
	saveFile, err := os.OpenFile("./miniapi.data", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	defer saveFile.Close()

	w := bufio.NewWriter(saveFile)

	if err == nil {
		fmt.Fprintf(w, "%s:%s\n", dataEntry.author, dataEntry.entry)
	}

	w.Flush()
}

func getEntries(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(404)
		fmt.Println("Something went bad")
		fmt.Fprintln(w, "Something went bad")
	} else {
		saveData, err := os.ReadFile("./miniapi.data")
		if err == nil {
			entries := strings.Split(string(saveData), "\n")
			for _, entry := range entries {
				if entry != "" {
					result := fmt.Sprint(strings.Split(entry, `:`)[1])
					fmt.Fprintln(w, result)
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/", getDate)
	http.HandleFunc("/add", addEntry)
	http.HandleFunc("/entries", getEntries)
	http.ListenAndServe(":4567", nil)
}
