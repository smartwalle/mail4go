package main

import (
	"fmt"
	"github.com/smartwalle/mail4go"
)

// gmail 需要先访问  https://www.google.com/settings/security/lesssecureapps，允许安全性较低的应用

func main() {
	// gmail
	var config = mail4go.NewMailConfig("username@gmail.com", "password", "smtp.gmail.com", "587")

	var e = mail4go.NewHTMLMessage("Title", "<a href='http://www.google.com'>Hello Google</a>")
	e.From = "From<smartwalle@gmail.com>"
	e.To = []string{"917996695@qq.com"}

	fmt.Println(mail4go.SendWithConfig(config, e))
}
