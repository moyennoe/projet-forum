package account

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type dbLogin struct {
	Id       int
	Pseudo   string
	Email    string
	Password string
}

var store = sessions.NewCookieStore([]byte("mysession"))

func initDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Fatal(err)
	}
	statement := `CREATE TABLE IF NOT EXISTS dblogin (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, pseudo TEXT NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL)`

	_, err = db.Exec(statement)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
func compareInputAndDB(email string, password string) bool {

	var u dbLogin
	db := initDatabase("dblogin.db")

	rows := selectAllFromTable(db, "dblogin")

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Pseudo, &u.Email, &u.Password)
		if u.Email == email && u.Password == password {
			return true
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return false
	// Fonction pour comparer un input de register avec la db déjà existante, pour
	// éviter les doubles comptes avec le même pseudo ou email
}

func insertIntoTypes(db *sql.DB, pseudo string, email string, password string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO dblogin (pseudo, email, password) VALUES (?, ?, ?)`, pseudo, email, password)

	return result.LastInsertId()
}
func selectAllFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

func Index(response http.ResponseWriter, request *http.Request) {
	tmp, _ := template.ParseFiles("html/index.html")
	tmp.Execute(response, nil)
}

func Login(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	pseudoSignup := request.Form.Get("signup-pseudo")
	emailSignup := request.Form.Get("signup-email")
	passwordSignup := request.Form.Get("signup-password")

	emailLogin := request.Form.Get("login-email")
	passwordLogin := request.Form.Get("login-password")

	if emailLogin != "" && passwordLogin != "" {
		boolLogin := compareInputAndDB(emailLogin, passwordLogin)
		if boolLogin == true {
			fmt.Println("emailLogin: ", emailLogin)
			fmt.Println("passwordLogin: ", passwordLogin)
			session, _ := store.Get(request, "mysession")
			session.Values["username"] = emailLogin
			session.Save(request, response)
			http.Redirect(response, request, "/welcome", http.StatusSeeOther)
		}
	}
	if pseudoSignup != "" && emailSignup != "" && passwordSignup != "" {
		fmt.Println("pseudoSignup: ", pseudoSignup)
		fmt.Println("emailSignup: ", emailSignup)
		fmt.Println("passwordSignup: ", passwordSignup)

		db := initDatabase("dblogin.db")
		defer db.Close()
		insertIntoTypes(db, pseudoSignup, emailSignup, passwordSignup)

		session, _ := store.Get(request, "mysession")
		session.Values["username"] = emailSignup
		session.Save(request, response)
		http.Redirect(response, request, "/welcome", http.StatusSeeOther)

	} else {
		data := map[string]interface{}{
			"err": "Email & Password didn't match",
		}
		tmp, _ := template.ParseFiles("html/index.html")
		tmp.Execute(response, data)
	}

}

func Welcome(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	username := session.Values["username"]
	fmt.Println("username :", username)
	data := map[string]interface{}{
		"username": username,
	}
	tmp, _ := template.ParseFiles("html/welcome.html")
	tmp.Execute(response, data)
}

func Logout(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	session.Options.MaxAge = -1
	session.Save(request, response)
	http.Redirect(response, request, "/index", http.StatusSeeOther)
}
