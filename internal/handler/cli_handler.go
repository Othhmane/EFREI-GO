package handler

import (
	"bufio"
	"fmt"
	"mini-crm/internal/service"
	"mini-crm/pkg/utils"
)

// CLIHandler gère les interactions CLI.
type CLIHandler struct {
	service *service.ContactService
	reader  *bufio.Reader
}

// NewCLIHandler crée un nouveau handler CLI.
func NewCLIHandler(service *service.ContactService, reader *bufio.Reader) *CLIHandler {
	return &CLIHandler{
		service: service,
		reader:  reader,
	}
}

// AddContactInteractive demande les infos et ajoute un contact.
func (h *CLIHandler) AddContactInteractive() {
	fmt.Print("ID: ")
	idStr, _ := utils.ReadLine(h.reader)
	id, err := utils.ParseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	fmt.Print("Nom: ")
	name, _ := utils.ReadLine(h.reader)

	fmt.Print("Email: ")
	email, _ := utils.ReadLine(h.reader)

	if err := h.service.AddContact(id, name, email); err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Println("✓ Contact ajouté avec succès.")
}

// ListContacts affiche tous les contacts.
func (h *CLIHandler) ListContacts() {
	contacts := h.service.GetAllContacts()
	if len(contacts) == 0 {
		fmt.Println("Aucun contact.")
		return
	}

	fmt.Println("\n=== Liste des contacts ===")
	for _, c := range contacts {
		fmt.Println("-", c)
	}
}

// DeleteContactInteractive supprime un contact par ID.
func (h *CLIHandler) DeleteContactInteractive() {
	fmt.Print("ID à supprimer: ")
	idStr, _ := utils.ReadLine(h.reader)
	id, err := utils.ParseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	if err := h.service.DeleteContact(id); err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Println("✓ Contact supprimé.")
}

// UpdateContactInteractive met à jour un contact.
func (h *CLIHandler) UpdateContactInteractive() {
	fmt.Print("ID à mettre à jour: ")
	idStr, _ := utils.ReadLine(h.reader)
	id, err := utils.ParseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	contact, ok := h.service.GetContactByID(id)
	if !ok {
		fmt.Println("ID introuvable.")
		return
	}

	fmt.Printf("Nouveau nom (laisser vide pour garder '%s'): ", contact.Name)
	name, _ := utils.ReadLine(h.reader)

	fmt.Printf("Nouvel email (laisser vide pour garder '%s'): ", contact.Email)
	email, _ := utils.ReadLine(h.reader)

	if err := h.service.UpdateContact(id, name, email); err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Println("✓ Contact mis à jour.")
}