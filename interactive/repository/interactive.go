package repository

import (
	"context"
	"github.com/mymikasa/skk/interactive/domain"
	"github.com/mymikasa/skk/interactive/repository/cache"
	"github.com/mymikasa/skk/interactive/repository/dao"
	"github.com/mymikasa/skk/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId int64, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bizId, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.LoggerV1
}

func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}

	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, bizId int64, cid, uid int64) error {
	err := c.dao.InsertCollectBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Cid:   cid,
		Uid:   uid,
	})

	if err != nil {
		return err
	}

	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}

	ie, err := c.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}

	//if err == nil {
	res := c.toDomain(ie)
	err = c.cache.Set(ctx, biz, bizId, res)

	if err != nil {
		c.l.Error("回写缓存失败",
			logger.String("biz", biz),
			logger.Int64("bizId", bizId),
			logger.Error(err))
	}
	return res, nil
	//}
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	intrs, err := c.dao.GetByIds(ctx, biz, ids)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Interactive, len(intrs))
	for _, intr := range intrs {
		result = append(result, c.toDomain(intr))
	}
	return result, nil
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao}
}

func (c *CachedInteractiveRepository) toDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizId:      ie.BizId,
		ReadCnt:    ie.ReadCnt,
		LikeCnt:    ie.LikeCnt,
		CollectCnt: ie.CollectCnt,
	}
}
