package domain

import (
	"fmt"
	"strings"
)

// Contact représente un contact du CRM.
type Contact struct {
	ID    int
	Name  string
	Email string
}

// NewContact crée un nouveau Contact après validation.
// Retourne une erreur si les données sont invalides.
func NewContact(id int, name, email string) (*Contact, error) {
	if id <= 0 {
		return nil, fmt.Errorf("l'ID doit être positif (reçu: %d)", id)
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("le nom est obligatoire")
	}
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("l'email est obligatoire")
	}
	return &Contact{
		ID:    id,
		Name:  strings.TrimSpace(name),
		Email: strings.TrimSpace(email),
	}, nil
}

// Update met à jour les champs du contact.
// Si un paramètre est vide, la valeur actuelle est conservée.
func (c *Contact) Update(name, email string) {
	if strings.TrimSpace(name) != "" {
		c.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		c.Email = strings.TrimSpace(email)
	}
}

// String retourne une représentation textuelle du contact.
func (c *Contact) String() string {
	return fmt.Sprintf("ID:%d | Nom:%s | Email:%s", c.ID, c.Name, c.Email)
}