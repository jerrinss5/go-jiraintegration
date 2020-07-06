package jiraintegration

import (
	"io/ioutil"

	"github.com/andygrunwald/go-jira"
	"github.com/sirupsen/logrus"
)

// getProjects ... is to fetch all the projects in the JIRA URL
func getProjects(jiraURL string, log *logrus.Entry) {
	jiraClient := getJIRAClient(jiraURL, log)
	req, _ := jiraClient.NewRequest("GET", "rest/api/2/project", nil)

	projects := new([]jira.Project)
	_, err := jiraClient.Do(req, projects)
	if err != nil {
		panic(err)
	}

	for _, project := range *projects {
		log.Infof("%s: %s\n", project.Key, project.Name)
	}
}

// createEpic ... is used to create epic for the specified jira URL
func createEpic(jiraURL string, log *logrus.Entry) {
	jiraDetails, err := getFromProperty("epic.json", log)
	if err != nil {
		log.Fatalf("Error occurred reading epic json file for description: " + err.Error())
	}
	epicName := make(map[string]interface{})
	// this value can be fetched from admin panel of fields within JIRA
	epicName["customfield_10006"] = jiraDetails.JiraEpicKey
	jiraClient := getJIRAClient(jiraURL, log)
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: jiraDetails.Description,
			Type: jira.IssueType{
				Name: jiraDetails.Type,
			},
			Project: jira.Project{
				Key: jiraDetails.JiraEpicKey,
			},
			Summary:  jiraDetails.Summary,
			Unknowns: epicName,
		},
	}
	issue, resp, err := jiraClient.Issue.Create(&i)
	if err != nil {
		log.Error("Some error occurred creating the JIRA ticket: " + err.Error())
		responeseByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error occurred reading the response body: " + err.Error())
		}
		log.Fatalf("JIRA Error Response body for creation: " + string(responeseByte))
	}
	log.Infof("Epic Created with details -> %s: and ID: %+v\n", issue.Key, issue.ID)
}
