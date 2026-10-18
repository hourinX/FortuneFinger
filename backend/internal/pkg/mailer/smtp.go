package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	mail "github.com/wneessen/go-mail"
)

type smtpClient struct {
	client *mail.Client
	config Config
}

func NewSMTPClient(config Config) (Client, error) {
	if config.Host == "" {
		return nil, errors.New("mailer: SMTP服务器地址不能为空")
	}
	if config.Port < 1 || config.Port > 65535 {
		return nil, errors.New("mailer: SMTP端口无效")
	}
	if config.Username == "" || config.Password == "" {
		return nil, errors.New("mailer: SMTP认证信息不能为空")
	}
	if config.SenderEmail == "" {
		return nil, errors.New("mailer: 发件人邮箱不能为空")
	}
	if config.Timeout <= 0 {
		return nil, errors.New("mailer: SMTP超时时间必须大于0")
	}

	options := []mail.Option{
		mail.WithPort(config.Port),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(config.Username),
		mail.WithPassword(config.Password),
		mail.WithTimeout(config.Timeout),
	}

	switch config.Security {
	case SecurityStartTLS:
		options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
	case SecuritySSLTLS:
		options = append(options, mail.WithSSL())
	case SecurityNone:
		options = append(options, mail.WithTLSPolicy(mail.NoTLS))
	default:
		return nil, errors.New("mailer: SMTP连接安全类型无效")
	}

	client, err := mail.NewClient(config.Host, options...)
	if err != nil {
		return nil, fmt.Errorf("mailer: 创建SMTP客户端失败: %w", err)
	}
	return &smtpClient{
		client: client,
		config: config,
	}, nil
}

func (c *smtpClient) Send(ctx context.Context, message Message) error {
	if len(message.To) == 0 {
		return errors.New("mailer: 收件人不能为空")
	}
	if message.Subject == "" {
		return errors.New("mailer: 邮件主题不能为空")
	}
	if message.TextBody == "" && message.HTMLBody == "" {
		return errors.New("mailer: 邮件正文不能为空")
	}

	msg := mail.NewMsg()
	if err := msg.FromFormat(c.config.SenderName, c.config.SenderEmail); err != nil {
		return fmt.Errorf("mailer: 设置发件人失败: %w", err)
	}
	if err := msg.To(message.To...); err != nil {
		return fmt.Errorf("mailer: 设置收件人失败: %w", err)
	}
	if c.config.ReplyTo != "" {
		if err := msg.ReplyTo(c.config.ReplyTo); err != nil {
			return fmt.Errorf("mailer: 设置回复邮箱失败: %w", err)
		}
	}

	msg.Subject(message.Subject)
	msg.SetDate()
	msg.SetMessageID()
	if message.TextBody != "" {
		msg.SetBodyString(mail.TypeTextPlain, message.TextBody)
		if message.HTMLBody != "" {
			msg.AddAlternativeString(mail.TypeTextHTML, message.HTMLBody)
		}
	} else {
		msg.SetBodyString(mail.TypeTextHTML, message.HTMLBody)
	}
	for _, embed := range message.Embeds {
		if err := msg.EmbedReader(embed.Name, bytes.NewReader(embed.Content), mail.WithFileContentID(embed.ContentID)); err != nil {
			return fmt.Errorf("mailer: 嵌入邮件资源失败: %w", err)
		}
	}

	if err := c.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("mailer: 发送邮件失败: %w", err)
	}
	return nil
}

func (c *smtpClient) Verify(ctx context.Context) error {
	if err := c.client.DialWithContext(ctx); err != nil {
		return fmt.Errorf("mailer: 验证SMTP连接失败: %w", err)
	}
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("mailer: 关闭SMTP连接失败: %w", err)
	}
	return nil
}
