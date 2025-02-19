package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Participant represents a database involved in the transaction
type Participant struct {
	Name string
	DB   *sql.DB
	Tx   *sql.Tx
}

// Prepare simulates the "prepare" phase in 2PC
func (p *Participant) Prepare(query string) error {
	var err error
	// Begin a transaction
	p.Tx, err = p.DB.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %v", p.Name, err)
	}

	// Simulate a write operation
	_, err = p.Tx.Exec(query)
	if err != nil {
		p.Tx.Rollback() // Rollback immediately on failure
		return fmt.Errorf("%s: failed to execute query: %v", p.Name, err)
	}

	fmt.Printf("%s: Prepared successfully\n", p.Name)
	return nil
}

// Commit simulates the "commit" phase
func (p *Participant) Commit() error {
	if err := p.Tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %v", p.Name, err)
	}
	fmt.Printf("%s: Committed successfully\n", p.Name)
	return nil
}

// Rollback simulates rolling back a transaction
func (p *Participant) Rollback() error {
	if err := p.Tx.Rollback(); err != nil {
		return fmt.Errorf("%s: failed to rollback transaction: %v", p.Name, err)
	}
	fmt.Printf("%s: Rolled back successfully\n", p.Name)
	return nil
}

// Coordinator manages the two-phase commit protocol
type Coordinator struct {
	Participants []*Participant
}

func (c *Coordinator) ExecuteTransaction(queries []string) {
	fmt.Println("Starting Two-Phase Commit...")

	// Phase 1: Prepare
	for i, participant := range c.Participants {
		if err := participant.Prepare(queries[i]); err != nil {
			fmt.Printf("Error in prepare phase: %v\n", err)
			// Rollback all participants if prepare fails
			c.Rollback()
			return
		}
	}

	// Phase 2: Commit
	fmt.Println("All participants prepared. Proceeding to commit...")
	for _, participant := range c.Participants {
		if err := participant.Commit(); err != nil {
			fmt.Printf("Error in commit phase: %v\n", err)
			return
		}
	}

	fmt.Println("Transaction completed successfully!")
}

func (c *Coordinator) Rollback() {
	fmt.Println("Rolling back all participants...")
	for _, participant := range c.Participants {
		participant.Rollback()
	}
}

func main() {
	// Database connections
	dbA, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbA")
	if err != nil {
		log.Fatalf("Failed to connect to Database A: %v", err)
	}
	defer dbA.Close()

	dbB, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/dbB")
	if err != nil {
		log.Fatalf("Failed to connect to Database B: %v", err)
	}
	defer dbB.Close()

	// Participants setup
	participantA := &Participant{Name: "Database-A", DB: dbA}
	participantB := &Participant{Name: "Database-B", DB: dbB}

	// Coordinator setup
	coordinator := &Coordinator{Participants: []*Participant{participantA, participantB}}

	// Queries for the transaction
	queries := []string{
		"INSERT INTO accounts (id, balance) VALUES (1, 100)",      // Database A query
		"INSERT INTO logs (event) VALUES ('Transaction Success')", // Database B query
	}

	// Execute the Two-Phase Commit
	coordinator.ExecuteTransaction(queries)
}
