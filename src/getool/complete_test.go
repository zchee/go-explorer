// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const completeTestFile = `
package main

import (
    p1 "github.com/user/repo1"
    "github.com/user/repo2"
    "github.com/user/repo3"
)
`

var completeTests = []struct {
	in  string
	out string
}{
	// Some of these tests depend on the content of the Go workspace. A failure
	// may be be because the workspace is not what's expected and not a failure
	// in the package.

	{"", "p1\nrepo2\nrepo3"},
	{"repo", "repo2\nrepo3"},
	{".", "./\n../"},
	{"..", "../"},
	{"/net/http Client", "Client."},
	{"/net/http client.postf", "Client.PostForm"},
	{"/github.com", "/github.com/"},
	{"../getool", "../getool/"},
	{"/go/", "/go/ast/\n/go/build/\n/go/doc/\n/go/format/\n/go/parser/\n/go/printer/\n/go/scanner/\n/go/token/"},
}

func TestComplete(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range completeTests {
		var buf bytes.Buffer
		doComplete(&Context{
			out: &buf,
			in:  strings.NewReader(completeTestFile),
			cwd: cwd,
			args: []string{
				tt.in[strings.LastIndex(tt.in, " ")+1:],
				"x " + tt.in,
				"",
			},
		})

		out := buf.String()
		if out != tt.out {
			t.Errorf("complete(%q) = %q, want %q", tt.in, out, tt.out)
		}
	}
}

var resolveTests = []struct {
	in  string
	out string
}{
	{"p1", "github.com/user/repo1"},
	{"repo2", "github.com/user/repo2"},
	{"/github.com/user/repo3", "github.com/user/repo3"},
	{".", "github.com/garyburd/getool"},
}

func TestResolve(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range resolveTests {
		var buf bytes.Buffer
		doResolve(&Context{
			out:  &buf,
			in:   strings.NewReader(completeTestFile),
			cwd:  cwd,
			args: []string{tt.in},
		})
		out := buf.String()
		if out != tt.out {
			t.Errorf("resolve(%q) = %q, want %q", tt.in, out, tt.out)
		}
	}
}
