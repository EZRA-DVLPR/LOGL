package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	// "github.com/EZRA-DVLPR/GameList/internal/scraper"
)

type Salary struct {
	Basic, HRA, TA float64
}

type Employee struct {
	FirstName, LastName, Email string
	Age                        int
	MonthlySalary              []Salary
}

func CreateDB() {
	fmt.Println("Creating the DB")
	data := Employee{
		FirstName: "Mark",
		LastName:  "Jones",
		Email:     "mark@gmail.com",
		Age:       25,
		MonthlySalary: []Salary{
			{
				Basic: 15000.00,
				HRA:   5000.00,
				TA:    2000.00,
			},
			{
				Basic: 16000.00,
				HRA:   5000.00,
				TA:    2100.00,
			},
			{
				Basic: 17000.00,
				HRA:   5000.00,
				TA:    2200.00,
			},
		},
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = os.WriteFile("test.json", file, 0644)
}

func ReadDB() {
	content, err := os.ReadFile("test.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload Employee
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	log.Printf("FirstName %s\n", payload.FirstName)
	log.Printf("LastName %s\n", payload.LastName)
	log.Printf("Email %s\n", payload.Email)
	log.Printf("Age %d\n", payload.Age)
	log.Println(payload.MonthlySalary)
}
