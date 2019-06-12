package gitsemvers

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/Songmu/gitmock"
)

var input = `dummy
foo-0.10.1
foo-0.9.0
foo-0.9.3
foo-0.8.4-pre
foo-0.8.4
foo-0.8.3+win
foo-0.8.2-pre.pre+win.win
foo-0.7.0-pre+win+invalid
v0.10.1
v0.9.0
v0.9.3
v0.8.4-pre
v0.8.4
v0.8.3+win
v0.8.2-pre.pre+win.win
v0.7.0-pre+win+invalid
`

func TestParseVersions(t *testing.T) {
	expect := []string{
		"v0.10.1",
		"v0.9.3",
		"v0.9.0",
		"v0.8.4",
	}
	sv := &Semvers{}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithVerPrefix(t *testing.T) {
	expect := []string{
		"foo-0.10.1",
		"foo-0.9.3",
		"foo-0.9.0",
		"foo-0.8.4",
	}
	sv := &Semvers{
		VerPrefix: "foo-",
	}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithPreRelease(t *testing.T) {
	expect := []string{
		"v0.10.1",
		"v0.9.3",
		"v0.9.0",
		"v0.8.4",
		"v0.8.4-pre",
	}
	sv := &Semvers{WithPreRelease: true}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithPreReleaseAndVerPrefix(t *testing.T) {
	expect := []string{
		"foo-0.10.1",
		"foo-0.9.3",
		"foo-0.9.0",
		"foo-0.8.4",
		"foo-0.8.4-pre",
	}
	sv := &Semvers{WithPreRelease: true, VerPrefix: "foo-"}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithBuildMetadata(t *testing.T) {
	expect := []string{
		"v0.10.1",
		"v0.9.3",
		"v0.9.0",
		"v0.8.4",
		"v0.8.3+win",
	}
	sv := &Semvers{WithBuildMetadata: true}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithBuildMetadataAndVerPrefix(t *testing.T) {
	expect := []string{
		"foo-0.10.1",
		"foo-0.9.3",
		"foo-0.9.0",
		"foo-0.8.4",
		"foo-0.8.3+win",
	}
	sv := &Semvers{WithBuildMetadata: true, VerPrefix: "foo-"}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithAllExtensions(t *testing.T) {
	expect := []string{
		"v0.10.1",
		"v0.9.3",
		"v0.9.0",
		"v0.8.4",
		"v0.8.4-pre",
		"v0.8.3+win",
		"v0.8.2-pre.pre+win.win",
	}
	sv := &Semvers{WithPreRelease: true, WithBuildMetadata: true}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestParseVersionsWithAllExtensionsAndVerPrefix(t *testing.T) {
	expect := []string{
		"foo-0.10.1",
		"foo-0.9.3",
		"foo-0.9.0",
		"foo-0.8.4",
		"foo-0.8.4-pre",
		"foo-0.8.3+win",
		"foo-0.8.2-pre.pre+win.win",
	}
	sv := &Semvers{WithPreRelease: true, WithBuildMetadata: true, VerPrefix: "foo-"}
	if !reflect.DeepEqual(sv.parseVersions(input), expect) {
		t.Errorf("something went wrong %+v", sv.parseVersions(input))
	}
}

func TestVersionStrings(t *testing.T) {
	gm, err := gitmock.New("")
	if err != nil {
		t.Fatal(err)
	}
	repoPath := gm.RepoPath()
	defer os.RemoveAll(repoPath)
	gm.Init()
	gm.Commit("--allow-empty", "-m", "initial commit")
	gm.Tag("0.0.1")
	gm.Tag("v0.0.2")
	gm.Tag("v0.0.2-pre")
	gm.Tag("v0.0.2-pre+win")
	gm.Tag("v0.0.2+win")
	gm.Tag("v0.0.2+win+invalid")

	sv := &Semvers{
		RepoPath: repoPath,
	}
	{
		expect := []string{
			"v0.0.2",
			"0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong")
		}
	}

	sv.WithPreRelease = true
	{
		expect := []string{
			"v0.0.2",
			"v0.0.2-pre",
			"0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong %+v", sv.VersionStrings())
		}
	}

	sv.WithPreRelease = false
	sv.WithBuildMetadata = true
	{
		expect := []string{
			"v0.0.2",
			"v0.0.2+win",
			"0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong %+v", sv.VersionStrings())
		}
	}

	sv.WithPreRelease = true
	sv.WithBuildMetadata = true
	{
		expect := []string{
			"v0.0.2",
			"v0.0.2+win",
			"v0.0.2-pre",
			"v0.0.2-pre+win",
			"0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong %+v", sv.VersionStrings())
		}
	}
}

func TestVersionStringsWithPrefix(t *testing.T) {
	gm, err := gitmock.New("")
	if err != nil {
		t.Fatal(err)
	}
	repoPath := gm.RepoPath()
	defer os.RemoveAll(repoPath)
	gm.Init()
	gm.Commit("--allow-empty", "-m", "initial commit")
	gm.Tag("0.0.1")
	gm.Tag("v0.0.2")
	gm.Tag("v0.0.2-pre")
	gm.Tag("v0.0.2-pre+win")
	gm.Tag("v0.0.2+win")
	gm.Tag("v0.0.2+win+invalid")
	gm.Tag("foo-0.0.1")
	gm.Tag("foo-0.0.2")
	gm.Tag("foo-0.0.2-pre")
	gm.Tag("foo-0.0.2-pre+win")
	gm.Tag("foo-0.0.2+win")
	gm.Tag("foo-0.0.2+win+invalid")

	sv := &Semvers{
		RepoPath:  repoPath,
		VerPrefix: "foo-",
	}

	{
		expect := []string{
			"foo-0.0.2",
			"foo-0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong")
		}
	}

	sv.WithPreRelease = true
	{
		expect := []string{
			"foo-0.0.2",
			"foo-0.0.2-pre",
			"foo-0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong %+v", sv.VersionStrings())
		}
	}

	sv.WithPreRelease = false
	sv.WithBuildMetadata = true
	{
		expect := []string{
			"foo-0.0.2",
			"foo-0.0.2+win",
			"foo-0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong %+v", sv.VersionStrings())
		}
	}

	sv.WithPreRelease = true
	sv.WithBuildMetadata = true
	{
		expect := []string{
			"foo-0.0.2",
			"foo-0.0.2+win",
			"foo-0.0.2-pre",
			"foo-0.0.2-pre+win",
			"foo-0.0.1",
		}
		if !reflect.DeepEqual(sv.VersionStrings(), expect) {
			t.Errorf("something went wrong %+v", sv.VersionStrings())
		}
	}
}

func TestCLI(t *testing.T) {
	gm, err := gitmock.New("")
	if err != nil {
		t.Fatal(err)
	}
	repoPath := gm.RepoPath()
	defer os.RemoveAll(repoPath)
	gm.Init()
	gm.Commit("--allow-empty", "-m", "initial commit")
	gm.Tag("0.0.1")
	gm.Tag("v0.0.2")
	gm.Tag("v0.0.2-pre")
	gm.Tag("v0.0.2-pre+win")
	gm.Tag("v0.0.2+win")
	gm.Tag("v0.0.2+win+invalid")
	gm.Tag("foo-0.0.1")
	gm.Tag("foo-0.0.2")
	gm.Tag("foo-0.0.2-pre")
	gm.Tag("foo-0.0.2-pre+win")
	gm.Tag("foo-0.0.2+win")
	gm.Tag("foo-0.0.2+win+invalid")

	{
		var b bytes.Buffer
		cli := CLI{OutStream: &b, ErrStream: ioutil.Discard}

		ret := cli.Run([]string{"--hoge"})
		if ret != exitCodeParseFlagErr {
			t.Errorf("exit code should be exitCodeParseFlagErr but: %d", ret)
		}

		cli.Run([]string{"--repo", repoPath})
		expect := "v0.0.2\n0.0.1\n"
		output := b.String()
		if output != expect {
			t.Errorf("expect: %s, but: %s", expect, output)
		}
	}

	{
		var b bytes.Buffer
		cli := CLI{OutStream: &b, ErrStream: ioutil.Discard}

		ret := cli.Run([]string{"--hoge"})
		if ret != exitCodeParseFlagErr {
			t.Errorf("exit code should be exitCodeParseFlagErr but: %d", ret)
		}

		cli.Run([]string{"--repo", repoPath, "--ver-prefix", "foo-"})
		expect := "foo-0.0.2\nfoo-0.0.1\n"
		output := b.String()
		if output != expect {
			t.Errorf("expect: %s, but: %s", expect, output)
		}
	}
}
