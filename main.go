package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// STRUCTURE Contact

type Contact struct {
	ID    int
	Name  string
	Email string
}

// Constructeur avec validation
func NewContact(id int, name, email string) (*Contact, error) {
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

// Mise à jour partielle
func (c *Contact) Update(name, email string) {
	if strings.TrimSpace(name) != "" {
		c.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		c.Email = strings.TrimSpace(email)
	}
}

// Affichage
func (c *Contact) String() string {
	return fmt.Sprintf("ID:%d | Nom:%s | Email:%s", c.ID, c.Name, c.Email)
}

// INTERFACE Storer (contrat de stockage)

type Storer interface {
	Add(c *Contact) error
	GetByID(id int) (*Contact, bool)
	GetAll() []*Contact
	Update(c *Contact) error
	Delete(id int) error
}

// MemoryStore (en mémoire)

type MemoryStore struct {
	contacts map[int]*Contact
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{contacts: make(map[int]*Contact)}
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
	res := make([]*Contact, 0, len(m.contacts))
	for _, c := range m.contacts {
		res = append(res, c)
	}
	return res
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

// JSONFileStore (fichier JSON persistant)

type JSONFileStore struct {
	filename string
	contacts map[int]*Contact
}

func NewJSONFileStore(filename string) (*JSONFileStore, error) {
	store := &JSONFileStore{
		filename: filename,
		contacts: make(map[int]*Contact),
	}
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("erreur chargement fichier: %w", err)
	}
	return store, nil
}

func (j *JSONFileStore) load() error {
	file, err := os.Open(j.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var contacts []*Contact
	dec := json.NewDecoder(file)
	if err := dec.Decode(&contacts); err != nil {
		return fmt.Errorf("erreur décodage JSON: %w", err)
	}
	for _, c := range contacts {
		j.contacts[c.ID] = c
	}
	return nil
}

func (j *JSONFileStore) save() error {
	list := make([]*Contact, 0, len(j.contacts))
	for _, c := range j.contacts {
		list = append(list, c)
	}
	file, err := os.Create(j.filename)
	if err != nil {
		return fmt.Errorf("erreur création fichier: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(list); err != nil {
		return fmt.Errorf("erreur encodage JSON: %w", err)
	}
	return nil
}

func (j *JSONFileStore) Add(c *Contact) error {
	if _, exists := j.contacts[c.ID]; exists {
		return fmt.Errorf("un contact avec l'ID %d existe déjà", c.ID)
	}
	j.contacts[c.ID] = c
	return j.save()
}

func (j *JSONFileStore) GetByID(id int) (*Contact, bool) {
	c, ok := j.contacts[id]
	return c, ok
}

func (j *JSONFileStore) GetAll() []*Contact {
	res := make([]*Contact, 0, len(j.contacts))
	for _, c := range j.contacts {
		res = append(res, c)
	}
	return res
}

func (j *JSONFileStore) Update(c *Contact) error {
	if _, ok := j.contacts[c.ID]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", c.ID)
	}
	j.contacts[c.ID] = c
	return j.save()
}

func (j *JSONFileStore) Delete(id int) error {
	if _, ok := j.contacts[id]; !ok {
		return fmt.Errorf("contact avec l'ID %d introuvable", id)
	}
	delete(j.contacts, id)
	return j.save()
}

// MAIN

func main() {
	// Choisis l'un des deux stores:
	// store := NewMemoryStore()
	store, err := NewJSONFileStore("contacts.json")
	if err != nil {
		fmt.Println("Erreur initialisation store:", err)
		return
	}

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
			fmt.Println("Choix invalidRéessayez.")
		}
	}
}

// UTILITAIRES I/O

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

// CRUD via Storer

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
