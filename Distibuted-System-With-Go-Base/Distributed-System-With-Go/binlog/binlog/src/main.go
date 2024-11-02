package main

import (
	"binlog/binlog"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Main function
// func main() {

// 	a := []int{2, 5}

// 	var b reflect.Value = reflect.ValueOf(&a)

// 	b = b.Elem()
// 	fmt.Println("Slice :", b)
// 	fmt.Println("Slice :", a)
// 	s := reflect.Indirect(b)
// 	t := b.Type()
// 	m := t.NumField()

// 	fmt.Println(s, " ... ", t, ".....", m)

// 	//use of ValueOf method

// 	b = reflect.Append(b, reflect.ValueOf(80))
// 	fmt.Println("Slice after appending data:", a)

// }

// package main

func main() {
	db, err := sql.Open("mysql", "testuser:Mysohapass@tcp(10.8.12.195:3306)/Test")
	if err != nil {
		log.Fatal(err)
	}
	go binlog.BinlogListener(db)

	go func() {
		for {
			fmt.Print("Thank for watching")

			if err != nil {
				log.Fatal(err)
			}
			// defer db.Close()

			// Begin the transaction
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			// Example query
			_, err = tx.Exec("INSERT INTO User (name) VALUES (?)", "lamkczxczxzhanh")
			if err != nil {
				// Rollback the transaction if an error occurs
				tx.Rollback()
				log.Fatal(err)
			}

			// Another query
			// _, err = tx.Exec("UPDATE User SET name = ? WHERE id = ", "vvxxcvc")
			// if err != nil {
			// 	// Rollback the transaction if an error occurs
			// 	tx.Rollback()
			// 	log.Fatal(err)
			// }
			// _, err = tx.Exec("DELETE FROM User WHERE id = (?)", 12)
			// if err != nil {
			// 	// Rollback the transaction if an error occurs
			// 	tx.Rollback()
			// 	log.Fatal(err)
			// }
			// Commit the transaction if all queries succeed
			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Transaction committed successfully")
	}()
	time.Sleep(5 * time.Minute)
}
