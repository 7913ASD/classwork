package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func main() {

	db, err := sql.Open("mysql", "root:123456@tcp(18.219.49.188:3306)/fss?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("连接成功")

	defer db.Close()

	res, err := db.Exec(`insert into student(name,sex,age,course) values (?,?,?,?)`, "su", "man", "18", "math")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.LastInsertId())

	stu, err := queryStudentByName(db, "su")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*stu)

	stu, err = queryStudentByName(db, "su1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*stu)

}

type Student struct {
	Id     int
	Name   string
	Sex    string
	Age    string
	Course string
}

func queryStudentByName(db *sql.DB, name string) (*Student, error) {
	var s Student
	row := db.QueryRow(`select * from student where name=?`, name)
	err := row.Scan(&s.Id, &s.Name, &s.Sex, &s.Course, &s.Age)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "mysql异常")
	}
	return &s, nil
}
