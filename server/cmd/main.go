package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../cmd/lib"
)

func main() {
	lib.InitDB("mydatabase.db")
	defer lib.CloseDB()

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/", homeHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// récupération des données via le formulaire
		nickname := r.FormValue("nickname")
		ageStr := r.FormValue("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}
		gender := r.FormValue("gender")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// enregistrement du user dans la base de données
		err = lib.InsertUser(nickname, age, gender, firstName, lastName, email, password)
		if err != nil {
			http.Error(w, "User registration error "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
