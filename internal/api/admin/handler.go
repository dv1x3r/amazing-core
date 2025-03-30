package admin

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/services/auth"

	"github.com/dv1x3r/w2go/w2"
	"github.com/gorilla/csrf"
)

//go:embed *.gotmpl
var templatesFS embed.FS

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(templatesFS, "*.gotmpl"))
}

func Handler(authService *auth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := authService.GetSessionUsername(w, r)
		if ok {
			if err := authService.RefreshSession(w, r); err != nil {
				w2.NewErrorResponse(err.Error()).Write(w, http.StatusInternalServerError)
				return
			}
		}

		data := map[string]any{"csrfToken": csrf.Token(r), "username": username}
		if err := tmpl.ExecuteTemplate(w, "admin.gotmpl", data); err != nil {
			w2.NewErrorResponse(err.Error()).Write(w, http.StatusInternalServerError)
		}
	})
}
