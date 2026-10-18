package mailer

import "time"

const (
	SecurityStartTLS = "1"
	SecuritySSLTLS   = "2"
	SecurityNone     = "3"
)

type Config struct {
	Host        string
	Port        int
	Security    string
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
	ReplyTo     string
	Timeout     time.Duration
}
