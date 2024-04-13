package repository

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"golang.org/x/net/context"
)

type interactionRepository struct {
	db orm.Database
}

func NewInteractionRepository(db orm.Database) domain.InteractionRepository {
	return &interactionRepository{
		db: db,
	}
}

func (repo *interactionRepository) IncrReadCount(c context.Context, biz string, id int64) error {
	// upset!
	//filter := map[string]interface{}{
	//	"read_cnt": "read_cnt + 1",
	//  "update_at": ...
	//}
	//repo.db.UpdateOne()
	return nil
}
