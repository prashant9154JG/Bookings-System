package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prashant9154/Booking_System/internal/config"
	"github.com/prashant9154/Booking_System/internal/forms"
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

// Reservation is a reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplates(w, r, "make-reservation.page.hbs", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name")
	form.Required("first_name", "last_name", "email", "phone")

	form.ValidEmail("email")

	form.MinLength("first_name", 3)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplates(w, r, "make-reservation.page.hbs", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	summary, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		log.Println("cannot get items from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = summary

	render.RenderTemplates(w, r, "reservation-summary.page.hbs", &models.TemplateData{
		Data: data,
	})
}
