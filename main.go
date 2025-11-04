package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Contact struct {
	ID    int
	Name  string
	Email string
}

var contacts = make(map[int]Contact)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu()
		fmt.Print("Votre choix: ")
		line, err := readLine(reader)
		if err != nil {
			fmt.Println("Erreur lecture:", err)
			continue
		}
		switch strings.TrimSpace(line) {
		case "1":
			// Ajouter un contact 
			addContactInteractive(reader)
		case "2":
			listContacts()
		case "3":
			deleteContactInteractive(reader)
		case "4":
			updateContactInteractive(reader)
		case "5":
			fmt.Println("Au revoir!")
			return
		default:
			fmt.Println("Choix invalide. Réessayez.")
		}
	}
}

func printMenu() {
	fmt.Println()
	fmt.Println("=== Mini-CRM ===")
	fmt.Println("1) Ajouter un contact")
	fmt.Println("2) Lister les contacts")
	fmt.Println("3) Supprimer un contact par ID")
	fmt.Println("4) Mettre à jour un contact")
	fmt.Println("5) Quitter")
}

func readLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}


func addContactInteractive(r *bufio.Reader) {
	fmt.Print("ID: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}
	// comma-ok idiom pour vérifier collision
	if _, exists := contacts[id]; exists {
		fmt.Println("Un contact avec cet ID existe déjà.")
		return
	}

	fmt.Print("Nom: ")
	name, _ := readLine(r)
	if name == "" {
		fmt.Println("Nom requis.")
		return
	}

	fmt.Print("Email: ")
	email, _ := readLine(r)
	if email == "" {
		fmt.Println("Email requis.")
		return
	}

	contacts[id] = Contact{ID: id, Name: name, Email: email}
	fmt.Println("Contact ajouté.")
}

func listContacts() {
	if len(contacts) == 0 {
		fmt.Println("Aucun contact.")
		return
	}
	fmt.Println("Contacts:")
	for id, c := range contacts {
		fmt.Printf("- ID:%d | Nom:%s | Email:%s\n", id, c.Name, c.Email)
	}
}

func deleteContactInteractive(r *bufio.Reader) {
	fmt.Print("ID à supprimer: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}
	if _, ok := contacts[id]; !ok { // comma ok idiom
		fmt.Println("ID introuvable.")
		return
	}
	delete(contacts, id)
	fmt.Println("Contact supprimé.")
}

func updateContactInteractive(r *bufio.Reader) {
	fmt.Print("ID à mettre à jour: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}
	c, ok := contacts[id] // comma ok idiom
	if !ok {
		fmt.Println("ID introuvable.")
		return
	}

	fmt.Printf("Nouveau nom (laisser vide pour garder: %s): ", c.Name)
	name, _ := readLine(r)
	if name != "" {
		c.Name = name
	}

	fmt.Printf("Nouvel email (laisser vide pour garder: %s): ", c.Email)
	email, _ := readLine(r)
	if email != "" {
		c.Email = email
	}

	contacts[id] = c
	fmt.Println("Contact mis à jour.")
}