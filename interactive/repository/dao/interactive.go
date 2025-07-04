package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizid int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	InsertCollectBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{db: db}
}

// IncrReadCnt 已读+1
func (dao *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()

	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("read_cnt + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

// BatchIncrReadCnt 批量增加阅读数
func (dao *GORMInteractiveDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)

		for i := 0; i < len(bizs); i++ {
			err := txDAO.IncrReadCnt(ctx, bizs[i], bizIds[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizid int64, uid int64) error {
	now := time.Now().UnixMilli()

	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			Biz:    biz,
			BizId:  bizid,
			Status: 1,
			Utime:  now,
			Ctime:  now,
		}).Error

		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` + 1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   bizid,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()

	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("uid = ? AND biz_id= ? AND biz=?", uid, bizId, biz).Updates(map[string]interface{}{
			"utime":  now,
			"status": 0,
		}).Error

		if err != nil {
			return err
		}

		return tx.Model(&Interactive{}).Where("biz=? AND biz_id=?").Updates(
			map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` + 1"),
				"utime":    now,
			}).Error
	})
}

func (dao *GORMInteractiveDAO) InsertCollectBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.Ctime = now
	cb.Utime = now

	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&cb).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			Biz:        cb.Biz,
			BizId:      cb.BizId,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ? AND status = ?",
			biz, bizId, uid, 1).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetCollectInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ?", biz, bizId, uid).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, bizId).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var res []Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id IN ?", biz, ids).
		Find(&res).Error
	return res, err
}

type UserLikeBiz struct {
	Id     int64  `gorm:"primary_key;autoIncrement"`
	Uid    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	Status int
	Utime  int64
	Ctime  int64
}

type UserCollectionBiz struct {
	Id    int64  `gorm:"primary_key;autoIncrement"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`

	Cid   int64 `gorm:"index"`
	Utime int64
	Ctime int64
}

type Interactive struct {
	Id    int64  `gorm:"primary_key"`
	Biz   string `gorm:"type:varchar(255); uniqueIndex:biz_type_ip"`
	BizId int64  `gorm:"uniqueIndex:biz_type_ip"`

	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Utime      int64
	Ctime      int64
}
