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
	headers  = map[string]string{
		"NanoDB":       "FHYmGdNwfiJngvF2z",
	}
)

func TestAlive(t *testing.T) {
	_, code, err := http.Do("HEAD", baseUrl, nil, nil)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log("alive")
}

func TestIllegalPath(t *testing.T) {
	_illegalPath(t, "..")
	_illegalPath(t, "/")
	_illegalPath(t, "-")
}
func _illegalPath(t *testing.T, path string) {
	body, code, err := http.Do("GET", baseUrl+path, nil, headers)
	if err != nil || code == 200 {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestUpdateFile(t *testing.T) {
	body, code, err := http.Do("POST", baseUrl+"novel/3382/chapter.json", testJson, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(string(body))
}

func TestGetDirnames(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl+"novel", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(string(body))
}

func TestGetFilenames(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl+"novel/3382", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(string(body))
}

func TestFile(t *testing.T) {
	body, code, err := http.Do("GET", baseUrl+"novel/3382/chapter.json", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(len(body))
}

func TestDeleteFile(t *testing.T) {
	body, code, err := http.Do("DELETE", baseUrl+"novel/3382/chapter.json", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(string(body))
}

func TestDeleteDir(t *testing.T) {
	body, code, err := http.Do("DELETE", baseUrl+"novel/3382", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(string(body))
}

func TestDeleteDB(t *testing.T) {
	body, code, err := http.Do("DELETE", baseUrl+"novel", nil, headers)
	if err != nil || code != 200 {
		t.Fatal(code, err)
	}
	t.Log(string(body))
}
