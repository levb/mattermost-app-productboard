package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func createNoteForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	field := func(n string) *apps.Field {
		return &apps.Field{
			Type:       apps.FieldTypeText,
			Name:       n,
			Label:      n,
			ModalLabel: "Note " + n,
			IsRequired: true,
		}
	}

	respondForm(w, &apps.Form{
		Title: "Creates a ProductBoard Note",
		Icon:  appURL(creq, "/static/icon.png"),
		Fields: []*apps.Field{
			field("title"),
			field("content"),
		},
		Call: &apps.Call{
			Path: "/create-note",
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
				OAuth2User:            apps.ExpandAll,
			},
		},
	})
}
