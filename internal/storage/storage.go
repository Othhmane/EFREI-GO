package storage

import "mini-crm/internal/domain"

// Storer définit l'interface pour les opérations de stockage des contacts
// Permet de découpler la logique métier de l'implémentation du stockage
type Storer interface {
	Add(c *domain.Contact) error
	GetByID(id int) (*domain.Contact, bool)
	GetAll() []*domain.Contact
	Update(c *domain.Contact) error
	Delete(id int) error
}