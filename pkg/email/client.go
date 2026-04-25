package email

import (
	"fmt"
	"html"
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
	subject := fmt.Sprintf("Booking Received - %s", c.bookingIDValue(booking))
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

	subject := fmt.Sprintf("New Booking Received - %s", c.bookingIDValue(booking))
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

// SendBookingApproval sends an approval email to the customer
func (c *Client) SendBookingApproval(to string, booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) error {
	subject := fmt.Sprintf("Booking Approved - %s", c.bookingIDValue(booking))
	htmlBody := c.buildBookingApprovalHTML(booking, packageName, slots, addons)

	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		c.logger.Error("Failed to send booking approval email",
			zap.String("booking_id", booking.ID),
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send booking approval: %w", err)
	}

	c.logger.Info("Booking approval email sent successfully",
		zap.String("booking_id", booking.ID),
		zap.String("to", to),
		zap.String("email_id", sent.Id),
	)

	return nil
}

// SendBookingRejection sends a rejection email to the customer
func (c *Client) SendBookingRejection(to string, booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) error {
	subject := fmt.Sprintf("Booking Status Update - %s", c.bookingIDValue(booking))
	htmlBody := c.buildBookingRejectionHTML(booking, packageName, slots, addons)

	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		c.logger.Error("Failed to send booking rejection email",
			zap.String("booking_id", booking.ID),
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send booking rejection: %w", err)
	}

	c.logger.Info("Booking rejection email sent successfully",
		zap.String("booking_id", booking.ID),
		zap.String("to", to),
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
	bookingSummary := c.buildBookingSummaryHTML(booking, packageName)
	customerDetails := c.buildCustomerDetailsHTML(booking, true)
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
        <h1>Booking Received - LQ Studio Photography</h1>

        <p>Dear %s,</p>
        <p>Thank you for booking with LQ Studio Photography. We have received your booking request and it is currently pending payment verification.</p>
        <p>Once your payment has been reviewed, we will update you on the next step for your booking.</p>

        <h2>Booking Summary</h2>
        %s

        <h2>Customer Details</h2>
        %s

        <h3>Scheduled Sessions</h3>
        %s

        %s

        <div class="total">
            <span class="label">Total Amount:</span>
            <span class="value">RM %s</span>
        </div>

        <h3>Payment Instructions</h3>
        <p>Please complete your payment and share your payment screenshot with LQ Studio for verification.</p>
        <p>Your booking will remain in pending status until payment verification is completed.</p>

        <div class="footer">
            <p><strong>Note:</strong> If you have any questions about your booking or payment, please contact LQ Studio directly.</p>
            <p>Thank you for choosing LQ Studio Photography.</p>
        </div>
    </div>
</body>
</html>
`,
		c.customerNameValue(booking),
		bookingSummary,
		customerDetails,
		slotList,
		addonsList,
		c.totalAmountValue(booking),
	)

	return html
}

// buildAdminNotificationHTML builds the HTML email for admin notification
func (c *Client) buildAdminNotificationHTML(booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) string {
	bookingSummary := c.buildBookingSummaryHTML(booking, packageName)
	customerDetails := c.buildCustomerDetailsHTML(booking, true)
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
        .total { font-size: 1.2em; font-weight: bold; color: #27ae60; margin-top: 15px; }
        .notes { background: #fff9e6; padding: 10px; margin-top: 10px; border-left: 4px solid #f39c12; }
    </style>
</head>
<body>
    <div class="container">
        <h1>New Booking Received</h1>

        <p>A new customer booking has been submitted and is awaiting follow-up.</p>

        <h2>Booking Summary</h2>
        %s

        <h2>Customer Details</h2>
        %s

        <h3>Scheduled Sessions</h3>
        %s

        %s

        <div class="total">
            <span class="label">Total Amount:</span>
            <span class="value">RM %s</span>
        </div>

        <div class="footer">
            <p>Please review the booking in the admin system and follow up with the customer as needed.</p>
        </div>
    </div>
</body>
</html>
`,
		bookingSummary,
		customerDetails,
		slotList,
		addonsList,
		c.totalAmountValue(booking),
	)

	return html
}

// buildBookingApprovalHTML builds the HTML email for booking approval notification
func (c *Client) buildBookingApprovalHTML(booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) string {
	bookingSummary := c.buildBookingSummaryHTML(booking, packageName)
	customerDetails := c.buildCustomerDetailsHTML(booking, true)
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
        h1 { color: #27ae60; border-bottom: 3px solid #2ecc71; padding-bottom: 10px; }
        h2 { color: #34495e; margin-top: 20px; }
        h3 { color: #7f8c8d; margin-top: 15px; }
        .detail { margin: 10px 0; }
        .label { font-weight: bold; color: #555; }
        .value { color: #333; }
        .slot-item, .addon-item { background: #f8f9fa; padding: 10px; margin: 5px 0; border-left: 4px solid #2ecc71; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; color: #7f8c8d; font-size: 0.9em; }
        .total { font-size: 1.2em; font-weight: bold; color: #27ae60; margin-top: 15px; }
        .success-message { background: #d4edda; border-left: 4px solid #28a745; padding: 15px; margin: 20px 0; color: #155724; }
    </style>
</head>
<body>
    <div class="container">
        <h1>✓ Booking Approved!</h1>

        <div class="success-message">
            <strong>Good news.</strong> Your booking has been approved and your session is now confirmed.
        </div>

        <p>Dear %s,</p>
        <p>Your payment has been verified, and your booking is now confirmed with LQ Studio Photography.</p>

        <h2>Booking Summary</h2>
        %s

        <h2>Customer Details</h2>
        %s

        <h3>Scheduled Sessions</h3>
        %s

        %s

        <div class="total">
            <span class="label">Total Amount:</span>
            <span class="value">RM %s</span>
        </div>

        <div class="footer">
            <p><strong>What's Next?</strong></p>
            <p>Please arrive at least 10 minutes before your scheduled session time and keep this email for reference.</p>
            <p>If you need help before your session, or if there are any changes to discuss, please contact LQ Studio as early as possible.</p>
            <p>We look forward to seeing you at the studio.</p>
        </div>
    </div>
</body>
</html>
`,
		c.customerNameValue(booking),
		bookingSummary,
		customerDetails,
		slotList,
		addonsList,
		c.totalAmountValue(booking),
	)

	return html
}

// buildBookingRejectionHTML builds the HTML email for booking rejection notification
func (c *Client) buildBookingRejectionHTML(booking *models.Booking, packageName string, slots []SlotInfo, addons []AddonInfo) string {
	bookingSummary := c.buildBookingSummaryHTML(booking, packageName)
	customerDetails := c.buildCustomerDetailsHTML(booking, true)
	slotList := c.buildSlotListHTML(slots)
	addonsList := c.buildAddonsListHTML(addons)

	// Extract admin notes if available
	adminNotes := "No specific reason provided."
	if booking.AdminNotes != "" {
		adminNotes = booking.AdminNotes
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        h1 { color: #e74c3c; border-bottom: 3px solid #c0392b; padding-bottom: 10px; }
        h2 { color: #34495e; margin-top: 20px; }
        h3 { color: #7f8c8d; margin-top: 15px; }
        .detail { margin: 10px 0; }
        .label { font-weight: bold; color: #555; }
        .value { color: #333; }
        .slot-item, .addon-item { background: #f8f9fa; padding: 10px; margin: 5px 0; border-left: 4px solid #e74c3c; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; color: #7f8c8d; font-size: 0.9em; }
        .total { font-size: 1.2em; font-weight: bold; color: #555; margin-top: 15px; }
        .warning-message { background: #f8d7da; border-left: 4px solid #dc3545; padding: 15px; margin: 20px 0; color: #721c24; }
        .info-box { background: #d1ecf1; border-left: 4px solid #17a2b8; padding: 15px; margin: 20px 0; color: #0c5460; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Booking Status Update</h1>

        <div class="warning-message">
            <strong>Booking Update:</strong> We are unable to approve this booking request at the moment.
        </div>

        <p>Dear %s,</p>
        <p>Thank you for your interest in LQ Studio Photography. After review, we are unable to proceed with this booking request in its current form.</p>

        <h2>Reason</h2>
        <p>%s</p>

        <h2>Booking Summary</h2>
        %s

        <h2>Customer Details</h2>
        %s

        <h3>Scheduled Sessions (Reference)</h3>
        %s

        %s

        <div class="total">
            <span class="label">Amount:</span>
            <span class="value">RM %s</span>
        </div>

        <div class="info-box">
            <h3>Next Step</h3>
            <p>If you have already made a payment, please contact LQ Studio directly so the team can assist you with the next arrangement, including any refund discussion if applicable.</p>
        </div>

        <div class="footer">
            <p>If you would like to discuss another date, package, or arrangement, please contact LQ Studio directly.</p>
            <p>We appreciate your understanding and hope to assist you again in the future.</p>
            <p>Best regards,<br>LQ Studio Photography</p>
        </div>
    </div>
</body>
</html>
`,
		c.customerNameValue(booking),
		c.escapeOrDefault(adminNotes, "No specific reason provided."),
		bookingSummary,
		customerDetails,
		slotList,
		addonsList,
		c.totalAmountValue(booking),
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
			c.escapeOrDefault(slot.Date, "Date to be confirmed"),
			c.escapeOrDefault(slot.ThemeName, "Theme to be confirmed"),
			c.escapeOrDefault(slot.Time, "Time to be confirmed"),
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
			c.escapeOrDefault(addon.Name, "Add-on"),
			c.nonZeroIntOrDefault(addon.Quantity, 1),
			c.escapeOrDefault(addon.Price, "0.00"),
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
		c.escapeOrDefault(notes, "No customer notes provided."),
	)
}

func (c *Client) buildBookingSummaryHTML(booking *models.Booking, packageName string) string {
	return fmt.Sprintf(`
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
        <div class="detail">
            <span class="label">Total Amount:</span>
            <span class="value">RM %s</span>
        </div>
`,
		c.bookingIDValue(booking),
		c.escapeOrDefault(packageName, "Package to be confirmed"),
		c.bookingStatusValue(booking),
		c.totalAmountValue(booking),
	)
}

func (c *Client) buildCustomerDetailsHTML(booking *models.Booking, includeNotes bool) string {
	if booking == nil {
		return `
        <div class="detail">
            <span class="value">Customer details are unavailable.</span>
        </div>
`
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
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
`,
		c.customerNameValue(booking),
		c.escapeOrDefault(booking.CustomerEmail, "Email not provided"),
		c.escapeOrDefault(booking.CustomerPhone, "Phone not provided"),
	))

	if includeNotes && strings.TrimSpace(booking.CustomerNotes) != "" {
		sb.WriteString(c.buildCustomerNotesHTML(booking.CustomerNotes))
	}

	return sb.String()
}

func (c *Client) bookingIDValue(booking *models.Booking) string {
	if booking == nil {
		return "unknown-booking"
	}
	return c.escapeOrDefault(booking.ID, "unknown-booking")
}

func (c *Client) customerNameValue(booking *models.Booking) string {
	if booking == nil {
		return "Customer"
	}
	return c.escapeOrDefault(booking.CustomerName, "Customer")
}

func (c *Client) bookingStatusValue(booking *models.Booking) string {
	if booking == nil {
		return "UNKNOWN"
	}
	return c.escapeOrDefault(string(booking.Status), "UNKNOWN")
}

func (c *Client) totalAmountValue(booking *models.Booking) string {
	if booking == nil {
		return "0.00"
	}
	return booking.TotalAmount.StringFixed(2)
}

func (c *Client) escapeOrDefault(value string, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		trimmed = fallback
	}
	return html.EscapeString(trimmed)
}

func (c *Client) nonZeroIntOrDefault(value int, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}
