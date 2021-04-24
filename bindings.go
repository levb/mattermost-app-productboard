package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func bindings(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	commandBindings := []*apps.Binding{
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
					Call: &apps.Call{
						Path: "/create-note",
					},
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
					Call: &apps.Call{
						Path: "/gdpr-purge",
					},
				},
			},
		},
	}

	respond(w, []apps.Binding{
		{
			Location: "/command",
			Bindings: []*apps.Binding{
				{
					Icon:        appURL(creq, "/static/icon.png"),
					Location:    "pb",
					Label:       "pb",
					Description: "ProductBoard integration.",
					Hint:        "[ configure | create | gdpr ]",
					Bindings:    commandBindings,
				},
			},
		},
	}, "")
}
