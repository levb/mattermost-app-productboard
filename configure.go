package main

import (
	_ "embed"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/pkg/errors"
)

func configure(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	token, _ := creq.Values["token"].(string)

	asUser := mmclient.AsActingUser(creq.Context)
	err := asUser.StoreOAuth2User(creq.Context.AppID, token)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to store token"))
		return
	}

	stored := ""
	err = asUser.GetOAuth2User(creq.Context.AppID, &stored)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve stored token"))
		return
	}
	respond(w, nil, "updated personal access token to `%s`", last4(stored))
}

func last4(in string) string {
	out := ""
	i := 0
	for ; i < len(in)-4; i++ {
		out += "*"
	}
	for ; i < len(in); i++ {
		out += in[i : i+1]
	}
	return out
}
