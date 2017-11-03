// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package integration

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"go-gitee/gitee"
)

var (
	client *gitee.Client

	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.
	auth bool
)

const msgEnvMissing = "Skipping test because the required environment variable (%v) is not present."
const envKeyGiteeUsername = "GITEE_USERNAME"
const envKeyGiteePassword = "GITEE_PASSWORD"
const envKeyClientID = "GITEE_CLIENT_ID"
const envKeyClientSecret = "GITEE_CLIENT_SECRET"
const InvalidTokenValue = "iamnotacroken"

func init() {
	vars := []string{envKeyGiteeUsername, envKeyGiteePassword, envKeyClientID, envKeyClientSecret}

	for _, v := range vars {
		value := os.Getenv(v)
		if value == "" {
			print("!!! " + fmt.Sprintf(msgEnvMissing, v) + " !!!\n\n")
		}
	}

	username, ok := os.LookupEnv(envKeyGiteeUsername)
	if !ok {
		print("!!! No OAuth token. Some tests won't run. !!!\n\n")
		client = gitee.NewClient(nil)
	}

	password, ok := os.LookupEnv(envKeyGiteePassword)
	if !ok {
		print("!!! No OAuth token. Some tests won't run. !!!\n\n")
		client = gitee.NewClient(nil)
		return
	}

	clientID, ok := os.LookupEnv(envKeyClientID)
	if !ok {
		print("!!! No OAuth token. Some tests won't run. !!!\n\n")
		client = gitee.NewClient(nil)
		return
	}

	clientSecret, ok := os.LookupEnv(envKeyClientSecret)
	if !ok {
		print("!!! No OAuth token. Some tests won't run. !!!\n\n")
		client = gitee.NewClient(nil)
		return
	}

	ctx := context.Background()
	conf := &oauth2.Config{
	    ClientID:     clientID,
	    ClientSecret: clientSecret,
	    Scopes:       []string{"user_info", "projects", "pull_requests", "issues", "notes", "keys", "hook", "groups", "gists"},
	    Endpoint: oauth2.Endpoint{
	        AuthURL:  "https://gitee.com/oauth/auth",
	        TokenURL: "https://gitee.com/oauth/token",
	    },
	}
	token,err := conf.PasswordCredentialsToken(ctx,username,password)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		client = gitee.NewClient(nil)
		return
	}

	tp := gitee.OAuthTransport{
		Token: token,
	}

	fmt.Printf("\n init oauth client ok.\n")
	client = gitee.NewClient(tp.Client())
	auth = true
}

func checkAuth(name string) bool {
	if !auth {
		fmt.Printf("No auth - skipping portions of %v\n", name)
	}
	return auth
}

func createRandomTestRepository(owner string, autoinit bool) (*gitee.Repository, error) {
	// create random repo name that does not currently exist
	var repoName string
	for {
		repoName = fmt.Sprintf("test-%d", rand.Int())
		_, resp, err := client.Repositories.Get(context.Background(), owner, repoName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				// found a non-existent repo, perfect
				break
			}

			return nil, err
		}
	}

	// create the repository
	repo, _, err := client.Repositories.Create(context.Background(), "", &gitee.Repository{Name: gitee.String(repoName), AutoInit: gitee.Bool(autoinit)})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
