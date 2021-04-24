package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func bindings(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	respond(w, []apps.Binding{
		{
			Location: apps.LocationCommand,
			Bindings: []*apps.Binding{
				{
					Icon:        appURL(creq, "/static/icon.png"),
					Location:    "pb",
					Label:       "pb",
					Description: "ProductBoard integration.",
					Hint:        "[ configure | create | gdpr ]",
					Bindings:    commandBindings(creq),
				},
			},
		},
		{
			Location: apps.LocationPostMenu,
			Bindings: []*apps.Binding{
				{
					Icon:        appURL(creq, "/static/icon.png"),
					Location:    "create-note-menu",
					Label:       "Create a ProductBoard Note",
					Description: "Use this post to create a Note in ProductBoard.",
					Hint:        " -- TODO --",
					Call:        createNoteCall("post-menu"),
				},
			},
		},
	}, "")
}

func commandBindings(creq *apps.CallRequest) []*apps.Binding {
	return []*apps.Binding{
		{
			Location:    "configure",
			Label:       "configure",
			Description: "Configure ProductBoard access.",
			Hint:        "[ flags ]",
			Call: &apps.Call{
				Path: "/configure",
			},
		},
		{
			Location:    "create",
			Label:       "create",
			Description: "Create an item in ProductBoard.",
			Hint:        "[ note ]",
			Bindings: []*apps.Binding{
				{
					Location:    "note",
					Label:       "note",
					Description: "Create a Note in ProductBoard.",
					Hint:        "[ flags ]",
					Call:        createNoteCall("command"),
				},
			},
		},
		{
			Location:    "gdpr",
			Label:       "gdpr",
			Description: "Manage GDPR compliance in ProductBoard.",
			Hint:        "[ purge ]",
			Bindings: []*apps.Binding{
				{
					Location:    "purge",
					Label:       "purge",
					Description: "Purge customer data in ProductBoard.",
					Hint:        "[ email ]",
					Call:        gdprPurgeCall,
				},
			},
		},
	}
}
