package implement

import (
	"gorm.io/gorm"
	"utopia-back/model"
)

type VideoDal struct {
	Db *gorm.DB
}

func (v *VideoDal) CreateVideo(video *model.Video) (id uint, err error) {
	res := v.Db.Create(&video)
	if res.Error != nil {
		return 0, res.Error
	}
	return video.ID, nil
}

func (v *VideoDal) IsVideoExist(videoId uint) (err error) {
	res := v.Db.Where("id = ?", videoId).First(&model.Video{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
