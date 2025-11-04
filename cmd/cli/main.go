package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mini-crm/internal/handler"
	"mini-crm/internal/service"
	"mini-crm/internal/storage"
	"mini-crm/pkg/utils"
)

func main() {
	store, err := storage.NewJSONFileStore("contacts.json")
	if err != nil {
		fmt.Println("Erreur initialisation store:", err)
		return
	}

	svc := service.NewContactService(store)

	reader := bufio.NewReader(os.Stdin)
	cliHandler := handler.NewCLIHandler(svc, reader)

	for {
		printMenu()
		fmt.Print("Votre choix: ")
		line, err := utils.ReadLine(reader)
		if err != nil {
			fmt.Println("Erreur lecture:", err)
			continue
		}
		switch strings.TrimSpace(line) {
		case "1":
			cliHandler.AddContactInteractive()
		case "2":
			cliHandler.ListContacts()
		case "3":
			cliHandler.DeleteContactInteractive()
		case "4":
			cliHandler.UpdateContactInteractive()
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