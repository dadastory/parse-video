package parser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRedBookPageVariants(t *testing.T) {
	for _, tc := range []struct {
		name, state string
		images      int
		video, live string
		failed      bool
	}{
		{name: "mobile album", state: `{"noteData":{"hasError":false,"data":{"noteData":{"title":"Public note","imageList":[{"url":"https://cdn.example/one.jpg","livePhoto":true,"stream":{"h264":[{"masterUrl":"https://cdn.example/motion.mp4"}]}},{"infoList":[{"imageScene":"WB_DFT","url":"https://cdn.example/two.webp"}]}]}}}}`, images: 2, live: "https://cdn.example/motion.mp4"},
		{name: "legacy album", state: `{"note":{"currentNoteId":"note","noteDetailMap":{"note":{"note":{"title":"Public note","imageList":[{"urlDefault":"https://cdn.example/legacy.jpg"}]}}}}}`, images: 1},
		{name: "mobile video", state: `{"noteData":{"data":{"noteData":{"title":"Public note","video":{"media":{"stream":{"h264":[{"masterUrl":"https://cdn.example/video.mp4"}]}}}}}}}`, video: "https://cdn.example/video.mp4"},
		{name: "empty shell", state: `{"global":{},"feed":{}}`, failed: true},
		{name: "platform error", state: `{"noteData":{"hasError":true,"data":{"noteData":{"imageList":[{"url":"https://cdn.example/error.jpg"}]}}}}`, failed: true},
		{name: "metadata only", state: `{"noteData":{"data":{"noteData":{"title":"Public note"}}}}`, failed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.UserAgent(), "iPhone") {
					t.Errorf("expected mobile page request, got %q", r.UserAgent())
				}
				fmt.Fprint(w, "<script>window.__INITIAL_STATE__ =\n"+tc.state+";</script>")
			}))
			defer server.Close()
			result, err := (redBook{}).parseShareUrl(server.URL)
			if tc.failed {
				if err == nil {
					t.Fatal("missing media must fail")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result.Title != "Public note" || len(result.Images) != tc.images || result.VideoUrl != tc.video {
				t.Fatalf("unexpected result: %+v", result)
			}
			if tc.live != "" && result.Images[0].LivePhotoUrl != tc.live {
				t.Fatal("missing live photo")
			}
		})
	}
}
