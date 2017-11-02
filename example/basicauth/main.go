// Copyright 2015 The go-gitee AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The basicauth command demonstrates using the gitee.BasicAuthTransport,
// including handling two-factor authentication. This won't currently work for
// accounts that use SMS to receive one-time passwords.
package main

import (
	// "bufio"
	"context"
	"fmt"
	// "os"
	// "strings"
	// "syscall"
	"golang.org/x/oauth2"

	"go-gitee/gitee"
)

func main() {

	ctx := context.Background()
	conf := &oauth2.Config{
	    ClientID:     "47e10a732882363062588a4cb26e0eea80eb4e3d32f60ce8193f7f96e467abac",
	    ClientSecret: "228f59e9306d14b95de611db8ffcb39524a9c6e499dc1c09eed1652a63781d57",
	    Scopes:       []string{"user_info", "projects"},
	    Endpoint: oauth2.Endpoint{
	        AuthURL:  "https://gitee.com/oauth/auth",
	        TokenURL: "https://gitee.com/oauth/token",
	    },
	}
	token,err := conf.PasswordCredentialsToken(ctx,"noreply@daiheimao.top","1qaz2wsx")

	tp := gitee.OAuthTransport{
		Token: token,
	}

	client := gitee.NewClient(tp.Client())
	
	user, _, err := client.Activity.ListStargazers(ctx, "")

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	fmt.Printf("\n%v\n", gitee.Stringify(user))
	fmt.Printf("\n%v\n", *user.Login)
}
