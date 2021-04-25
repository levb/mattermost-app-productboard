package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func bindings(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	user, _ := userFromContext(creq)
	isConnectedAPI := user != nil && user.AccessToken != ""
	isConnectedGDPR := user != nil && user.GDPRToken != ""
	isAdmin := creq.Context.ActingUser != nil && creq.Context.ActingUser.IsSystemAdmin()

	bindings := []apps.Binding{}
	if isConnectedAPI {
		bindings = append(bindings, apps.Binding{
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
		})
	}

	commandBindings := []*apps.Binding{}
	if !isConnectedAPI || !isConnectedGDPR {
		commandBindings = append(commandBindings, &apps.Binding{
			Location:    "connect",
			Label:       "connect",
			Description: "Connect to ProductBoard (configure access tokens).",
			Hint:        "[ flags ]",
			Call:        connectCall,
		})
	}
	if isConnectedAPI || isConnectedGDPR {
		commandBindings = append(commandBindings, &apps.Binding{
			Location:    "disconnect",
			Label:       "disconnect",
			Description: "Disconnect from ProductBoard (clear access tokens).",
			Call:        disconnectCall,
		})
	}
	if isConnectedAPI {
		commandBindings = append(commandBindings, &apps.Binding{
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
		})
	}
	if isConnectedGDPR {
		commandBindings = append(commandBindings, &apps.Binding{
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
		})
	}
	if isAdmin {
		commandBindings = append(commandBindings, &apps.Binding{
			Location:    "subscribe",
			Label:       "subscribe",
			Description: "Subscribe to ProductBoard notifications (webhooks).",
			Call:        subscribeCall,
		})
	}

	bindings = append(bindings, apps.Binding{
		Location: apps.LocationCommand,
		Bindings: []*apps.Binding{
			{
				Icon:        appURL(creq, "/static/icon.png"),
				Location:    "pb",
				Label:       "pb",
				Description: "ProductBoard integration.",
				Hint:        "[ connect | create | gdpr ]",
				Bindings:    commandBindings,
			},
		},
	})

	respond(w, bindings, "")
}
