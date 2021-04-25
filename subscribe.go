package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func subscribe(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	webhookURL := creq.Context.MattermostSiteURL + creq.Context.AppPath + apps.PathWebhook + "/test" +
		"?secret=" + creq.Context.App.WebhookSecret

	respond(w, nil, "ProductBoard Slack webhooks don't seem to work for non slack.com URLs. If there was a way to direct them here, the webhook URL would be:\n`%s`", webhookURL)
}
