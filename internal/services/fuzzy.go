package services

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func DisplayFuzzyFind(items []Item) ([]int, error) {
	idxs, err := fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			if items[i].Content != nil {
				return "📁 " + items[i].Name
			}
			return "🎵 " + items[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			if items[i].Content != nil {
				return fmt.Sprintf("Playlist: %s\nAuthor: %s\nViews: %s\nContent:\n%s",
					items[i].Name,
					items[i].Owner,
					items[i].Views,
					strings.Join(items[i].Content, "\n"),
				)
			}
			return fmt.Sprintf("Song: %s\nAuthor: %s\nViews: %s", items[i].Name, items[i].Owner, items[i].Views)
		}), fuzzyfinder.WithHeader("Search for a song"))
	return idxs, err
}
