package jiraintegration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// getFromProperty ... reads from json file the description and returns unmarshalled data
func getFromProperty(fileName string, log *logrus.Entry) (JiraDetails, error) {
	var jiraType JiraDetails
	jiraTypeData, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("Some error occurred reading json file: " + err.Error())
		return jiraType, err
	}

	err = json.Unmarshal(jiraTypeData, &jiraType)
	if err != nil {
		log.Errorf("Some error occurred unmarshalling data: " + err.Error())
		return jiraType, err
	}

	return jiraType, nil
}
