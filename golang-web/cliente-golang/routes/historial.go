package routes

import (
	"html/template"
	"log"
	"net/http"
)

//GET
func addEntradaHistorialFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/addEntrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", Page{Title: "Añadir entrada historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func addAnaliticaHistorialFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/addAnalitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", Page{Title: "Añadir analitica historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
