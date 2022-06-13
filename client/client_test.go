package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aexvir/lnk/proto"
)

func TestClientCreateLink(t *testing.T) {
	tests := map[string]struct {
		target string // used as flag to control server behaviour on this test
		slug   string

		wantLink string
		wantErr  string
	}{
		"server error": {
			target:  "servererr",
			wantErr: "request failed",
		},
		"successful request": {
			target:   "random",
			wantLink: "random",
		},
		"successful request with custom slug": {
			target:   "specific",
			slug:     "test",
			wantLink: "test",
		},
	}

	downstream := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var payload proto.CreateLinkReq
				err := json.NewDecoder(r.Body).Decode(&payload)
				require.NoError(t, err, "the client should always serialize the payload correctly")

				switch payload.Target {
				case "random":
					w.WriteHeader(http.StatusOK)
					err = json.NewEncoder(w).Encode(proto.LinkId{Slug: "random"})
					if err != nil {
						t.Fatal(err)
					}
				case "specific":
					w.WriteHeader(http.StatusOK)
					err = json.NewEncoder(w).Encode(proto.LinkId{Slug: *payload.Slug})
					if err != nil {
						t.Fatal(err)
					}
				case "malformed":
					w.WriteHeader(http.StatusNotFound)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
		),
	)
	defer downstream.Close()

	client, err := NewLnkClient(WithBaseUrl(downstream.URL))
	require.NoError(t, err, "nothing to fail here for")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var slug *string
			if test.slug != "" {
				slug = &test.slug
			}
			gotLink, gotErr := client.CreateLink(test.target, slug)

			if test.wantErr != "" {
				require.Error(t, gotErr, "should have failed on this request")
				assert.Contains(t, gotErr.Error(), test.wantErr, "got an error, but not the correct one")
				return
			}

			require.NoError(t, gotErr, "shouldn't error here")
			require.NotNil(t, gotLink, "there should be a non-nil link here")
			assert.Equal(t, test.wantLink, gotLink.Slug, "unexpected link returned")
		})
	}
}

func TestClientGetLink(t *testing.T) {
	testlink := proto.LinkDetails{
		Slug:   "exists",
		Target: "http://google.com",
		Hits:   42,
		Stats: []*proto.DailyHits{
			{
				Date: "2022-06-12",
				Hits: 42,
			},
		},
	}

	tests := map[string]struct {
		slug string

		wantLink proto.LinkDetails
		wantErr  string
	}{
		"empty slug": {
			wantErr: "unexpected response status code",
		},
		"correct slug but missing": {
			slug:    "missing",
			wantErr: "not found",
		},
		"correct slug but malformed response": {
			slug:    "malformed",
			wantErr: "error decoding response",
		},
		"correct slug": {
			slug:     "exists",
			wantLink: testlink,
		},
	}

	downstream := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/links/exists":
					w.WriteHeader(http.StatusOK)
					err := json.NewEncoder(w).Encode(testlink)
					require.NoError(t, err, "shouldn't fail encoding the response")
				case "/api/links/malformed":
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("hehe"))
				case "/api/links/missing":
					w.WriteHeader(http.StatusNotFound)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
		),
	)
	defer downstream.Close()

	client, err := NewLnkClient(WithBaseUrl(downstream.URL))
	require.NoError(t, err, "nothing to fail here for")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotLink, gotErr := client.GetLink(test.slug)

			if test.wantErr != "" {
				require.Error(t, gotErr, "should have failed on this request")
				assert.Contains(t, gotErr.Error(), test.wantErr, "got an error, but not the correct one")
				return
			}

			require.NoError(t, gotErr, "shouldn't error here")
			require.NotNil(t, gotLink, "there should be a non-nil link here")
			assert.Equal(t, test.wantLink, *gotLink, "the returned link is not matching expectations")
		})
	}

}
