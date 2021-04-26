package main

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func gdprPurge(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	email, _ := creq.Values["email"].(string)
	u, err := userFromContext(creq)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	q := url.Values{}
	q.Add("email", email)
	pbReq, err := http.NewRequest(http.MethodDelete, "https://api.productboard.com/v1/customers/delete_all_data?"+q.Encode(), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	pbReq.Header.Add("Private-Token", u.GDPRToken)

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
