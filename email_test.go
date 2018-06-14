package mail4go

import (
	"fmt"
	"testing"
)

func Test_SendEmail(t *testing.T) {
	var config = &MailConfig{}
	config.username = "smartwalle@163.com"
	config.host = "smtp.163.com"
	config.password = "test123456"
	config.port = "25"

	var e = NewHtmlMessage("title", "<a href='http://www.google.com'>Google</a>")
	e.From = "From<smartwalle@163.com>"
	e.To = []string{"917996695@qq.com"}

	fmt.Println(SendMail(config, e))
}
