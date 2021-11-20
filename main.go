package main

//Using postgresSQL and Go lang APIs for basic CRUD operations.

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// update the postgres Database
const (
	user     = "postgres"
	host     = "localhost"
	dbname   = "postgres"
	password = "system"
	port     = 5432
)

var psqlInfo string = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func main() {

	r := gin.Default()
	r.GET("/employee", readEmployeeHandler)
	r.POST("/employee", createEmplyeeHandler)
	r.PUT("/employee", updateEmployeeHandler)
	r.DELETE("/employee/:id", deleteEmployeeHAndler)

	r.Run(":8000")
}

func readEmployeeHandler(c *gin.Context) {

	db, err := sql.Open("postgres", psqlInfo)
	fmt.Print("CONNECT")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	//Creating employee struct

	type employee struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		City string `json:"city"`
	}

	emp := employee{}

	id := c.Request.URL.Query().Get("id")

	log.Println(id)
	i, _ := strconv.Atoi(id)

	row := db.QueryRow("Select * FROM employee where id=$1", i)

	//error in row

	err = row.Scan(&emp.Id, &emp.Name, &emp.City)

	if err != nil {
		log.Println(err)
		c.JSON(500, "Row does not exist")
		return
	}

	c.JSON(200, emp)

}

func createEmplyeeHandler(c *gin.Context) {

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	// Getting data from POST request body
	decoder := json.NewDecoder(c.Request.Body)

	type body_struct struct {
		Name string
		City string
	}

	var one body_struct

	err = decoder.Decode(&one)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Print(one.Name)
	fmt.Print(one.City)

	//Adding row to the table
	rowadd, err := db.Exec("Insert into employee (name,city) VALUES ($1,$2)", one.Name, one.City)
	if err != nil {
		log.Println(err)
		c.JSON(500, "Error in addition in table!")
		return
	}

	output, err := rowadd.LastInsertId()
	if err != nil {
		log.Println(err)
		c.JSON(500, "Conversion")
		return
	}
	c.JSON(200, fmt.Sprintf("Added = %v ", output))
}

func updateEmployeeHandler(c *gin.Context) {

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	// Getting data from PUT request body
	decoder := json.NewDecoder(c.Request.Body)

	type body_struct struct {
		Id   int
		Name string
		City string
	}

	var one body_struct

	err = decoder.Decode(&one)
	if err != nil {
		log.Println(err)
		return
	}

	var query string
	var params = make([]interface{}, 3)

	//Updating existing row
	if one.Name == "" {
		query = "Update employee set city = $1 where id = $2"
		params = []interface{}{one.City, one.Id}
	} else if one.City == "" {
		query = "Update employee set name = $1 where id = $2"
		params = []interface{}{one.Name, one.Id}
	} else {
		query = "Update employee set name = $1, city= $2 where id = $3"
		params = []interface{}{one.Name, one.City, one.Id}
	}
	_, err = db.Exec(query, params...)
	if err != nil {
		log.Println(err)
		c.JSON(500, "Error in updation in table!")
		return
	}

	c.JSON(200, fmt.Sprintf("Updated = %v ", one.Id))

}

func deleteEmployeeHAndler(c *gin.Context) {

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	deleteId := c.Params.ByName("id")

	//Delete row
	_, err = db.Exec("Delete from employee where id = $1", deleteId)
	if err != nil {
		log.Println(err)
		c.JSON(500, "Error in deletion from table!")
		return
	}

	c.JSON(200, fmt.Sprintf("Deleted = %v ", deleteId))

}
