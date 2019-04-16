package api

// StreamMessageSubscription gives the graphql query to establish a subscription for a livestream's chat messages
func StreamMessageSubscription() string {
	return `subscription StreamMessageSubscription($streamer: String!) {
		streamMessageReceived(streamer: $streamer) {
		  type
		  ... on ChatGift {
			id
			gift
			amount
			recentCount
			expireDuration
			...VStreamChatSenderInfoFrag
		  }
		  ... on ChatHost {
			id
			viewer
			...VStreamChatSenderInfoFrag
		  }
		  ... on ChatSubscription {
			id
			month
			...VStreamChatSenderInfoFrag
		  }
		  ... on ChatChangeMode {
			mode
		  }
		  ... on ChatText {
			id
			content
			...VStreamChatSenderInfoFrag
		  }
		  ... on ChatFollow {
			id
			...VStreamChatSenderInfoFrag
		  }
		  ... on ChatDelete {
			ids
		  }
		  ... on ChatBan {
			id
			...VStreamChatSenderInfoFrag
		  }
		  ... on ChatModerator {
			id
			...VStreamChatSenderInfoFrag
			add
		  }
		  ... on ChatEmoteAdd {
			id
			...VStreamChatSenderInfoFrag
			emote
		  }
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
		}
	  }`
}
