// Package elasticsearch provides support for access the elasticsearch.
package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
	"time"
)

// Config is the required properties to use the elasticsearch.
type Config struct {
	URL   string
	Sniff bool
	Gzip  bool
}

// Open knows how to open an elasticsearch connection based on the config.
func Open(cfg Config) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(cfg.URL),
		elastic.SetSniff(cfg.Sniff),
		elastic.SetGzip(cfg.Gzip),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// StatusCheck returns nil if it can successfully talk to the elasticsearch.
// It returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, client *elastic.Client, url string) error {
	var pingError error
	for attempts := 1; ; attempts++ {
		_, _, pingError = client.Ping(url).Do(ctx)
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil

}