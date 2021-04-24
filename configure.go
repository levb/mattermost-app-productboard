package main

import (
	_ "embed"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
)

type User struct {
	AccessToken string
	GDPRToken   string
}

func (u User) asList() string {
	out := ""
	if u.AccessToken != "" {
		out += "- Personal Access Token: `" + last4(u.AccessToken) + "`\n"
	}
	if u.GDPRToken != "" {
		out += "- GDPR Public API Token: `" + last4(u.GDPRToken) + "`\n"
	}
	if out == "" {
		return "empty"
	}
	return out
}

func userFromContext(creq *apps.CallRequest) (*User, error) {
	m, _ := creq.Context.OAuth2.User.(map[string]interface{})
	if len(m) == 0 {
		return nil, errors.Errorf("empty or wrong type %T for User record in Context", creq.Context.OAuth2.User)
	}

	u := User{}
	for k, v := range m {
		switch k {
		case "AccessToken":
			u.AccessToken = v.(string)
		case "GDPRToken":
			u.GDPRToken = v.(string)
		default:
			return nil, errors.Errorf("unrecognized key: %s", k)
		}
	}
	return &u, nil
}

func configure(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	accessToken, _ := creq.Values["access_token"].(string)
	gdprToken, _ := creq.Values["gdpr_token"].(string)
	clear, _ := creq.Values["clear"].(bool)
	asUser := mmclient.AsActingUser(creq.Context)

	u := User{}
	var err error
	action := ""
	defer func() {
		if err != nil {
			return
		}
		err = asUser.GetOAuth2User(creq.Context.AppID, &u)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve updated user record"))
			return
		}
		respond(w, nil, "%s:\n%s", action, u.asList())
	}()

	if clear {
		if accessToken != "" || gdprToken != "" {
			respondError(w, http.StatusBadRequest, errors.New("can not specify --clear and any other options in the same command"))
			return
		}
		err := asUser.StoreOAuth2User(creq.Context.AppID, u)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to clear the user tokens"))
			return
		}
		action = "Cleared all API access tokens"
		return
	}

	err = asUser.GetOAuth2User(creq.Context.AppID, &u)
	if errors.Cause(err) == utils.ErrNotFound {
		err = nil
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve previous user record"))
		return
	}

	changed := false
	if accessToken != "" {
		u.AccessToken = accessToken
		changed = true
	}
	if gdprToken != "" {
		u.GDPRToken = gdprToken
		changed = true
	}
	if !changed {
		action = "No change"
		return
	}

	err = asUser.StoreOAuth2User(creq.Context.AppID, u)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to store user record"))
		return
	}
	action = "Updated"
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
