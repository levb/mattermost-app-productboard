// TODO:
//
// - Webhooks: make them work, if possible, or remove
// - add command and form icons
// - /invite @user --create-note
// - /invite --list
// - /uninvite @user --create-note
// 		- /disconnect --invite @from-user (same as ^^, but from the other end)
// - /create note --as @user
// 		- /create note --as default if only 1 is available
//
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/felixge/httpsnoop"
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

//go:embed icon.png
var iconData []byte

//go:embed manifest.json
var manifestData []byte

func main() {
	// Serve its own manifest as HTTP for convenience in dev. mode.
	withLog("/manifest.json", handleData("application/json", manifestData))

	// Serve the icon for the App.
	withLog("/static/icon.png", handleData("image/png", iconData))

	// Serve the Channel Header and Command bindings for the App.
	withLog("/bindings", call(bindings))

	// `connect` and `disconnect` commands - acces tokens
	withLog("/connect/form", call(connectForm))
	withLog("/connect/submit", call(connect))
	withLog("/disconnect/form", call(disconnectForm))
	withLog("/disconnect/submit", call(disconnect))

	// `create note` command - creates a note.
	withLog("/create-note/form", call(createNoteForm))
	withLog("/create-note/submit", call(createNote))

	// `subscribe` command - creates a webhook subscription.
	withLog("/subscribe/form", call(subscribeForm))
	withLog("/subscribe/submit", call(subscribe))

	// `gdpr` command - purge customer data.
	withLog("/gdpr-purge/form", call(gdprPurgeForm))
	withLog("/gdpr-purge/submit", call(gdprPurge))

	withLog("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	log.Printf("listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func withLog(path string, f http.HandlerFunc) {
	http.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(f, w, r)
			log.Printf(
				"%s %s (code=%d dt=%s written=%d)",
				r.Method, r.URL, m.Code, m.Duration, m.Written)
		},
	)
}

func call(f func(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		creq := apps.CallRequest{}
		err := json.NewDecoder(req.Body).Decode(&creq)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		f(w, req, &creq)
	}
}

func handleData(ct string, data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		writeData(w, ct, http.StatusOK, data)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("failed to encode output: %v", err)
		return
	}
	writeData(w, "application/json", statusCode, data)
}

func respond(w http.ResponseWriter, v interface{}, format string, args ...interface{}) {
	writeJSON(w, http.StatusOK, apps.CallResponse{
		Type:     apps.CallResponseTypeOK,
		Data:     v,
		Markdown: md.Markdownf(format, args...),
	})
}

func respondError(w http.ResponseWriter, statusCode int, err error) {
	writeJSON(w, statusCode, apps.CallResponse{
		Type:      apps.CallResponseTypeError,
		Data:      err,
		ErrorText: err.Error(),
	})
}

func respondForm(w http.ResponseWriter, f *apps.Form) {
	writeJSON(w, http.StatusOK, apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: f,
	})
}

func writeData(w http.ResponseWriter, ct string, statusCode int, data []byte) {
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(data)))
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		log.Printf("failed to encode output: %v", err)
		return
	}
}

func appURL(creq *apps.CallRequest, path string) string {
	return creq.Context.MattermostSiteURL + creq.Context.AppPath + path
}

func permalink(creq *apps.CallRequest) string {
	if creq.Context.PostID == "" {
		return ""
	}
	u, _ := url.Parse(creq.Context.MattermostSiteURL)
	u.Path = path.Join(u.Path, "_redirect", "pl", creq.Context.PostID)
	return u.String()
}
