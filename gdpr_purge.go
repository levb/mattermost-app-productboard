package main

import (
	_ "embed"
	"net/http"
	"net/url"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/pkg/errors"
)

func gdprPurge(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	email, _ := creq.Values["email"].(string)
	token, _ := creq.Context.OAuth2.User.(string)

	u, _ := url.Parse("https://api.productboard.com/v1/customers/delete_all_data")
	u.Query().Add("email", email)
	pbReq, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	pbReq.Header.Add("Authorization", "Bearer "+token)
	pbReq.Header.Add("Private-Token", token)

	client := &http.Client{}
	resp, err := client.Do(pbReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		respondError(w, resp.StatusCode, errors.Errorf("response status %s", resp.Status))
		return
	}

	respond(w, nil, "Request to delete all data for `%s` has been accepted.", email)
}
