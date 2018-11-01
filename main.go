package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Engineer struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ImgURL       string `json:"imgUrl"`
	Title        string `json:"title"`
	Department   string `json:"department"`
	HourlyWage   int    `json:"hourlyWage"`
	HoursPerWeek int    `json:"hoursPerWeek"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	AuthID       string `json:"auth_id"`
}

type Project struct {
	ID            uint   `json:"id"`
	ProjectName   string `json:"projectName"`
	Client        string `json:"client"`
	LogoClient    string `json:"logoClient"`
	Revenue       int    `json:"revenue"`
	FrontendHours int    `json:"frontendHours"`
	BackendHours  int    `json:"backendHours"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	AuthID        string `json:"auth_id"`
}

var db *gorm.DB
var err error

func main() {
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	DB_URL := "host=localhost user=bsinkule dbname=sp-go sslmode=disable"
	dbURL := os.Getenv(DB_URL)
	db, err = gorm.Open(
		"postgres", dbURL)
	// "host=localhost user=bsinkule dbname=sp-go sslmode=disable")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Engineer{})
	db.AutoMigrate(&Project{})

	router.HandleFunc("/engineers/", GetEngineers).Methods("GET")
	router.HandleFunc("/engineers/{id}", GetEngineer).Methods("GET")
	router.HandleFunc("/engineers/", CreateEngineer).Methods("POST")
	router.HandleFunc("/engineers/{id}", DeleteEngineer).Methods("DELETE")
	router.HandleFunc("/engineers/{id}", UpdateEngineer).Methods("PUT")

	router.HandleFunc("/projects/", GetProjects).Methods("GET")
	router.HandleFunc("/projects/{id}", GetProject).Methods("GET")
	router.HandleFunc("/projects/", CreateProject).Methods("POST")
	router.HandleFunc("/projects/{id}", DeleteProject).Methods("DELETE")
	router.HandleFunc("/projects/{id}", UpdateProject).Methods("PUT")

	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(router)))
	// log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(router)))
}

func GetEngineers(w http.ResponseWriter, r *http.Request) {
	var engineers []Engineer
	db.Order("id desc").Find(&engineers)
	json.NewEncoder(w).Encode(&engineers)
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	var projects []Project
	db.Order("id desc").Find(&projects)
	json.NewEncoder(w).Encode(&projects)
}

func GetEngineer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var engineer Engineer
	db.First(&engineer, params["id"])
	json.NewEncoder(w).Encode(&engineer)
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var project Project
	db.First(&project, params["id"])
	json.NewEncoder(w).Encode(&project)
}

func UpdateEngineer(w http.ResponseWriter, r *http.Request) {
	// if HandleCORS(w, r) {
	// 	return
	// }

	params := mux.Vars(r)
	var engineer Engineer
	db.First(&engineer, params["id"])
	json.NewDecoder(r.Body).Decode(&engineer)
	db.Save(&engineer)
	json.NewEncoder(w).Encode(&engineer)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var project Project
	db.First(&project, params["id"])
	json.NewDecoder(r.Body).Decode(&project)
	db.Save(&project)
	json.NewEncoder(w).Encode(&project)
}

func CreateEngineer(w http.ResponseWriter, r *http.Request) {

	// b, err := ioutil.ReadAll(r.Body)
	// log.Println("Request Body:", string(b))
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// reader := bytes.NewReader(b)

	var engineer Engineer
	json.NewDecoder(r.Body).Decode(&engineer)
	db.Create(&engineer)
	json.NewEncoder(w).Encode(&engineer)

	// log.Printf("engineers struct %+v \n", engineer)

	// var engineer Engineer
	// err = json.Unmarshal(b, &engineer)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// db.Create(&engineer)

	// output, err := json.Marshal(engineer)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	// w.Header().Set("content-type", "application/json")
	// w.Write(output)

}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var project Project
	json.NewDecoder(r.Body).Decode(&project)
	db.Create(&project)
	json.NewEncoder(w).Encode(&project)
}

func DeleteEngineer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var engineer Engineer
	db.First(&engineer, params["id"])
	db.Delete(&engineer)

	var engineers []Engineer
	db.Find(&engineers)
	json.NewEncoder(w).Encode(&engineers)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var project Project
	db.First(&project, params["id"])
	db.Delete(&project)

	var projects []Project
	db.Find(&projects)
	json.NewEncoder(w).Encode(&projects)
}

// // HandleCORS adds CORS headers if a CORS request was made
// // Read more here: http://www.html5rocks.com/en/tutorials/cors/
// func HandleCORS(w http.ResponseWriter, r *http.Request) bool {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "PUT")
// 	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization")
// 	w.Header().Set("Access-Control-Max-Age", "300") // In seconds - 5 mins

// 	// Short-circuit pre-flight OPTIONS requests with a 200 success response
// 	method := r.Header.Get("Access-Control-Request-Method")
// 	if method != "" && r.Method == "OPTIONS" {
// 		w.WriteHeader(200)
// 		return true
// 	}
// 	// Do not short-circuit the response
// 	return false
// }
