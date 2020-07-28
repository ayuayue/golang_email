package mailers

import "fmt"

// GerReceiveMail 获取接收邮件
func GerReceiveMail(host, port, name, password string) []*MailItem {
	return GetFolderMail(fmt.Sprintf("%s:%s", host, port), name, password, "INBOX")

}

// GetReceiveMailMessage 获取邮件信息
func GetReceiveMailMessage(host, port, name, password string, id uint32) *MailItem {
	return GetMessage(fmt.Sprintf("%s:%s", host, port), name, password, "INBOX", id)
}
