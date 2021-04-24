package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func createNoteCall(state string) *apps.Call {
	return &apps.Call{
		Path: "/create-note",
		Expand: &apps.Expand{
			ActingUserAccessToken: apps.ExpandAll,
			OAuth2User:            apps.ExpandAll,
			Post:                  apps.ExpandSummary,
		},
		State: state,
	}
}

func createNoteForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	title, _ := creq.Values["title"].(string)
	content, _ := creq.Values["content"].(string)
	email, _ := creq.Values["email"].(string)
	tags, _ := creq.Values["tags"].(string)
	if creq.Context.PostID != "" {
		post := creq.Context.Post
		if title == "" {
			title = "Mattermost: " + post.Message
			if len(title) > 72 {
				title = title[:72]
			}
		}
		if content != "" {
			content += "\n\n"
		}
		content = fmt.Sprintf("Note created from a message in Mattermost, %s:\n\n%s", permalink(creq), post.Message)
	}

	form := &apps.Form{
		Title: "Creates a ProductBoard Note",
		Icon:  appURL(creq, "/static/icon.png"),
		Fields: []*apps.Field{
			{
				Type:       apps.FieldTypeText,
				Name:       "title",
				Label:      "title",
				ModalLabel: "Title",
				IsRequired: true,
				Value:      title,
			},
			{
				Type:        apps.FieldTypeText,
				TextSubtype: "textarea",
				Name:        "content",
				Label:       "content",
				ModalLabel:  "Text",
				IsRequired:  true,
				Value:       content,
			},
			{
				Type:             apps.FieldTypeText,
				Name:             "email",
				Label:            "email",
				ModalLabel:       "Customer Email",
				Description:      "The email address of the user (customer) associated with this note",
				AutocompleteHint: "e.g. `name@example.test`",
				Value:            email,
			},
			{
				Type:             apps.FieldTypeText,
				Name:             "tags",
				Label:            "tags",
				ModalLabel:       "Tags",
				Description:      "Tags for the note, comma-separated. Use `` to enter spaces.",
				AutocompleteHint: "e.g. `3.0,for triage,large customer`",
				Value:            tags,
			},
		},
		Call: createNoteCall("form"),
	}

	if strings.HasPrefix(string(creq.Context.Location), string(apps.LocationCommand)) {
		form.Fields = append(form.Fields, &apps.Field{
			Type:        apps.FieldTypeBool,
			Name:        "interactive",
			Label:       "interactive",
			Description: "Review and edit interactively before submitting.",
		})
	}

	respondForm(w, form)
}
