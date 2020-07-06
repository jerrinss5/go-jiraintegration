package jiraintegration

// JiraDetails ... Structure to hold epic related details
type JiraDetails struct {
	Summary        string `json:"summary"`
	Description    string `json:"description"`
	JiraEpicKey    string `json:"jiraEpicKey"`
	JiraProjectKey string `json:"jiraProjectKey"`
	ConsumerKey    string `json:"consumerKey"`
	JiraPrivateKey string `json:"jiraPrivateKey"`
	EpicName       string `json:"epicName"`
	Type           string `json:"type"`
}
