package translation

import (
	"github.com/aexvir/lnk/internal/storage"
	"github.com/aexvir/lnk/proto"
)

// DbLinkToProto translates a storage link model to its proto link model counterpart.
func DbLinkToProto(link *storage.Link) *proto.LinkDetails {
	stats := make([]*proto.DailyHits, 0, len(link.Histogram))
	for date, count := range link.Histogram {
		stats = append(stats, &proto.DailyHits{Date: date, Hits: count})
	}

	return &proto.LinkDetails{
		Slug:   link.Slug,
		Target: link.Target,
		Hits:   link.Hits,
		Stats:  stats,
	}
}
