package main

import (
	"fmt"
	"github.com/smartwalle/mail4go"
)

func main() {
	// gmail
	var config = mail4go.NewMailConfig("xxxx@gmail.com", "xxx", "smtp.gmail.com", "587")

	var e = mail4go.NewHtmlMessage("title", "<a href='http://www.google.com'>Google</a>")
	e.From = "From<hoteldelins@gmail.com>"
	e.To = []string{"917996695@qq.com"}

	fmt.Println(mail4go.SendWithConfig(config, e))
}
