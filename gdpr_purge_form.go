package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func gdprPurgeForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respondForm(w, &apps.Form{
		Title: "Purge all data for a customer (GDPR)",
		Icon:  appURL(creq, "/static/icon.png"),
		Fields: []*apps.Field{
			{
				Type:       apps.FieldTypeText,
				Name:       "email",
				Label:      "email",
				ModalLabel: "Customer Email",
				IsRequired: true,
			},
			{
				Type:       apps.FieldTypeBool,
				Name:       "confirm",
				Label:      "confirm",
				ModalLabel: "Confirm purging customer data",
				IsRequired: true,
			},
		},
		Call: &apps.Call{
			Path: "/gdpr-purge",
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
				OAuth2User:            apps.ExpandAll,
			},
		},
	})
}
