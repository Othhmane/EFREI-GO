package service

import (
	"mini-crm/internal/domain"
	"mini-crm/internal/notifier"
	"mini-crm/internal/storage"
)

// ContactService gère la logique métier des contacts.
type ContactService struct {
	store     storage.Storer
	notifiers []notifier.Notifier // ← Ajout de la slice de notificateurs
}

// NewContactService crée un nouveau service de contacts.
func NewContactService(store storage.Storer, notifiers []notifier.Notifier) *ContactService {
	return &ContactService{
		store:     store,
		notifiers: notifiers,
	}
}

// AddContact crée et ajoute un nouveau contact.
func (s *ContactService) AddContact(id int, name, email string) error {
	contact, err := domain.NewContact(id, name, email)
	if err != nil {
		return err
	}

	// Ajouter le contact
	if err := s.store.Add(contact); err != nil {
		return err
	}

	// Envoyer des notifications
	s.notifyAll("Nouveau contact ajouté : " + name)

	return nil
}

// notifyAll envoie un message à tous les notificateurs.
func (s *ContactService) notifyAll(message string) {
	for _, n := range s.notifiers {
		if err := n.Send(message); err != nil {
			// Log l'erreur mais ne bloque pas l'opération
			println("Erreur notification:", err.Error())
		}
	}
}

// GetContactByID récupère un contact par son ID.
func (s *ContactService) GetContactByID(id int) (*domain.Contact, bool) {
	return s.store.GetByID(id)
}

// GetAllContacts retourne tous les contacts.
func (s *ContactService) GetAllContacts() []*domain.Contact {
	return s.store.GetAll()
}

// UpdateContact met à jour un contact existant.
func (s *ContactService) UpdateContact(id int, name, email string) error {
	contact, ok := s.store.GetByID(id)
	if !ok {
		return nil
	}
	contact.Update(name, email)

	// Envoyer des notifications
	s.notifyAll("Contact mis à jour : " + contact.Name)

	return s.store.Update(contact)
}

// DeleteContact supprime un contact par ID.
func (s *ContactService) DeleteContact(id int) error {
	// Récupérer le contact avant de le supprimer (pour la notification)
	contact, ok := s.store.GetByID(id)
	if !ok {
		return s.store.Delete(id) // Retourne l'erreur du store
	}

	// Supprimer
	if err := s.store.Delete(id); err != nil {
		return err
	}

	// Notifier
	s.notifyAll("Contact supprimé : " + contact.Name)

	return nil
}