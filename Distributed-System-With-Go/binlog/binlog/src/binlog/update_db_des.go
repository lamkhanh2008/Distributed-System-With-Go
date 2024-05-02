package binlog

// import (
// 	"fmt"

// 	"github.com/siddontang/go-mysql/canal"
// )

// func action_in_db_des(e *canal.RowsEvent, user interface{}) {
// 	switch e.Action {
// 	case canal.UpdateAction:

// 		fmt.Printf("User %d name changed from %s to %s \n", user.Id, olduser.Name, olduser.Name)
// 	case canal.InsertAction:
// 		fmt.Printf("User %d is created with name %s\n", user.Id, user.Name)
// 	case canal.DeleteAction:
// 		fmt.Printf("User %d is deleted with name %s \n", user.Id, user.Name)
// 	default:
// 		fmt.Printf("Unknown Action")
// 	}
// }
