package main

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/pkg/errors"
)

func disconnect(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	asUser := mmclient.AsActingUser(creq.Context)

	err := asUser.StoreOAuth2User(creq.Context.AppID, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to clear the user tokens"))
		return
	}

	user := User{}
	err = asUser.GetOAuth2User(creq.Context.AppID, &user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve updated user record"))
		return
	}

	respond(w, nil, "Cleared all API access tokens:\n%s", user.asList())
}
