package subscriber

import (
	"context"
)

// Filter is a filter for addresses.
// It is used to filter transactions.
// It could be in future bloom filter.
type Filter interface {
	Test(ctx context.Context, address string) bool
}

var _ Filter = (*filter)(nil)

type filter struct {
	subscriber Subscriber
}

func NewFilter(subscriber Subscriber) Filter {
	return &filter{
		subscriber: subscriber,
	}
}

func (f *filter) Test(ctx context.Context, address string) bool {
	exists, err := f.subscriber.Test(ctx, address)
	if err != nil {
		return false
	}

	return exists
}
