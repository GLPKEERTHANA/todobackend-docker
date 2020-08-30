package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
)

type Delete struct {
	User_id int `json:"user_id"`
	Feed_id int `json:"feed_id"`
}
type DeleteFeed struct {
	Success string `json:"success"`
	Error   string `json:"error"`
}
type Feed struct {
	Feed_id     int    `json:"feed_id"`
	Feed        string `json:"feed"`
	Feed_status string `json:"feed_status"`
	User_id     int    `json:"user_id"`
}
type Users struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
type User struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserData struct {
	Userdata Users `json:"userData"`
}
type UserID struct {
	User_id int `json:"user_id"`
}
type FeedData struct {
	Feeddata []Feed `json:"feedData"`
}
type LoginData struct {
	Feeddata []Users `json:"feedData"`
}
type UserFeed struct {
	User_id int    `json:"user_id"`
	Feed    string `json:"feed"`
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
	var user Users
	fmt.Println("-------------------------------------Login called--------------------------------------------------------")
	InfoLogger.Println("Starting the application...")
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb") username:password
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	
	if err != nil {
		panic("failed to connect database")
		ErrorLogger.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	var login User
	_ = json.NewDecoder(r.Body).Decode(&login)
	var username string = login.Username
	var password string = login.Password
	db.Where("Username = ? AND Password = ?", username, password).Find(&user)
	json1 := UserData{}
	json1.Userdata = user
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(json1)
	InfoLogger.Println("Login accepted for " + user.Username)
}
func feed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Called to get incomplete tasks of a user")
	vars := mux.Vars(r)
	user_id := vars["id"]
	fmt.Println("-------------------feed called---------------------------------------------")
	//var login UserID
	//	_ = json.NewDecoder(r.Body).Decode(&login)
	//	var user_id int = login.User_id
	var status string = "F"
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		panic("failed to connect database")
		ErrorLogger.Println(err.Error())
	}
	defer db.Close()
	feed := []Feed{}
	db.Where("user_id = ? AND feed_status=?", user_id, status).Find(&feed)

	json1 := FeedData{}
	json1.Feeddata = feed
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(json1)
	InfoLogger.Println("Sent the information about incomplete tasks of a user")
}

func feedDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Called to get completed tasks of a user")
	/*var login UserID
	_ = json.NewDecoder(r.Body).Decode(&login)
	var user_id int = login.User_id*/
	vars := mux.Vars(r)
	user_id := vars["id"]
	fmt.Println("-------------------feedDone called---------------------------------------------")
	var status string = "T"
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		panic("failed to connect database")
		ErrorLogger.Println(err.Error())
	}
	defer db.Close()
	feed := []Feed{}
	db.Where("user_id = ? AND feed_status=?", user_id, status).Find(&feed)

	json1 := FeedData{}
	json1.Feeddata = feed
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(json1)
	InfoLogger.Println("Sent the information about incomplete tasks of a user")
}
func feedusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Called to get the usernames of all users")
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		panic("failed to connect database")
		ErrorLogger.Println(err.Error())
	}
	defer db.Close()
	var username string = "admin"
	user := []Users{}
	db.Where("Username  <>  ? ", username).Find(&user)
	json1 := LoginData{}
	json1.Feeddata = user
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(json1)
	InfoLogger.Println("Sent the usernames of all users")
}

func feedUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var login UserFeed
	_ = json.NewDecoder(r.Body).Decode(&login)
	var userid int = login.User_id
	var feed string = login.Feed
	InfoLogger.Println("New Task addition " + feed)
	var status string = "F"

	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		ErrorLogger.Println(err.Error())
		panic("failed to connect database")

	}
	InfoLogger.Println("New Task addition complete.. ")
	feed1 := Feed{Feed: feed, Feed_status: status, User_id: userid}
	db.NewRecord(feed1)
	db.Create(&feed1)
}

func feedDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Task Deletion called")
	var login Delete
	_ = json.NewDecoder(r.Body).Decode(&login)
	var userid int = login.User_id
	var feedid int = login.Feed_id

	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		panic("failed to connect database")

	}
	db.Where("user_id = ? AND feed_id=?", userid, feedid).Delete(Feed{})
	json1 := DeleteFeed{}
	w.WriteHeader(200)
	if err != nil {
		ErrorLogger.Println(err.Error() + " Delete ERROR ...")
		json1.Error = "Delete error"
		json.NewEncoder(w).Encode(json1)
	}
	json1.Success = "Feed deleted"
	json.NewEncoder(w).Encode(json1)
	InfoLogger.Println("Feed deletion done... ")

}
func feedstatus(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Task Status updation called")
	w.Header().Set("Content-Type", "application/json")
	var login Delete
	_ = json.NewDecoder(r.Body).Decode(&login)
	userid := login.User_id
	feedid := login.Feed_id
	fmt.Println("--------------------------------------------------------")
	fmt.Println(userid)
	fmt.Println(feedid)

	fmt.Println("--------------------------------------------------------")
	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/reactdb")
	db, err := gorm.Open("mysql", "sql12362860:nxBBF29dcx@tcp(sql12.freemysqlhosting.net:3306)/sql12362860")
	if err != nil {
		panic("failed to connect database")
		ErrorLogger.Println(err.Error())
	}
	feed := Feed{}
	db.Where("user_id = ? AND feed_id=? ", userid, feedid).Find(&feed)
	var status string = feed.Feed_status
	json1 := DeleteFeed{}
	w.WriteHeader(200)
	if err != nil {
		ErrorLogger.Println(err.Error() + " Task updation error.. ")
		json1.Error = "Feed error"
		json.NewEncoder(w).Encode(json1)
	} else {

		var status1 string = ""
		if status == "T" {
			status1 = "F"
			db.Table("feeds").Where("feed_id = ? and user_id=?", feedid, userid).Updates(map[string]interface{}{"feed_status": status1})

		} else {
			status1 = "T"
			db.Table("feeds").Where("feed_id = ? and user_id=?", feedid, userid).Updates(map[string]interface{}{"feed_status": status1})

		}
		json1.Success = "Feed updated"
		json.NewEncoder(w).Encode(json1)
		InfoLogger.Println("Task status updation done... ")
	}
}

func main() {
	mux := mux.NewRouter()
	port:=os.Getenv("PORT")
	log.Println("Server started on: http://localhost:8080")
	//mux.HandleFunc("/login", login).Methods("POST")
	mux.HandleFunc("/todo/login", login).Methods("POST")
	mux.HandleFunc("/todo/feed/{id}", feed).Methods("GET")
	mux.HandleFunc("/todo/feedDone/{id}", feedDone).Methods("GET")
	mux.HandleFunc("/todo/users", feedusers).Methods("GET")
	mux.HandleFunc("/todo", feedUpdate).Methods("POST")
	mux.HandleFunc("/todo/{id}", feedDelete).Methods("DELETE")
	mux.HandleFunc("/todo/{id}", feedstatus).Methods("PUT")
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

}
