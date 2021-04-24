package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

type PBCreateNoteRequest struct {
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	CustomerEmail string   `json:"customer_email,omitempty"`
	DisplayURL    string   `json:"display_url,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

type PBCreateNoteResponse struct {
	Links struct {
		HTML string `json:"html"`
	} `json:"links"`
}

func createNote(w http.ResponseWriter, req *http.Request, creq *apps.CallRequest) {
	interactive, _ := creq.Values["interactive"].(string)
	title, _ := creq.Values["title"].(string)
	content, _ := creq.Values["content"].(string)
	email, _ := creq.Values["email"].(string)
	tagStr, _ := creq.Values["tags"].(string)
	tags := strings.Split(tagStr, ",")
	if len(tags) == 0 || len(tags) == 1 && tags[0] == "" {
		tags = nil
	}

	if creq.State == "post-menu" || interactive == "true" || title == "" || content == "" {
		createNoteForm(w, req, creq)
		return
	}

	user, err := userFromContext(creq)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(&PBCreateNoteRequest{
		Title:         title,
		Content:       content,
		DisplayURL:    permalink(creq),
		CustomerEmail: email,
		Tags:          tags,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	pbReq, err := http.NewRequest("POST", "https://api.productboard.com/notes", bytes.NewReader(data))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	pbReq.Header.Add("Authorization", "Bearer "+user.AccessToken)
	pbReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(pbReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respondError(w, resp.StatusCode, errors.Errorf("response status %s", resp.Status))
		return
	}

	pb := PBCreateNoteResponse{}
	err = json.NewDecoder(resp.Body).Decode(&pb)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Post feedback to the channel, as the acting user
	asUser := mmclient.AsActingUser(creq.Context)
	asUser.CreatePost(&model.Post{
		UserId:    creq.Context.ActingUserID,
		ChannelId: creq.Context.ChannelID,
		RootId:    creq.Context.RootPostID,
		Message: fmt.Sprintf(
			"[%s](%s) has been submitted to ProductBoard for processing by a Product Manager.",
			title, pb.Links.HTML),
	})

	respond(w, nil, "[%s](%s) created.", title, pb.Links.HTML)
}
