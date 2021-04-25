package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func connectForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respondForm(w, &apps.Form{
		Title: "Connect to ProductBoard (configure access)",
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
		},
		Call: &apps.Call{
			Path: "/connect",
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
			},
		},
	})
}
