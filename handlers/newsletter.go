package handlers

import (
	"canvas/model"
	"canvas/storage"
	"canvas/views"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type signupper interface {
	SignupForNewsletter(ctx context.Context, email model.Email) (string, error)
}

func NewsLetterSignup(mux chi.Router, s *storage.Database) {
	mux.Post("/newsletter/signup", func(w http.ResponseWriter, r *http.Request) {

		email := model.Email(r.FormValue("email"))
		if !email.IsValid() {
			http.Error(w, "email is invalid", http.StatusBadRequest)
		}

		if _, err := s.SignupForNewsLetter(r.Context(), email); err != nil {
			http.Error(w, "error signing up, refresh to try again", http.StatusBadRequest)
		}

		http.Redirect(w, r, "/newsletter/thanks", http.StatusFound)
	})
}

func NewsletterThanks(mux chi.Router) {
	mux.Get("/newsletter/thanks", func(w http.ResponseWriter, r *http.Request) {
		_ = views.NewsletterThanksPage("/newsletter/thanks").Render(w)
	})
}
