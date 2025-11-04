package main

import (
	"bufio"
	"encoding/json" // Pour encoder/décoder JSON
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ============================================================
// 1. STRUCTURE Contact
// ============================================================

type Contact struct {
	ID    int
	Name  string
	Email string
}

// Constructeur avec validation
func NewContact(id int, name, email string) (*Contact, error) {
	if id <= 0 {
		return nil, fmt.Errorf("l'ID doit être positif")
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

// Méthode Update pour modifier un contact
func (c *Contact) Update(name, email string) {
	if strings.TrimSpace(name) != "" {
		c.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		c.Email = strings.TrimSpace(email)
	}
}

// Méthode String pour affichage (implémente fmt.Stringer)
func (c *Contact) String() string {
	return fmt.Sprintf("ID:%d | Nom:%s | Email:%s", c.ID, c.Name, c.Email)
}

// ============================================================
// 2. INTERFACE Storer (contrat de stockage)
// ============================================================

// Storer définit les opérations de stockage pour les contacts.
// Principe: on découple la logique métier (CRUD) de l'implémentation du stockage.
type Storer interface {
	Add(c *Contact) error
	GetByID(id int) (*Contact, bool)
	GetAll() []*Contact
	Update(c *Contact) error
	Delete(id int) error
}

// ============================================================
// 3. STRUCT MemoryStore (stockage en mémoire)
// ============================================================

type MemoryStore struct {
	contacts map[int]*Contact
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		contacts: make(map[int]*Contact),
	}
}

func (m *MemoryStore) Add(c *Contact) error {
	if _, exists := m.contacts[c.ID]; exists {
		return fmt.Errorf("un contact avec l'ID %d existe déjà", c.ID)
	}
	m.contacts[c.ID] = c
	return nil
}

func (m *MemoryStore) GetByID(id int) (*Contact, bool) {
	c, ok := m.contacts[id]
	return c, ok
}

func (m *MemoryStore) GetAll() []*Contact {
	result := make([]*Contact, 0, len(m.contacts))
	for _, c := range m.contacts {
		result = append(result, c)
	}
	return result
}

func (m *MemoryStore) Update(c *Contact) error {
	if _, ok := m.contacts[c.ID]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", c.ID)
	}
	m.contacts[c.ID] = c
	return nil
}

func (m *MemoryStore) Delete(id int) error {
	if _, ok := m.contacts[id]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", id)
	}
	delete(m.contacts, id)
	return nil
}

// ============================================================
// 4. STRUCT JSONFileStore (stockage persistant dans un fichier JSON)
// ============================================================

// JSONFileStore stocke les contacts dans un fichier JSON.
// Implémente l'interface Storer.
// Principe: chaque opération (Add, Update, Delete) sauvegarde immédiatement dans le fichier.
type JSONFileStore struct {
	filename string
	contacts map[int]*Contact
}

// NewJSONFileStore crée un JSONFileStore et charge les contacts depuis le fichier.
// Si le fichier n'existe pas, il sera créé au premier ajout.
func NewJSONFileStore(filename string) (*JSONFileStore, error) {
	store := &JSONFileStore{
		filename: filename,
		contacts: make(map[int]*Contact),
	}

	// Charger les contacts existants depuis le fichier
	if err := store.load(); err != nil {
		// Si le fichier n'existe pas, ce n'est pas une erreur (on démarre vide)
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("erreur chargement fichier: %w", err)
		}
	}

	return store, nil
}

// load charge les contacts depuis le fichier JSON.
func (j *JSONFileStore) load() error {
	// Ouvrir le fichier en lecture
	file, err := os.Open(j.filename)
	if err != nil {
		return err // Retourne l'erreur (peut être os.IsNotExist)
	}
	defer file.Close()

	// Décoder le JSON dans une slice temporaire
	var contacts []*Contact
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&contacts); err != nil {
		return fmt.Errorf("erreur décodage JSON: %w", err)
	}

	// Remplir la map depuis la slice
	for _, c := range contacts {
		j.contacts[c.ID] = c
	}

	return nil
}

// save sauvegarde tous les contacts dans le fichier JSON.
// Appelée après chaque modification (Add, Update, Delete).
func (j *JSONFileStore) save() error {
	// Convertir la map en slice pour l'encodage JSON
	contacts := make([]*Contact, 0, len(j.contacts))
	for _, c := range j.contacts {
		contacts = append(contacts, c)
	}

	// Créer/écraser le fichier
	file, err := os.Create(j.filename)
	if err != nil {
		return fmt.Errorf("erreur création fichier: %w", err)
	}
	defer file.Close()

	// Encoder en JSON avec indentation (lisible)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(contacts); err != nil {
		return fmt.Errorf("erreur encodage JSON: %w", err)
	}

	return nil
}

// Add ajoute un contact et sauvegarde dans le fichier.
func (j *JSONFileStore) Add(c *Contact) error {
	if _, exists := j.contacts[c.ID]; exists {
		return fmt.Errorf("un contact avec l'ID %d existe déjà", c.ID)
	}
	j.contacts[c.ID] = c
	return j.save() // Sauvegarde immédiate
}

// GetByID récupère un contact par ID (comma-ok idiom).
func (j *JSONFileStore) GetByID(id int) (*Contact, bool) {
	c, ok := j.contacts[id]
	return c, ok
}

// GetAll retourne tous les contacts.
func (j *JSONFileStore) GetAll() []*Contact {
	result := make([]*Contact, 0, len(j.contacts))
	for _, c := range j.contacts {
		result = append(result, c)
	}
	return result
}

// Update met à jour un contact et sauvegarde.
func (j *JSONFileStore) Update(c *Contact) error {
	if _, ok := j.contacts[c.ID]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", c.ID)
	}
	j.contacts[c.ID] = c
	return j.save() // Sauvegarde immédiate
}

// Delete supprime un contact et sauvegarde.
func (j *JSONFileStore) Delete(id int) error {
	if _, ok := j.contacts[id]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", id)
	}
	delete(j.contacts, id)
	return j.save() // Sauvegarde immédiate
}

// ============================================================
// 5. FONCTION PRINCIPALE (choix du store au démarrage)
// ============================================================

func main() {
	// Choix du store: décommenter celui que tu veux utiliser

	// Option 1: Stockage en mémoire (volatile)
	// store := NewMemoryStore()

	// Option 2: Stockage persistant dans un fichier JSON
	store, err := NewJSONFileStore("contacts.json")
	if err != nil {
		fmt.Println("Erreur initialisation store:", err)
		return
	}

	// Reader pour les entrées utilisateur
	reader := bufio.NewReader(os.Stdin)

	// Boucle principale du menu
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
			addContactInteractive(reader, store)
		case "2":
			listContacts(store)
		case "3":
			deleteContactInteractive(reader, store)
		case "4":
			updateContactInteractive(reader, store)
		case "5":
			fmt.Println("Au revoir!")
			return
		default:
			fmt.Println("Choix invalide. Réessayez.")
		}
	}
}

// ============================================================
// 6. AFFICHAGE DU MENU
// ============================================================

func printMenu() {
	fmt.Println()
	fmt.Println("=== Mini-CRM ===")
	fmt.Println("1) Ajouter un contact")
	fmt.Println("2) Lister les contacts")
	fmt.Println("3) Supprimer un contact par ID")
	fmt.Println("4) Mettre à jour un contact")
	fmt.Println("5) Quitter")
}

// ============================================================
// 7. UTILITAIRES I/O
// ============================================================

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

// ============================================================
// 8. OPÉRATIONS CRUD (reçoivent l'interface Storer)
// ============================================================

func addContactInteractive(r *bufio.Reader, store Storer) {
	fmt.Print("ID: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	fmt.Print("Nom: ")
	name, _ := readLine(r)

	fmt.Print("Email: ")
	email, _ := readLine(r)

	contact, err := NewContact(id, name, email)
	if err != nil {
		fmt.Println("Erreur création contact:", err)
		return
	}

	if err := store.Add(contact); err != nil {
		fmt.Println("Erreur ajout:", err)
		return
	}

	fmt.Println("✓ Contact ajouté avec succès.")
}

func listContacts(store Storer) {
	contacts := store.GetAll()
	if len(contacts) == 0 {
		fmt.Println("Aucun contact.")
		return
	}

	fmt.Println("\n=== Liste des contacts ===")
	for _, c := range contacts {
		fmt.Println("-", c)
	}
}

func deleteContactInteractive(r *bufio.Reader, store Storer) {
	fmt.Print("ID à supprimer: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	if err := store.Delete(id); err != nil {
		fmt.Println("Erreur suppression:", err)
		return
	}

	fmt.Println("✓ Contact supprimé.")
}

func updateContactInteractive(r *bufio.Reader, store Storer) {
	fmt.Print("ID à mettre à jour: ")
	idStr, _ := readLine(r)
	id, err := parseInt(idStr)
	if err != nil {
		fmt.Println("ID invalide:", err)
		return
	}

	c, ok := store.GetByID(id)
	if !ok {
		fmt.Println("ID introuvable.")
		return
	}

	fmt.Printf("Nouveau nom (laisser vide pour garder '%s'): ", c.Name)
	name, _ := readLine(r)

	fmt.Printf("Nouvel email (laisser vide pour garder '%s'): ", c.Email)
	email, _ := readLine(r)

	c.Update(name, email)

	if err := store.Update(c); err != nil {
		fmt.Println("Erreur mise à jour:", err)
		return
	}

	fmt.Println("✓ Contact mis à jour.")
}