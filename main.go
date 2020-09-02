package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
)

type ErrorMesage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Error struct {
	Status int         `json:"status"`
	Type   string      `json:"type"`
	Errors ErrorMesage `"json:errors"`
}
type Feed struct {
	Feed_id     int    `json:"feed_id"`
	Feed        string `json:"feed"`
	Feed_status string `json:"feed_status"`
	User_id     int    `json:"user_id"`
}

type DeleteFeed struct {
	Success string `json:"success"`
	Error   string `json:"error"`
}

type Users struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UserData struct {
	Userdata Users `json:"userData"`
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func login(w http.ResponseWriter, r *http.Request) {
	var ErrorObject Error
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	defer db.Close()
	var user Users
	w.Header().Set("Content-Type", "application/json")
	var login Users
	_ = json.NewDecoder(r.Body).Decode(&login)
	var username string = login.Username
	var password string = login.Password

	var count int
	db.Where("Username = ? AND Password = ?", username, password).Find(&user).Count(&count)
	if count > 0 {
		retunedObject := UserData{}
		retunedObject.Userdata = user
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(retunedObject)
	} else {
		retunedObject := DeleteFeed{}
		retunedObject.Error = "Wrong username and password"
		json.NewEncoder(w).Encode(retunedObject)
	}
	InfoLogger.Println("Login accepted for " + user.Username)
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user Users
	var ErrorObject Error
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	defer db.Close()
	var login Users
	_ = json.NewDecoder(r.Body).Decode(&login)
	var username string = login.Username
	var password string = login.Password
	var email string = login.Email
	var name string = login.Name

	username_regex := regexp.MustCompile(`^[A-Za-z0-9_]{4,10}`)
	var usernamevalid = username_regex.MatchString(username)

	password_regex := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_]{4,20}`)
	var passwordvalid = password_regex.MatchString(password)

	email_regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var emailvalid = email_regex.MatchString(email)

	retunedObject := DeleteFeed{}
	if usernamevalid == false {
		retunedObject.Error = "Invalid username"
		json.NewEncoder(w).Encode(retunedObject)
	} else if emailvalid == false {
		retunedObject.Error = "Invalid email"
		json.NewEncoder(w).Encode(retunedObject)

	} else if passwordvalid == false {
		retunedObject.Error = "Invalid password"
		json.NewEncoder(w).Encode(retunedObject)
	} else {
		var count int
		db.Where("Username = ? OR Email = ?", username, email).Find(&user).Count(&count)
		if count == 0 {
			feed1 := Users{Username: username, Password: password, Name: name, Email: email}
			db.NewRecord(feed1)
			db.Create(&feed1)
			db.Where("Username = ? AND Password = ?", username, password).Find(&user).Count(&count)
			retunedObject := UserData{}
			retunedObject.Userdata = user
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(retunedObject)
		} else {
			retunedObject.Error = "username or email exists"
			json.NewEncoder(w).Encode(retunedObject)
		}
	}
	fmt.Println(login)

	InfoLogger.Println("Login accepted for " + username)
}
func ErrorObjectInitialisation(title string, message string, status int, errorType string) Error {
	var e Error
	var errorMessage ErrorMesage
	errorMessage.Title = title
	errorMessage.Message = message
	e.Status = status
	e.Type = errorType
	e.Errors = errorMessage
	return e
}
func taskStatusFalse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Called to get incomplete tasks of a user")
	vars := mux.Vars(r)
	user_id := vars["userId"]
	var user Users
	var usercount int
	var ErrorObject Error
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	db.Where("user_id = ?", user_id).Find(&user).Count(&usercount)
	if usercount == 0 {
		ErrorObject = ErrorObjectInitialisation("UserId Not Found", "UserId Not Found", 404, "UserId Not Found")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorObject)
		ErrorLogger.Println(err)

	} else {
		feed := []Feed{}
		var status string = "F"
		defer db.Close()
		db.Where("user_id = ? AND feed_status=?", user_id, status).Find(&feed)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(feed)
		InfoLogger.Println("Sent the information about incomplete tasks of a user")
	}

}
func taskStatusTrue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Called to get completed tasks of a user")
	vars := mux.Vars(r)
	user_id := vars["userId"]
	var user Users
	var usercount int
	var ErrorObject Error
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	db.Where("user_id = ?", user_id).Find(&user).Count(&usercount)
	if usercount == 0 {
		ErrorObject = ErrorObjectInitialisation("UserId Not Found", "UserId Not Found", 404, "UserId Not Found")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorObject)
		ErrorLogger.Println(err)

	} else {
		feed := []Feed{}
		var status string = "T"
		defer db.Close()
		db.Where("user_id = ? AND feed_status=?", user_id, status).Find(&feed)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(feed)
		InfoLogger.Println("Sent the information about completed tasks of a user")
	}

}

func taskUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Called to get the usernames of all users")
	var ErrorObject Error
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	defer db.Close()
	var username string = "admin"
	user := []Users{}
	db.Where("Username  <>  ? ", username).Find(&user)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(user)
	InfoLogger.Println("Sent the usernames of all users")
}

func feedUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var login Feed
	_ = json.NewDecoder(r.Body).Decode(&login)
	var user_id int = login.User_id
	var feed string = login.Feed
	InfoLogger.Println("New Task addition " + feed)
	var status string = "F"
	var usercount int
	var user Users
	var ErrorObject Error
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	db.Where("user_id = ?", user_id).Find(&user).Count(&usercount)
	if usercount == 0 {
		ErrorObject = ErrorObjectInitialisation("UserId Not Found", "UserId Not Found", 404, "UserId Not Found")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorObject)
		ErrorLogger.Println(err)

	}
	feed1 := Feed{Feed: feed, Feed_status: status, User_id: user_id}
	db.NewRecord(feed1)
	db.Create(&feed1)
	ErrorObject = ErrorObjectInitialisation("Task Inserted", "Task Inserted", 201, "Task Inserted")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(ErrorObject)
	InfoLogger.Println("New Task addition complete.. ")

}

func feedDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Task Deletion called")
	var ErrorObject Error
	vars := mux.Vars(r)
	feedid := vars["feedId"]
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	defer db.Close()
	var count int
	var user Feed
	db.Where("feed_id=?", feedid).Find(&user).Count(&count)
	if count == 0 {
		ErrorObject = ErrorObjectInitialisation("Task Not Found", "Task Not Found", 404, "Task Not Nound")
		ErrorLogger.Println(err)
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	db.Where("feed_id=?", feedid).Delete(Feed{})
	ErrorObject = ErrorObjectInitialisation("Task Deletion Success", "Task Deletion Success", 200, "Task Deletion Success")
	ErrorLogger.Println(err)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(ErrorObject)
	InfoLogger.Println("Feed deletion done... ")

}

func feedstatus(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Task Status updation called")
	w.Header().Set("Content-Type", "application/json")
	var ErrorObject Error
	vars := mux.Vars(r)
	feedid := vars["feedId"]

	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorObject = ErrorObjectInitialisation("Internal Server Error", "Internal Server Error", 500, "Internal Server Error")
		ErrorLogger.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	defer db.Close()
	var count int
	var user Feed
	db.Where("feed_id=? ", feedid).Find(&user).Count(&count)
	if count == 0 {
		ErrorObject = ErrorObjectInitialisation("Task Not Found", "Task Not Found", 404, "Task Not Nound")
		ErrorLogger.Println(err)
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorObject)
		return
	}
	feed := Feed{}
	db.Where("feed_id=? ", feedid).Find(&feed)
	var status string = feed.Feed_status
	var status1 string = ""
	if status == "T" {
		status1 = "F"
		db.Table("feeds").Where("feed_id = ? ", feedid).Updates(map[string]interface{}{"feed_status": status1})

	} else {
		status1 = "T"
		db.Table("feeds").Where("feed_id = ? ", feedid).Updates(map[string]interface{}{"feed_status": status1})

	}
	ErrorObject = ErrorObjectInitialisation("Task Updation Success", "Task Updation Success", 200, "Task Updation Success")
	ErrorLogger.Println(err)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(ErrorObject)
	InfoLogger.Println("Task status updation done... ")

}

func main() {
	mux := mux.NewRouter()
	port := os.Getenv("PORT")
	log.Println("Server started on: http://localhost:8080")
	mux.HandleFunc("/todo/users/signup", signup).Methods("POST")
	mux.HandleFunc("/todo/users/login", login).Methods("POST")
	mux.HandleFunc("/todo/task/{feedId}", feedDelete).Methods("DELETE")
	mux.HandleFunc("/todo/task/{feedId}", feedstatus).Methods("PUT")
	mux.HandleFunc("/todo/task/statusFalse/{userId}", taskStatusFalse).Methods("GET")
	mux.HandleFunc("/todo/task/statusTrue/{userId}", taskStatusTrue).Methods("GET")
	mux.HandleFunc("/todo/users", taskUsers).Methods("GET")
	mux.HandleFunc("/todo/task", feedUpdate).Methods("POST")
	handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete},
		AllowedHeaders: []string{"*"},
		Debug:          true,
	})

	handler = c.Handler(handler)

	log.Fatal(http.ListenAndServe(":"+port, handler))
	//log.Fatal(http.ListenAndServe(":8080", handler))
}
