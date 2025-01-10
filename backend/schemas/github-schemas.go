package schemas

import (
	"errors"
	"time"
)

type GithubAction string

const (
	UpdateCommitInRepo      GithubAction = "UpdateCommitInRepo"
	UpdatePullRequestInRepo GithubAction = "UpdatePullRequestInRepo"
)

type GithubReaction string

// GitHubTokenResponse represents the response from Github when a token is requested.
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUserInfo struct {
	Login     string `json:"login"`
	Id        uint64 `json:"id"         gorm:"primaryKey"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
	HtmlURL   string `json:"html_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type GithubUserEmail struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visibility string `json:"visibility"`
}

// Errors Messages.
var (
	ErrGithubSecretNotSet   = errors.New("GITHUB_SECRET is not set")
	ErrGithubClientIdNotSet = errors.New("GITHUB_CLIENT_ID is not set")
)

type GithubCommit struct {
	URL         string `json:"url"`
	Sha         string `json:"sha"`
	NodeID      string `json:"node_id"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Commit      struct {
		URL    string `json:"url"`
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			URL string `json:"url"`
			Sha string `json:"sha"`
		} `json:"tree"`
		CommentCount int `json:"comment_count"`
		Verification struct {
			Verified   bool        `json:"verified"`
			Reason     string      `json:"reason"`
			Signature  interface{} `json:"signature"`
			Payload    interface{} `json:"payload"`
			VerifiedAt interface{} `json:"verified_at"`
		} `json:"verification"`
	} `json:"commit"`
	Author struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		URL string `json:"url"`
		Sha string `json:"sha"`
	} `json:"parents"`
}

type GithubActionUpdateCommitInRepo struct {
	RepoName string `json:"repo_name"`
}

type GithubActionUpdateCommitInRepoStorage struct {
	Time time.Time `json:"time"`
}

type GithubPullRequest struct {
	URL               string `json:"url"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	HTMLURL           string `json:"html_url"`
	DiffURL           string `json:"diff_url"`
	PatchURL          string `json:"patch_url"`
	IssueURL          string `json:"issue_url"`
	CommitsURL        string `json:"commits_url"`
	ReviewCommentsURL string `json:"review_comments_url"`
	ReviewCommentURL  string `json:"review_comment_url"`
	CommentsURL       string `json:"comments_url"`
	StatusesURL       string `json:"statuses_url"`
	Number            int    `json:"number"`
	State             string `json:"state"`
	Locked            bool   `json:"locked"`
	Title             string `json:"title"`
	User              struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"user"`
	Body   string `json:"body"`
	Labels []struct {
		ID          int    `json:"id"`
		NodeID      string `json:"node_id"`
		URL         string `json:"url"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Default     bool   `json:"default"`
	} `json:"labels"`
	Milestone struct {
		URL         string `json:"url"`
		HTMLURL     string `json:"html_url"`
		LabelsURL   string `json:"labels_url"`
		ID          int    `json:"id"`
		NodeID      string `json:"node_id"`
		Number      int    `json:"number"`
		State       string `json:"state"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Creator     struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"creator"`
		OpenIssues   int       `json:"open_issues"`
		ClosedIssues int       `json:"closed_issues"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		ClosedAt     time.Time `json:"closed_at"`
		DueOn        time.Time `json:"due_on"`
	} `json:"milestone"`
	ActiveLockReason string    `json:"active_lock_reason"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ClosedAt         time.Time `json:"closed_at"`
	MergedAt         time.Time `json:"merged_at"`
	MergeCommitSha   string    `json:"merge_commit_sha"`
	Assignee         struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"assignee"`
	Assignees []struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"assignees"`
	RequestedReviewers []struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"requested_reviewers"`
	RequestedTeams []struct {
		ID                  int         `json:"id"`
		NodeID              string      `json:"node_id"`
		URL                 string      `json:"url"`
		HTMLURL             string      `json:"html_url"`
		Name                string      `json:"name"`
		Slug                string      `json:"slug"`
		Description         string      `json:"description"`
		Privacy             string      `json:"privacy"`
		Permission          string      `json:"permission"`
		NotificationSetting string      `json:"notification_setting"`
		MembersURL          string      `json:"members_url"`
		RepositoriesURL     string      `json:"repositories_url"`
		Parent              interface{} `json:"parent"`
	} `json:"requested_teams"`
	Head struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
		User  struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Repo struct {
			ID       int    `json:"id"`
			NodeID   string `json:"node_id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			Owner    struct {
				Login             string `json:"login"`
				ID                int    `json:"id"`
				NodeID            string `json:"node_id"`
				AvatarURL         string `json:"avatar_url"`
				GravatarID        string `json:"gravatar_id"`
				URL               string `json:"url"`
				HTMLURL           string `json:"html_url"`
				FollowersURL      string `json:"followers_url"`
				FollowingURL      string `json:"following_url"`
				GistsURL          string `json:"gists_url"`
				StarredURL        string `json:"starred_url"`
				SubscriptionsURL  string `json:"subscriptions_url"`
				OrganizationsURL  string `json:"organizations_url"`
				ReposURL          string `json:"repos_url"`
				EventsURL         string `json:"events_url"`
				ReceivedEventsURL string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"owner"`
			Private          bool        `json:"private"`
			HTMLURL          string      `json:"html_url"`
			Description      string      `json:"description"`
			Fork             bool        `json:"fork"`
			URL              string      `json:"url"`
			ArchiveURL       string      `json:"archive_url"`
			AssigneesURL     string      `json:"assignees_url"`
			BlobsURL         string      `json:"blobs_url"`
			BranchesURL      string      `json:"branches_url"`
			CollaboratorsURL string      `json:"collaborators_url"`
			CommentsURL      string      `json:"comments_url"`
			CommitsURL       string      `json:"commits_url"`
			CompareURL       string      `json:"compare_url"`
			ContentsURL      string      `json:"contents_url"`
			ContributorsURL  string      `json:"contributors_url"`
			DeploymentsURL   string      `json:"deployments_url"`
			DownloadsURL     string      `json:"downloads_url"`
			EventsURL        string      `json:"events_url"`
			ForksURL         string      `json:"forks_url"`
			GitCommitsURL    string      `json:"git_commits_url"`
			GitRefsURL       string      `json:"git_refs_url"`
			GitTagsURL       string      `json:"git_tags_url"`
			GitURL           string      `json:"git_url"`
			IssueCommentURL  string      `json:"issue_comment_url"`
			IssueEventsURL   string      `json:"issue_events_url"`
			IssuesURL        string      `json:"issues_url"`
			KeysURL          string      `json:"keys_url"`
			LabelsURL        string      `json:"labels_url"`
			LanguagesURL     string      `json:"languages_url"`
			MergesURL        string      `json:"merges_url"`
			MilestonesURL    string      `json:"milestones_url"`
			NotificationsURL string      `json:"notifications_url"`
			PullsURL         string      `json:"pulls_url"`
			ReleasesURL      string      `json:"releases_url"`
			SSHURL           string      `json:"ssh_url"`
			StargazersURL    string      `json:"stargazers_url"`
			StatusesURL      string      `json:"statuses_url"`
			SubscribersURL   string      `json:"subscribers_url"`
			SubscriptionURL  string      `json:"subscription_url"`
			TagsURL          string      `json:"tags_url"`
			TeamsURL         string      `json:"teams_url"`
			TreesURL         string      `json:"trees_url"`
			CloneURL         string      `json:"clone_url"`
			MirrorURL        string      `json:"mirror_url"`
			HooksURL         string      `json:"hooks_url"`
			SvnURL           string      `json:"svn_url"`
			Homepage         string      `json:"homepage"`
			Language         interface{} `json:"language"`
			ForksCount       int         `json:"forks_count"`
			StargazersCount  int         `json:"stargazers_count"`
			WatchersCount    int         `json:"watchers_count"`
			Size             int         `json:"size"`
			DefaultBranch    string      `json:"default_branch"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			IsTemplate       bool        `json:"is_template"`
			Topics           []string    `json:"topics"`
			HasIssues        bool        `json:"has_issues"`
			HasProjects      bool        `json:"has_projects"`
			HasWiki          bool        `json:"has_wiki"`
			HasPages         bool        `json:"has_pages"`
			HasDownloads     bool        `json:"has_downloads"`
			Archived         bool        `json:"archived"`
			Disabled         bool        `json:"disabled"`
			Visibility       string      `json:"visibility"`
			PushedAt         time.Time   `json:"pushed_at"`
			CreatedAt        time.Time   `json:"created_at"`
			UpdatedAt        time.Time   `json:"updated_at"`
			Permissions      struct {
				Admin bool `json:"admin"`
				Push  bool `json:"push"`
				Pull  bool `json:"pull"`
			} `json:"permissions"`
			AllowRebaseMerge    bool        `json:"allow_rebase_merge"`
			TemplateRepository  interface{} `json:"template_repository"`
			TempCloneToken      string      `json:"temp_clone_token"`
			AllowSquashMerge    bool        `json:"allow_squash_merge"`
			AllowAutoMerge      bool        `json:"allow_auto_merge"`
			DeleteBranchOnMerge bool        `json:"delete_branch_on_merge"`
			AllowMergeCommit    bool        `json:"allow_merge_commit"`
			SubscribersCount    int         `json:"subscribers_count"`
			NetworkCount        int         `json:"network_count"`
			License             struct {
				Key     string `json:"key"`
				Name    string `json:"name"`
				URL     string `json:"url"`
				SpdxID  string `json:"spdx_id"`
				NodeID  string `json:"node_id"`
				HTMLURL string `json:"html_url"`
			} `json:"license"`
			Forks      int `json:"forks"`
			OpenIssues int `json:"open_issues"`
			Watchers   int `json:"watchers"`
		} `json:"repo"`
	} `json:"head"`
	Base struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
		User  struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Repo struct {
			ID       int    `json:"id"`
			NodeID   string `json:"node_id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			Owner    struct {
				Login             string `json:"login"`
				ID                int    `json:"id"`
				NodeID            string `json:"node_id"`
				AvatarURL         string `json:"avatar_url"`
				GravatarID        string `json:"gravatar_id"`
				URL               string `json:"url"`
				HTMLURL           string `json:"html_url"`
				FollowersURL      string `json:"followers_url"`
				FollowingURL      string `json:"following_url"`
				GistsURL          string `json:"gists_url"`
				StarredURL        string `json:"starred_url"`
				SubscriptionsURL  string `json:"subscriptions_url"`
				OrganizationsURL  string `json:"organizations_url"`
				ReposURL          string `json:"repos_url"`
				EventsURL         string `json:"events_url"`
				ReceivedEventsURL string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"owner"`
			Private          bool        `json:"private"`
			HTMLURL          string      `json:"html_url"`
			Description      string      `json:"description"`
			Fork             bool        `json:"fork"`
			URL              string      `json:"url"`
			ArchiveURL       string      `json:"archive_url"`
			AssigneesURL     string      `json:"assignees_url"`
			BlobsURL         string      `json:"blobs_url"`
			BranchesURL      string      `json:"branches_url"`
			CollaboratorsURL string      `json:"collaborators_url"`
			CommentsURL      string      `json:"comments_url"`
			CommitsURL       string      `json:"commits_url"`
			CompareURL       string      `json:"compare_url"`
			ContentsURL      string      `json:"contents_url"`
			ContributorsURL  string      `json:"contributors_url"`
			DeploymentsURL   string      `json:"deployments_url"`
			DownloadsURL     string      `json:"downloads_url"`
			EventsURL        string      `json:"events_url"`
			ForksURL         string      `json:"forks_url"`
			GitCommitsURL    string      `json:"git_commits_url"`
			GitRefsURL       string      `json:"git_refs_url"`
			GitTagsURL       string      `json:"git_tags_url"`
			GitURL           string      `json:"git_url"`
			IssueCommentURL  string      `json:"issue_comment_url"`
			IssueEventsURL   string      `json:"issue_events_url"`
			IssuesURL        string      `json:"issues_url"`
			KeysURL          string      `json:"keys_url"`
			LabelsURL        string      `json:"labels_url"`
			LanguagesURL     string      `json:"languages_url"`
			MergesURL        string      `json:"merges_url"`
			MilestonesURL    string      `json:"milestones_url"`
			NotificationsURL string      `json:"notifications_url"`
			PullsURL         string      `json:"pulls_url"`
			ReleasesURL      string      `json:"releases_url"`
			SSHURL           string      `json:"ssh_url"`
			StargazersURL    string      `json:"stargazers_url"`
			StatusesURL      string      `json:"statuses_url"`
			SubscribersURL   string      `json:"subscribers_url"`
			SubscriptionURL  string      `json:"subscription_url"`
			TagsURL          string      `json:"tags_url"`
			TeamsURL         string      `json:"teams_url"`
			TreesURL         string      `json:"trees_url"`
			CloneURL         string      `json:"clone_url"`
			MirrorURL        string      `json:"mirror_url"`
			HooksURL         string      `json:"hooks_url"`
			SvnURL           string      `json:"svn_url"`
			Homepage         string      `json:"homepage"`
			Language         interface{} `json:"language"`
			ForksCount       int         `json:"forks_count"`
			StargazersCount  int         `json:"stargazers_count"`
			WatchersCount    int         `json:"watchers_count"`
			Size             int         `json:"size"`
			DefaultBranch    string      `json:"default_branch"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			IsTemplate       bool        `json:"is_template"`
			Topics           []string    `json:"topics"`
			HasIssues        bool        `json:"has_issues"`
			HasProjects      bool        `json:"has_projects"`
			HasWiki          bool        `json:"has_wiki"`
			HasPages         bool        `json:"has_pages"`
			HasDownloads     bool        `json:"has_downloads"`
			Archived         bool        `json:"archived"`
			Disabled         bool        `json:"disabled"`
			Visibility       string      `json:"visibility"`
			PushedAt         time.Time   `json:"pushed_at"`
			CreatedAt        time.Time   `json:"created_at"`
			UpdatedAt        time.Time   `json:"updated_at"`
			Permissions      struct {
				Admin bool `json:"admin"`
				Push  bool `json:"push"`
				Pull  bool `json:"pull"`
			} `json:"permissions"`
			AllowRebaseMerge    bool        `json:"allow_rebase_merge"`
			TemplateRepository  interface{} `json:"template_repository"`
			TempCloneToken      string      `json:"temp_clone_token"`
			AllowSquashMerge    bool        `json:"allow_squash_merge"`
			AllowAutoMerge      bool        `json:"allow_auto_merge"`
			DeleteBranchOnMerge bool        `json:"delete_branch_on_merge"`
			AllowMergeCommit    bool        `json:"allow_merge_commit"`
			SubscribersCount    int         `json:"subscribers_count"`
			NetworkCount        int         `json:"network_count"`
			License             struct {
				Key     string `json:"key"`
				Name    string `json:"name"`
				URL     string `json:"url"`
				SpdxID  string `json:"spdx_id"`
				NodeID  string `json:"node_id"`
				HTMLURL string `json:"html_url"`
			} `json:"license"`
			Forks      int `json:"forks"`
			OpenIssues int `json:"open_issues"`
			Watchers   int `json:"watchers"`
		} `json:"repo"`
	} `json:"base"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Issue struct {
			Href string `json:"href"`
		} `json:"issue"`
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
		ReviewComment struct {
			Href string `json:"href"`
		} `json:"review_comment"`
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		Statuses struct {
			Href string `json:"href"`
		} `json:"statuses"`
	} `json:"_links"`
	AuthorAssociation string      `json:"author_association"`
	AutoMerge         interface{} `json:"auto_merge"`
	Draft             bool        `json:"draft"`
}

type GithubActionUpdatePullRequestInRepo struct {
	RepoName string `json:"repo_name"`
}

type GithubActionUpdatePullRequestInRepoStorage struct {
	Time time.Time `json:"time"`
}
