package handler

import (
	"net/http"

	"github.com/prashant9154/Booking_System/pkg/config"
	"github.com/prashant9154/Booking_System/pkg/models"
	"github.com/prashant9154/Booking_System/pkg/render"
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

	render.RenderTemplates(w, "home.page.hbs", &models.TemplateData{})
}

// About is a about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic here
	stringMap := map[string]string{}

	stringMap["test"] = "Hello Again!"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//send data to the template
	render.RenderTemplates(w, "about.page.hbs", &models.TemplateData{
		StringMap: stringMap,
	})
}
