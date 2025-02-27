package utils

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)



func SendEmail(receiver string, table_name string, reservation_datetime string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("CONFIG_AUTH_EMAIL"))
	m.SetHeader("To", receiver)
	log.Printf(receiver, table_name, reservation_datetime)
	m.SetHeader("Subject", "Welcome To WeReserve")

	htmlBody :=  `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Reservation Details</title>
	</head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0;">
		<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border: 1px solid #dddddd; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">
			<!-- Header -->
			<div style="background-color: #007bff; color: #ffffff; text-align: center; padding: 20px; font-size: 24px; font-weight: bold;">
				Reservation Details
			</div>

			<!-- Content -->
			<div style="padding: 20px;">
				<p style="margin: 0 0 10px;">Hello,</p>
				<p style="margin: 0 0 20px;">Here are the details of your reservation:</p>

				<!-- Table -->
				<table style="width: 100%; border-collapse: collapse; margin-top: 20px;">
					<thead>
						<tr>
							<th style="border: 1px solid #dddddd; padding: 10px; text-align: left; background-color: #f9f9f9; font-weight: bold;">Table Name</th>
							<th style="border: 1px solid #dddddd; padding: 10px; text-align: left; background-color: #f9f9f9; font-weight: bold;">Reservation Date & Time</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td style="border: 1px solid #dddddd; padding: 10px;">` + table_name + `</td>
							<td style="border: 1px solid #dddddd; padding: 10px;">` + reservation_datetime + `</td>
						</tr>
					</tbody>
				</table>
			</div>

			<!-- Footer -->
			<div style="text-align: center; padding: 15px; background-color: #f9f9f9; font-size: 12px; color: #666666;">
				This is an automated email. Please do not reply directly to this message.
			</div>
		</div>
	</body>
	</html>
	`

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(
		viper.GetString("CONFIG_SMTP_HOST"),
		viper.GetInt("CONFIG_SMTP_PORT"),
		viper.GetString("CONFIG_AUTH_EMAIL"),
		viper.GetString("CONFIG_AUTH_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)

	}

	return nil
}

