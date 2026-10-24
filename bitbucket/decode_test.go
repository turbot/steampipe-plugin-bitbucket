package bitbucket

import (
	"encoding/json"
	"testing"
)

// Bodies below are trimmed from real api.bitbucket.org responses.

func TestRepositoryPermissionDecodesTheGrantAndItsRepository(t *testing.T) {
	const body = `{"values":[{"type":"repository_permission","permission":"admin",
	  "repository":{"type":"repository","full_name":"acme/widget","name":"widget","uuid":"{8c7ebd36}"},
	  "user":{"display_name":"Ada","uuid":"{3a349a75}","account_id":"712020:5adf"}}]}`

	list := new(RepositoryPermissionList)
	if err := json.Unmarshal([]byte(body), list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p := list.Permissions[0]
	if p.Permission != "admin" {
		t.Errorf("permission = %q, want admin", p.Permission)
	}
	if p.Repository.FullName != "acme/widget" {
		t.Errorf("repository full_name = %q, want acme/widget", p.Repository.FullName)
	}
	if p.User.UUID != "{3a349a75}" || p.User.AccountId != "712020:5adf" {
		t.Errorf("user did not decode: %+v", p.User)
	}
}

func TestDeployKeyDistinguishesNeverUsedFromUnused(t *testing.T) {
	const body = `{"values":[
	  {"id":1,"label":"ci","comment":"ci@acme","key":"ssh-rsa AAAA","type":"deploy_key",
	   "created_on":"2024-01-02T03:04:05.000000+00:00","last_used":null},
	  {"id":2,"label":"deploy","key":"ssh-rsa BBBB","type":"deploy_key",
	   "created_on":"2024-01-02T03:04:05.000000+00:00","last_used":"2026-09-01T00:00:00.000000+00:00"}]}`

	list := new(DeployKeyList)
	if err := json.Unmarshal([]byte(body), list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(list.Keys) != 2 {
		t.Fatalf("got %d keys, want 2", len(list.Keys))
	}
	if list.Keys[0].LastUsed != nil {
		t.Error("a key that has never been used must decode as null, not as the zero time")
	}
	if list.Keys[1].LastUsed == nil {
		t.Error("last_used did not decode")
	}
	if list.Keys[0].CreatedOn == nil {
		t.Error("created_on did not decode, so a key cannot be aged for rotation")
	}
}
