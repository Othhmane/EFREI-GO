package storage

import (
	"encoding/json"
	"fmt"
	"mini-crm/internal/domain"
	"os"
)

//  stocke les contacts dans un fichier JSON (persistant)
type JSONFileStore struct {
	filename string
	contacts map[int]*domain.Contact
}

//  crée un JSONFileStore et charge les contacts depuis le fichier
func NewJSONFileStore(filename string) (*JSONFileStore, error) {
	store := &JSONFileStore{
		filename: filename,
		contacts: make(map[int]*domain.Contact),
	}

	// Charger les contacts existants
	if err := store.load(); err != nil {
		// Si le fichier n'existe pas, ce n'est pas une erreur
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("erreur chargement fichier: %w", err)
		}
	}

	return store, nil
}

// load charge les contacts depuis le fichier JSON
func (j *JSONFileStore) load() error {
	file, err := os.Open(j.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var contacts []*domain.Contact
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&contacts); err != nil {
		return fmt.Errorf("erreur décodage JSON: %w", err)
	}

	// Remplir la map
	for _, c := range contacts {
		j.contacts[c.ID] = c
	}

	return nil
}

// save sauvegarde tous les contacts dans le fichier JSON.
func (j *JSONFileStore) save() error {
	// Convertir la map en slice
	contacts := make([]*domain.Contact, 0, len(j.contacts))
	for _, c := range j.contacts {
		contacts = append(contacts, c)
	}

	// Créer/écraser le fichier
	file, err := os.Create(j.filename)
	if err != nil {
		return fmt.Errorf("erreur création fichier: %w", err)
	}
	defer file.Close()

	// Encoder en JSON avec indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(contacts); err != nil {
		return fmt.Errorf("erreur encodage JSON: %w", err)
	}

	return nil
}

// Add ajoute un contact et sauvegarde.
func (j *JSONFileStore) Add(c *domain.Contact) error {
	if _, exists := j.contacts[c.ID]; exists {
		return fmt.Errorf("un contact avec l'ID %d existe déjà", c.ID)
	}
	j.contacts[c.ID] = c
	return j.save()
}

// GetByID récupère un contact par ID.
func (j *JSONFileStore) GetByID(id int) (*domain.Contact, bool) {
	c, ok := j.contacts[id]
	return c, ok
}

// GetAll retourne tous les contacts.
func (j *JSONFileStore) GetAll() []*domain.Contact {
	result := make([]*domain.Contact, 0, len(j.contacts))
	for _, c := range j.contacts {
		result = append(result, c)
	}
	return result
}

// Update met à jour un contact et sauvegarde.
func (j *JSONFileStore) Update(c *domain.Contact) error {
	if _, ok := j.contacts[c.ID]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", c.ID)
	}
	j.contacts[c.ID] = c
	return j.save()
}

// Delete supprime un contact et sauvegarde.
func (j *JSONFileStore) Delete(id int) error {
	if _, ok := j.contacts[id]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", id)
	}
	delete(j.contacts, id)
	return j.save()
}