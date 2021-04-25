package main

import (
	_ "embed"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
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
	respond(w, nil, "%s:\n%s", "Cleared all API access tokens", user.asList())

}
