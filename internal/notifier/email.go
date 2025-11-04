package notifier

import "fmt"

// EmailNotifier envoie des notifications par email.
type EmailNotifier struct {
	From string
	To   string
}

// NewEmailNotifier crée un nouveau notificateur email.
func NewEmailNotifier(from, to string) *EmailNotifier {
	return &EmailNotifier{
		From: from,
		To:   to,
	}
}

// Send envoie un email (simulation).
func (e *EmailNotifier) Send(message string) error {
	fmt.Printf("📧 [EMAIL] De: %s | À: %s | Message: %s\n", e.From, e.To, message)
	// Ici tu pourrais utiliser une vraie lib comme smtp ou sendgrid
	return nil
}