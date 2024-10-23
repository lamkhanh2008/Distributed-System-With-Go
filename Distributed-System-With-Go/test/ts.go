package main

import "fmt"

type MyStruct struct {
	Date2  int32
	PeerID int
}

func Test(arr [2]int) {
	arr[0] = 1
}

func main() {
	arr := [2]int{4, 3}
	Test(arr)
	fmt.Println(arr)
	// Tạo danh sách các đối tượng MyStruct
	// list := []MyStruct{
	// 	{Date2: 20230826, PeerID: 5},
	// 	{Date2: 20230826, PeerID: 5},
	// 	{Date2: 20230825, PeerID: 2},
	// 	{Date2: 20230826, PeerID: 1},
	// 	{Date2: 20230825, PeerID: 3},
	// }

	// // Sắp xếp danh sách theo Date2 tăng dần, nếu Date2 bằng nhau thì theo PeerID tăng dần
	// sort.Slice(list, func(i, j int) bool {
	// 	if list[i].Date2 == list[j].Date2 {
	// 		return list[i].PeerID > list[j].PeerID
	// 	}
	// 	return list[i].Date2 > list[j].Date2
	// })

	// // In danh sách đã sắp xếp
	// for _, item := range list {
	// 	fmt.Printf("Date2: %d, PeerID: %d\n", item.Date2, item.PeerID)
	// }
}
