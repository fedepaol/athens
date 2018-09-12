package credentials

type pullParams struct {
	RepoURL      string
	RepoProtocol string
	Operation    string
}

type credentials struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}
