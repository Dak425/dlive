package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dak425/dlive/pkg/api"
)

func main() {
	c := api.Client{
		URL:  api.DefaultURL,
		Auth: "",
	}

	resp, err := c.GlobalInformation()
	//resp, err := c.LivestreamPage("TheHighlord", true, false)
	//resp, err := c.LivestreamChatRoomInfo("TheHighlord", false, "50")
	//resp, err := c.LivestreamProfileVideos("TheHighlord", "10")
	//resp, err := c.LivestreamProfileReplays("TheHighlord", "10")
	//resp, err := c.LivestreamProfileFollowers("TheHighlord", api.SortAlpha, "20", false)
	//resp, err := c.LivestreamProfileFollowing("TheHighlord", "AZ", "20", false)
	//resp, err := c.LivestreamProfileWallet("TheHighlord", "20", false)
	//resp, err := c.TopContributors("TheHighlord", api.ContributionSummaryMonth, "10", false)
	//resp, err := c.StreamChatBannedUsers("TheHighlord", "20", "")
	//resp, err := c.StreamChatModerators("TheHighlord", "20", "")

	// Requires Login (Authentication Token)
	//resp, err := c.AllowedActions("BoneCloset", "TheHighlord")

	if err != nil {
		log.Fatal(err)
	}

	prettyResponse, err := json.MarshalIndent(resp, "", " ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(prettyResponse))
}
