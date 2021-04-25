package main

import (
	_ "embed"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
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
