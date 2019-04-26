package api

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

// StreamDonateMutation returns the graphql mutation for donating LINO to a streamer
func StreamDonateMutation() string {
	return `mutation StreamDonate($stream: DonateInput!) {
		donate(stream: $stream) {
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

// SendStreamChatMessageMutation returns the graphql mutation for sending a message to a streamer's chat
func SendStreamChatMessageMutation() string {
	return `mutation SendStreamChatMessage($stream: SendStreamchatMessageInput!) {
		sendStreamchatMessage(stream: $stream) {
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

// EmoteSaveMutation returns the graphql mutation for saving a sticker emote for the logged in user
func EmoteSaveMutation() string {
	return `mutation EmoteSave($stream: SaveEmoteInput!) {
  saveEmote(stream: $stream) {
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

// EmoteDeleteMutation returns the graphql mutation for removing a sticker from the list of saved stickers for the logged in user
func EmoteDeleteMutation() string {
	return `mutation EmoteDelete($stream: DeleteEmoteInput!) {
  deleteEmote(stream: $stream) {
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
