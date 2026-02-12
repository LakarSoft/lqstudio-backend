package email

import (
	"fmt"
	"lqstudio-backend/internal/models"
	"strings"

	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

// Client handles email sending via Resend service
type Client struct {
	client  *resend.Client
	from    string
	adminTo string
	logger  *zap.Logger
}

// NewClient creates a new email client
func NewClient(apiKey string, from string, adminTo string, logger *zap.Logger) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("email API key is required")
	}
	if from == "" {
		return nil, fmt.Errorf("email from address is required")
	}

	client := resend.NewClient(apiKey)

	return &Client{
		client:  client,
		from:    from,
		adminTo: adminTo,
		logger:  logger,
	}, nil
}

// SendBookingConfirmation sends a confirmation email to the customer
func (c *Client) SendBookingConfirmation(to string, booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) error {
	subject := fmt.Sprintf("Booking Confirmation - %s", booking.ID)
	htmlBody := c.buildCustomerConfirmationHTML(booking, packageName, slots, addons)

	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		c.logger.Error("Failed to send booking confirmation email",
			zap.String("booking_id", booking.ID),
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send booking confirmation: %w", err)
	}

	c.logger.Info("Booking confirmation email sent successfully",
		zap.String("booking_id", booking.ID),
		zap.String("to", to),
		zap.String("email_id", sent.Id),
	)

	return nil
}

// SendAdminNotification sends a booking notification to admin
func (c *Client) SendAdminNotification(booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) error {
	if c.adminTo == "" {
		c.logger.Warn("Admin email not configured, skipping admin notification",
			zap.String("booking_id", booking.ID),
		)
		return nil
	}

	subject := fmt.Sprintf("New Booking Received - %s", booking.ID)
	htmlBody := c.buildAdminNotificationHTML(booking, packageName, slots, addons)

	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{c.adminTo},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		c.logger.Error("Failed to send admin notification email",
			zap.String("booking_id", booking.ID),
			zap.String("to", c.adminTo),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send admin notification: %w", err)
	}

	c.logger.Info("Admin notification email sent successfully",
		zap.String("booking_id", booking.ID),
		zap.String("to", c.adminTo),
		zap.String("email_id", sent.Id),
	)

	return nil
}

// SlotInfo represents slot information for email templates
type SlotInfo struct {
	ThemeName string
	Date      string
	Time      string
}

// AddonInfo represents addon information for email templates
type AddonInfo struct {
	Name     string
	Quantity int
	Price    string
}

// buildCustomerConfirmationHTML builds the HTML email for customer confirmation
func (c *Client) buildCustomerConfirmationHTML(booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) string {
	slotList := c.buildSlotListHTML(slots)
	addonsList := c.buildAddonsListHTML(addons)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        h1 { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; }
        h2 { color: #34495e; margin-top: 20px; }
        h3 { color: #7f8c8d; margin-top: 15px; }
        .detail { margin: 10px 0; }
        .label { font-weight: bold; color: #555; }
        .value { color: #333; }
        .slot-item, .addon-item { background: #f8f9fa; padding: 10px; margin: 5px 0; border-left: 4px solid #3498db; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; color: #7f8c8d; font-size: 0.9em; }
        .total { font-size: 1.2em; font-weight: bold; color: #27ae60; margin-top: 15px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Booking Confirmation - LQ Studio Photography</h1>

        <p>Dear %s,</p>
        <p>Thank you for booking with LQ Studio Photography! Your booking has been confirmed and is pending payment verification.</p>

        <h2>Booking Details</h2>
        <div class="detail">
            <span class="label">Booking ID:</span>
            <span class="value">%s</span>
        </div>
        <div class="detail">
            <span class="label">Package:</span>
            <span class="value">%s</span>
        </div>
        <div class="detail">
            <span class="label">Status:</span>
            <span class="value">%s</span>
        </div>

        <h3>Scheduled Sessions</h3>
        %s

        %s

        <div class="total">
            <span class="label">Total Amount:</span>
            <span class="value">RM %s</span>
        </div>

        <h3>Payment Instructions</h3>
        <p>Please complete your payment and upload the payment screenshot through our booking system. Your booking will be confirmed once payment is verified.</p>

        <div class="footer">
            <p><strong>Note:</strong> If you have any questions about your booking, please don't hesitate to contact us.</p>
            <p>Thank you for choosing LQ Studio Photography!</p>
        </div>
    </div>
</body>
</html>
`,
		booking.CustomerName,
		booking.ID,
		packageName,
		string(booking.Status),
		slotList,
		addonsList,
		booking.TotalAmount.StringFixed(2),
	)

	return html
}

// buildAdminNotificationHTML builds the HTML email for admin notification
func (c *Client) buildAdminNotificationHTML(booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) string {
	slotList := c.buildSlotListHTML(slots)
	addonsList := c.buildAddonsListHTML(addons)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        h1 { color: #c0392b; border-bottom: 3px solid #e74c3c; padding-bottom: 10px; }
        h2 { color: #34495e; margin-top: 20px; }
        h3 { color: #7f8c8d; margin-top: 15px; }
        .detail { margin: 10px 0; }
        .label { font-weight: bold; color: #555; }
        .value { color: #333; }
        .slot-item, .addon-item { background: #f8f9fa; padding: 10px; margin: 5px 0; border-left: 4px solid #e74c3c; }
        .action-button { display: inline-block; padding: 12px 24px; background: #3498db; color: white; text-decoration: none; border-radius: 4px; margin-top: 15px; }
        .total { font-size: 1.2em; font-weight: bold; color: #27ae60; margin-top: 15px; }
        .notes { background: #fff9e6; padding: 10px; margin-top: 10px; border-left: 4px solid #f39c12; }
    </style>
</head>
<body>
    <div class="container">
        <h1>New Booking Received</h1>

        <h2>Booking Information</h2>
        <div class="detail">
            <span class="label">Booking ID:</span>
            <span class="value">%s</span>
        </div>
        <div class="detail">
            <span class="label">Status:</span>
            <span class="value">%s</span>
        </div>

        <h2>Customer Details</h2>
        <div class="detail">
            <span class="label">Name:</span>
            <span class="value">%s</span>
        </div>
        <div class="detail">
            <span class="label">Email:</span>
            <span class="value">%s</span>
        </div>
        <div class="detail">
            <span class="label">Phone:</span>
            <span class="value">%s</span>
        </div>

        <h2>Package & Sessions</h2>
        <div class="detail">
            <span class="label">Package:</span>
            <span class="value">%s</span>
        </div>

        <h3>Scheduled Sessions</h3>
        %s

        %s

        <div class="total">
            <span class="label">Total Amount:</span>
            <span class="value">RM %s</span>
        </div>

        %s

        <a href="http://localhost:8080/api/admin/bookings/%s" class="action-button">View Booking Details</a>
    </div>
</body>
</html>
`,
		booking.ID,
		string(booking.Status),
		booking.CustomerName,
		booking.CustomerEmail,
		booking.CustomerPhone,
		packageName,
		slotList,
		addonsList,
		booking.TotalAmount.StringFixed(2),
		c.buildCustomerNotesHTML(booking.CustomerNotes),
		booking.ID,
	)

	return html
}

// buildSlotListHTML builds the HTML for the list of slots
func (c *Client) buildSlotListHTML(slots []SlotInfo) string {
	if len(slots) == 0 {
		return "<p>No sessions scheduled.</p>"
	}

	var sb strings.Builder
	for i, slot := range slots {
		sb.WriteString(fmt.Sprintf(`
        <div class="slot-item">
            <strong>Session %d:</strong> %s<br>
            <strong>Theme:</strong> %s<br>
            <strong>Time:</strong> %s
        </div>
`,
			i+1,
			slot.Date,
			slot.ThemeName,
			slot.Time,
		))
	}

	return sb.String()
}

// buildAddonsListHTML builds the HTML for the list of addons
func (c *Client) buildAddonsListHTML(addons []AddonInfo) string {
	if len(addons) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("<h3>Selected Add-ons</h3>\n")

	for _, addon := range addons {
		sb.WriteString(fmt.Sprintf(`
        <div class="addon-item">
            <strong>%s</strong><br>
            Quantity: %d | Price: RM %s
        </div>
`,
			addon.Name,
			addon.Quantity,
			addon.Price,
		))
	}

	return sb.String()
}

// buildCustomerNotesHTML builds the HTML for customer notes
func (c *Client) buildCustomerNotesHTML(notes string) string {
	if notes == "" {
		return ""
	}

	return fmt.Sprintf(`
        <h3>Customer Notes</h3>
        <div class="notes">
            %s
        </div>
`,
		notes,
	)
}
