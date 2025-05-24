package client

import (
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
	"strings"
)

// These are the file types supported by (device asset) FileUpload.
var fileExtToContentType = map[string]string{
	".gif": "image/gif",
	".jpg": "image/jpeg",
	".png": "image/png",
	".mp3": "audio/mpeg",
	".mp4": "audio/mp4",
	".wav": "audio/wave",
	".caf": "audio/x-caf",
}

var acceptableTypes = func() string {
	keys := []string{}
	for key := range maps.Keys(fileExtToContentType) {
		keys = append(keys, key[1:])
	}
	return strings.Join(keys, ", ")
}()

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFileProtect(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))

	ext := filepath.Ext(filename)
	contentType, ok := fileExtToContentType[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file type, got: '%s', acceptable types: '%s'", ext[1:], acceptableTypes)
	}

	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
