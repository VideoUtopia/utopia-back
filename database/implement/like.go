package implement

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"utopia-back/database/abstract"
	"utopia-back/model"
)

type LikeDal struct{ Db *gorm.DB }

func (l *LikeDal) GetLikeUserId(videoId uint) (user []uint, err error) {
	res := l.Db.Model(&model.Like{}).Where("video_id = ? AND status = ?", videoId, true).Find(&user)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = res.Error
	}
	return user, err
}

func (l *LikeDal) Like(userId uint, videoId uint) (err error) {
	res := l.Db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"status": true}),
		}).Create(&model.Like{
		UserID:  userId,
		VideoID: videoId,
		Status:  true,
	})
	return res.Error
}

func (l *LikeDal) UnLike(userId uint, videoId uint) (err error) {
	res := l.Db.Model(&model.Like{}).Where("user_id = ? AND video_id = ?", userId, videoId).Update("status", false)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = res.Error
	}
	return err
}

func (l *LikeDal) IsLike(userId uint, videoId uint) (isLike bool, err error) {
	var like model.Like
	res := l.Db.Where("user_id = ? AND video_id = ?", userId, videoId).First(&like)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = res.Error
	}
	return like.Status, err
}

func (l *LikeDal) GetLikeCount(videoId uint) (count int64, err error) {
	var likeCount int64
	res := l.Db.Model(&model.Like{}).Where("video_id = ?", videoId).Count(&likeCount)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = res.Error
	}
	return likeCount, err
}

var _ abstract.LikeDal = (*LikeDal)(nil)
