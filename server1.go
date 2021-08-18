package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UsersArray struct {
	Id            string   `json:"Id"`
	SecretCode    string   `json:"SecretCode"`
	Name          string   `json:"Name"`
	Address       string   `json:"Address"`
	PhoneNumber   string   `json:"PhoneNumber"`
	UserType      string   `json:"UserType"`
	Requested     []string `json:"Requested"`
	PendingReq    []string `json:"PendingReq"`
	ConnectedUser []string `json:"ConnectedUser"`
}

type UsersArrayConstraint struct {
	Id            string   `json:"Id"`
	Name          string   `json:"Name"`
	Address       string   `json:"Address"`
	PhoneNumber   string   `json:"PhoneNumber"`
	UserType      string   `json:"UserType"`
	Requested     []string `json:"Requested"`
	PendingReq    []string `json:"PendingReq"`
	ConnectedUser []string `json:"ConnectedUser"`
}
type Userstemp struct {
	SecretCode string `json:"SecretCode"`
	Id         string `json:"Id"`
}

var count = 2
var usersmap map[string]UsersArray
var p sync.WaitGroup
var m sync.Mutex

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	var user UsersArray
	user = usersmap[login.SecretCode]
	json.NewEncoder(w).Encode(user)
}

func increment(wg *sync.WaitGroup, m *sync.Mutex) {
	m.Lock()
	count = count + 1
	m.Unlock()
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var users UsersArray
	b := 999999
	a := 111111
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &users)
	//p.Add(1)
	increment(&p, &m)
	//p.Wait()

	users.Id = strconv.Itoa(count)
	rand.Seed(time.Now().UnixNano())
	users.SecretCode = strconv.Itoa(a + rand.Intn(b-a+1))
	usersmap[users.SecretCode] = users
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var users UsersArray
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &users)
	userSd := users.SecretCode
	for m, n := range usersmap {
		if n.SecretCode == userSd {
			if users.Name != " " {
				n.Name = users.Name
			}
			if users.Address != " " {
				n.Address = users.Address
			}
			if users.PhoneNumber != " " {
				n.PhoneNumber = users.PhoneNumber
			}
			if users.UserType != " " {
				n.UserType = users.UserType
			}
			if users.Requested != nil {
				n.Requested = users.Requested
			}
			if users.PendingReq != nil {
				n.PendingReq = users.PendingReq
			}
			if users.ConnectedUser != nil {
				n.ConnectedUser = users.ConnectedUser
			}
			usersmap[m] = n
		}
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	var user UsersArray
	user = usersmap[login.Id]
	json.NewEncoder(w).Encode(user)
}

func GetAllusers(w http.ResponseWriter, r *http.Request) {
	for k := range usersmap {
		json.NewEncoder(w).Encode(usersmap[k])
	}
}

func GetAllDonors(w http.ResponseWriter, r *http.Request) {

	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	PatientSd := login.SecretCode
	PatientUser := usersmap[PatientSd]
	if PatientUser.UserType == "Patient" {
		for k := range usersmap {
			if usersmap[k].UserType == "Donor" {
				var DonorUser UsersArrayConstraint
				found := 0
				DonorUser.Id = usersmap[k].Id
				DonorUser.Name = usersmap[k].Name
				for _, queryDonor := range usersmap[k].ConnectedUser {
					if queryDonor == PatientSd {
						found = 1
					}
				}
				if found == 1 {
					DonorUser.Address = usersmap[k].Address
					DonorUser.PhoneNumber = usersmap[k].PhoneNumber
					DonorUser.UserType = usersmap[k].UserType
				}
				if found == 0 {
					DonorUser.Address = " "
					DonorUser.PhoneNumber = " "
					DonorUser.UserType = " "
				}
				DonorUser.Requested = usersmap[k].Requested
				DonorUser.PendingReq = usersmap[k].PendingReq
				DonorUser.ConnectedUser = usersmap[k].ConnectedUser
				json.NewEncoder(w).Encode(DonorUser)

			}
		}

	}

}

func GetAllPatients(w http.ResponseWriter, r *http.Request) {
	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	DonorSd := login.SecretCode
	DonorUser := usersmap[DonorSd]
	if DonorUser.UserType == "Donor" {
		for k := range usersmap {
			if usersmap[k].UserType == "Patient" {
				var PatientUser UsersArrayConstraint
				found := 0
				PatientUser.Id = usersmap[k].Id
				PatientUser.Name = usersmap[k].Name
				for _, queryPatient := range usersmap[k].ConnectedUser {
					if queryPatient == DonorSd {
						found = 1
					}
				}
				if found == 1 {
					PatientUser.Address = usersmap[k].Address
					PatientUser.PhoneNumber = usersmap[k].PhoneNumber
					PatientUser.UserType = usersmap[k].UserType
				}
				if found == 0 {
					PatientUser.Address = " "
					PatientUser.PhoneNumber = " "
					PatientUser.UserType = " "
				}
				PatientUser.Requested = usersmap[k].Requested
				PatientUser.PendingReq = usersmap[k].PendingReq
				PatientUser.ConnectedUser = usersmap[k].ConnectedUser
				json.NewEncoder(w).Encode(PatientUser)
			}
		}
	}
}

func SendRequest(w http.ResponseWriter, r *http.Request) {

	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	PatientSd := login.SecretCode
	DonorId := login.Id
	PatientId := usersmap[PatientSd].Id
	PatientUser := usersmap[login.SecretCode]
	PatientUser.PendingReq = append(PatientUser.PendingReq, DonorId)
	var DonorUser UsersArray
	for k, v := range usersmap {
		if usersmap[k].Id == DonorId {
			DonorUser = v
		}
	}
	DonorUser.Requested = append(DonorUser.Requested, PatientId)

}

func AcceptRequest(w http.ResponseWriter, r *http.Request) {

	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	PatientId := login.Id
	var PatientUser UsersArray
	DonorSd := login.SecretCode
	DonorId := usersmap[DonorSd].Id
	DonorUser := usersmap[DonorSd]
	for j, _ := range DonorUser.Requested {
		if DonorUser.Requested[j] == PatientId {
			DonorUser.Requested = append(DonorUser.Requested[:j], DonorUser.Requested[j+1:]...)
			DonorUser.ConnectedUser = append(DonorUser.ConnectedUser, PatientId)
		}
	}

	for k, v := range usersmap {
		if usersmap[k].Id == PatientId {
			PatientUser = v
		}
	}
	for i, _ := range PatientUser.PendingReq {
		if PatientUser.PendingReq[i] == DonorId {
			PatientUser.PendingReq = append(PatientUser.PendingReq[:i], PatientUser.PendingReq[i+1:]...)
			PatientUser.ConnectedUser = append(PatientUser.ConnectedUser, DonorId)
		}
	}
}

func CancelConnection(w http.ResponseWriter, r *http.Request) {
	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	PatientSd := login.SecretCode
	DonorId := login.Id
	PatientId := usersmap[PatientSd].Id
	PatientUser := usersmap[PatientSd]
	var DonorUser UsersArray
	for i, _ := range PatientUser.ConnectedUser {
		if PatientUser.ConnectedUser[i] == DonorId {
			PatientUser.ConnectedUser = append(PatientUser.ConnectedUser[:i], PatientUser.ConnectedUser[i+1:]...)
		}
	}
	for k, v := range usersmap {
		if usersmap[k].Id == DonorId {
			DonorUser = v
		}
	}
	for j, _ := range DonorUser.ConnectedUser {
		if DonorUser.ConnectedUser[j] == PatientId {

			DonorUser.ConnectedUser = append(DonorUser.ConnectedUser[:j], DonorUser.ConnectedUser[j+1:]...)
		}
	}

}

func CancelRequest(w http.ResponseWriter, r *http.Request) {
	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	PatientSd := login.SecretCode
	DonorId := login.Id
	PatientId := usersmap[PatientSd].Id
	PatientUser := usersmap[PatientSd]
	var DonorUser UsersArray
	for i, _ := range PatientUser.PendingReq {
		if PatientUser.PendingReq[i] == DonorId {
			PatientUser.PendingReq = append(PatientUser.PendingReq[:i], PatientUser.PendingReq[i+1:]...)
		}
	}
	for k, v := range usersmap {
		if usersmap[k].Id == DonorId {
			DonorUser = v
		}
	}
	for j, _ := range DonorUser.Requested {
		if DonorUser.Requested[j] == PatientId {
			DonorUser.Requested = append(DonorUser.Requested[:j], DonorUser.Requested[j+1:]...)
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var login Userstemp
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	var user UsersArray
	user = usersmap[login.SecretCode]
	delete(usersmap, login.SecretCode)
	json.NewEncoder(w).Encode(user)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login", LoginUser)
	http.HandleFunc("/createUser", CreateUser)
	http.HandleFunc("/GetAllusers", GetAllusers)
	http.HandleFunc("/getUser", GetUser)
	http.HandleFunc("/GetAllDonors", GetAllDonors)
	http.HandleFunc("/GetAllPatients", GetAllPatients)
	http.HandleFunc("/SendRequest", SendRequest)
	http.HandleFunc("/acceptRequest", AcceptRequest)
	http.HandleFunc("/cancelConnection", CancelConnection)
	http.HandleFunc("/cancelRequest", CancelRequest)
	http.HandleFunc("/deleteUser", DeleteUser)
	http.HandleFunc("/UpdateUser", UpdateUser)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	usersmap = make(map[string]UsersArray)
	usersmap["898982"] = UsersArray{Id: "1", SecretCode: "898982", Name: "Kireeti", Address: "kakinada", PhoneNumber: "9381305447", UserType: "Donor", Requested: []string{""}, PendingReq: []string{""}, ConnectedUser: []string{"767262"}}
	usersmap["767262"] = UsersArray{Id: "2", SecretCode: "767262", Name: "hello", Address: "Bangalore", PhoneNumber: "9059800515", UserType: "Patient", Requested: []string{""}, PendingReq: []string{""}, ConnectedUser: []string{"898982"}}
	handleRequests()
}