package mailers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// MailItem item
type MailItem struct {
	Subject string
	Fid     string
	ID      uint32
	From    string
	To      string
	Body    string
	Date    string
}

// MailPageList list
type MailPageList struct {
	MailItems []*MailItem
}

// CheckEmailPasseword 校验邮箱密码
func CheckEmailPasseword(server string, email string, password string) bool {
	if !strings.Contains(server, ":") {
		return false
	}
	var c *client.Client
	serverSlice := strings.Split(server, ":")
	port, _ := strconv.Atoi(serverSlice[1])
	if port != 993 && port != 143 {
		return false
	}
	// 登录
	c = connect(server, email, password)
	// 退出
	defer c.Logout()
	if c == nil {
		return false
	}
	return true

}

// connect 连接
func connect(server, email, password string) *client.Client {
	var c *client.Client
	var err error
	serverSlice := strings.Split(server, ":")
	uri := serverSlice[0]
	port, _ := strconv.Atoi(serverSlice[1])
	if port != 993 && port != 143 {
		return nil
	}
	if port == 993 {
		c, err = client.DialTLS(fmt.Sprintf("%s:%d", uri, port), nil)
	} else {
		c, err = client.Dial(fmt.Sprintf("%s:%d", uri, port))
	}
	if err != nil {
		return nil
	}

	// 登录
	if err := c.Login(email, password); err != nil {
		return nil
	}

	return c
}
func Connect(server, email, password string) bool {
	var c *client.Client
	var err error
	serverSlice := strings.Split(server, ":")
	uri := serverSlice[0]
	port, _ := strconv.Atoi(serverSlice[1])
	if port != 993 && port != 143 {
		return false
	}
	if port == 993 {
		c, err = client.DialTLS(fmt.Sprintf("%s:%d", uri, port), nil)
	} else {
		c, err = client.Dial(fmt.Sprintf("%s:%d", uri, port))
	}
	if err != nil {
		return false
	}

	// 登录
	if err := c.Login(email, password); err != nil {
		return false
	}

	return true
}

// GetMailNum 获取邮件总数
func GetMailNum(server, email, password string) map[string]int {
	var c *client.Client
	c = connect(server, email, password)
	// 退出
	defer c.Logout()
	if c == nil {
		return nil
	}
	//列出邮箱
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()
	//存储邮件夹
	var folders = make(map[string]int)
	for m := range mailboxes {
		folders[m.Name] = 0
	}
	for k := range folders {
		mbox, _ := c.Select(k, true)
		if mbox != nil {
			folders[k] = int(mbox.Messages)
		}
	}
	return folders
}

// GetFolders 获取邮件夹
func GetFolders(server, email, password, folder string) map[string]int {
	var c *client.Client
	c = connect(server, email, password)
	// 退出
	defer c.Logout()
	if c == nil {
		return nil
	}
	//类出邮箱
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()
	//存储邮件夹
	var folders = make(map[string]int)
	for m := range mailboxes {
		folders[m.Name] = 0
	}
	for k := range folders {
		if k == folder {
			mbox, _ := c.Select(k, true)
			if mbox != nil {
				folders[k] = int(mbox.Messages)
			}
			break
		}
	}
	return folders
}

// GetFolderMail 获取邮件夹邮件
func GetFolderMail(server, email, password, folder string) []*MailItem {
	var c *client.Client
	c = connect(server, email, password)
	// 退出
	defer c.Logout()
	if c == nil {
		return nil
	}
	mbox, _ := c.Select(folder, true)
	// fmt.Println(mbox.Messages)
	to := mbox.Messages
	from := uint32(1)
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	fetchItem := imap.FetchItem(imap.FetchEnvelope)
	items := make([]imap.FetchItem, 0)
	items = append(items, fetchItem)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()
	var mailPagelist = new(MailPageList)
	dec := GetDecoder()

	for msg := range messages {
		// log.Println(msg.Envelope.Date)
		ret, err := dec.Decode(msg.Envelope.Subject)
		if err != nil {
			ret, _ = dec.DecodeHeader(msg.Envelope.Subject)
		}
		var mailitem = new(MailItem)
		mailitem.Subject = ret
		mailitem.ID = msg.SeqNum
		mailitem.Fid = folder
		mailitem.Date = msg.Envelope.Date.String()

		from := ""
		for _, s := range msg.Envelope.Sender {
			from += s.Address()
		}
		mailitem.From = from
		// fmt.Println(mailitem)
		mailPagelist.MailItems = append(mailPagelist.MailItems, mailitem)

	}
	return mailPagelist.MailItems
}

// GetMessage 获取消息
func GetMessage(server, email, password, folder string, id uint32) *MailItem {
	var c *client.Client
	c = connect(server, email, password)
	// 退出
	defer c.Logout()
	if c == nil {
		return nil
	}
	// Select INBOX
	mbox, err := c.Select(folder, false)
	if err != nil {
		log.Fatal(err)
	}
	// Get lastest message by id
	if mbox.Messages == 0 {
		log.Fatal("no message")
	}
	seqset := new(imap.SeqSet)
	seqset.AddNum(id)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Fatal(err)
		}
	}()
	msg := <-messages
	if msg == nil {
		log.Fatal("no message")
	}
	r := msg.GetBody(section)

	if r == nil {
		log.Fatal("no body")
	}
	var mailitem = new(MailItem)

	//Create a new mail reader
	mr, _ := mail.CreateReader(r)

	//Print some info about this message
	header := mr.Header
	date, _ := header.Date()

	mailitem.Date = date.String()

	var f string
	dec := GetDecoder()

	if from, err := header.AddressList("From"); err != nil {
		for _, address := range from {
			fromStr := address.String()
			temp, _ := dec.DecodeHeader(fromStr)
			f += " " + temp
		}
	}
	mailitem.From = f
	// log.Println("From: ", mailitem.From)

	var t string
	if to, err := header.AddressList("To"); err != nil {
		// log.Println("To: ", to)
		for _, address := range to {
			toStr := address.String()
			temp, _ := dec.DecodeHeader(toStr)
			t += " " + temp
		}
	}
	mailitem.To = t
	subject, _ := header.Subject()
	s, err := dec.Decode(subject)
	if err != nil {
		s, _ = dec.DecodeHeader(subject)
	}
	// log.Println("Subject: ", s)
	mailitem.Subject = s
	var bodyMap = make(map[string]string)
	bodyMap["text/plain"] = ""
	bodyMap["text/html"] = ""

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {

		}
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text

			b, _ := ioutil.ReadAll(p.Body)
			ct := p.Header.Get("Content-Type")
			if strings.Contains(ct, "text/plain") {
				bodyMap["text/plain"] += Encoding(string(b), ct)
			} else {
				bodyMap["text/html"] += Encoding(string(b), ct)
			}

		case *mail.AttachmentHeader:
			// has attachment

			filename, _ := h.Filename()
			log.Println("Attachment: ", filename)
		}
	}
	if bodyMap["text/html"] != "" {
		mailitem.Body = bodyMap["text/html"]
	} else {
		mailitem.Body = bodyMap["text/plain"]
	}
	return mailitem
}

// GetDecoder 编码处理
func GetDecoder() *mime.WordDecoder {
	dec := new(mime.WordDecoder)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		charset = strings.ToLower(charset)
		switch charset {
		case "gb2312":
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}
			//ret:=bytes.NewReader(content)
			//ret:=transform.NewReader(bytes.NewReader(content), simplifiedchinese.HZGB2312.NewEncoder())

			utf8str := ConvertToStr(string(content), "gbk", "utf-8")
			t := bytes.NewReader([]byte(utf8str))
			//ret:=utf8.DecodeRune(t)
			//log.Println(ret)
			return t, nil
		case "gbk":
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}
			//ret:=bytes.NewReader(content)
			//ret:=transform.NewReader(bytes.NewReader(content), simplifiedchinese.HZGB2312.NewEncoder())

			utf8str := ConvertToStr(string(content), "gbk", "utf-8")
			t := bytes.NewReader([]byte(utf8str))
			//ret:=utf8.DecodeRune(t)
			//log.Println(ret)
			return t, nil
		case "gb18030":
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}
			//ret:=bytes.NewReader(content)
			//ret:=transform.NewReader(bytes.NewReader(content), simplifiedchinese.HZGB2312.NewEncoder())

			utf8str := ConvertToStr(string(content), "gbk", "utf-8")
			t := bytes.NewReader([]byte(utf8str))
			//ret:=utf8.DecodeRune(t)
			//log.Println(ret)
			return t, nil
		default:
			return nil, fmt.Errorf("unhandle charset:%s", charset)

		}
	}
	return dec
}

// ConvertToStr 任意编码转特定编码
func ConvertToStr(src string, srcCode string, tagCode string) string {
	result := mahonia.NewDecoder(srcCode).ConvertString(src)
	//srcCoder := mahonia.NewDecoder(srcCode)
	//srcResult := srcCoder.ConvertString(src)
	//tagCoder := mahonia.NewDecoder(tagCode)
	//_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	//result := string(cdata)
	return result
}

// func ConvertToStr(src string, srcCode string, tagCode string) string {
// 	// 原编码
// 	srcCoder := mahonia.NewDecoder(srcCode)
// 	// 原编码写入原字符
// 	srcResult := srcCoder.ConvertString(src)
// 	//目标编码
// 	tagCoder := mahonia.NewDecoder(tagCode)
// 	// 转为目标字节
// 	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
// 	result := string(cdata)
// 	return result
// }

// Encoding 转换编码
func Encoding(html string, ct string) string {
	e, name := DetermineEncoding(html)
	if name != "utf-8" {
		html = ConvertToStr(html, "gbk", "utf-8")
		e = unicode.UTF8
	}
	r := strings.NewReader(html)

	utf8Reader := transform.NewReader(r, e.NewDecoder())
	//将其他编码的reader转换为常用的utf8reader
	all, _ := ioutil.ReadAll(utf8Reader)
	return string(all)
}

// DetermineEncoding 确定编码
func DetermineEncoding(html string) (encoding.Encoding, string) {
	e, name, _ := charset.DetermineEncoding([]byte(html), "")
	return e, name
}
