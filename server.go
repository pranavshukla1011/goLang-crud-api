package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type subjects struct {
	Math         int `json:"math" `
	Science       int `json:"science" `
	SocialScience int `json:"social_science" `
	English       int `json:"english" `
	Hindi         int `json:"hindi" `
}

type Student struct {
	Id        string   `json:"id,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	RollNo    string   `json:"roll_no,omitempty"`
	Dob       string   `json:"dob,omitempty"`
	Email     string   `json:"email,omitempty" `
	Phone      int     `json:"phone,omitempty"`
	Sub       subjects `json:"sub,omitempty" `
}

type studentsHandler struct {
	sync.Mutex
	store map[string]Student
}

//helper functions
func isEmailValid(e string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(e)
}

// api functions
func (h *studentsHandler) students (w http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			h.getAllstudents(w,r);
			return
		case "POST":
			h.addNewstudent(w,r);
			return
		default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

//creating store
func newstudentsHandler()*studentsHandler{
	return &studentsHandler{
		store: map[string]Student{},
	}
}

//server init
func main(){
	studentsHandler := newstudentsHandler()
	http.HandleFunc("/students", studentsHandler.students)
	http.HandleFunc("/students/", studentsHandler.student)
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		panic(err)
	}
}


func (h *studentsHandler) student (w http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET_STUDENT":
			h.getstudent(w,r);
			return
		case "UPDATE_STUDENT":
			h.updatestudent(w,r);
			return
		case "DELETE_STUDENT":
			h.deletestudent(w,r)
		default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}
//-----------------------------------/students----------------------------

//get all students
func (h *studentsHandler) getAllstudents (w http.ResponseWriter, r *http.Request) {
	students := make([]Student, len(h.store))

	h.Lock()
	i := 0
	for _, student := range h.store{
		students[i] = student
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(students)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write((jsonBytes))
}

//add new student
func (h *studentsHandler) addNewstudent (w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
		}
	
	ct := r.Header.Get("content-type")
	if ct != "application/json"{
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("err:  content-type: application/json but got '%s'", ct)))
	}
	var student Student
	e := json.Unmarshal(bodyBytes, &student)
	
	if e != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return 
	}

	if student.FirstName == ""{
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("err:  first-name but got empty string")))
		return
	}

	if student.LastName == ""{
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("err:  last-name but got empty string")))
		return
	}

	if student.RollNo == ""{
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("err: need roll number but got empty")))
		return
	}
	
	if !isEmailValid(student.Email){
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("err: enter valid email")))
		return
	}

	student.Id = student.RollNo
	h.Lock()

	h.store[student.Id]=  student

	defer h.Unlock()
}

//---------------------------/students/{id}---------------------------
//get student
func (h *studentsHandler) getstudent (w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	student, ok:= h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("invalid path")))
		return
	}
	
	jsonBytes, err := json.Marshal(student)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write((jsonBytes))
}

//update student

func (h *studentsHandler) updatestudent (w http.ResponseWriter, r *http.Request) {

	println("if no_updation found, please check your field names:\n first_name\nlast_name\n,omitempty\nroll_no\nemail\nphone\nsub { \nmath\nsciencesocial_science\nenglishhindi\n}")
	parts := strings.Split(r.URL.String(), "/") //check path is valid
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body) //check body is valid
	defer r.Body.Close()

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	ct := r.Header.Get("content-type")
	if ct != "application/json"{
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("err:  content-type: application/json but got '%s'", ct)))
	}
	
	h.Lock()
	student, ok:= h.store[parts[2]] //current student info
	defer h.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("invalid path")))
		return
	}

	var updateStudent Student
	e := json.Unmarshal(bodyBytes, &updateStudent) //updated student info

	if e != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return 
	}

	if updateStudent.FirstName != ""{
		student.FirstName = updateStudent.FirstName
	}

	if updateStudent.LastName != ""{
		student.LastName = updateStudent.LastName
	}

	if updateStudent.RollNo != ""{
		student.RollNo = updateStudent.RollNo
	}

	if updateStudent.Dob != ""{
		student.Dob = updateStudent.Dob
	}

	if updateStudent.Phone != 0{
		student.Phone = updateStudent.Phone
	}

	if updateStudent.Sub.English != 0{
		student.Sub.English = updateStudent.Sub.English
	}

	if updateStudent.Sub.Hindi != 0{
		student.Sub.Hindi = updateStudent.Sub.Hindi
	}

	if updateStudent.Sub.Math != 0{
		student.Sub.Math = updateStudent.Sub.Math
	}

	if updateStudent.Sub.Science != 0{
		student.Sub.Science = updateStudent.Sub.Science
	}

	if updateStudent.Sub.SocialScience != 0{
		student.Sub.SocialScience = updateStudent.Sub.SocialScience
	}

	if !isEmailValid(student.Email){
		student.Email = updateStudent.Email
	}

	h.store[parts[2]] = student; //updating in db

	w.WriteHeader(http.StatusOK)
	fmt.Println(student)
	w.Write([]byte(fmt.Sprintf("'%s' fields updated", parts[2])))
}

//delete student
func (h *studentsHandler) deletestudent (w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	_, ok:= h.store[parts[2]] //current student info
	defer h.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("invalid path")))
		return
	}
	delete(h.store, parts[2])
	println("deleted field %s", parts[2]);
}
