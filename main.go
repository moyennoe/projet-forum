package main

import (
	"fmt"
	"net/http"

	account "./code"
)

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", account.Index)
	http.HandleFunc("/index", account.Index)
	http.HandleFunc("/login", account.Login)
	http.HandleFunc("/welcome", account.Welcome)
	http.HandleFunc("/logout", account.Logout)
	fmt.Println("8080")
	http.ListenAndServe(":8080", nil)
}
