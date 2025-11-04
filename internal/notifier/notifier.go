package notifier

// Notifier définit l'interface pour envoyer des notifications.
type Notifier interface {
	Send(message string) error
}