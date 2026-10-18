package html

import (
	_ "embed"
	"fmt"
)

const WordmarkContentID = "fortunefinger-wordmark"

//go:embed assets/fortune-finger-wordmark.png
var wordmark []byte

type TestData struct {
	SentAt  string
	Channel string
	Status  string
}

func TestTemplate(data TestData) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>FortuneFinger 测试邮件</title>
</head>
<body style="margin:0;padding:0;background:#ffffff;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI','Microsoft YaHei',Arial,sans-serif;color:#111111;">
  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background:#ffffff;">
    <tr>
      <td align="center" style="padding:32px 20px 52px;">
        <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="max-width:620px;">
          <tr>
            <td style="padding:0 0 24px;">
              <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                <tr>
                  <td align="left" valign="middle">
                    <img src="cid:%s" width="136" alt="FortuneFinger" style="display:block;width:136px;max-width:100%%;height:auto;border:0;">
                  </td>
                  <td align="right" valign="middle" style="color:#ef3f3f;font-size:11px;font-weight:700;letter-spacing:1.5px;white-space:nowrap;">SMTP VERIFIED</td>
                </tr>
              </table>
            </td>
          </tr>
          <tr><td style="height:2px;background:#1f1f1f;font-size:0;line-height:0;">&nbsp;</td></tr>
          <tr>
            <td style="padding:38px 0 36px;">
              <h1 style="margin:0 0 24px;color:#111111;font-size:27px;font-weight:700;line-height:1.35;">邮件服务配置成功</h1>
              <p style="margin:0 0 14px;color:#303030;font-size:15px;line-height:1.9;">FortuneFinger 已成功连接您的邮件服务器并完成消息投递测试。</p>
              <p style="margin:0;color:#303030;font-size:15px;line-height:1.9;">您的验证码与业务通知功能现已就绪，可以正常发送。这是一封由系统自动触发的确认邮件。</p>
            </td>
          </tr>
          <tr><td style="height:1px;background:#e8e8e8;font-size:0;line-height:0;">&nbsp;</td></tr>
          <tr>
            <td style="padding:25px 0;">
              <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                <tr>
                  <td style="padding:7px 0;color:#8c8c8c;font-size:13px;">发送时间</td>
                  <td align="right" style="padding:7px 0;color:#111111;font-size:13px;font-weight:600;">%s</td>
                </tr>
                <tr>
                  <td style="padding:7px 0;color:#8c8c8c;font-size:13px;">发信通道</td>
                  <td align="right" style="padding:7px 0;color:#111111;font-size:13px;font-weight:600;">%s</td>
                </tr>
                <tr>
                  <td style="padding:7px 0;color:#8c8c8c;font-size:13px;">连接状态</td>
                  <td align="right" style="padding:7px 0;color:#111111;font-size:13px;font-weight:600;"><span style="display:inline-block;width:6px;height:6px;margin-right:8px;background:#2f9e44;border-radius:50%%;vertical-align:2px;"></span>%s</td>
                </tr>
              </table>
            </td>
          </tr>
          <tr><td style="height:1px;background:#e8e8e8;font-size:0;line-height:0;">&nbsp;</td></tr>
          <tr>
            <td style="padding:28px 0 0;color:#999999;font-size:12px;line-height:1.8;">
              <p style="margin:0 0 6px;">这是一封系统自动发出的测试邮件，无需回复。</p>
              <p style="margin:0;">© 2026 FortuneFinger. All rights reserved.</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, WordmarkContentID, data.SentAt, data.Channel, data.Status)
}

func Wordmark() []byte {
	return wordmark
}
