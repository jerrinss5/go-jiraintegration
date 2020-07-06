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

func createIssueStory(jiraURL string, log *logrus.Entry) {
	jiraDetails, err := getFromProperty("issue.json", log)
	summary := jiraDetails.Summary
	description := jiraDetails.Description
	if err != nil {
		log.Fatalf("Error occurred reading issue json file: " + err.Error())
	}
	// summary has a field length restriction of 254
	// https://confluence.atlassian.com/jirasoftwareserver071/advanced-searching-fields-reference-800707156.html
	if len(summary) > 254 {
		summary = summary[0:244]
	}
	// description has a character limit restriction of 32766
	// https://community.atlassian.com/t5/Jira-questions/JIRA-Character-Limits-on-Comments-and-Description-32-767/qaq-p/259057
	if len(description) > 32766 {
		description = description[0:32766]
	}

	jiraClient := getJIRAClient(jiraURL, log)
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: description,
			Type: jira.IssueType{
				Name: jiraDetails.Type,
			},
			Project: jira.Project{
				Key: jiraDetails.JiraProjectKey,
			},
			Summary: summary,
		},
	}

	issue, resp, err := jiraClient.Issue.Create(&i)
	if err != nil {
		log.Error("Some error occurred creating the JIRA ticket: " + err.Error())
		responeseByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error occurred reading the response body: " + err.Error())
			return
		}
		log.Info("JIRA Error Response body for creation: " + string(responeseByte))
		return
	}

	log.Infof("Story created with key: %s: ID: %+v\n", issue.Key, issue.ID)
}
