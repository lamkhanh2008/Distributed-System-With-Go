package binlog

import (
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
)

func BinlogListener(db *sql.DB) {
	c, err := getDefaultCanal()

	if err == nil {
		_, err := c.GetMasterPos()
		if err == nil {
			c.SetEventHandler(&binlogHandler{BinlogParser: BinlogParser{dbSource: db}})
			c.RunFrom(mysql.Position{Name: "mysql-bin.000014", Pos: 179812055})

			// c.RunFrom(coords)
		}
	}
}

func getDefaultCanal() (*canal.Canal, error) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", "10.8.12.195", 3306)
	cfg.User = "testuser"
	cfg.Password = "Mysohapass"
	cfg.Flavor = "mysql"
	cfg.Dump.ExecutionPath = ""
	// cfg.Dump.TableDB = "Test"
	// cfg.Dump.Tables = []string{"canal_test"}
	return canal.NewCanal(cfg)
}

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r, " ", string(debug.Stack()))
		}
	}()

	var n int
	var k int
	k = 1
	if e.Action == canal.UpdateAction {
		n = 1
		k = 2
	}
	// fmt.Print(e, "--xxx")
	fmt.Print("------------", "e.Rows: ", e.Rows, " len e.rows: ", len(e.Rows), " e.action: ", e.Action, " k: ", k)
	for i := n; i < len(e.Rows); i += k {
		key := e.Table.Schema + "." + e.Table.Name
		switch key {
		case User{}.SchemaName() + "." + User{}.TableName():
			user := User{}
			h.GetBinLogData(&user, e, i)
			switch e.Action {
			case canal.UpdateAction:
				olduser := User{}
				// tx, err := h.dbSource.Begin()

				// Example query
				update, err := h.dbSource.Prepare("UPDATE User_test SET name = ? WHERE Id = ?")
				if err != nil {
					// Rollback the transaction if an error occurs
					log.Fatal(err)
				}
				_, err = update.Exec(user.Name, user.Id)
				if err != nil {
					// Rollback the transaction if an error occurs
					log.Fatal(err)
				}
				h.GetBinLogData(&olduser, e, i-1)

				fmt.Printf("User %d name changed from %s to %s \n", user.Id, olduser.Name, olduser.Name)
			case canal.InsertAction:
				insertVote := "INSERT INTO User_test(id, name) VALUES(?,?)"
				_, err := h.dbSource.Exec(insertVote, user.Id, user.Name)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("User %d is created with name %s\n", user.Id, user.Name)
			case canal.DeleteAction:
				fmt.Printf("User %d is deleted with name %s \n", user.Id, user.Name)
			default:
				fmt.Printf("Unknown Action")
			}
		}
	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}
