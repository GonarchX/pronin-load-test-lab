package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"load-test-lab/internal/domain/model"
	logger "load-test-lab/pkg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	db     *sql.DB
	errLog = logger.NewErrorFile("server")
)

const employeesDirectory = "employees"

func main() {
	initDB()

	if err := os.MkdirAll(employeesDirectory, 0755); err != nil {
		errLog.Fatalf("unable to create employees directory: %v", err)
	}

	defer db.Close()

	http.HandleFunc("GET /v0/employees", readHandler)
	http.HandleFunc("POST /v0/employees", writeHandler)
	http.HandleFunc("POST /v0/employees/broken", brokenWriteHandler)

	flushToFilePeriodInSeconds, _ := strconv.Atoi(os.Getenv("FLUSH_PERIOD_IN_SECONDS"))
	flushToFilePeriod := time.Duration(flushToFilePeriodInSeconds) * time.Second
	go func() {
		ticker := time.NewTicker(flushToFilePeriod)
		for {
			<-ticker.C
			saveEmployeesToFile()
			newEmployees = []model.Employee{}
		}
	}()

	errLog.Fatal(http.ListenAndServe(":9080", nil))
}

func initDB() {
	host := os.Getenv("DATABASE_HOST")
	port, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASS")
	dbname := os.Getenv("DATABASE_DBNAME")

	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)
	log.Println(psqlInfo) // just for debug

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		errLog.Fatalf("Error opening database: %q", err)
	}

	err = db.Ping()
	if err != nil {
		errLog.Fatalf("Error connecting to the database: %q", err)
	}

	migrate(err)
}

func migrate(err error) {
	_, err = db.Exec(
		`CREATE TABLE employees
				(
					id     SERIAL PRIMARY KEY,
					name   VARCHAR(100)   NOT NULL,
					salary DECIMAL(10, 2) NOT NULL
				);`)
	if err != nil {
		log.Printf("Failed to create table 'employees': %v", err)
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var employee model.Employee
	err := db.QueryRow("SELECT id, name, salary FROM employees WHERE id=$1", id).Scan(&employee.ID, &employee.Name, &employee.Salary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employee)
}

var newEmployees []model.Employee

func writeHandler(w http.ResponseWriter, r *http.Request) {
	var employee model.Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO employees (name, salary) VALUES ($1, $2) RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, employee.Name, employee.Salary).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employee.ID = id
	newEmployees = append(newEmployees, employee)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

func brokenWriteHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Internal error", http.StatusInternalServerError)
}

func saveEmployeesToFile() {
	if len(newEmployees) == 0 {
		return
	}

	file, err := json.MarshalIndent(newEmployees, "", " ")
	if err != nil {
		log.Println("Unable to marshal newEmployees:", err)
		return
	}

	filename := fmt.Sprintf("%s/employees_%d.json", employeesDirectory, time.Now().Unix())
	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		log.Println("Unable to write file:", err)
		return
	}
	log.Println("Employees saved to", filename)
}
