package repository

import (
	"errors"
	"fmt"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/LXJ0000/go-backend/script"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type interactionRepository struct {
	dao   orm.Database
	cache cache.RedisCache
}

func NewInteractionRepository(dao orm.Database, cache cache.RedisCache) domain.InteractionRepository {
	return &interactionRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *interactionRepository) GetByIDs(c context.Context, biz string, ids []int64) (map[int64]domain.Interaction, error) {
	mp := map[int64]domain.Interaction{}
	var interactions []domain.Interaction

	if err := repo.dao.Raw(c).Model(&domain.Interaction{}).Where("biz = ? and biz_id in (?)", biz, ids).Find(&interactions).Error; err != nil {
		return nil, err
	}

	for _, interaction := range interactions {
		mp[interaction.BizID] = interaction
	}
	return mp, nil
}

// BatchIncrReadCount 批量增加read_cnt 需保证 len(biz) == len(id)
func (repo *interactionRepository) BatchIncrReadCount(c context.Context, biz []string, id []int64) error {
	fn := func(tx *gorm.DB) error {
		//now := time.Now().UnixMicro()
		dao := orm.NewDatabase(tx)
		update := map[string]interface{}{
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
			//"updated_at": now,
		}

		for i := 0; i < len(biz); i++ {
			i := i // 1.22 可不写
			create := &domain.Interaction{
				BizID:   id[i],
				Biz:     biz[i],
				ReadCnt: 1,
			}
			//create.CreatedAt = now
			//create.UpdatedAt = now

			if err := dao.UpsertOne(c, &domain.Interaction{}, update, create); err != nil {
				slog.Error("IncrReadCount Fail", "Error", err.Error(), "biz", biz[i], "biz_id", id[i])
			}
			go func() { // TODO new lua script or pipeline
				if err := repo.cacheIncrCnt(context.Background(), biz[i], id[i], "read_cnt"); err != nil {
					slog.Warn("Redis Op Fail With CacheIncrReadCnt", "Error", err.Error(), "biz", biz[i], "bizID", id[i])
				}
			}()
		}
		return nil
	}
	_ = repo.dao.Transaction(c, fn)
	return nil
}

func (repo *interactionRepository) IncrReadCount(c context.Context, biz string, id int64) error {
	//now := time.Now().UnixMicro()
	update := map[string]interface{}{
		"read_cnt": gorm.Expr("`read_cnt` + 1"),
		//"updated_at": now,
	}
	create := &domain.Interaction{
		BizID:   id,
		Biz:     biz,
		ReadCnt: 1,
	}
	//create.CreatedAt = now
	//create.UpdatedAt = now
	if err := repo.dao.UpsertOne(c, &domain.Interaction{}, update, create); err != nil {
		return err
	}
	go func() {
		if err := repo.cacheIncrCnt(context.Background(), biz, id, "read_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheIncrReadCnt", "Error", err.Error(), "biz", biz, "bizID", id)
		}
	}()
	// 数据库操作成功即认为业务处理成功
	return nil
}

func (repo *interactionRepository) Like(c context.Context, biz string, bizID, userID int64) error {
	//now := time.Now().UnixMicro()
	updateInteraction := map[string]interface{}{
		"like_cnt": gorm.Expr("`like_cnt` + 1"),
		//"updated_at": now,
	}
	createInteraction := &domain.Interaction{
		BizID:   bizID,
		Biz:     biz,
		LikeCnt: 1,
	}
	//createInteraction.CreatedAt = now
	//createInteraction.UpdatedAt = now
	updateUserLike := map[string]interface{}{
		"status": true,
		//"updated_at": now,
	}
	createUserLike := &domain.UserLike{
		BizID:  bizID,
		Biz:    biz,
		UserID: userID,
		Status: true,
	}
	//createUserLike.CreatedAt = now
	//createUserLike.UpdatedAt = now
	fn := func(tx *gorm.DB) error {
		dao := orm.NewDatabase(tx)
		if err := dao.UpsertOne(c, &domain.Interaction{}, updateInteraction, createInteraction); err != nil {
			return err
		}
		return dao.UpsertOne(c, &domain.UserLike{}, updateUserLike, createUserLike)
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheIncrCnt(context.Background(), biz, bizID, "like_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheIncrLikeCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) CancelLike(c context.Context, biz string, bizID, userID int64) error {
	fn := func(tx *gorm.DB) error {
		//now := time.Now().UnixMicro()
		dao := orm.NewDatabase(tx)
		//1. 更新 UserLike status = false
		createUserLike := &domain.UserLike{
			UserID: userID,
			BizID:  bizID,
			Biz:    biz,
		}
		//createUserLike.CreatedAt = now
		//createUserLike.UpdatedAt = now
		if err := dao.UpdateOne(c,
			&domain.UserLike{},
			createUserLike,
			map[string]interface{}{
				"status": false,
				//"updated_at": now,
			},
		); err != nil {
			return err
		}
		//2. 更新 interaction like_cnt - 1
		createLikeCnt := &domain.Interaction{
			BizID: bizID,
			Biz:   biz,
		}
		//createUserLike.CreatedAt = now
		//createUserLike.UpdatedAt = now
		if err := dao.UpdateOne(c,
			&domain.Interaction{},
			createLikeCnt,
			map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				//"updated_at": now,
			},
		); err != nil {
			return err
		}
		return nil
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheDecrCnt(context.Background(), biz, bizID, "like_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheDecrLikeCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) Stat(c context.Context, biz string, bizID, userID int64) (domain.Interaction, domain.UserInteractionStat, error) {
	var isLike, isCollect bool
	var err error
	var interaction domain.Interaction
	eg := errgroup.Group{}
	eg.Go(func() error {
		isLike, err = repo.isLike(c, biz, bizID, userID)
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		isCollect, err = repo.isCollect(c, biz, bizID, userID)
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		res, err := repo.cache.HGetAll(c, key(biz, bizID))
		if err == nil && len(res) > 0 {
			interaction.CollectCnt, _ = strconv.Atoi(res["collect_cnt"])
			interaction.ReadCnt, _ = strconv.Atoi(res["read_cnt"])
			interaction.LikeCnt, _ = strconv.Atoi(res["like_cnt"])
			return nil
		}
		var interaction domain.Interaction
		filter := map[string]interface{}{
			"biz_id": bizID,
			"biz":    biz,
		}
		if err := repo.dao.FindOne(c, &domain.Interaction{}, filter, &interaction); err != nil {
			return nil
		}
		go func() {
			if err := repo.cache.HSet(context.Background(),
				key(biz, bizID),
				"read_cnt", interaction.ReadCnt,
				"collect_cnt", interaction.CollectCnt,
				"like_cnt", interaction.LikeCnt,
			); err != nil {
				slog.Warn("Redis Op Fail With HSet", "Error", err.Error(), "biz", biz, "bizID", bizID, "Key", key(biz, bizID))
			}
		}()
		return nil
	})
	if err := eg.Wait(); err != nil {
		return domain.Interaction{}, domain.UserInteractionStat{}, err
	}
	return interaction,
		domain.UserInteractionStat{
			Liked:     isLike,
			Collected: isCollect,
		}, nil
}

func (repo *interactionRepository) Collect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	//now := time.Now().UnixMicro()
	updateInteraction := map[string]interface{}{
		"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
		//"updated_at":  now,
	}
	createInteraction := &domain.Interaction{
		BizID:      bizID,
		Biz:        biz,
		CollectCnt: 1,
	}
	//createInteraction.CreatedAt = now
	//createInteraction.UpdatedAt = now
	updateUserCollect := map[string]interface{}{
		"status": true,
		//"updated_at": now,
	}
	createUserCollect := &domain.UserCollect{
		BizID:        bizID,
		Biz:          biz,
		UserID:       userID,
		CollectionID: collectionID,
		Status:       true,
	}
	//createUserCollect.CreatedAt = now
	//createUserCollect.UpdatedAt = now
	fn := func(tx *gorm.DB) error {
		dao := orm.NewDatabase(tx)
		if err := dao.UpsertOne(c, &domain.Interaction{}, updateInteraction, createInteraction); err != nil {
			return err
		}
		return dao.UpsertOne(c, &domain.UserCollect{}, updateUserCollect, createUserCollect)
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheIncrCnt(context.Background(), biz, bizID, "collect_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheIncrCollectCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	fn := func(tx *gorm.DB) error {
		//now := time.Now().UnixMicro()
		dao := orm.NewDatabase(tx)
		//1. 更新 UserLike status = false
		createUserCollect := &domain.UserCollect{
			UserID:       userID,
			BizID:        bizID,
			CollectionID: collectionID,
			Biz:          biz,
		}
		//createUserCollect.CreatedAt = now
		//createUserCollect.UpdatedAt = now
		if err := dao.UpdateOne(c,
			&domain.UserCollect{},
			createUserCollect,
			map[string]interface{}{
				"status": false,
				//"updated_at": now,
			},
		); err != nil {
			return err
		}
		//2. 更新 interaction like_cnt - 1
		createCollectCnt := &domain.Interaction{
			BizID: bizID,
			Biz:   biz,
		}
		//createUserCollect.CreatedAt = now
		//createUserCollect.UpdatedAt = now
		if err := dao.UpdateOne(c,
			&domain.Interaction{},
			createCollectCnt,
			map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` - 1"),
				//"updated_at":  now,
			},
		); err != nil {
			return err
		}
		return nil
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheDecrCnt(context.Background(), biz, bizID, "collect_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheDecrCollectCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) isLike(c context.Context, biz string, bizID, userID int64) (bool, error) {
	var item domain.UserLike
	filter := map[string]interface{}{
		"user_id": userID,
		"biz_id":  bizID,
		"biz":     biz,
	}
	if err := repo.dao.FindOne(c, &domain.UserLike{}, filter, &item); err != nil {
		return false, err
	}
	return item.Status, nil
}

func (repo *interactionRepository) isCollect(c context.Context, biz string, bizID, userID int64) (bool, error) {
	var item domain.UserCollect
	filter := map[string]interface{}{
		"user_id": userID,
		"biz_id":  bizID,
		"biz":     biz,
	}
	if err := repo.dao.FindOne(c, &domain.UserCollect{}, filter, &item); err != nil {
		return false, err
	}
	return item.Status, nil
}

func (repo *interactionRepository) cacheIncrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.LuaWithReturnInt(c, script.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, 1)
	return err
}

func (repo *interactionRepository) cacheDecrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.LuaWithReturnInt(c, script.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, -1)
	return err
}

func key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}
