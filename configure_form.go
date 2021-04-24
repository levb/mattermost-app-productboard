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
				Type:       "text",
				Name:       "token",
				ModalLabel: "Personal Access Token",
				Label:      "token",
				IsRequired: true,
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
