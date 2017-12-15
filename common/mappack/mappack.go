package mappack

import (
	"archive/zip"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	bettererrors "github.com/xtuc/better-errors"
)

type MappackInMemoryArchive struct {
	Zip   zip.ReadCloser
	Files map[string][]byte
}

func UnzipAndGetHandles(filename string) (*MappackInMemoryArchive, error) {
	mappackInMemoryArchive := &MappackInMemoryArchive{
		Files: make(map[string][]byte),
	}

	reader, err := zip.OpenReader(filename)

	if err != nil {
		better := bettererrors.
			New("Could not open archive").
			With(err).
			SetContext("filename", filename)

		return nil, better
	}

	for _, file := range reader.File {
		fd, err := file.Open()

		if err != nil {
			better := bettererrors.
				New("Could not open file in archive").
				With(err).
				SetContext("filename", file.Name)

			return nil, better
		}

		content, readErr := ioutil.ReadAll(fd)

		if readErr != nil {
			berror := bettererrors.
				New("Could not read file in archive").
				With(err).
				SetContext("filename", file.Name)

			return nil, berror
		}

		mappackInMemoryArchive.Files[file.Name] = content
		fd.Close()
	}

	return mappackInMemoryArchive, nil
}

func (m *MappackInMemoryArchive) Open(name string) ([]byte, error) {
	if content, hasFile := m.Files[name]; hasFile {
		return content, nil
	}

	berror := bettererrors.
		New("File not found").
		SetContext("filename", name)

	return nil, berror
}

func (m *MappackInMemoryArchive) Close() {
	m.Zip.Close()
}

func (m *MappackInMemoryArchive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	initialExt := filepath.Ext(r.URL.Path)

	if strings.HasSuffix(r.URL.Path, "model.json") {
		r.URL.Path += ".gz"
		w.Header().Set("Content-Encoding", "gzip")
	}

	content, err := m.Open(r.URL.Path)

	if err != nil {
		http.NotFound(w, r)
	} else {
		ctype := mime.TypeByExtension(initialExt)

		w.Header().Set("Content-Type", ctype)
		w.Header().Set("Content-Size", strconv.Itoa(len(content)))
		w.Write(content)
	}
}
