package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

type redBook struct{ requestClient }

func (r redBook) parseShareUrl(shareUrl string) (*VideoParseInfo, error) {
	client := r.newClient()
	videoRes, err := client.R().
		SetHeader(HttpHeaderUserAgent, DefaultUserAgent).
		Get(shareUrl)
	if err != nil {
		return nil, err
	}

	if videoRes.IsError() {
		return nil, ErrNoMedia
	}

	re := regexp.MustCompile(`(?s)window.__INITIAL_STATE__\s*=\s*(.*?)</script>`)
	findRes := re.FindSubmatch(videoRes.Body())
	if len(findRes) < 2 {
		return nil, errors.New("parse video json info from html fail")
	}

	jsonBytes := findRes[1]
	if gjson.GetBytes(jsonBytes, "noteData.hasError").Bool() {
		return nil, ErrNoMedia
	}

	nodeId := gjson.GetBytes(jsonBytes, "note.currentNoteId").String()
	data := gjson.GetBytes(jsonBytes, "noteData.data.noteData")
	if !data.Exists() {
		data = gjson.GetBytes(jsonBytes, fmt.Sprintf("note.noteDetailMap.%s.note", nodeId))
	}
	if !data.Exists() {
		return nil, ErrNoMedia
	}

	videoUrl := data.Get("video.media.stream.h264.0.masterUrl").String()

	// 获取图集图片地址
	imagesObjArr := data.Get("imageList").Array()
	images := make([]ImgInfo, 0, len(imagesObjArr))
	if len(videoUrl) <= 0 {
		for _, imageItem := range imagesObjArr {
			imageUrl := redBookImageURL(imageItem)
			if len(imageUrl) <= 0 {
				continue
			}
			imgId := strings.Split(imageUrl[strings.LastIndex(imageUrl, "/")+1:], "!")[0]
			// 如果链接中带有 spectrum/ , 替换域名时需要带上
			spectrumStr := ""
			if strings.Contains(imageUrl, "spectrum") {
				spectrumStr = "spectrum/"
			}
			newUrl := fmt.Sprintf("https://ci.xiaohongshu.com/notes_pre_post/%s%s?imageView2/format/jpg", spectrumStr, imgId)
			imgInfo := ImgInfo{
				Url: newUrl,
			}
			// 如果原图片网址中没有 notes_pre_post 关键字，不支持替换域名，使用原域名
			if !strings.Contains(imageUrl, "notes_pre_post") {
				imgInfo.Url = imageUrl
			}
			if imageItem.Get("livePhoto").Bool() {
				for _, livePhotoItem := range imageItem.Get("stream.h264").Array() {
					if livePhotoUrl := livePhotoItem.Get("masterUrl").String(); len(livePhotoUrl) > 0 {
						imgInfo.LivePhotoUrl = livePhotoUrl
					}
				}
			}
			images = append(images, imgInfo)
		}
	}

	parseInfo := &VideoParseInfo{
		Title:    data.Get("title").String(),
		VideoUrl: data.Get("video.media.stream.h264.0.masterUrl").String(),
		CoverUrl: redBookImageURL(data.Get("imageList.0")),
		Images:   images,
	}
	parseInfo.Author.Uid = data.Get("user.userId").String()
	parseInfo.Author.Name = data.Get("user.nickname").String()
	parseInfo.Author.Avatar = data.Get("user.avatar").String()

	if parseInfo.VideoUrl == "" && len(images) == 0 {
		return nil, ErrNoMedia
	}
	return parseInfo, nil
}

func redBookImageURL(image gjson.Result) string {
	for _, field := range []string{"urlDefault", "url"} {
		if value := image.Get(field).String(); value != "" {
			return value
		}
	}
	var fallback string
	for _, info := range image.Get("infoList").Array() {
		value := info.Get("url").String()
		if value == "" {
			continue
		}
		if info.Get("imageScene").String() == "WB_DFT" {
			return value
		}
		if fallback == "" {
			fallback = value
		}
	}
	return fallback
}
