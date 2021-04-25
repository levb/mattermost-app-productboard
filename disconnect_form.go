package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func disconnectForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respondForm(w, &apps.Form{
		Title: "Disconnect ProductBoard access)",
		Icon:  appURL(creq, "/static/icon.png"),
		Call: &apps.Call{
			Path: "/disconnect",
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
			},
		},
		Fields: []*apps.Field{},
	})
}
