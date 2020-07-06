package jiraintegration

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

func getJIRAClient(jiraURL string, log *logrus.Entry) *jira.Client {
	ctx := context.Background()

	jiraConsumerKey, jiraPrivateKey := getJIRASecretsFromFile(log)

	keyDERBlock, _ := pem.Decode([]byte(jiraPrivateKey))
	if keyDERBlock == nil {
		log.Fatal("unable to decode key PEM block")
	}
	if !(keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY")) {
		log.Fatalf("unexpected key DER block type: %s", keyDERBlock.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	if err != nil {
		log.Fatalf("unable to parse PKCS1 private key. %v", err)
	}

	config := oauth1.Config{
		ConsumerKey: jiraConsumerKey,
		CallbackURL: "oob", /* for command line usage */
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: jiraURL + "plugins/servlet/oauth/request-token",
			AuthorizeURL:    jiraURL + "plugins/servlet/oauth/authorize",
			AccessTokenURL:  jiraURL + "plugins/servlet/oauth/access-token",
		},
		Signer: &oauth1.RSASigner{
			PrivateKey: privateKey,
		},
	}
	jiraClient, err := jira.NewClient(getJIRAHTTPClient(ctx, &config, jiraURL, log), jiraURL)
	if err != nil {
		log.Fatalf("unable to create new JIRA client. %v", err)
	}
	return jiraClient
}

func getJIRASecretsFromFile(log *logrus.Entry) (string, string) {
	jiraSecrets, err := getFromProperty("jiraproperties.json", log)
	if err != nil {
		log.Error("Some error occurred from the property file: " + err.Error())
	}
	jiraPrivateKey := jiraSecrets.JiraPrivateKey

	i := 0
	jiraPrivateKeyFormated := "-----BEGIN RSA PRIVATE KEY-----\n"

	for i < len(jiraPrivateKey) {
		if (i + 64) < len(jiraPrivateKey) {
			jiraPrivateKeyFormated += jiraPrivateKey[i:i+64] + "\n"
		} else {
			jiraPrivateKeyFormated += jiraPrivateKey[i:] + "\n"
		}

		i = i + 64
	}
	jiraPrivateKeyFormated += "-----END RSA PRIVATE KEY-----"
	return jiraSecrets.ConsumerKey, jiraPrivateKeyFormated
}

func getJIRAHTTPClient(ctx context.Context, config *oauth1.Config, jiraURL string, log *logrus.Entry) *http.Client {
	cacheFile, err := jiraTokenCacheFile(jiraURL)
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := jiraTokenFromFile(cacheFile)
	if err != nil {
		tok = getJIRATokenFromWeb(config)
		saveJIRAToken(cacheFile, tok, log)
	}
	return config.Client(ctx, tok)
}

func getJIRATokenFromWeb(config *oauth1.Config) *oauth1.Token {
	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		log.Fatalf("Unable to get request token. %v", err)
	}
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatalf("Unable to get authorization url. %v", err)
	}
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authorizationURL.String())

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code. %v", err)
	}

	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, code)
	if err != nil {
		log.Fatalf("Unable to get access token. %v", err)
	}
	return oauth1.NewToken(accessToken, accessSecret)
}

func jiraTokenCacheFile(jiraURL string) (string, error) {
	tokenCacheDir := "data/.credentials"
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(jiraURL+".json")), nil
}

func jiraTokenFromFile(file string) (*oauth1.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth1.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveJIRAToken(file string, token *oauth1.Token, log *logrus.Entry) {
	log.Infof("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
