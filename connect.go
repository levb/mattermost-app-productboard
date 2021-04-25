package main

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
)

func connect(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	accessToken, _ := creq.Values["access_token"].(string)
	gdprToken, _ := creq.Values["gdpr_token"].(string)
	asUser := mmclient.AsActingUser(creq.Context)

	user := User{}
	err := asUser.GetOAuth2User(creq.Context.AppID, &user)
	if errors.Cause(err) == utils.ErrNotFound {
		err = nil
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve previous user record"))
		return
	}

	changed := false
	if accessToken != "" {
		user.AccessToken = accessToken
		changed = true
	}
	if gdprToken != "" {
		user.GDPRToken = gdprToken
		changed = true
	}
	if !changed {
		respond(w, nil, user.asList())
		return
	}

	err = asUser.StoreOAuth2User(creq.Context.AppID, user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to store user record"))
		return
	}

	err = asUser.GetOAuth2User(creq.Context.AppID, &user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve updated user record"))
		return
	}
	respond(w, nil, "Updated:\n%s", user.asList())
}
