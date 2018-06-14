package mail4go

import (
	"fmt"
	"testing"
)

func Test_SendEmail(t *testing.T) {
	var config = &MailConfig{}
	config.username = "developer_mail@163.com"
	config.host = "smtp.163.com"
	config.password = "rkrntactzdinzcjk"
	config.port = "25"

	var e = NewHtmlMessage("title", "<a href='http://www.google.com'>Google</a>")
	e.From = "From<developer_mail@163.com>"
	e.To = []string{"917996695@qq.com"}

	fmt.Println(SendMail(config, e))
}
