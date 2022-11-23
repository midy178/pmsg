// Copyright 2022 The pmsg Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/lenye/pmsg/pkg/version"
)

var UserAgent string

var ErrHttpRequest = errors.New("http request error")

const (
	contentTypeJson = "application/json;charset=utf-8"
)

func DefaultUserAgent() string {
	return fmt.Sprintf("%s/%s (%s; %s) %s/%s", version.AppName, version.Version, runtime.GOOS, runtime.GOARCH, version.BuildGit, version.BuildTime)
}

func userAgent() string {
	if UserAgent != "" {
		return UserAgent
	}
	return DefaultUserAgent()
}

// Get http get
func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent())

	return http.DefaultClient.Do(req)
}

// Post http post
func Post(url, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("User-Agent", userAgent())

	return http.DefaultClient.Do(req)
}

func fileToBody(bodyWriter *multipart.Writer, formName, fileName string) error {
	fileWriter, err := bodyWriter.CreateFormFile(formName, filepath.Base(fileName))
	if err != nil {
		return fmt.Errorf("multipart.Writer.CreateFormFile failed, %w", err)
	}

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("open file failed, %w", err)
	}

	if _, err := io.Copy(fileWriter, f); err != nil {
		return fmt.Errorf("file io.Copy failed, %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("close file  failed, %w", err)
	}
	return nil
}

// MultipartForm 保存文件或其他字段信息
type MultipartForm struct {
	params map[string][]string
	files  map[string]string
}

func NewMultipartForm() *MultipartForm {
	return &MultipartForm{
		params: make(map[string][]string),
		files:  make(map[string]string),
	}
}

func (t *MultipartForm) AddFile(formName, fileName string) *MultipartForm {
	t.files[formName] = fileName
	return t
}

func (t *MultipartForm) AddParam(key, value string) *MultipartForm {
	if param, ok := t.params[key]; ok {
		t.params[key] = append(param, value)
	} else {
		t.params[key] = []string{value}
	}
	return t
}

// PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(url string, form *MultipartForm) (*http.Response, error) {
	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	for formName, fileName := range form.files {
		if err := fileToBody(bodyWriter, formName, fileName); err != nil {
			return nil, err
		}
	}
	for k, v := range form.params {
		for _, vv := range v {
			if err := bodyWriter.WriteField(k, vv); err != nil {
				return nil, fmt.Errorf("multipart.Writer.WriteField failed, %w", err)
			}
		}
	}
	contentType := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		return nil, fmt.Errorf("multipart.Writer.Close failed, %w", err)
	}
	return Post(url, contentType, bodyBuf)
}

// PostFile 上传文件
func PostFile(url, formName, fileName string) (*http.Response, error) {
	form := NewMultipartForm().AddFile(formName, fileName)
	return PostMultipartForm(url, form)
}
