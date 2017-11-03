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
	"testing"
	"reflect"

	"go-gitee/gitee"
)

func TestUsers_Update(t *testing.T) {
	if !checkAuth("TestUsers_Get") {
		return
	}

	u, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		t.Fatalf("Users.Get('') returned error: %v", err)
	}

	if *u.Login == "" {
		t.Errorf("wanted non-empty values for user.Login")
	}

	// save original address
	var address gitee.UserAddress
	if u.Address != nil {
		address = *u.Address
	}

	// update address to test value
	randAddress := fmt.Sprintf("test-%d", rand.Int())
	testLoc := gitee.UserAddress{Address : &randAddress}

	u, _, err = client.Users.EditAddress(context.Background(), &testLoc)
	if err != nil {
		t.Fatalf("Users.Update returned error: %v", err)
	}

	// refetch user and check address value
	u, _, err = client.Users.Get(context.Background(), "")
	if err != nil {
		t.Fatalf("Users.Get('') returned error: %v", err)
	}

	if !reflect.DeepEqual(u.Address.Address, testLoc.Address) {
		t.Errorf("Users.Get('') has address: %v, want: %v", *u.Address, testLoc)
	}

	// set address back to the original value
	u.Address = &address
	_, _, err = client.Users.Edit(context.Background(), u)
	if err != nil {
		t.Fatalf("Users.Edit returned error: %v", err)
	}
}

func TestUsers_Emails(t *testing.T) {
	if !checkAuth("TestUsers_Emails") {
		return
	}

	email, _, err := client.Users.GetEmail(context.Background())
	if err != nil {
		t.Fatalf("Users.GetEmail() returned error: %v", err)
	}

	// create random address not currently in user's emails
	randEmail := &gitee.UserEmail{ Email : gitee.String(fmt.Sprintf("test-%d@example.com", rand.Int())) }
		

	// Add new address
	_, _, err = client.Users.AddEmail(context.Background(), randEmail)
	if err != nil {
		t.Fatalf("Users.AddEmails() returned error: %v", err)
	}

	// List emails again and verify new email is present
	email, _, err = client.Users.GetEmail(context.Background())
	if err != nil {
		t.Fatalf("Users.GetEmail() returned error: %v", err)
	}

	if *email.UnconfirmedEmail != *randEmail.Email {
		t.Fatalf("Users.GetEmail() does not contain new address: %v", *email.UnconfirmedEmail)
	}

	// Remove new address
	_,_, err = client.Users.DeleteEmail(context.Background())
	if err != nil {
		t.Fatalf("Users.DeleteEmail() returned error: %v", err)
	}

	// List emails again and verify new email was removed
	email, _, err = client.Users.GetEmail(context.Background())
	if err != nil {
		t.Fatalf("Users.GetEmail() returned error: %v", err)
	}

	if email.UnconfirmedEmail != nil {
		t.Fatalf("Users.GetEmail() still contains address %v after removing it", *email.UnconfirmedEmail)
	}
}

func TestUsers_Keys(t *testing.T) {
	keys, _, err := client.Users.ListKeys(context.Background(), nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	if len(keys) == 0 {
		t.Errorf("Users.ListKeys('') returned no keys")
	}

	// the rest of the tests requires auth
	if !checkAuth("TestUsers_Keys") {
		return
	}

	// TODO: make this integration test work for any authenticated user.
	keys, _, err = client.Users.ListKeys(context.Background(), nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	// ssh public key for testing (fingerprint: a7:22:ad:8c:36:9f:68:65:eb:ae:a1:e4:59:73:c1:76)
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCy/RIqaMFj2wjkOEjx9EAU0ReLAIhodga82/feo5nnT9UUkHLbL9xrIavfdLHx28lD3xYgPfAoSicUMaAeNwuQhmuerr2c2LFGxzrdXP8pVsQ+Ol7y7OdmFPfe0KrzoZaLJs9aSiZ4VKyY4z5Se/k2UgcJTdgQVlLfw/P96aqCx8yUu94BiWqkDqYEvgWKRNHrTiIo1EXeVBCCcfgNZe1suFfNJUJSUU2T3EG2bpwBbSOCjE3FyH8+Lz3K3BOGzm3df8E7Regj9j4YIcD8cWJYO86jLJoGgQ0L5MSOq+ishNaHQXech22Ix03D1lVMjCvDT7S/C94Z1LzhI2lhvyff"
	for _, k := range keys {
		if k.Key != nil && *k.Key == key {
			t.Fatalf("Test key already exists for user. Please manually remove it first.")
		}
	}

	// Add new key
	_, _, err = client.Users.CreateKey(context.Background(), &gitee.Key{
		Title: gitee.String("go-gitee test key"),
		Key:   gitee.String(key),
	})
	if err != nil {
		t.Fatalf("Users.CreateKey() returned error: %v", err)
	}

	// List keys again and verify new key is present
	keys, _, err = client.Users.ListKeys(context.Background(), nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	var id int
	for _, k := range keys {
		if k.Key != nil && *k.Key == key {
			id = *k.ID
			break
		}
	}

	if id == 0 {
		t.Fatalf("Users.ListKeys('') does not contain added test key")
	}

	// Verify that fetching individual key works
	k, _, err := client.Users.GetKey(context.Background(), id)
	if err != nil {
		t.Fatalf("Users.GetKey(%q) returned error: %v", id, err)
	}
	if *k.Key != key {
		t.Fatalf("Users.GetKey(%q) returned key %v, want %v", id, *k.Key, key)
	}

	// Remove test key
	_, err = client.Users.DeleteKey(context.Background(), id)
	if err != nil {
		t.Fatalf("Users.DeleteKey(%d) returned error: %v", id, err)
	}

	// List keys again and verify test key was removed
	keys, _, err = client.Users.ListKeys(context.Background(), nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	for _, k := range keys {
		if k.Key != nil && *k.Key == key {
			t.Fatalf("Users.ListKeys('') still contains test key after removing it")
		}
	}
}
