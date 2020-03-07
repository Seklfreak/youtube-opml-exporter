package pkg

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/api/youtube/v3"
)

const (
	documentFmt string = `
<opml version="1.1">
    <body>
        <outline text="YouTube Subscriptions" title="YouTube Subscriptions">
%s
        </outline>
    </body>
</opml>
`
	itemFmt string = `
            <outline text="%s" title="%s" type="rss"
                     xmlUrl="https://www.youtube.com/feeds/videos.xml?channel_id=%s"/>
`
)

func ExportHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refreshToken")
	if refreshToken == "" {
		http.Error(w, "please specify refresh token", http.StatusBadRequest)
		return
	}

	yt, err := ytServiceFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var items []*youtube.Subscription

	var pageToken string
	for {
		logger.Info("making YouTube API request", zap.String("page_token", pageToken))

		resp, err := yt.Subscriptions.
			List("snippet").
			MaxResults(50).
			Order("alphabetical").
			Mine(true).
			PageToken(pageToken).
			Do()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items = append(items, resp.Items...)

		pageToken = resp.NextPageToken
		if resp.NextPageToken == "" {
			break
		}
	}

	logger.Info("found items", zap.Int("amount", len(items)))

	var result string
	for _, item := range items {
		if item == nil || item.Snippet == nil || item.Snippet.ResourceId == nil ||
			item.Snippet.ResourceId.ChannelId == "" {
			logger.Warn("skipping item because it does not have all required fields", zap.Any("item", item))
			continue
		}

		result += fmt.Sprintf(itemFmt, item.Snippet.Title, item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
	}

	result = fmt.Sprintf(documentFmt, result)

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, result)
}
