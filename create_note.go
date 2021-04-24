package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

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
	title, _ := creq.Values["title"].(string)
	content, _ := creq.Values["content"].(string)
	token, _ := creq.Context.OAuth2.User.(string)

	permalink := ""
	if creq.Context.PostID != "" {
		u, _ := url.Parse(creq.Context.MattermostSiteURL)
		u.Path = path.Join(u.Path, "_redirect", "pl", creq.Context.PostID)
		permalink = u.String()
	}

	data, err := json.Marshal(&PBCreateNoteRequest{
		Title:      title,
		Content:    content,
		DisplayURL: permalink,
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
	pbReq.Header.Add("Authorization", "Bearer "+token)
	pbReq.Header.Add("Private-Token", token)
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

	noteLink := ""
	if permalink == "" {
		noteLink = fmt.Sprintf("[Note](%s)", pb.Links.HTML)
	} else {
		noteLink = fmt.Sprintf("[Note](%s) from [post](%s)", pb.Links.HTML, permalink)
	}
	message := fmt.Sprintf(
		"Thanks! A new %s has been submitted to ProductBoard for processing by a PM.\n"+
			"Contact a Product Manager if you need an urgent reply.",
		noteLink)

	// Post feedback to the channel, as the acting user
	asUser := mmclient.AsActingUser(creq.Context)
	asUser.CreatePost(&model.Post{
		UserId:    creq.Context.ActingUserID,
		ChannelId: creq.Context.ChannelID,
		RootId:    creq.Context.RootPostID,
		Message:   message,
	})

	respond(w, nil, message)
}
