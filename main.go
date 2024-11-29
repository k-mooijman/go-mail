package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	username := "mailtest@mooijman.info"
	password := "JdhvVTn7tq"
	//server := "mail.mooijman.info"
	server := "mail218.hostingdiscounter.nl"
	port := 993

	imapClient, err := connectToServer(username, password, server, port)
	if err != nil {
		log.Fatal(err)
	}
	defer imapClient.Logout()

	if err := fetchEmails(imapClient); err != nil {
		log.Fatal(err)
	}
	if err := mailBoxes(imapClient); err != nil {
		log.Fatal(err)
	}

	//to := "kasper@mooijman.info"
	//subject := "Test Email"
	//body := "This is a test email sent from a Go-based email client."

	//if err := sendEmail(username, password, server, 587, to, subject, body); err != nil {
	//	log.Fatal(err)
	//}

}

func connectToServer(username, password, server string, port int) (*client.Client, error) {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", server, port), nil)
	if err != nil {
		return nil, err
	}

	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	return c, nil
}

func fetchEmails(imapClient *client.Client) error {
	// Select the mailbox you want to read
	mailbox, err := imapClient.Select("INBOX", false)
	if err != nil {
		return err
	}

	//// Define the range of emails to fetch
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(1, mailbox.Messages)
	//section := &imap.BodySectionName{}
	//items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}
	items := []imap.FetchItem{imap.FetchItem("BODY.PEEK[]"), imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate}

	//Fetch the required message attributes
	messages := make(chan *imap.Message, 10)

	go func() {
		if err := imapClient.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {
		//fmt.Println("Subject:", msg)
		fmt.Println("Subject:", msg.Envelope.Subject)
		fmt.Println("from:", msg.Envelope.From)
		fmt.Println("ID:", msg.Envelope.MessageId)
		//fmt.Println("flags:", msg.Flags)
		fmt.Println("-------------------")
		for flag := range msg.Flags {
			flagString := imap.CanonicalFlag(msg.Flags[flag])
			fmt.Println(" flag :", flagString)
		}
		fmt.Println("-------------------")
	}

	return nil
}

func sendEmail(username, password, server string, port int, to, subject, body string) error {
	auth := smtp.PlainAuth("", username, password, server)

	msg := "From: " + username + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	err := smtp.SendMail(server+":"+strconv.Itoa(port), auth, username, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func mailBoxes(imapClient *client.Client) error {
	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- imapClient.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")

	var list []*imap.MailboxInfo
	for m := range mailboxes {
		list = append(list, m)
		log.Printf("* %s    -- %s \n", m.Name, m.Attributes)
	}

	if err := <-done; err != nil {
		log.Fatalf("Error listing mailboxes: %v", err)
	}

	return nil
}
