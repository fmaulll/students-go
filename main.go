package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	Id        int
	Name      string
	StudentId string
	Major     string
}

func dbConnect() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:Password1234@tcp(127.0.0.1:3306)/student-portal")
	if err != nil {
		panic(err.Error())
	}

	return db
}

func getStudents(context *gin.Context) {

	db := dbConnect()

	selDB, err := db.Query("SELECT * FROM students")
	if err != nil {
		panic(err.Error())
	}

	emp := Student{}
	res := []Student{}
	for selDB.Next() {
		var id int
		var name, studentId, major string
		err = selDB.Scan(&id, &name, &studentId, &major)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.StudentId = studentId
		emp.Major = major
		// jsonRes, err := json.Marshal(Student{Id: id, Name: name, StudentId: studentId, Major: major})
		// if err != nil {
		// 	return
		// }
		res = append(res, emp)

	}
	context.IndentedJSON(http.StatusOK, res)
	defer db.Close()
}

func addStudent(context *gin.Context) {
	var newStudent Student
	db := dbConnect()

	if err := context.BindJSON(&newStudent); err != nil {
		return
	}

	selDB, err := db.Query("INSERT INTO students (name, studentId, major) VALUES " + "(" + "'" + newStudent.Name + "'" + ", " + "'" + newStudent.StudentId + "'" + ", " + "'" + newStudent.Major + "'" + ")")
	if err != nil {
		panic(err.Error())
	}
	context.IndentedJSON(http.StatusCreated, newStudent)
	defer selDB.Close()
}

func deleteStudent(context *gin.Context) {
	id := context.Param("id")
	db := dbConnect()

	selDB, err := db.Query("DELETE FROM students WHERE (`id` = " + id + ")")
	if err != nil {
		panic(err.Error())
	}

	selDB.Close()
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student deleted!"})
}

func getStudent(context *gin.Context) {
	id := context.Param("id")
	db := dbConnect()

	selDB, err := db.Query("SELECT * FROM students WHERE (`id` = " + id + ")")

	if err != nil {
		panic(err.Error())
	}
	selDB.Close()
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student deleted!"})

}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/students", getStudents)
	router.POST("/students", addStudent)
	router.GET("/students/delete/:id", deleteStudent)
	router.GET("/students/:id", getStudent)
	router.Run("localhost:9000")
}
