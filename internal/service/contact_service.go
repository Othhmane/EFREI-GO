package service

import (
	"mini-crm/internal/domain"
	"mini-crm/internal/storage"
)

type ContactService struct {
	store storage.Storer
}

func NewContactService(store storage.Storer) *ContactService {
	return &ContactService{store: store}
}

func (s *ContactService) AddContact(id int, name, email string) error {
	contact, err := domain.NewContact(id, name, email)
	if err != nil {
		return err
	}
	return s.store.Add(contact)
}

func (s *ContactService) GetContactByID(id int) (*domain.Contact, bool) {
	return s.store.GetByID(id)
}

func (s *ContactService) GetAllContacts() []*domain.Contact {
	return s.store.GetAll()
}

func (s *ContactService) UpdateContact(id int, name, email string) error {
	contact, ok := s.store.GetByID(id)
	if !ok {
		return nil
	}
	contact.Update(name, email)
	return s.store.Update(contact)
}

func (s *ContactService) DeleteContact(id int) error {
	return s.store.Delete(id)
}