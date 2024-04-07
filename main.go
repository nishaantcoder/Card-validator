package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/validate", validateHandler)

	port := 3000
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := strings.TrimSpace(r.FormValue("cardNumber"))

	isValid := luhnAlgorithm(cardNumber)

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]interface{}{"CardNumber": cardNumber, "Result": getResultMessage(isValid)})
}

func luhnAlgorithm(cardNumber string) bool {
	var ss string
	for _, r := range cardNumber {
		if !unicode.IsSpace(r) {
			ss += string(r)
		}
	}
	var sum int64 = 0
	parity := len(ss) % 2

	cardNumWithoutChecksum := ss[:len(ss)-1]

	for i, v := range cardNumWithoutChecksum {
		item, err := strconv.Atoi(string(v))

		if err != nil {
			fmt.Println(err)
			return false
		}
		if int64(i)%2 != int64(parity) {
			sum += int64(item)
		} else if item > 4 {
			sum += int64(2*item - 9)
		} else {
			sum += int64(2 * item)
		}
	}

	checkDigit, err := strconv.Atoi(ss[len(ss)-1:])

	if err != nil {
		fmt.Println(err)
		return false
	}
	SumMod := sum % 10

	if SumMod == int64(0) {
		return SumMod == int64(checkDigit)
	}
	return int64(10)-SumMod == int64(checkDigit)
}

func getResultMessage(isValid bool) string {
	if isValid {
		return "Valid Credit Card!"
	}
	return "Invalid Credit Card!"
}
