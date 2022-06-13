package svc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/aexvir/lnk/api"
	"github.com/aexvir/lnk/internal/logging"
)

type LinkStore interface {
	CreateLink(target string, slug *string) (string, error)
	GetLink(slug string) (*api.Link, error)
	DeleteLink(slug string) error
	AllLinks() []*api.Link

	GetTarget(slug string) (string, error)
	RegisterHit(slug string)
}

// LinkRedirectHandler fetches a shortened link by slug and if there is a link
// with that slug, it redirects to that link's target url.
func LinkRedirectHandler(store LinkStore) http.HandlerFunc {
	log := logging.NewLogger("lnk.redirect")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respond(w, http.StatusMethodNotAllowed, "only get requests allowed")
			return
		}

		slug := path.Base(r.URL.Path)
		log.Write("visit", "slug: %s", slug)

		target, err := store.GetTarget(slug)
		if err != nil {
			respond(w, http.StatusNotFound, err.Error())
			return
		}

		store.RegisterHit(slug)
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	}
}

// LinkMgmtHandler performs different actions depending on the request method.
// The actions match the expected restful verb.
func LinkMgmtHandler(store LinkStore) http.HandlerFunc {
	log := logging.NewLogger("lnk.mgmt")

	return func(w http.ResponseWriter, r *http.Request) {
		slug := path.Base(r.URL.Path)

		switch r.Method {
		case http.MethodGet:
			log.Write("get", "slug: %s", slug)

			if slug == "links" { // /api/links
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(store.AllLinks())
				if err != nil {
					log.Write("error", err.Error())
					respond(w, http.StatusInternalServerError, err.Error())
					return
				}
				return
			}

			link, err := store.GetLink(slug)
			if err != nil {
				log.Write("error", err.Error())
				respond(w, http.StatusNotFound, "link not found: %s", slug)
				return
			}

			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(link)
			if err != nil {
				log.Write("error", err.Error())
				respond(w, http.StatusInternalServerError, err.Error())
			}

		case http.MethodPost:
			var payload api.CreateLinkReq
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				log.Write("error", err.Error())
				respond(w, http.StatusBadRequest, "invalid payload: %s", err.Error())
				return
			}

			slug, err := store.CreateLink(payload.Target, payload.Slug)
			if err != nil {
				log.Write("error", err.Error())
				respond(w, http.StatusInternalServerError, err.Error())
				return
			}

			ret := api.LinkResp{Link: slug}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(ret)
			if err != nil {
				log.Write("error", err.Error())
				respond(w, http.StatusInternalServerError, err.Error())
				return
			}

			log.Write("post", "created slug: %s", slug)

		case http.MethodDelete:
			err := store.DeleteLink(slug)
			if err != nil {
				log.Write("error", err.Error())
				respond(w, http.StatusBadRequest, err.Error())
			}

			respond(w, http.StatusOK, "deleted")
			log.Write("delete", "slug %s", slug)

		default:
			respond(w, http.StatusMethodNotAllowed, "")
		}
	}
}

func respond(w http.ResponseWriter, status int, msg string, args ...any) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf(msg, args...)))
}
