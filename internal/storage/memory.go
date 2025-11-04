package storage

import (
	"fmt"
	"mini-crm/internal/domain"
)

// MemoryStore stocke les contacts en mémoire (volatile).
type MemoryStore struct {
	contacts map[int]*domain.Contact
}

// NewMemoryStore crée une nouvelle instance de MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		contacts: make(map[int]*domain.Contact),
	}
}

// Add ajoute un contact dans le store.
func (m *MemoryStore) Add(c *domain.Contact) error {
	if _, exists := m.contacts[c.ID]; exists {
		return fmt.Errorf("un contact avec l'ID %d existe déjà", c.ID)
	}
	m.contacts[c.ID] = c
	return nil
}

// GetByID récupère un contact par son ID.
func (m *MemoryStore) GetByID(id int) (*domain.Contact, bool) {
	c, ok := m.contacts[id]
	return c, ok
}

// GetAll retourne tous les contacts.
func (m *MemoryStore) GetAll() []*domain.Contact {
	result := make([]*domain.Contact, 0, len(m.contacts))
	for _, c := range m.contacts {
		result = append(result, c)
	}
	return result
}

// Update met à jour un contact existant.
func (m *MemoryStore) Update(c *domain.Contact) error {
	if _, ok := m.contacts[c.ID]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", c.ID)
	}
	m.contacts[c.ID] = c
	return nil
}

// Delete supprime un contact par ID.
func (m *MemoryStore) Delete(id int) error {
	if _, ok := m.contacts[id]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", id)
	}
	delete(m.contacts, id)
	return nil
}