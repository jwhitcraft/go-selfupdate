package selfupdate

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"testing"
	"runtime"
)

var testHash = sha256.New()

func TestUpdaterFetchMustReturnNonNilReaderCloser(t *testing.T) {
	mr := &mockRequester{}
	mr.handleRequest(
		func(url string) (io.ReadCloser, error) {
			return nil, nil
		})
	updater := createUpdater(mr)
	err := updater.BackgroundRun()
	if err != nil {
		equals(t, "Fetch was expected to return non-nil ReadCloser", err.Error())
	} else {
		t.Log("Expected an error")
		t.Fail()
	}
}

func TestUpdaterWithEmptyPayloadNoErrorNoUpdate(t *testing.T) {
	mr := &mockRequester{}
	mr.handleRequest(
		func(url string) (io.ReadCloser, error) {
			equals(t, "http://updates.yourdomain.com/myapp/" + runtime.GOOS + "-amd64.json", url)
			return newTestReaderCloser("{}"), nil
		})
	updater := createUpdater(mr)

	err := updater.BackgroundRun()
	if err != nil {
		t.Errorf("Error occurred: %#v", err)
	}
}

func TestUpdaterWithEmptyPayloadNoErrorNoUpdateEscapedPath(t *testing.T) {
	mr := &mockRequester{}
	mr.handleRequest(
		func(url string) (io.ReadCloser, error) {
			equals(t, "http://updates.yourdomain.com/myapp%2Bfoo/" + runtime.GOOS + "-amd64.json", url)
			return newTestReaderCloser("{}"), nil
		})
	updater := createUpdaterWithEscapedCharacters(mr)

	err := updater.BackgroundRun()
	if err != nil {
		t.Errorf("Error occurred: %#v", err)
	}
}

func TestUpdaterHasUpdateWillReturnEmptyString(t *testing.T) {
	mr := &mockRequester{}
	mr.handleRequest(
		func(url string) (io.ReadCloser, error) {
			equals(t, "http://updates.yourdomain.com/myapp/" + runtime.GOOS + "-amd64.json", url)
			return newTestReaderCloser(`{
    "Version": "1.2",
    "Sha256": "UFYeBigrXaRqJyPgCO9f8Pr22rAi0ZoovbTSmO6vxaY="
}`), nil
		})
	updater := createUpdater(mr)

	newVersion, err := updater.HasUpdate()
	if newVersion != "" {
		t.Errorf("Expected newVersion to be empty, but got %s", newVersion)
	}
	if err != nil {
		t.Errorf("Error occurred: %#v", err)
	}
}

func TestUpdaterHasUpdateWillReturnNewVersion(t *testing.T) {
	mr := &mockRequester{}
	mr.handleRequest(
		func(url string) (io.ReadCloser, error) {
			equals(t, "http://updates.yourdomain.com/myapp/" + runtime.GOOS + "-amd64.json", url)
			return newTestReaderCloser(`{
    "Version": "1.3",
    "Sha256": "UFYeBigrXaRqJyPgCO9f8Pr22rAi0ZoovbTSmO6vxaY="
}`), nil
		})
	updater := createUpdater(mr)

	newVersion, err := updater.HasUpdate()
	if newVersion != "1.3" {
		t.Errorf("Expected newVersion to be 1.3, but got %s", newVersion)
	}
	if err != nil {
		t.Errorf("Error occurred: %#v", err)
	}
}

func createUpdater(mr *mockRequester) *Updater {
	return &Updater{
		CurrentVersion: "1.2",
		ApiURL:         "http://updates.yourdomain.com/",
		BinURL:         "http://updates.yourdownmain.com/",
		DiffURL:        "http://updates.yourdomain.com/",
		Dir:            "update/",
		CmdName:        "myapp", // app name
		Requester:      mr,
	}
}

func createUpdaterWithEscapedCharacters(mr *mockRequester) *Updater {
	return &Updater{
		CurrentVersion: "1.2+foobar",
		ApiURL:         "http://updates.yourdomain.com/",
		BinURL:         "http://updates.yourdownmain.com/",
		DiffURL:        "http://updates.yourdomain.com/",
		Dir:            "update/",
		CmdName:        "myapp+foo", // app name
		Requester:      mr,
	}
}

func equals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Log(fmt.Sprintf("Expected: %#v %#v\n", expected, actual))
		t.Fail()
	}
}

type testReadCloser struct {
	buffer *bytes.Buffer
}

func newTestReaderCloser(payload string) io.ReadCloser {
	return &testReadCloser{buffer: bytes.NewBufferString(payload)}
}

func (trc *testReadCloser) Read(p []byte) (n int, err error) {
	return trc.buffer.Read(p)
}

func (trc *testReadCloser) Close() error {
	return nil
}
