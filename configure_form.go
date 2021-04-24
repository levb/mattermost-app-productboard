package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func configureForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respondForm(w, &apps.Form{
		Title: "Configures ProductBoard access credentials",
		Icon:  appURL(creq, "/static/icon.png"),
		Fields: []*apps.Field{
			{
				Type:       apps.FieldTypeText,
				Name:       "access_token",
				ModalLabel: "Personal Access Token",
				Label:      "access-token",
			},
			{
				Type:       apps.FieldTypeText,
				Name:       "gdpr_token",
				ModalLabel: "GDPR Public API Access Token",
				Label:      "gdpr-token",
			},
			{
				Type:       apps.FieldTypeBool,
				Name:       "clear",
				ModalLabel: "Clear All Tokens",
				Label:      "clear",
			},
		},
		Call: &apps.Call{
			Path: "/configure",
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
			},
		},
	})
}
