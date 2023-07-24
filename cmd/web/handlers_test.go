package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/windevkay/hardpassv2/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestPasswordView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct{
		name		string
		urlPath		string
		wantCode	int
		wantBody	string
	}{
		{
			name: "Valid ID",
			urlPath: "/password/view/1",
			wantCode: http.StatusOK,
			wantBody: "test",
		},
		{
			name: "Non-existent ID",
			urlPath: "/password/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			urlPath: "/password/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			urlPath: "/password/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			urlPath: "/password/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			urlPath: "/password/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests { 
		t.Run(tt.name, func(t *testing.T) {
			// do a post to login first to get a valid cookie for subseqnet requests
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			} 
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)
	
	const (
		validName = "Bob"
		validPasword = "validPa$$word"
		validEmail = "bob@example.com"
		formTag = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct{
		name string 
		userName string 
		userEmail string 
		userPassword string 
		csrfToken string 
		wantCode int 
		wantFormTag string
	}{
		{
			name: "Valid Submission",
			userName: validName,
			userEmail: validEmail,
			userPassword: validPasword,
			csrfToken: validCSRFToken,
			wantCode: http.StatusSeeOther,
		},
		{
			name: "Invalid CSRF Token",
			userName: validName,
			userEmail: validEmail,
			userPassword: validPasword,
			csrfToken: "wrongToken",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Empty name",
			userName: " ",
			userEmail: validEmail,
			userPassword: validPasword,
			csrfToken: validCSRFToken,
			wantCode: http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name: "Empty email",
			userName: validName,
			userEmail: "",
			userPassword: validPasword,
			csrfToken: validCSRFToken,
			wantCode: http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name: "Empty password",
			userName: validName,
			userEmail: validEmail,
			userPassword: "",
			csrfToken: validCSRFToken,
			wantCode: http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name: "Invalid email",
			userName: validName,
			userEmail: "bob@example",
			userPassword: validPasword,
			csrfToken: validCSRFToken,
			wantCode: http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name: "Short password",
			userName: validName,
			userEmail: validEmail,
			userPassword: "pa$$",
			csrfToken: validCSRFToken,
			wantCode: http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		// {
		// 	name: "Duplicate email",
		// 	userName: validName,
		// 	userEmail: "dupe@example.com",
		// 	userPassword: validPasword,
		// 	csrfToken: validCSRFToken,
		// 	wantCode: http.StatusUnprocessableEntity,
		// 	wantFormTag: formTag,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			form := url.Values{}
			form.Add("name", tt.userName) 
			form.Add("email", tt.userEmail) 
			form.Add("password", tt.userPassword) 
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}