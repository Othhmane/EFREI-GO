package notifier

import "fmt"

// SmsNotifier envoie des notifications par SMS.
type SmsNotifier struct {
	PhoneNumber string
}

// NewSmsNotifier crée un nouveau notificateur SMS.
func NewSmsNotifier(phoneNumber string) *SmsNotifier {
	return &SmsNotifier{
		PhoneNumber: phoneNumber,
	}
}

// Send envoie un SMS (simulation).
func (s *SmsNotifier) Send(message string) error {
	fmt.Printf("📱 [SMS] À: %s | Message: %s\n", s.PhoneNumber, message)
	// Ici tu pourrais utiliser une API comme Twilio
	return nil
}