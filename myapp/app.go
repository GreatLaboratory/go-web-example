package myapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// User 유저 인터페이스
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var userMap map[int]*User
var lastID int

type fooHandler struct{}

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello World") // Fprint는 writer에 프린트를 하라는 뜻
}

func getAllUserInfoHandler(rw http.ResponseWriter, r *http.Request) {
	if len(userMap) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "No Users")
		return
	}
	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}
	data, _ := json.Marshal(users)
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, string(data))
}

func getUserInfoHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		return
	}
	user, ok := userMap[id]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "No User Id:", id)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(rw, string(data))
}

func createUserHandler(rw http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		return
	}

	// Create User
	lastID++
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(rw, string(data))
}

func deleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		return
	}
	_, ok := userMap[id]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "No User Id:", id)
		return
	}
	delete(userMap, id)

	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, "Deleted User Id:", id)
}

func updateUserHandler(rw http.ResponseWriter, r *http.Request) {
	updateUser := new(User)
	err := json.NewDecoder(r.Body).Decode(updateUser) // user interface의 데이터가 아닌 경우는 err로 예외처리된다.
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "Bad Request : ", err)
		return
	}

	user, ok := userMap[updateUser.ID]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "No User Id:", updateUser.ID)
		return
	}

	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(rw, string(data))
}

func (f *fooHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) { // mux.Handle의 두번째 인자인 http.Handler 인터페이스에서 구현해야하는 함수 이름이 ServeHTTP이다.
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user) // user interface의 데이터가 아닌 경우는 err로 예외처리된다.
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "Bad Request : ", err)
		return
	}
	user.CreatedAt = time.Now()
	// fmt.Printf("Hi. My name is %s", user.FirstName+" "+user.LastName)

	data, _ := json.Marshal(user) // CreatedAt이 업데이트된 user struct를 다시 response json형태로 바꿔주는 작업
	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprint(rw, string(data)) // 현재 data는 byte[]이므로 이걸 스트링으로 변환해야한다.
}
func barHandler(rw http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Bar"
	}
	fmt.Fprintf(rw, "Hello %s!", name)
}

func uploadsHandler(rw http.ResponseWriter, r *http.Request) {
	// 1. 클라가 보내주는 파일을 read하는 파트
	uploadFile, header, err := r.FormFile("upload_file")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		return
	}
	defer uploadFile.Close()

	// 2. 새로운 파일을 저장할 공간을 만드는 파트
	dirname := "./uploads"
	os.Mkdir(dirname, 0777)
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
		return
	}

	// 3. 생긴 공간에 받은 파일을 copy하는 파트
	io.Copy(file, uploadFile)
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, filepath)
}

// NewHTTPHandler -> http.Handler를 리턴하는 함수
func NewHTTPHandler() http.Handler {
	userMap = make(map[int]*User)
	lastID = 0
	mux := mux.NewRouter()
	// mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/bar", barHandler)
	mux.Handle("/foo", &fooHandler{})
	mux.HandleFunc("/uploads", uploadsHandler)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("public")))) // mux로 하려면 이렇게 굉장히 까다롭게 경로를 지정해줘야 한다.
	mux.HandleFunc("/users", getAllUserInfoHandler).Methods("GET")
	mux.HandleFunc("/user", createUserHandler).Methods("POST")
	mux.HandleFunc("/user", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/user/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	mux.HandleFunc("/user/{id:[0-9]+}", getUserInfoHandler).Methods("GET")

	return mux
}
