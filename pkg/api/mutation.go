package api

import "net/url"

type LoginWithWalletArgs struct {
	Payload       string `json:"payload"`
	SignedPayload string `json:"signed_payload"`
}

// LoginWithWalletMutation returns the graphql mutation for logging with a LINO wallet account
func LoginWithWalletMutation() string {
	return `mutation LoginWithWallet($payload: String!, $signedPayload: String!) {
		loginWithWallet(payload: $payload, signedPayload: $signedPayload) {
		  ...LoginWithThirdParty
		  __typename
		}
	  }
	  
	  fragment LoginWithThirdParty on LoginResponse {
		me {
		  id
		  private {
			accessToken
			__typename
		  }
		  __typename
		}
		accessToken
		err {
		  code
		  message
		  __typename
		}
		__typename
	  }
	  `
}

// VideoPermLinkMutation returns the graphql mutation for uploading a video to DLive
func VideoPermLinkMutation() string {
	return `mutation VideoPermlink {
		videoPermlinkGenerate {
		  permlink
		  permlinkToken
		  err {
			code
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type FollowUserArgs struct {
	Streamer string `json:"streamer"`
}

// FollowUserMutation returns the graphql mutation for becoming a follower of a streamer
func FollowUserMutation() string {
	return `mutation FollowUser($streamer: String!) {
		follow(streamer: $streamer) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type UnfollowUserArgs FollowUserArgs

// UnfollowUserMutation returns the graphql mutation for unfollowing a streamer
func UnfollowUserMutation() string {
	return `mutation UnfollowUser($streamer: String!) {
		unfollow(streamer: $streamer) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type DonateInput struct {
	Count    int    `json:"count"`
	PermLink string `json:"permlink"`
	Type     string `json:"type"`
}

type StreamDonateArgs struct {
	Input DonateInput `json:"input"`
}

// StreamDonateMutation returns the graphql mutation for donating LINO to a streamer
func StreamDonateMutation() string {
	return `mutation StreamDonate($input: DonateInput!) {
		donate(input: $input) {
		  id
		  recentCount
		  expireDuration
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type SendStreamChatMessageInput struct {
	Message     string `json:"message"`
	RoomRole    string `json:"roomRole"`
	Streamer    string `json:"streamer"`
	Subscribing bool   `json:"subscribing"`
}

type SendStreamChatMessageArgs struct {
	Input SendStreamChatMessageInput `json:"input"`
}

// SendStreamChatMessageMutation returns the graphql mutation for sending a message to a streamer's chat
func SendStreamChatMessageMutation() string {
	return `mutation SendStreamChatMessage($input: SendStreamchatMessageInput!) {
		sendStreamchatMessage(input: $input) {
		  err {
			code
			__typename
		  }
		  message {
			type
			... on ChatText {
			  id
			  content
			  ...VStreamChatSenderInfoFrag
			  __typename
			}
			__typename
		  }
		  __typename
		}
	  }
	  
	  fragment VStreamChatSenderInfoFrag on SenderInfo {
		subscribing
		role
		roomRole
		sender {
		  id
		  username
		  displayname
		  avatar
		  partnerStatus
		  __typename
		}
		__typename
	  }
	  `
}

type SetAllowStickerArgs struct {
	Allow bool `json:"allow"`
}

// SetAllowStickerMutation returns the graphql mutation for enabling or disabling stickers in a streamer's chat
func SetAllowStickerMutation() string {
	return `mutation SetAllowSticker($allow: Boolean!) {
		allowEmoteSet(allow: $allow) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type SetChatIntervalArgs struct {
	Seconds int `json:"seconds"`
}

// SetChatIntervalMutation returns the graphql mutation for setting the how often viewers can send messages to a streamer's chat
func SetChatIntervalMutation() string {
	return `mutation SetChatInterval($seconds: Int!) {
		chatIntervalSet(seconds: $seconds) {
		  err {
			code
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type DeleteChatArgs struct {
	Streamer string `json:"streamer"`
	ID       string `json:"id"`
}

// DeleteChatMutation returns the graphql mutation for deleting a message from a streamer's chat
func DeleteChatMutation() string {
	return `mutation DeleteChat($streamer: String!, $id: String!) {
		chatDelete(streamer: $streamer, id: $id) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type AddModeratorArgs struct {
	Username string `json:"username"`
}

// AddModeratorMutation returns the graphql mutation for setting a user in a streamer's chat as a moderator
func AddModeratorMutation() string {
	return `mutation AddModerator($username: String!) {
		moderatorAdd(username: $username) {
		  err {
			code
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type RemoveModeratorArgs AddModeratorArgs

// RemoveModeratorMutation returns the graphql mutation for removing a user as a moderator in the given streamer's chat
func RemoveModeratorMutation() string {
	return `mutation RemoveModerator($username: String!) {
		moderatorRemove(username: $username) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type UnbanStreamChatUserArgs struct {
	Streamer string `json:"streamer"`
	Username string `json:"username"`
}

// UnbanStreamChatUserMutation returns the graphql mutation for unbanning a user in a streamer's chat
func UnbanStreamChatUserMutation() string {
	return `mutation UnbanStreamChatUser($streamer: String!, $username: String!) {
		streamchatUserUnban(streamer: $streamer, username: $username) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type BanStreamChatUserArgs UnbanStreamChatUserArgs

// BanStreamChatUserMutation returns the graphql mutation for banning a user from a streamer's chat
// Couldn't get in dev console, is based on UnbanStreamChatUserMutation
func BanStreamChatUserMutation() string {
	return `mutation BanStreamChatUser($streamer: String!, $username: String!) {
		streamchatUserBan(streamer: $streamer, username: $username) {
		  err {
			code
			message
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type SetStreamTemplateInput struct {
	AgeRestriction bool    `json:"ageRestriction"`
	CategoryID     int     `json:"categoryID"`
	DisableAlert   bool    `json:"disableAlert"`
	LanguageID     int     `json:"languageID"`
	ThumbnailURL   url.URL `json:"thumbnailUrl"`
	Title          string  `json:"title"`
}

type SetStreamTemplateArgs struct {
	Template SetStreamTemplateInput `json:"template"`
}

// SetStreamTemplateMutation returns the graphql mutation for saving a user's stream metadata
func SetStreamTemplateMutation() string {
	return `mutation SetStreamTemplate($template: SetStreamTemplateInput!) {
		streamTemplateSet(template: $template) {
		  err {
			code
			__typename
		  }
		  __typename
		}
	  }
	  `
}

// GenerateStreamKeyMutation returns the graphql mutation for generating the key needed to stream data to a livestream profile
func GenerateStreamKeyMutation() string {
	return `mutation generateStreamKey {
		streamKeyGenerate {
		  url
		  key
		  err {
			code
			__typename
		  }
		  __typename
		}
	  }
	  `
}

type SaveEmoteInput struct {
	Level   string `json:"level"`
	MyLevel string `json:"myLevel"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

type EmoteSaveArgs struct {
	Input SaveEmoteInput `json:"input"`
}

// EmoteSaveMutation returns the graphql mutation for saving a sticker emote for the logged in user
func EmoteSaveMutation() string {
	return `mutation EmoteSave($input: SaveEmoteInput!) {
  saveEmote(input: $input) {
    emote {
      name
      username
      sourceURL
      mimeType
      level
      type
      __typename
    }
    err {
      code
      message
      __typename
    }
    __typename
  }
}
`
}

type DeleteEmoteInput struct {
	Level string `json:"level"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type EmoteDeleteArgs struct {
	Input DeleteEmoteInput `json:"input"`
}

// EmoteDeleteMutation returns the graphql mutation for removing a sticker from the list of saved stickers for the logged in user
func EmoteDeleteMutation() string {
	return `mutation EmoteDelete($input: DeleteEmoteInput!) {
  deleteEmote(input: $input) {
    err {
      code
      message
      __typename
    }
    __typename
  }
}
`
}

type DeletePastBroadcastArgs struct {
	Permlink string `json:"permlink"`
}

// DeleteLastBroadcastMutation returns the graphql mutation for deleting a stream replay
func DeleteLastBroadcastMutation() string {
	return `mutation DeletePastbroadcast($permlink: String!) {
		  pastbroadcastDelete(permlink: $permlink) {
		    err {
		      code
		      __typename
		    }
		    __typename
		  }
		}
		`
}
