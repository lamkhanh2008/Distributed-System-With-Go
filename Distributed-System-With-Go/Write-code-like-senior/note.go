package main

//Go is pass by value: truyền 1 giá trị vào hàm thì hàm đó sẽ copy giá trị đó
//Use struct: chú ý lỗi nil pointer, có nghĩa là hàm struct con trỏ được gọi đến rỗng. Nhưng hàm đó lại đang chạy hàm implement nó

type Worker interface {
	Run() error
	Wait()
}

type implement struct {
}
