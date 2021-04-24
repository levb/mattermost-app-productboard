package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func createNoteForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respondForm(w, &apps.Form{
		Title: "Creates a ProductBoard Note",
		Icon:  appURL(creq, "/static/icon.png"),
		Fields: []*apps.Field{
			{
				Type:       apps.FieldTypeText,
				Name:       "title",
				Label:      "title",
				ModalLabel: "Title",
				IsRequired: true,
			},
			{
				Type:       apps.FieldTypeText,
				Name:       "content",
				Label:      "content",
				ModalLabel: "Text",
				IsRequired: true,
			},
			{
				Type:             apps.FieldTypeText,
				Name:             "email",
				Label:            "email",
				ModalLabel:       "Customer Email",
				Description:      "The email address of the user (customer) associated with this note",
				AutocompleteHint: "e.g. `name@example.test`",
			},
			{
				Type:             apps.FieldTypeText,
				Name:             "tags",
				Label:            "tags",
				ModalLabel:       "Tags",
				Description:      "Tags for the note, comma-separated. Use `` to enter spaces.",
				AutocompleteHint: "e.g. `3.0,for triage,large customer`",
			},
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
