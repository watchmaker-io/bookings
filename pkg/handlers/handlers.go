package handlers

import (
	"github.com/watchmaker-io/bookings/pkg/config"
	"github.com/watchmaker-io/bookings/pkg/models"
	"github.com/watchmaker-io/bookings/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	remoteIp := request.RemoteAddr
	m.App.Session.Put(request.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(writer, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, request *http.Request) {
	stringMap := map[string]string{
		"test": "Hello, again",
	}

	remoteIp := m.App.Session.GetString(request.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(writer, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
