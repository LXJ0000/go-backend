package domain

import "context"

type RankUsecase interface {
	TopN(c context.Context) error
}

type RankRepository interface {
}
