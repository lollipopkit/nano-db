package main

import (
	_ "embed"
	"testing"

	"github.com/lollipopkit/gommon/http"
)

const (
	baseUrl = "http://localhost:3770/"
)
var (
	//go:embed test.json
	testJson []byte
	headers = map[string]string{
		"Cookie": "n=dGVzdA==; s=76af6d77b277376fde1baa23addae763da77b;",
		"Content-Type": "application/json",
	}
)

func TestAlive(t *testing.T) {
	_, code, err := http.Do("HEAD", baseUrl, nil, nil)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Log("alive")
}

func TestStatus(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl, nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("status: %s", string(body))
}

func TestUpdateFile(t *testing.T) {
	body, code, err := http.Do("POST", baseUrl + "novel/3382/chapter.json", testJson, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("update file: %s", string(body))
}

func TestSearchDir(t *testing.T) {
	reqBody := map[string]string{
		"path": "list.0.list.0.id",
		"regex": "51",
	}
	body, code, err := http.Do("POST", baseUrl + "novel/3382", reqBody, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("search dir: %s", string(body))
}

func TestSearchDB(t *testing.T) {
	reqBody := map[string]string{
		"path": "list.0.list.0.id",
		"regex": "51",
	}
	body, code, err := http.Do("POST", baseUrl + "novel", reqBody, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("search db: %s", string(body))
}

func TestGetDirnames(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl + "novel", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("dir: %s", string(body))
}

func TestGetFilenames(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl + "novel/3382", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("files: %s", string(body))
}

func TestFile(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl + "novel/3382/chapter.json", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("file: %s", string(body))
}

func TestDeleteFile(t *testing.T) {
	body, code, err := http.Do("DELETE", baseUrl + "novel/3382/chapter.json", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("delete file: %s", string(body))
}

func TestDeleteDir(t *testing.T) {
	body, code, err := http.Do("DELETE", baseUrl + "novel/3382", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("delete dir: %s", string(body))
}

func TestDeleteDB(t *testing.T) {
	body, code, err := http.Do("DELETE", baseUrl + "novel", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(err)
	}
	t.Logf("delete db: %s", string(body))
}
