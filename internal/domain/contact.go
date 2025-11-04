package domain

import (
	"fmt"
	"strings"
)

type Contact struct {
	ID    int
	Name  string
	Email string
}

// création un nouveau Contact après validation
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

// Si un paramètre est vide, la valeur actuelle est conservée
func (c *Contact) Update(name, email string) {
	if strings.TrimSpace(name) != "" {
		c.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		c.Email = strings.TrimSpace(email)
	}
}

//  retourne contact
func (c *Contact) String() string {
	return fmt.Sprintf("ID:%d | Nom:%s | Email:%s", c.ID, c.Name, c.Email)
}