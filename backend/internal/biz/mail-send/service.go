package mailsend

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	sysconfig "fortunefinger/backend/internal/biz/sys-config"
	"fortunefinger/backend/internal/pkg/crypto/aesgcm"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/mailer"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

const mailConfigCode = "mail"

type Input struct {
	To          []string
	Subject     string
	TextBody    string
	HTMLBody    string
	HTMLBuilder func(ConnectionInfo) string
	Embeds      []mailer.Embed
}

type ConnectionInfo struct {
	Channel string
	Status  string
}

type mailConfig struct {
	Enabled     uint8  `json:"enabled"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Security    string `json:"security"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	SenderName  string `json:"senderName"`
	SenderEmail string `json:"senderEmail"`
	ReplyTo     string `json:"replyTo"`
	Timeout     int    `json:"timeout"`
}

type Service struct {
	svcCtx        *svc.ServiceContext
	configService *sysconfig.Service
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{
		svcCtx:        svcCtx,
		configService: sysconfig.NewService(svcCtx),
	}
}

func (s *Service) Send(ctx context.Context, input Input) error {
	storedConfig, err := s.configService.GetConfigByCode(ctx, mailConfigCode)
	if err != nil {
		s.svcCtx.Logger.Error("[MailSend] 查询邮件配置失败", zap.Error(err))
		return err
	}
	if storedConfig == nil {
		err := errorx.New(errorx.ServiceUnavailable, "邮件服务尚未配置")
		return err
	}

	var config mailConfig
	if err := json.Unmarshal([]byte(storedConfig.Config), &config); err != nil {
		s.svcCtx.Logger.Error("[MailSend] 解析邮件配置失败",
			zap.Uint64("config_id", storedConfig.ID),
			zap.Error(err),
		)
		return errorx.Wrap(errorx.InternalServer, "邮件配置无效", err)
	}
	if config.Enabled != 1 {
		err := errorx.New(errorx.ServiceUnavailable, "邮件服务未启用")
		s.svcCtx.Logger.Warn("[MailSend] 邮件服务未启用",
			zap.Uint64("config_id", storedConfig.ID),
			zap.Error(err),
		)
		return err
	}

	cipher, err := aesgcm.New()
	if err != nil {
		s.svcCtx.Logger.Error("[MailSend] 初始化邮件认证信息解密器失败", zap.Error(err))
		return errorx.Wrap(errorx.InternalServer, "邮件认证信息解密失败", err)
	}
	username, err := cipher.Decrypt(config.Username)
	if err != nil {
		s.svcCtx.Logger.Error("[MailSend] 解密SMTP Login失败",
			zap.Uint64("config_id", storedConfig.ID),
			zap.Error(err),
		)
		return errorx.Wrap(errorx.InternalServer, "邮件认证信息解密失败", err)
	}
	password, err := cipher.Decrypt(config.Password)
	if err != nil {
		s.svcCtx.Logger.Error("[MailSend] 解密SMTP Key失败",
			zap.Uint64("config_id", storedConfig.ID),
			zap.Error(err),
		)
		return errorx.Wrap(errorx.InternalServer, "邮件认证信息解密失败", err)
	}

	client, err := mailer.NewSMTPClient(mailer.Config{
		Host:        config.Host,
		Port:        config.Port,
		Security:    config.Security,
		Username:    string(username),
		Password:    string(password),
		SenderName:  config.SenderName,
		SenderEmail: config.SenderEmail,
		ReplyTo:     config.ReplyTo,
		Timeout:     time.Duration(config.Timeout) * time.Second,
	})
	if err != nil {
		s.svcCtx.Logger.Error("[MailSend] 创建SMTP客户端失败",
			zap.Uint64("config_id", storedConfig.ID),
			zap.Error(err),
		)
		return errorx.Wrap(errorx.InternalServer, "邮件配置无效", err)
	}
	htmlBody := input.HTMLBody
	if input.HTMLBuilder != nil {
		htmlBody = input.HTMLBuilder(ConnectionInfo{
			Channel: resolveChannel(config.Host),
			Status:  resolveStatus(config.Security),
		})
	}
	if err := client.Send(ctx, mailer.Message{
		To:       input.To,
		Subject:  input.Subject,
		TextBody: input.TextBody,
		HTMLBody: htmlBody,
		Embeds:   input.Embeds,
	}); err != nil {
		s.svcCtx.Logger.Error("[MailSend] 发送邮件失败",
			zap.Uint64("config_id", storedConfig.ID),
			zap.Int("recipient_count", len(input.To)),
			zap.Error(err),
		)
		return errorx.Wrap(errorx.ServiceUnavailable, "邮件发送失败，请稍后重试", err)
	}
	return nil
}

func resolveChannel(host string) string {
	value := strings.ToLower(host)
	switch {
	case strings.Contains(value, "brevo.com"):
		return "Brevo SMTP"
	case strings.Contains(value, "gmail.com"):
		return "Gmail SMTP"
	case strings.Contains(value, "office365.com"), strings.Contains(value, "outlook.com"):
		return "Outlook SMTP"
	default:
		return host
	}
}

func resolveStatus(security string) string {
	switch security {
	case mailer.SecurityStartTLS:
		return "畅通（STARTTLS）"
	case mailer.SecuritySSLTLS:
		return "畅通（SSL/TLS）"
	case mailer.SecurityNone:
		return "畅通（无加密）"
	default:
		return "畅通"
	}
}
