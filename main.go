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


// Principe: Je valide AVANT de créer l'objet, pour garantir qu'un Contact est toujours valide.
func NewContact(id int, name, email string) (*Contact, error) {
	// Validation: ID doit être strictement positif
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


// Si un paramètre est vide, je conserve la valeur actuelle.
func (c *Contact) Update(name, email string) {
	if strings.TrimSpace(name) != "" {
		c.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		c.Email = strings.TrimSpace(email)
	}
}

// String retourne une représentation textuelle du Contact (pour affichage)
func (c *Contact) String() string {
	return fmt.Sprintf("ID:%d | Nom:%s | Email:%s", c.ID, c.Name, c.Email)
}


// on peut modifier directement les champs sans réassigner dans la map
var contacts = make(map[int]*Contact)


func main() {
	// Reader tamponné pour lire les entrées utilisateur
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("Votre choix: ")

		// Lecture du choix utilisateur
		line, err := readLine(reader)
		if err != nil {
			fmt.Println("Erreur lecture:", err)
			continue
		}

		switch strings.TrimSpace(line) {
		case "1":
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


// readLine lit une ligne depuis le reader et retire les espaces de début/fin.
// Gère l'erreur de lecture proprement.
func readLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

// Retourne l'erreur si la conversion échoue
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

///////////////////////////////////////////////// 8. OPÉRATIONS CRUD /////////////////////////////////////////////////

// addContactInteractive demande les infos au clavier et ajoute un contact.
// Utilise le constructeur NewContact pour valider les données.
func addContactInteractive(r *bufio.Reader) {
	fmt.Print("ID: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	if _, exists := contacts[id]; exists {
		fmt.Println("Un contact avec cet ID existe déjà.")
		return
	}

	fmt.Print("Nom: ")
	name, _ := readLine(r)

	fmt.Print("Email: ")
	email, _ := readLine(r)

	// Créer le contact via le constructeur (validation automatique)
	contact, err := NewContact(id, name, email)
	if err != nil {
		fmt.Println("Erreur création contact:", err)
		return
	}

	// Stocker le pointeur dans la map
	contacts[id] = contact
	fmt.Println("✓ Contact ajouté avec succès.")
}

// listContacts affiche tous les contacts présents dans la map
// Note: l'ordre d'itération d'une map n'est pas garanti
func listContacts() {
	if len(contacts) == 0 {
		fmt.Println("Aucun contact.")
		return
	}

	fmt.Println("\n=== Liste des contacts ===")
	// Itération sur la map (clé = id, valeur = pointeur vers Contact)
	for _, c := range contacts {
		// Utilise la méthode String() du Contact (fmt.Stringer)
		fmt.Println("-", c)
	}
}

// deleteContactInteractive supprime un contact par ID après vérification.
func deleteContactInteractive(r *bufio.Reader) {
	fmt.Print("ID à supprimer: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	// Vérifier l'existence avec comma-ok idiom
	if _, ok := contacts[id]; !ok {
		fmt.Println("ID introuvable.")
		return
	}

	// Suppression dans la map (O(1))
	delete(contacts, id)
	fmt.Println("✓ Contact supprimé.")
}

// updateContactInteractive met à jour un contact existant.
// Utilise la méthode Update() attachée au type Contact.
func updateContactInteractive(r *bufio.Reader) {
	fmt.Print("ID à mettre à jour: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	// Récupérer le contact (comma-ok idiom)
	// Ici, c est un *pointeur* vers Contact, donc on peut modifier directement
	c, ok := contacts[id]
	if !ok {
		fmt.Println("ID introuvable.")
		return
	}

	fmt.Printf("Nouveau nom (laisser vide pour garder '%s'): ", c.Name)
	name, _ := readLine(r)

	// Proposer de modifier l'email
	fmt.Printf("Nouvel email (laisser vide pour garder '%s'): ", c.Email)
	email, _ := readLine(r)

	// Appeler la méthode Update du Contact (logique encapsulée)
	c.Update(name, email)

	// Pas besoin de réassigner dans la map: c est un pointeur,
	// les modifications sont directement reflétées dans contacts[id]
	fmt.Println("✓ Contact mis à jour.")
}
