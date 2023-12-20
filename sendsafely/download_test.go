/*
   Copyright 2022 Ryan SVIHLA

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// sendsafely package decrypts files, combines file parts into whole files, and handles api access to the sendsafely rest api
package sendsafely

import (
	"crypto/rand"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This is the default happy path test, no errors
func TestCalculateExecutionCalls(t *testing.T) {
	req := calculateExecutionCalls(1)
	if len(req) != 1 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)
	}
}

func TestCalculateExecutionCallsBeforeBoundary(t *testing.T) {
	req := calculateExecutionCalls(24)
	if len(req) != 1 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 24 {
		t.Errorf("expected end segment to be 24 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)
	}
}

func TestCalculateExecutionCallsAtBoundary(t *testing.T) {
	req := calculateExecutionCalls(25)
	if len(req) != 1 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 25 {
		t.Errorf("expected end segment to be 25 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)
	}
}
func TestCalculateExecutionCallsOverTheBoundary(t *testing.T) {
	req := calculateExecutionCalls(26)
	if len(req) != 2 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 25 {
		t.Errorf("expected end segment to be 25 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)

	}

	part2 := req[1]
	if part2.EndSegment != 26 {
		t.Errorf("expected end segment to be 26 but was '%v'", part2.EndSegment)
	}
	if part2.StartSegment != 26 {
		t.Errorf("expected end segment to be 26 but was '%v'", part2.StartSegment)
	}
}

func TestNoParts(t *testing.T) {
	req := calculateExecutionCalls(0)
	if len(req) != 0 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
}

func TestHumanProvidesAccurateAggregation(t *testing.T) {
	b := Human(1024)
	if b != "1024 bytes" {
		t.Errorf("expected 1023 bytes but was %v", b)
	}
	b = Human(1025)
	if b != "1.00 kb" {
		t.Errorf("expected 1.00 kb but was %v", b)
	}

	b = Human(1048576)
	if b != "1024.00 kb" {
		t.Errorf("expected 1024.00 kb but was %v", b)
	}
	b = Human(1048577)
	if b != "1.00 mb" {
		t.Errorf("expected 1.00 mb but was %v", b)
	}

	b = Human(1048576 * 1024)
	if b != "1024.00 mb" {
		t.Errorf("expected 1024.00 kb but was %v", b)
	}
	b = Human(1048577 * 1024)
	if b != "1.00 gb" {
		t.Errorf("expected 1.00 gb but was %v", b)
	}
}

type MockClient struct {
	RetrieveByPackagePackage           Package
	RetrieveByPackageErr               error
	GetDownloadUrlsForFileErr          error
	GetDownloadUrlsForFileDownloadUrls []DownloadURL
	PackageIDs                         []string
	Packages                           []Package
	FileIds                            []string
	KeyCodes                           []string
	Starts                             []int
	Ends                               []int
}

func (m *MockClient) RetrievePackageByID(packageID string) (Package, error) {
	m.PackageIDs = append(m.PackageIDs, packageID)
	return m.RetrieveByPackagePackage, m.RetrieveByPackageErr
}

func (m *MockClient) GetDownloadUrlsForFile(p Package, fileID, keyCode string, start, end int) ([]DownloadURL, error) {
	m.Packages = append(m.Packages, p)
	m.FileIds = append(m.FileIds, fileID)
	m.KeyCodes = append(m.KeyCodes, keyCode)
	m.Starts = append(m.Starts, start)
	m.Ends = append(m.Ends, end)
	return m.GetDownloadUrlsForFileDownloadUrls, m.GetDownloadUrlsForFileErr
}

type MockDownloader struct {
	FileNames        []string
	Urls             []string
	Err              error
	SubDirToDownload string
	Pass             string
	KeyCode          string
}

func (m *MockDownloader) DownloadFile(fileName, url string) error {
	//file1 := filepath.Join(m.SubDirToDownload, "00010101T000000_", fileName)
	tmpFileName := strings.TrimSuffix(fileName, ".encrypted")
	token := make([]byte, 128)
	_, err := rand.Read(token)
	if err != nil {
		log.Fatalf("unable to readBytes %v", err)
	}
	if err := os.WriteFile(tmpFileName, token, 0600); err != nil {
		log.Fatalf("unable to write %v", err)
	}
	_, err = EncryptFile(tmpFileName, m.Pass+m.KeyCode)
	if err != nil {
		log.Fatalf("unable to encrypt %v", err)
	}
	m.FileNames = append(m.FileNames, fileName)
	m.Urls = append(m.Urls, url)
	return m.Err
}

func TestDownloadFiles(t *testing.T) {
	expectedKeyCode := "keyCode"
	expectedPackageID := "packageID1213"

	a := DownloadArgs{
		DownloadDir:      t.TempDir(),
		KeyCode:          expectedKeyCode,
		PackageID:        expectedPackageID,
		Verbose:          false,
		SubDirToDownload: filepath.Join(t.TempDir(), "testpackages"),
		MaxFileSizeByte:  1000000000,
		SkipList:         []string{},
	}

	mockClient := &MockClient{}
	p := Package{}
	p.ServerSecret = "serverSecretPassword"

	p.Files = []File{
		{FileID: "fileID1", FileName: "filename1.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID2", FileName: "filename2.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID3", FileName: "filename3.txt", Parts: 1, FileSize: 10},
	}
	mockClient.RetrieveByPackagePackage = p
	downloadURL := DownloadURL{}
	downloadURL.URL = "http://localhost:1999/filename1.txt"
	downloadURL.Part = 0
	mockClient.GetDownloadUrlsForFileDownloadUrls = []DownloadURL{downloadURL}
	mockDownloader := &MockDownloader{}
	mockDownloader.Pass = p.ServerSecret
	mockDownloader.KeyCode = expectedKeyCode
	mockDownloader.SubDirToDownload = a.SubDirToDownload
	_, invalidFiles, err := DownloadFilesFromPackage(mockClient, mockDownloader, a)
	if err != nil {
		t.Fatal(err)
	}
	if len(invalidFiles) > 0 {
		t.Errorf("expected no invalid files but had %v", len(invalidFiles))
	}
	if len(mockDownloader.FileNames) != 3 {
		t.Errorf("expected 3 entries but had %v: output %#v", len(mockDownloader.FileNames), mockDownloader.FileNames)
	}
}

func TestSkipFilesOnDownload(t *testing.T) {
	expectedKeyCode := "keyCode"
	expectedPackageID := "packageID1213"

	a := DownloadArgs{
		DownloadDir:      t.TempDir(),
		KeyCode:          expectedKeyCode,
		PackageID:        expectedPackageID,
		Verbose:          false,
		SubDirToDownload: filepath.Join(t.TempDir(), "testpackages"),
		MaxFileSizeByte:  1000000000,
		SkipList:         []string{"fileID1", "fileID2", "fileID3"},
	}

	mockClient := &MockClient{}
	p := Package{}
	p.ServerSecret = "serverSecretPassword"

	p.Files = []File{
		{FileID: "fileID1", FileName: "filename1.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID2", FileName: "filename2.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID3", FileName: "filename3.txt", Parts: 1, FileSize: 10},
	}
	mockClient.RetrieveByPackagePackage = p
	downloadURL := DownloadURL{}
	downloadURL.URL = "http://localhost:1999/filename1.txt"
	downloadURL.Part = 0
	mockClient.GetDownloadUrlsForFileDownloadUrls = []DownloadURL{downloadURL}
	mockDownloader := &MockDownloader{}
	mockDownloader.Pass = p.ServerSecret
	mockDownloader.KeyCode = expectedKeyCode
	mockDownloader.SubDirToDownload = a.SubDirToDownload
	_, invalidFiles, err := DownloadFilesFromPackage(mockClient, mockDownloader, a)
	if err != nil {
		t.Fatal(err)
	}
	if len(invalidFiles) > 0 {
		t.Errorf("expected no invalid files but had %v", len(invalidFiles))
	}
	if len(mockDownloader.FileNames) != 0 {
		t.Errorf("expected no entries but had %v", len(mockDownloader.FileNames))
	}
}

func TestSkipFilesOnDownloadWithOnlyOneSkip(t *testing.T) {
	expectedKeyCode := "keyCode"
	expectedPackageID := "packageID1213"

	a := DownloadArgs{
		DownloadDir:      t.TempDir(),
		KeyCode:          expectedKeyCode,
		PackageID:        expectedPackageID,
		Verbose:          false,
		SubDirToDownload: filepath.Join(t.TempDir(), "testpackages"),
		MaxFileSizeByte:  1000000000,
		SkipList:         []string{"fileID3"},
	}

	mockClient := &MockClient{}
	p := Package{}
	p.ServerSecret = "serverSecretPassword"

	p.Files = []File{
		{FileID: "fileID1", FileName: "filename1.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID2", FileName: "filename2.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID3", FileName: "filename3.txt", Parts: 1, FileSize: 10},
	}
	mockClient.RetrieveByPackagePackage = p
	downloadURL := DownloadURL{}
	downloadURL.URL = "http://localhost:1999/filename1.txt"
	downloadURL.Part = 0
	mockClient.GetDownloadUrlsForFileDownloadUrls = []DownloadURL{downloadURL}
	mockDownloader := &MockDownloader{}
	mockDownloader.Pass = p.ServerSecret
	mockDownloader.KeyCode = expectedKeyCode
	mockDownloader.SubDirToDownload = a.SubDirToDownload
	_, invalidFiles, err := DownloadFilesFromPackage(mockClient, mockDownloader, a)
	if err != nil {
		t.Fatal(err)
	}
	if len(invalidFiles) > 0 {
		t.Errorf("expected no invalid files but had %v", len(invalidFiles))
	}
	if len(mockDownloader.FileNames) != 2 {
		t.Errorf("expected 2 entries but had %v: output %#v", len(mockDownloader.FileNames), mockDownloader.FileNames)
	}
}

func TestSkipFilesOverLimit(t *testing.T) {
	expectedKeyCode := "keyCode"
	expectedPackageID := "packageID1213"

	a := DownloadArgs{
		DownloadDir:      t.TempDir(),
		KeyCode:          expectedKeyCode,
		PackageID:        expectedPackageID,
		Verbose:          false,
		SubDirToDownload: filepath.Join(t.TempDir(), "testpackages"),
		MaxFileSizeByte:  9,
		SkipList:         []string{},
	}

	mockClient := &MockClient{}
	p := Package{}
	p.ServerSecret = "serverSecretPassword"

	p.Files = []File{
		{FileID: "fileID1", FileName: "filename1.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID2", FileName: "filename2.txt", Parts: 1, FileSize: 10},
		{FileID: "fileID3", FileName: "filename3.txt", Parts: 1, FileSize: 10},
	}
	mockClient.RetrieveByPackagePackage = p
	downloadURL := DownloadURL{}
	downloadURL.URL = "http://localhost:1999/filename1.txt"
	downloadURL.Part = 0
	mockClient.GetDownloadUrlsForFileDownloadUrls = []DownloadURL{downloadURL}
	mockDownloader := &MockDownloader{}
	mockDownloader.Pass = p.ServerSecret
	mockDownloader.KeyCode = expectedKeyCode
	mockDownloader.SubDirToDownload = a.SubDirToDownload
	_, invalidFiles, err := DownloadFilesFromPackage(mockClient, mockDownloader, a)
	if err != nil {
		t.Fatal(err)
	}
	if len(invalidFiles) > 0 {
		t.Errorf("expected no invalid files but had %v", len(invalidFiles))
	}
	if len(mockDownloader.FileNames) != 0 {
		t.Errorf("expected no entries but had %v", len(mockDownloader.FileNames))
	}
}
