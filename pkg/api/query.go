package api

// GlobalInformationQuery returns the graphql query string for retrieving global information about Dlive
func GlobalInformationQuery() string {
	return `query GlobalInformation {
		globalInfo {
			languages {
				id
				backendID
				language
				code
				__typename
			}
			__typename
		}
	}
	`
}

// MeGlobalQuery returns the graphql query for retrieving about the current user
func MeGlobalQuery() string {
	return `query MeGlobal {
		me {
			...MeGlobalFrag
			__typename
		}
	}
	
	fragment MeGlobalFrag on User {
		id
		username
		...VDliveAvatarFrag
		displayname
		partnerStatus
		role
		private {
			accessToken
			insecure
			email
			phone
			nextDisplayNameChangeTime
			language
			showSubSettingTab
			__typename
		}
		...SettingsSubscribeFrag
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment SettingsSubscribeFrag on User {
		id
		subSetting {
			badgeColor
			badgeText
			textColor
			__typename
		}
		__typename
	}
	`
}

type MeDashboardArgs struct {
	IsLoggedIn bool `json:"isLoggedIn"`
}

// MeDashboardQuery gives the graphql query to obtain information about the authenticated user's dashboard (settings, stats, chatroom)
func MeDashboardQuery() string {
	return `query MeDashboard($isLoggedIn: Boolean!) {
		me {
			...MeDashboardFrag
			__typename
		}
	}
	
	fragment MeDashboardFrag on User {
		id
		...DashboardStreamSettingsFrag
		...DashboardHostSettingFrag
		...DashboardStatsFrag
		...DashboardStreamChatroomFrag
		__typename
	}
	
	fragment DashboardStreamSettingsFrag on User {
		livestream {
			id
			permlink
			...VVideoPlayerFrag
			__typename
		}
		hostingLivestream {
			id
			permlink
			creator {
				username
				...VDliveAvatarFrag
				...VDliveNameFrag
				__typename
			}
			...VVideoPlayerFrag
			__typename
		}
		private {
			streamTemplate {
				title
				ageRestriction
				thumbnailUrl
				disableAlert
				category {
					id
					backendID
					title
					__typename
				}
				language {
					id
					backendID
					code
					language
					__typename
				}
				__typename
			}
			filterWords
			__typename
		}
		__typename
	}
	
	fragment VVideoPlayerFrag on Livestream {
		disableAlert
		category {
			id
			title
			__typename
		}
		language {
			language
			__typename
		}
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment DashboardHostSettingFrag on User {
		id
		hostingLivestream {
			creator {
				username
				...VDliveAvatarFrag
				...VDliveNameFrag
				__typename
			}
			__typename
		}
		__typename
	}
	
	fragment DashboardStatsFrag on User {
		id
		livestream {
			watchingCount
			totalReward
			__typename
		}
		followers {
			totalCount
			__typename
		}
		private {
			subscribers {
				totalCount
				__typename
			}
			__typename
		}
		wallet {
			totalEarning
			__typename
		}
		__typename
	}
	
	fragment DashboardStreamChatroomFrag on User {
		...MeLivestreamChatroomFrag
		__typename
	}
	
	fragment MeLivestreamChatroomFrag on User {
		id
		username
		role
		...MeEmoteFrag
		__typename
	}
	
	fragment MeEmoteFrag on User {
		id
		role @include(if: $isLoggedIn)
		emote {
			...EmoteMineFrag
			...EmoteChannelFrag
			__typename
		}
		__typename
	}
	
	fragment EmoteMineFrag on AllEmotes {
		mine {
			list {
				name
				username
				sourceURL
				mimeType
				level
				type
				__typename
			}
			__typename
		}
		__typename
	}
	
	fragment EmoteChannelFrag on AllEmotes {
		channel {
			list {
				name
				username
				sourceURL
				mimeType
				level
				type
				__typename
			}
			__typename
		}
		__typename
	}
	`
}

type MeLivestreamArgs struct {
	IsLoggedIn bool `json:"isLoggedIn"`
}

// MeLivestreamQuery provides the graphql query for obtaining data related to the current user's livestream
func MeLivestreamQuery() string {
	return `query MeLivestream($isLoggedIn: Boolean!) {
		me {
			...MeLivestreamFrag
			__typename
		}
	}
	
	fragment MeLivestreamFrag on User {
		id
		...MeLivestreamChatroomFrag
		__typename
	}
	
	fragment MeLivestreamChatroomFrag on User {
		id
		username
		role
		...MeEmoteFrag
		__typename
	}
	
	fragment MeEmoteFrag on User {
		id
		role @include(if: $isLoggedIn)
		emote {
			...EmoteMineFrag
			...EmoteChannelFrag
			__typename
		}
		__typename
	}
	
	fragment EmoteMineFrag on AllEmotes {
		mine {
			list {
				name
				username
				sourceURL
				mimeType
				level
				type
				__typename
			}
			__typename
		}
		__typename
	}
	
	fragment EmoteChannelFrag on AllEmotes {
		channel {
			list {
				name
				username
				sourceURL
				mimeType
				level
				type
				__typename
			}
			__typename
		}
		__typename
	}
	`
}

// MeBalanceQuery gives the graphql query to obtain information about authenticated user's balance
func MeBalanceQuery() string {
	return `query MeBalance {
		me {
			...MeBalanceFrag
			__typename
		}
	}
	
	fragment MeBalanceFrag on User {
		id
		wallet {
			balance
			__typename
		}
		__typename
	}
	`
}

type MeSubscribingArgs struct {
	First int    `json:"first"`
	After string `json:"after"`
}

// MeSubscribingQuery returns the graphql query to get the list of users the currently authenticated user is subbed to
func MeSubscribingQuery() string {
	return `query MeSubscribing($first: Int!, $after: String) {
		me {
			...MeSubscribingFrag
			__typename
		}
	}
	
	fragment MeSubscribingFrag on User {
		id
		private {
			subscribing(first: $first, after: $after) {
				totalCount
				pageInfo {
					startCursor
					endCursor
					hasNextPage
					hasPreviousPage
					__typename
				}
				list {
					streamer {
						username
						displayname
						avatar
						partnerStatus
						__typename
					}
					tier
					status
					lastBilledDate
					subscribedAt
					month
					__typename
				}
				__typename
			}
			__typename
		}
		__typename
	}
	`
}

// MePartnerProgressQuery returns the graphql query to get information about the currently authenticated user's partner progress
func MePartnerProgressQuery() string {
	return `query MePartnerProgress {
		me {
			...MePartnerProgressFrag
			__typename
		}
	}
	
	fragment MePartnerProgressFrag on User {
		id
		followers {
			totalCount
			__typename
		}
		private {
			previousStats {
				partnerStats {
					streamingHours
					streamingDays
					donationReceived
					__typename
				}
				contentBonus
				__typename
			}
			partnerProgress {
				partnerStatus
				current {
					followerCount
					streamingHours
					streamingDays
					donationReceived
					lockPoint
					__typename
				}
				target {
					followerCount
					streamingHours
					streamingDays
					donationReceived
					lockPoint
					__typename
				}
				eligible
				__typename
			}
			__typename
		}
		__typename
	}
	`
}

type LivestreamPageArgs struct {
	DisplayName string `json:"displayname"`
	Add         bool   `json:"add"`
	IsLoggedIn  bool   `json:"isLoggedIn"`
}

// LivestreamPageQuery gives the graphql query for obtaining data about a user's livestream
func LivestreamPageQuery() string {
	return `query LivestreamPage($displayname: String!, $add: Boolean!, $isLoggedIn: Boolean!) {
		userByDisplayName(displayname: $displayname) {
			id
			...VDliveAvatarFrag
			...VDliveNameFrag
			...VFollowFrag
			...VSubscriptionFrag
			banStatus
			about
			avatar
			myRoomRole @include(if: $isLoggedIn)
			isMe @include(if: $isLoggedIn)
			isSubscribing @include(if: $isLoggedIn)
			livestream {
				id
				permlink
				watchTime(add: $add)
				...LivestreamInfoFrag
				...VVideoPlayerFrag
				__typename
			}
			hostingLivestream {
				id
				creator {
					...VDliveAvatarFrag
					displayname
					username
					__typename
				}
				...VVideoPlayerFrag
				__typename
			}
			...LivestreamProfileFrag
			__typename
		}
	}
	
	fragment LivestreamInfoFrag on Livestream {
		category {
			title
			imgUrl
			id
			backendID
			__typename
		}
		title
		watchingCount
		totalReward
		...VDonationGiftFrag
		...VPostInfoShareFrag
		__typename
	}
	
	fragment VDonationGiftFrag on Post {
		permlink
		creator {
			username
			__typename
		}
		__typename
	}
	
	fragment VPostInfoShareFrag on Post {
		permlink
		title
		content
		category {
			id
			backendID
			title
			__typename
		}
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment LivestreamProfileFrag on User {
		isMe @include(if: $isLoggedIn)
		canSubscribe
		private @include(if: $isLoggedIn) {
			subscribers {
				totalCount
				__typename
			}
			__typename
		}
		videos {
			totalCount
			__typename
		}
		pastBroadcasts {
			totalCount
			__typename
		}
		followers {
			totalCount
			__typename
		}
		following {
			totalCount
			__typename
		}
		...ProfileAboutFrag
		__typename
	}
	
	fragment ProfileAboutFrag on User {
		id
		about
		__typename
	}
	
	fragment VVideoPlayerFrag on Livestream {
		disableAlert
		category {
			id
			title
			__typename
		}
		language {
			language
			__typename
		}
		__typename
	}
	
	fragment VFollowFrag on User {
		id
		username
		displayname
		isFollowing @include(if: $isLoggedIn)
		isMe @include(if: $isLoggedIn)
		followers {
			totalCount
			__typename
		}
		__typename
	}
	
	fragment VSubscriptionFrag on User {
		id
		username
		displayname
		isSubscribing @include(if: $isLoggedIn)
		canSubscribe
		isMe @include(if: $isLoggedIn)
		__typename
	}
	`
}

type LivestreamPageRefetchArgs LivestreamPageArgs

// LivestreamPageRefetchQuery gives the graphql query refreshing data about a specific streamer's page
func LivestreamPageRefetchQuery() string {
	return `query LivestreamPageRefetch($displayname: String!, $add: Boolean!, $isLoggedIn: Boolean!) {
		userByDisplayName(displayname: $displayname) {
			id
			username
			myRoomRole @include(if: $isLoggedIn)
			isFollowing @include(if: $isLoggedIn)
			isSubscribing @include(if: $isLoggedIn)
			livestream {
				id
				permlink
				watchTime(add: $add)
				...LivestreamInfoFrag
				...VVideoPlayerFrag
				__typename
			}
			hostingLivestream {
				id
				permlink
				creator {
					...VDliveAvatarFrag
					displayname
					username
					__typename
				}
				...VVideoPlayerFrag
				__typename
			}
			__typename
		}
	}
	
	fragment LivestreamInfoFrag on Livestream {
		category {
			title
			imgUrl
			id
			backendID
			__typename
		}
		title
		watchingCount
		totalReward
		...VDonationGiftFrag
		...VPostInfoShareFrag
		__typename
	}
	
	fragment VDonationGiftFrag on Post {
		permlink
		creator {
			username
			__typename
		}
		__typename
	}
	
	fragment VPostInfoShareFrag on Post {
		permlink
		title
		content
		category {
			id
			backendID
			title
			__typename
		}
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VVideoPlayerFrag on Livestream {
		disableAlert
		category {
			id
			title
			__typename
		}
		language {
			language
			__typename
		}
		__typename
	}
	`
}

type LivestreamChatRoomInfoArgs struct {
	DisplayName string `json:"displayname"`
	IsLoggedIn  bool   `json:"isLoggedIn"`
	Limit       int    `json:"limit"`
}

// LivestreamChatRoomInfoQuery gives the graphql query to get data related to the chat of a livestream page
func LivestreamChatRoomInfoQuery() string {
	return `query LivestreamChatroomInfo($displayname: String!, $isLoggedIn: Boolean!, $limit: Int!) {
		userByDisplayName(displayname: $displayname) {
			id
			...VLivestreamChatroomFrag
			__typename
		}
	}
	
	fragment VLivestreamChatroomFrag on User {
		id
		isFollowing @include(if: $isLoggedIn)
		role @include(if: $isLoggedIn)
		myRoomRole @include(if: $isLoggedIn)
		isSubscribing @include(if: $isLoggedIn)
		...VStreamChatroomHeaderFrag
		...VStreamChatroomListFrag
		...StreamChatroomInputFrag
		chats(count: 50) {
			type
			... on ChatGift {
				id
				gift
				amount
				...VStreamChatSenderInfoFrag
				__typename
			}
			... on ChatHost {
				id
				viewer
				...VStreamChatSenderInfoFrag
				__typename
			}
			... on ChatSubscription {
				id
				month
				...VStreamChatSenderInfoFrag
				__typename
			}
			... on ChatText {
				id
				content
				...VStreamChatSenderInfoFrag
				__typename
			}
			... on ChatModerator {
				id
				add
				...VStreamChatSenderInfoFrag
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
			... on ChatFollow {
				id
				...VStreamChatSenderInfoFrag
				__typename
			}
			... on ChatEmoteAdd {
				id
				emote
				...VStreamChatSenderInfoFrag
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
			__typename
		}
		__typename
	}
	
	fragment VStreamChatroomHeaderFrag on User {
		id
		username
		displayname
		livestream {
			id
			permlink
			__typename
		}
		...VTopContributorsFrag
		__typename
	}
	
	fragment VTopContributorsFrag on User {
		id
		displayname
		livestream {
			id
			__typename
		}
		__typename
	}
	
	fragment VStreamChatroomListFrag on User {
		...VStreamChatRowStreamerFrag
		...PinnedGiftsFrag
		__typename
	}
	
	fragment VStreamChatRowStreamerFrag on User {
		displayname
		...VStreamChatRowSenderInfoStreamerFrag
		...VStreamChatProfileCardStreamerFrag
		...StreamChatTextRowStreamerFrag
		__typename
	}
	
	fragment VStreamChatRowSenderInfoStreamerFrag on User {
		id
		subSetting {
			badgeText
			badgeColor
			textColor
			__typename
		}
		__typename
	}
	
	fragment VStreamChatProfileCardStreamerFrag on User {
		id
		username
		myRoomRole @include(if: $isLoggedIn)
		role
		__typename
	}
	
	fragment StreamChatTextRowStreamerFrag on User {
		id
		username
		myRoomRole @include(if: $isLoggedIn)
		emote @include(if: $isLoggedIn) {
			channel {
				list {
					name
					username
					sourceURL
					mimeType
					level
					type
					__typename
				}
				__typename
			}
			__typename
		}
		__typename
	}
	
	fragment PinnedGiftsFrag on User {
		id
		recentDonations(limit: $limit) {
			user {
				...VDliveAvatarFrag
				...VDliveNameFrag
				__typename
			}
			...PinnedGiftItemFrag
			__typename
		}
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment PinnedGiftItemFrag on DonationBlock {
		user {
			id
			username
			displayname
			...VDliveAvatarFrag
			...VDliveNameFrag
			__typename
		}
		count
		type
		updatedAt
		expiresAt
		expirationTime
		__typename
	}
	
	fragment StreamChatroomInputFrag on User {
		chatMode
		chatInterval
		myRoomRole @include(if: $isLoggedIn)
		livestream {
			permlink
			creator {
				username
				__typename
			}
			__typename
		}
		...StreamChatMemberManageTabFrag
		...StreamChatModeSettingsFrag
		...EmoteBoardStreamerFrag
		__typename
	}
	
	fragment StreamChatMemberManageTabFrag on User {
		id
		username
		displayname
		myRoomRole @include(if: $isLoggedIn)
		__typename
	}
	
	fragment StreamChatModeSettingsFrag on User {
		id
		chatMode
		allowEmote
		chatInterval
		__typename
	}
	
	fragment EmoteBoardStreamerFrag on User {
		id
		username
		partnerStatus
		myRoomRole @include(if: $isLoggedIn)
		emote @include(if: $isLoggedIn) {
			channel {
				list {
					name
					username
					sourceURL
					mimeType
					level
					type
					__typename
				}
				__typename
			}
			__typename
		}
		__typename
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

type LivestreamLanguagesArgs struct {
	CategoryID int `json:"categoryID"`
}

// LivestreamLanguagesQuery returns the graphql query to get data about the available languages to set a stream to
func LivestreamLanguagesQuery() string {
	return `query LivestreamsLanguages($categoryID: Int) {
		languages(categoryID: $categoryID) {
			...LanguageFrag
			__typename
		}
	}
	
	fragment LanguageFrag on Language {
		id
		backendID
		language
		__typename
	}
	`
}

type LivestreamProfileVideoArgs struct {
	DisplayName string `json:"displayname"`
	SortedBy    string `json:"sortedBy"`
	First       int    `json:"first"`
	After       string `json:"after"`
}

// LivestreamProfileVideoQuery returns the graphql query for getting the videos of a specified streamer
func LivestreamProfileVideoQuery() string {
	return `query LivestreamProfileVideo($displayname: String!, $sortedBy: VideoSortOrder, $first: Int, $after: String) {
		userByDisplayName(displayname: $displayname) {
			id
			videos(sortedBy: $sortedBy, first: $first, after: $after) {
				pageInfo {
					endCursor
					hasNextPage
					__typename
				}
				list {
					...ProfileVideoSnapFrag
					__typename
				}
				__typename
			}
			username
			__typename
		}
	}
	
	fragment ProfileVideoSnapFrag on Video {
		permlink
		thumbnailUrl
		title
		totalReward
		createdAt
		viewCount
		length
		creator {
			displayname
			__typename
		}
		__typename
	}
	`
}

type LivestreamProfileReplayArgs struct {
	DisplayName string `json:"displayname"`
	First       int    `json:"first"`
	After       string `json:"after"`
}

// LivestreamProfileReplayQuery returns the graphql query for getting the replays of a specified streamer
func LivestreamProfileReplayQuery() string {
	return `query LivestreamProfileReplay($displayname: String!, $first: Int, $after: String) {
		userByDisplayName(displayname: $displayname) {
			id
			pastBroadcasts(first: $first, after: $after) {
				pageInfo {
					endCursor
					hasNextPage
					__typename
				}
				list {
					...ProfileReplaySnapFrag
					__typename
				}
				__typename
			}
			username
			__typename
		}
	}
	
	fragment ProfileReplaySnapFrag on PastBroadcast {
		permlink
		thumbnailUrl
		title
		totalReward
		createdAt
		viewCount
		playbackUrl
		creator {
			displayname
			__typename
		}
		__typename
	}
	`
}

type LivestreamProfileWalletArgs struct {
	DisplayName string `json:"displayname"`
	First       int    `json:"first"`
	After       string `json:"after"`
	IsLoggedIn  bool   `json:"isLoggedIn"`
}

// LivestreamProfileWalletQuery returns the graphql query for getting balance, earnings, and transactions of a streamer
func LivestreamProfileWalletQuery() string {
	return `query LivestreamProfileWallet($displayname: String!, $first: Int, $after: String, $isLoggedIn: Boolean!) {
		userByDisplayName(displayname: $displayname) {
			id
			username
			displayname
			isMe @include(if: $isLoggedIn)
			wallet {
				balance
				totalEarning
				__typename
			}
			transactions(first: $first, after: $after) {
				totalCount
				pageInfo {
					startCursor
					endCursor
					hasNextPage
					hasPreviousPage
					__typename
				}
				list {
					seq
					txType
					createdAt
					description
					amount
					balance
					__typename
				}
				__typename
			}
			__typename
		}
	}
	`
}

type LivestreamProfileFollowersArgs struct {
	DisplayName string `json:"displayname"`
	SortedBy    string `json:"sortedBy"`
	First       int    `json:"first"`
	After       string `json:"after"`
	IsLoggedIn  bool   `json:"isLoggedIn"`
}

// LivestreamProfileFollowersQuery returns the graphql query for getting a streamer's followers
func LivestreamProfileFollowersQuery() string {
	return `query LivestreamProfileFollowers($displayname: String!, $sortedBy: RelationSortOrder, $first: Int, $after: String, $isLoggedIn: Boolean!) {
		userByDisplayName(displayname: $displayname) {
			id
			displayname
			followers(sortedBy: $sortedBy, first: $first, after: $after) {
				pageInfo {
					endCursor
					hasNextPage
					__typename
				}
				list {
					...VDliveAvatarFrag
					...VDliveNameFrag
					...VFollowFrag
					__typename
				}
				__typename
			}
			__typename
		}
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment VFollowFrag on User {
		id
		username
		displayname
		isFollowing @include(if: $isLoggedIn)
		isMe @include(if: $isLoggedIn)
		followers {
			totalCount
			__typename
		}
		__typename
	}
	`
}

type LivestreamProfileFollowingArgs LivestreamProfileFollowersArgs

// LivestreamProfileFollowingQuery returns the graphql query for getting the users a streamer follows
func LivestreamProfileFollowingQuery() string {
	return `query LivestreamProfileFollowing($displayname: String!, $sortedBy: RelationSortOrder, $first: Int, $after: String, $isLoggedIn: Boolean!) {
		userByDisplayName(displayname: $displayname) {
			id
			displayname
			following(sortedBy: $sortedBy, first: $first, after: $after) {
				pageInfo {
					endCursor
					hasNextPage
					__typename
				}
				list {
					...VDliveAvatarFrag
					...VDliveNameFrag
					...VFollowFrag
					__typename
				}
				__typename
			}
			__typename
		}
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment VFollowFrag on User {
		id
		username
		displayname
		isFollowing @include(if: $isLoggedIn)
		isMe @include(if: $isLoggedIn)
		followers {
			totalCount
			__typename
		}
		__typename
	}
	`
}

type TopContributorsArgs struct {
	DisplayName string `json:"displayname"`
	Rule        string `json:"rule"`
	First       int    `json:"first"`
	After       string `json:"after"`
	QueryStream bool   `json:"queryStream"`
}

// TopContributorsQuery gives the graphql query to get data about users who are the top contributors for a livestream page
func TopContributorsQuery() string {
	return `query TopContributors($displayname: String!, $rule: ContributionSummaryRule, $first: Int, $after: String, $queryStream: Boolean!) {
		userByDisplayName(displayname: $displayname) {
			id
			...TopContributorsOfStreamerFrag @skip(if: $queryStream)
			livestream @include(if: $queryStream) {
				...TopContributorsOfLivestreamFrag
				__typename
			}
			__typename
		}
	}
	
	fragment TopContributorsOfStreamerFrag on User {
		id
		topContributions(rule: $rule, first: $first, after: $after) {
			pageInfo {
				endCursor
				hasNextPage
				__typename
			}
			list {
				amount
				contributor {
					id
					...VDliveNameFrag
					...VDliveAvatarFrag
					__typename
				}
				__typename
			}
			__typename
		}
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment TopContributorsOfLivestreamFrag on Livestream {
		id
		topContributions(first: $first, after: $after) {
			pageInfo {
				endCursor
				hasNextPage
				__typename
			}
			list {
				amount
				contributor {
					id
					...VDliveNameFrag
					...VDliveAvatarFrag
					__typename
				}
				__typename
			}
			__typename
		}
		__typename
	}
	`
}

type HomePageLivestreamArgs struct {
	First            int    `json:"first"`
	After            string `json:"after"`
	LanguageID       int    `json:"languageID"`
	CategoryID       int    `json:"categoryID"`
	ShowNSFW         bool   `json:"showNSFW"`
	UserLanguageCode string `json:"user_language_code"`
}

// HomePageLivestreamQuery gives the graphql query to get data about the live streams that would be shown on the homepage
func HomePageLivestreamQuery() string {
	return `query HomePageLivestream($first: Int, $after: String, $languageID: Int, $categoryID: Int, $showNSFW: Boolean, $userLanguageCode: String) {
		livestreams(stream: {first: $first, after: $after, languageID: $languageID, categoryID: $categoryID, showNSFW: $showNSFW, order: TRENDING, userLanguageCode: $userLanguageCode}) {
			...VCategoryLivestreamFrag
			__typename
		}
	}
	
	fragment VCategoryLivestreamFrag on LivestreamConnection {
		pageInfo {
			endCursor
			hasNextPage
			__typename
		}
		list {
			permlink
			ageRestriction
			...VLivestreamSnapFrag
			__typename
		}
		__typename
	}
	
	fragment VLivestreamSnapFrag on Livestream {
		id
		creator {
			username
			displayname
			...VDliveAvatarFrag
			...VDliveNameFrag
			__typename
		}
		title
		totalReward
		watchingCount
		thumbnailUrl
		lastUpdatedAt
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

// HomePageLeaderboardQuery gives the graphql query to get data about the streamers with the biggest gains in LINO
func HomePageLeaderboardQuery() string {
	return `query HomePageLeaderboard {
		leaderboard {
			...LeaderboardFrag
			__typename
		}
	}
	
	fragment LeaderboardFrag on LeaderboardConnection {
		list {
			user {
				displayname
				wallet {
					lastDayEarning
					__typename
				}
				...VDliveAvatarFrag
				...VDliveNameFrag
				__typename
			}
			change
			__typename
		}
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

type HomePageCategoriesArgs struct {
	First      int    `json:"first"`
	After      string `json:"after"`
	LanguageID int    `json:"languageID"`
}

// HomePageCategoriesQuery gives the graphql query to get data about the available stream categories
func HomePageCategoriesQuery() string {
	return `query HomePageCategories($first: Int, $after: String, $languageID: Int) {
		categories(stream: {first: $first, after: $after, languageID: $languageID}) {
			...HomeCategoriesFrag
			__typename
		}
	}
	
	fragment HomeCategoriesFrag on CategoryConnection {
		pageInfo {
			endCursor
			hasNextPage
			__typename
		}
		list {
			...VCategoryCardFrag
			__typename
		}
		__typename
	}
	
	fragment VCategoryCardFrag on Category {
		id
		backendID
		title
		imgUrl
		watchingCount
		__typename
	}
	`
}

type HomePageCarouselsArgs struct {
	Count            int    `json:"count"`
	UserLanguageCode string `json:"userLanguageCode"`
}

// HomePageCarouselsQuery returns the graphql query to get data used to populate the home page stream carousels
func HomePageCarouselsQuery() string {
	return `query HomePageCarousels($count: Int, $userLanguageCode: String) {
		carousels(count: $count, userLanguageCode: $userLanguageCode) {
			type
			item {
				... on Livestream {
					id
					permlink
					...VLivestreamSnapFrag
					__typename
				}
				... on Poster {
					thumbnailURL
					redirectLink
					__typename
				}
				__typename
			}
			__typename
		}
	}
	
	fragment VLivestreamSnapFrag on Livestream {
		id
		creator {
			username
			displayname
			...VDliveAvatarFrag
			...VDliveNameFrag
			__typename
		}
		title
		totalReward
		watchingCount
		thumbnailUrl
		lastUpdatedAt
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

type BrowsePageSearchCategoriesArgs struct {
	Text  string `json:"text"`
	First int    `json:"first"`
	After string `json:"after"`
}

// BrowsePageSearchCategoriesQuery returns the graphql query to get all the categories you can filter on while browsing streams
func BrowsePageSearchCategoriesQuery() string {
	return `query BrowsePageSearchCategory($text: String!, $first: Int, $after: String) {
		search(text: $text) {
			trendingCategories(first: $first, after: $after) {
				...HomeCategoriesFrag
				__typename
			}
			__typename
		}
	}
	
	fragment HomeCategoriesFrag on CategoryConnection {
		pageInfo {
			endCursor
			hasNextPage
			__typename
		}
		list {
			...VCategoryCardFrag
			__typename
		}
		__typename
	}
	
	fragment VCategoryCardFrag on Category {
		id
		backendID
		title
		imgUrl
		watchingCount
		__typename
	}
	`
}

type FollowingPageLivestreamsArgs struct {
	First int    `json:"first"`
	After string `json:"after"`
}

// FollowingPageLivestreamsQuery returns the graphql query to get the streamers the currently authenticated user is following
func FollowingPageLivestreamsQuery() string {
	return `query FollowingPageLivestreams($first: Int, $after: String) {
		livestreamsFollowing(first: $first, after: $after) {
			...FollowingLivestreamsFrag
			__typename
		}
	}
	
	fragment FollowingLivestreamsFrag on LivestreamConnection {
		pageInfo {
			endCursor
			hasNextPage
			__typename
		}
		list {
			...VLivestreamSnapFrag
			__typename
		}
		__typename
	}
	
	fragment VLivestreamSnapFrag on Livestream {
		id
		creator {
			username
			displayname
			...VDliveAvatarFrag
			...VDliveNameFrag
			__typename
		}
		title
		totalReward
		watchingCount
		thumbnailUrl
		lastUpdatedAt
		__typename
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

type FollowingPageVideosArgs FollowingPageLivestreamsArgs

// FollowingPageVideosQuery returns the graphql query for getting videos uploaded by users the authenticated user is following
func FollowingPageVideosQuery() string {
	return `query FollowingPageVideos($first: Int, $after: String) {
		videosFollowing(first: $first, after: $after) {
			...FollowingVideosFrag
			__typename
		}
	}
	
	fragment FollowingVideosFrag on VideoConnection {
		pageInfo {
			endCursor
			hasNextPage
			__typename
		}
		list {
			...FollowingVideosSnapFrag
			__typename
		}
		__typename
	}
	
	fragment FollowingVideosSnapFrag on Video {
		creator {
			username
			displayname
			...VDliveNameFrag
			__typename
		}
		permlink
		title
		totalReward
		thumbnailUrl
		createdAt
		viewCount
		length
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

type SearchPageArgs struct {
	Text       string `json:"text"`
	First      int    `json:"first"`
	After      string `json:"after"`
	IsLoggedIn bool   `json:"isLoggedIn"`
}

// SearchPageQuery returns the graphql query for searching streamers, active streams, and videos for a specific term
func SearchPageQuery() string {
	return `query SearchPage($text: String!, $first: Int, $after: String, $isLoggedIn: Boolean!) {
		search(text: $text) {
			...SearchFrag
			__typename
		}
	}
	
	fragment SearchFrag on SearchResult {
		users(first: $first, after: $after) {
			...SearchUsersFrag
			__typename
		}
		livestreams(first: $first, after: $after) {
			list {
				...SearchItemLivestreamFrag
				__typename
			}
			__typename
		}
		videos(first: $first, after: $after) {
			list {
				...SearchItemVideoFrag
				__typename
			}
			__typename
		}
		__typename
	}
	
	fragment SearchItemLivestreamFrag on Livestream {
		creator {
			...VDliveNameFrag
			__typename
		}
		title
		totalReward
		watchingCount
		thumbnailUrl
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	
	fragment SearchItemVideoFrag on VideoPB {
		... on Video {
			creator {
				...VDliveNameFrag
				__typename
			}
			permlink
			title
			totalReward
			thumbnailUrl
			createdAt
			viewCount
			length
			content
			__typename
		}
		... on PastBroadcast {
			creator {
				...VDliveNameFrag
				__typename
			}
			permlink
			title
			totalReward
			thumbnailUrl
			createdAt
			viewCount
			length
			content
			__typename
		}
		__typename
	}
	
	fragment SearchUsersFrag on UserConnection {
		list {
			displayname
			avatar
			...VFollowFrag
			__typename
		}
		__typename
	}
	
	fragment VFollowFrag on User {
		id
		username
		displayname
		isFollowing @include(if: $isLoggedIn)
		isMe @include(if: $isLoggedIn)
		followers {
			totalCount
			__typename
		}
		__typename
	}
	`
}

type StreamChatBannedUsersArgs struct {
	DisplayName string `json:"displayname"`
	First       int    `json:"first"`
	After       string `json:"after"`
	Search      string `json:"search"`
}

// StreamChatBannedUsersQuery returns the graphql query for getting a list of users banned in a specific streamer's chat
func StreamChatBannedUsersQuery() string {
	return `query StreamChatBannedUsers($displayname: String!, $first: Int, $after: String, $search: String) {
		userByDisplayName(displayname: $displayname) {
			id
			chatBannedUsers(first: $first, after: $after, search: $search) {
				pageInfo {
					endCursor
					hasNextPage
					__typename
				}
				list {
					username
					...VDliveAvatarFrag
					...VDliveNameFrag
					__typename
				}
				__typename
			}
			__typename
		}
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

type StreamChatModeratorsArgs StreamChatBannedUsersArgs

// StreamChatModeratorsQuery returns the graphql query for getting a list of users that are moderators in a specific streamer's chat
func StreamChatModeratorsQuery() string {
	return `query StreamChatModerators($displayname: String!, $first: Int, $after: String, $search: String) {
		userByDisplayName(displayname: $displayname) {
			id
			chatModerators(first: $first, after: $after, search: $search) {
				pageInfo {
					endCursor
					hasNextPage
					__typename
				}
				list {
					username
					...VDliveAvatarFrag
					...VDliveNameFrag
					__typename
				}
				__typename
			}
			__typename
		}
	}
	
	fragment VDliveAvatarFrag on User {
		avatar
		__typename
	}
	
	fragment VDliveNameFrag on User {
		displayname
		partnerStatus
		__typename
	}
	`
}

type AllowedActionsArgs struct {
	Username string `json:"username"`
	Streamer string `json:"streamer"`
}

// AllowedActionsQuery returns the graphql query for getting a list of actions one user may take upon another on a given streamer's page
func AllowedActionsQuery() string {
	return `query AllowedActions($username: String!, $streamer: String!) {
		user(username: $username) {
			id
			allowedActionsIn(streamer: $streamer)
			__typename
		}
	}
	`
}
