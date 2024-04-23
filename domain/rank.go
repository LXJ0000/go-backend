package domain

import "context"

type RankUsercase interface {
	TopN(c context.Context) error
}
