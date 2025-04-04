package notifications

import (
	"fmt"
	"net/smtp"
	"os"
)



// SendSMS → ارسال پیامک با Ghasedak یا Kavenegar
func SendSMS(phoneNumber, message string) error {
	// اینجا باید کد ارسال پیامک از طریق API سرویس پیامکی رو بذاری
	fmt.Printf("📲 ارسال پیامک به %s: %s\n", phoneNumber, message)
	return nil
}

// SendEmail → ارسال ایمیل با SMTP
func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	// اطلاعات سرور SMTP (مثلاً Gmail)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// تنظیمات احراز هویت SMTP
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// قالب‌بندی پیام
	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	// ارسال ایمیل
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}

	fmt.Printf("📧 ایمیل ارسال شد به %s\n", to)
	return nil
}
