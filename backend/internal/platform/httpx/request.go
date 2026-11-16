package httpx

import (
	"encoding/json"
	"mime"
	"net/http"
)

// maxBodyBytes caps every JSON request body.
const maxBodyBytes = 1 << 20

// DecodeJSON reads a JSON request body into v. Bodies not sent as
// application/json are rejected, which also keeps HTML forms from other sites
// from reaching state-changing endpoints.
func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return ErrUnsupportedMediaType()
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(v); err != nil {
		return ErrMalformedRequest()
	}
	return nil
}
