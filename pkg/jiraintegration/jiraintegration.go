package jiraintegration

import (
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

// Environment ... properties file loaded via the file
type Environment struct {
	JiraSecretName string `json:"jiraSecretName"`
	JiraURL        string `json:"jiraURL"`
	JiraProjectKey string `json:"jiraProjectKey"`
	JiraEpicKey    string `json:"jiraEpicKey"`
}

var log *logrus.Entry

func init() {
	logger := logrus.New()

	logger.Out = os.Stderr // default log to stderr
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	log = logger.WithFields(logrus.Fields{"package": "jiraintegration"})
}

// Run ... runs the Jira integration service after fetching the property file location
func Run() {
	var jiraURL, jiraType string
	log.Infof("GO JIRA integration example")
	flag.StringVarP(&jiraURL, "jiraurl", "j", "", "URL of the JIRA server")
	flag.StringVarP(&jiraType, "type", "t", "", "Takes in either type epic or project")
	// flag.StringVarP(&summary, "summary", "tl", "", "Summary for the epic")
	// flag.StringVarP(&desc, "desc", "d", "", "Description for the epic")
	// flag.StringVarP(&epicName, "name", "n", "", "Name for the epic")
	flag.PrintDefaults()
	flag.Parse()

	if jiraType == "" {
		log.Fatalf("JIRA Type is required and specified via t flag. Exiting!")
	}

	switch jiraType {
	case "epic":
		log.Tracef("Epic JIRA Type is selected!")
		createEpic(jiraURL, log)
	case "project":
		log.Infof("Project JIRA Type is selected!")
	case "fetch":
		log.Infof("Get project details is selected")
		getProjects(jiraURL, log)
	default:
		log.Fatalf("Invalid project type provided: " + jiraType)
	}

}
