package mailer

import "context"

type Client interface {
	Send(ctx context.Context, message Message) error
	Verify(ctx context.Context) error
}

type Message struct {
	To       []string
	Subject  string
	TextBody string
	HTMLBody string
	Embeds   []Embed
}

type Embed struct {
	Name      string
	ContentID string
	Content   []byte
}
