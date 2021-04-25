package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

var subscribeCall = &apps.Call{
	Path: "/subscribe",
	Expand: &apps.Expand{
		App:              apps.ExpandAll,
		AdminAccessToken: apps.ExpandAll,
	},
}

func subscribeForm(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respondForm(w, &apps.Form{
		Title:  "DOES NOT WORK: Subscribe to receive ProductBoard updates (webhooks))",
		Icon:   appURL(creq, "/static/icon.png"),
		Call:   subscribeCall,
		Fields: []*apps.Field{},
	})
}
