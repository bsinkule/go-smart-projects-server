package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

type Engineer struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	ImgURL       string  `json:"imgUrl"`
	Title        string  `json:"title"`
	Department   string  `json:"department"`
	HourlyWage   float32 `json:"hourlyWage"`
	HoursPerWeek float32 `json:"hoursPerWeek"`
	StartDate    string  `json:"startDate"`
	EndDate      string  `json:"endDate"`
	AuthID       string  `json:"auth_id"`
}

type Project struct {
	ID            uint    `json:"id"`
	ProjectName   string  `json:"projectName"`
	Client        string  `json:"client"`
	LogoClient    string  `json:"logoClient"`
	Revenue       float32 `json:"revenue"`
	FrontendHours float32 `json:"frontendHours"`
	BackendHours  float32 `json:"backendHours"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
	AuthID        string  `json:"auth_id"`
}

var db *gorm.DB
var err error

func main() {
	router := mux.NewRouter()

	db, err = gorm.Open(
		"postgres",
		// "host="+os.Getenv("localhost")+" user="+os.Getenv("bsinkule")+
		// 	" dbname="+os.Getenv("spittle")+" sslmode=disable")

		// os.ExpandEnv("host=${localhost} user=${bsinkule} dbname=${spittle} sslmode=disable"))
		"host=localhost user=bsinkule dbname=sp-go sslmode=disable")

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

	handler := cors.Default().Handler(router)

	// log.Fatal(http.ListenAndServe(":"+os.ExpandEnv("path=${5432}"), handler))
	log.Fatal(http.ListenAndServe(":5000", handler))
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
	var engineer Engineer
	json.NewDecoder(r.Body).Decode(&engineer)
	db.Create(&engineer)
	json.NewEncoder(w).Encode(&engineer)
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
