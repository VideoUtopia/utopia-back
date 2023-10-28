package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"net/http"
	"strconv"
	"utopia-back/config"
	utils "utopia-back/pkg/util"
	"utopia-back/service/abstract"
	v1 "utopia-back/service/implement/v1"
)

type VideoController struct {
	VideoService abstract.VideoService
}

func NewVideoController() *VideoController {
	return &VideoController{VideoService: v1.NewVideoService()}
}

type uploadVideoTokenData struct {
	Token string `json:"token"`
}

type uploadCallbackReq struct {
	Key       string `json:"key" validate:"required"`
	IsImage   string `json:"is_image" validate:"required"`
	AuthorId  string `json:"author_id"` // todo 可更改为JWT-Token，增强安全性
	VideoType string `json:"video_type"`
	CoverUrl  string `json:"cover_url"`
	Describe  string `json:"describe"`
}

type uploadCallbackData struct {
	ImageUrl string `json:"image_url"`
}

// 上传是否为图片
const callbackIsImage = "YES"

func (v *VideoController) UploadVideoToken(c *gin.Context) {
	upToken := utils.GetCallbackToken()
	c.JSON(http.StatusOK, &ResponseWithData{
		Code: SuccessCode,
		Msg:  "ok",
		Data: uploadVideoTokenData{
			Token: upToken,
		},
	})
}

func (v *VideoController) UploadVideoCallback(c *gin.Context) {
	var (
		r                   uploadCallbackReq
		err                 error
		authorId, videoType uint
	)

	// 请求处理失败，返回错误信息
	defer func() {
		if err != nil {
			c.JSON(http.StatusOK, &ResponseWithoutData{
				Code: ErrorCode,
				Msg:  err.Error(),
			})
		}
	}()
	// 校验是否为七牛云调用
	isQiNiu, err := qbox.VerifyCallback(utils.GetMac(), c.Request)
	if !isQiNiu || err != nil {
		if err == nil {
			err = errors.New("非七牛云服务发送")
		}
		return
	}
	// 接收参数并绑定
	if err = c.BindJSON(&r); err != nil {
		return
	}
	// 参数校验
	authorId, videoType, err = callbackReqValidate(r)
	if err != nil {
		return
	}
	url := config.V.GetString("qiniu.kodoApi") + "/" + r.Key
	if r.IsImage == callbackIsImage {
		c.JSON(http.StatusOK, &ResponseWithData{
			Code: SuccessCode,
			Msg:  "ok",
			Data: uploadCallbackData{
				ImageUrl: url,
			},
		})
		return
	}
	err = v.VideoService.UploadVideoCallback(authorId, url, r.CoverUrl, r.Describe, videoType)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &ResponseWithoutData{
		Code: SuccessCode,
		Msg:  "ok",
	})

}

func callbackReqValidate(r uploadCallbackReq) (authorId uint, videoType uint, err error) {
	if err = utils.Validate.Struct(r); err != nil {
		return
	}
	if r.IsImage == callbackIsImage {
		return
	}
	// 上传视频，校验参数
	aid, err1 := strconv.ParseUint(r.AuthorId, 10, 64)
	tid, err2 := strconv.ParseUint(r.VideoType, 10, 64)
	if err1 != nil || err2 != nil || r.CoverUrl == "null" {
		err = errors.New("参数传递错误")
		return
	}
	authorId, videoType = uint(aid), uint(tid)
	return
}
