package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prashant9154/Booking_System/internal/config"
	"github.com/prashant9154/Booking_System/internal/models"
	"github.com/prashant9154/Booking_System/internal/render"
)

// Repo is the repository used by handlers
var Repo *Repository

// Repositiry is a Repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is a home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplates(w, r, "home.page.hbs", &models.TemplateData{})
}

// About is a about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic here
	stringMap := map[string]string{}

	stringMap["test"] = "Hello Again!"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//send data to the template
	render.RenderTemplates(w, r, "about.page.hbs", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Contact is a Contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "contact.page.hbs", &models.TemplateData{})
}

// Generals is a Generals page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "generals-quarter.page.hbs", &models.TemplateData{})
}

// Majors is a Majors page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "majors-suite.page.hbs", &models.TemplateData{})
}

// Availability is a Search avaialability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "search-availability.page.hbs", &models.TemplateData{})
}

// PostAvailability is a Search avaialability page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is: %s and End date is: %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for the availability and JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")

	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation is a make reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "make-reservation.page.hbs", &models.TemplateData{})
}
