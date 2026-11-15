package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Project struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Image       string   `json:"image"`
	Stack       []string `json:"stack"`
	Intro       string   `json:"intro"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	Repository  string   `json:"repository"`
	Video       string   `json:"video"`
}

var templates *template.Template
var projets []Project

func chargerProjets() {
	fichier, err := os.ReadFile("./codes/projets.json")
	if err != nil {
		log.Fatal("Impossible de lire projets.json :", err)
	}

	err = json.Unmarshal(fichier, &projets)
	if err != nil {
		log.Fatal("Erreur dans le format du JSON :", err)
	}
}

func afficherPage(w http.ResponseWriter, nomPage string, donnees interface{}) {
	err := templates.ExecuteTemplate(w, nomPage, donnees)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		log.Println("Erreur template :", err)
	}
}

func pageAccueil(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		afficherPage(w, "erreur.html", nil)
		return
	}
	afficherPage(w, "index.html", nil)
}

func pageProjets(w http.ResponseWriter, r *http.Request) {
	afficherPage(w, "projets.html", projets)
}

func pageDetailProjet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/projects/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		afficherPage(w, "erreur.html", nil)
		return
	}

	for _, p := range projets {
		if p.ID == id {
			afficherPage(w, "detail_projet.html", p)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	afficherPage(w, "erreur.html", nil)
}

func pageContact(w http.ResponseWriter, r *http.Request) {
	afficherPage(w, "contact.html", nil)
}

func main() {
	chargerProjets()

	templates = template.Must(template.ParseFiles(
		"./index.html",
		"./pages/projets.html",
		"./pages/detail_projet.html",
		"./pages/contact.html",
		"./pages/erreur.html",
	))

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", pageAccueil)
	http.HandleFunc("/projects", pageProjets)
	http.HandleFunc("/projects/", pageDetailProjet)
	http.HandleFunc("/contact", pageContact)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
