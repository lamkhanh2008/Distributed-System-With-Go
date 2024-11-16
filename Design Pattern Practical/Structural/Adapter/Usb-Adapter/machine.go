package main

import "fmt"

type Client struct {
}

type Computer interface {
	InsertIntoLightNingPort()
}

func (c *Client) InsertLightningConnectorIntoComputer(com Computer) {
	fmt.Println("Client inserts lightning connector into computer")
	com.InsertIntoLightNingPort()
}

type MacBook struct {
}

func (m *MacBook) InsertIntoLightNingPort() {
	fmt.Println("MacBook pro")
}

type Windows struct{}

func (w *Windows) insertIntoUSBPort() {
	fmt.Println("USB connector is plugged into windows machine.")
}

type WindowsAdapter struct {
	windowMachine *Windows
}

func (w *WindowsAdapter) InsertIntoLightNingPort() {
	fmt.Println("Windows pro")
	w.windowMachine.insertIntoUSBPort()
}

func main() {
	mac := MacBook{}
	win := WindowsAdapter{&Windows{}}
	client := Client{}
	client.InsertLightningConnectorIntoComputer(&mac)
	client.InsertLightningConnectorIntoComputer(&win)
}

//example: client cần sạc macbook dùng dây sạc lighning nhưng win thì k có -> cần có 1 usb để chuyển cổng win sang lightning

// tư tưởng: ví dụ có 1 lỗ tròn, 2 vật dạng 1 là hình trụ mặt là đáy hình vuông và hình trụ đáy là hình tròn, khi đưa hình tròn vào lỗ tròn thì sẽ đi qua còn hình vuioong thì k phù hợp -> cần có 1 adapter để biến đổi hình vuông sang hình tròn để đi qua.
