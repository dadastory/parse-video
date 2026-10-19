package parser

// 视频渠道来源
const (
	SourceDouYin       = "douyin"       // 抖音
	SourceKuaiShou     = "kuaishou"     // 快手
	SourcePiPiXia      = "pipixia"      // 皮皮虾
	SourceHuoShan      = "huoshan"      // 火山
	SourceWeiBo        = "weibo"        // 微博
	SourceWeiShi       = "weishi"       // 微视
	SourceLvZhou       = "lvzhou"       // 绿洲
	SourceZuiYou       = "zuiyou"       // 最右
	SourceQuanMin      = "quanmin"      // 度小视(原 全民小视频)
	SourceXiGua        = "xigua"        // 西瓜
	SourceLiShiPin     = "lishipin"     // 梨视频
	SourcePiPiGaoXiao  = "pipigaoxiao"  // 皮皮搞笑
	SourceHuYa         = "huya"         // 虎牙
	SourceAcFun        = "acfun"        // A站
	SourceDouPai       = "doupai"       // 逗拍
	SourceMeiPai       = "meipai"       // 美拍
	SourceQuanMinKGe   = "quanminkge"   // 全民K歌
	SourceSixRoom      = "sixroom"      // 六间房
	SourceXinPianChang = "xinpianchang" // 新片场
	SourceHaoKan       = "haokan"       // 好看视频
	SourceRedBook      = "redbook"      // 小红书
	SourceBiliBili     = "bilibili"     // 哔哩哔哩
	SourceTwitter      = "twitter"      // X/Twitter
	SourceQQVideo      = "qqvideo"      // 腾讯视频
	SourceCCTV         = "cctv"         // 央视网
	SourceSohu         = "sohu"         // 搜狐视频
)

// http 相关
const (
	HttpHeaderUserAgent   = "User-Agent" //http header
	HttpHeaderReferer     = "Referer"
	HttpHeaderContentType = "Content-Type"
	HttpHeaderCookie      = "Cookie"

	// DefaultUserAgent 默认UserAgent
	DefaultUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 26_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Mobile/15E148 Safari/604.1"
)

// videoShareUrlParser 根据视频分享地址解析
type videoShareUrlParser interface {
	parseShareUrl(shareUrl string) (*VideoParseInfo, error)
}

// videoIdParser 根据视频ID解析
type videoIdParser interface {
	parseVideoID(videoId string) (*VideoParseInfo, error)
}

// VideoParseInfo 视频解析信息
type VideoParseInfo struct {
	Author struct {
		Uid    string `json:"uid"`    // 作者id
		Name   string `json:"name"`   // 作者名称
		Avatar string `json:"avatar"` // 作者头像
	} `json:"author"`
	Title    string    `json:"title"`     // 描述
	VideoUrl string    `json:"video_url"` // 视频播放地址
	MusicUrl string    `json:"music_url"` // 音乐播放地址
	CoverUrl string    `json:"cover_url"` // 视频封面地址
	Images   []ImgInfo `json:"images"`    // 图集图片地址列表
}

type ImgInfo struct {
	Url          string `json:"url"`            // 图片url
	LivePhotoUrl string `json:"live_photo_url"` // livePhoto 视频地址
}

// BatchParseItem 批量解析时, 单条解析格式
type BatchParseItem struct {
	ParseInfo *VideoParseInfo // 视频解析信息
	Error     error           // 错误, 如果单条解析失败时, 记录error信息
}

// 视频渠道信息
type videoSourceInfo struct {
	VideoShareUrlDomain []string            // 视频分享地址域名
	VideoShareUrlParser videoShareUrlParser // 视频分享地址解析方法
	VideoIdParser       videoIdParser       // 视频id解析方法, 有些渠道可能没有id解析方法
}

// 视频渠道映射信息
var videoSourceInfoMapping = sourceInfoMapping(requestClient{})

func sourceInfoMapping(client requestClient) map[string]videoSourceInfo {
	return map[string]videoSourceInfo{
		SourceDouYin: {
			VideoShareUrlDomain: []string{"v.douyin.com", "www.iesdouyin.com", "www.douyin.com"},
			VideoShareUrlParser: douYin{requestClient: client},
			VideoIdParser:       douYin{requestClient: client},
		},
		SourceKuaiShou: {
			VideoShareUrlDomain: []string{"v.kuaishou.com", "www.kuaishou.com"},
			VideoShareUrlParser: kuaiShou{requestClient: client},
		},
		SourceZuiYou: {
			VideoShareUrlDomain: []string{"share.xiaochuankeji.cn"},
			VideoShareUrlParser: zuiYou{requestClient: client},
		},
		SourceXiGua: {
			VideoShareUrlDomain: []string{"v.ixigua.com"},
			VideoShareUrlParser: xiGua{requestClient: client},
			VideoIdParser:       xiGua{requestClient: client},
		},
		SourcePiPiXia: {
			VideoShareUrlDomain: []string{"h5.pipix.com"},
			VideoShareUrlParser: piPiXia{requestClient: client},
			VideoIdParser:       piPiXia{requestClient: client},
		},
		SourceWeiShi: {
			VideoShareUrlDomain: []string{"isee.weishi.qq.com"},
			VideoShareUrlParser: weiShi{requestClient: client},
			VideoIdParser:       weiShi{requestClient: client},
		},
		SourceHuoShan: {
			VideoShareUrlDomain: []string{"share.huoshan.com"},
			VideoShareUrlParser: huoShan{requestClient: client},
			VideoIdParser:       huoShan{requestClient: client},
		},
		SourceLiShiPin: {
			VideoShareUrlDomain: []string{"www.pearvideo.com"},
			VideoShareUrlParser: liShiPin{requestClient: client},
			VideoIdParser:       liShiPin{requestClient: client},
		},
		SourcePiPiGaoXiao: {
			VideoShareUrlDomain: []string{"h5.pipigx.com"},
			VideoShareUrlParser: piPiGaoXiao{requestClient: client},
			VideoIdParser:       piPiGaoXiao{requestClient: client},
		},
		SourceQuanMin: {
			VideoShareUrlDomain: []string{"xspshare.baidu.com"},
			VideoShareUrlParser: quanMin{requestClient: client},
			VideoIdParser:       quanMin{requestClient: client},
		},
		SourceHuYa: {
			VideoShareUrlDomain: []string{"v.huya.com"},
			VideoShareUrlParser: huYa{requestClient: client},
			VideoIdParser:       huYa{requestClient: client},
		},
		SourceAcFun: {
			VideoShareUrlDomain: []string{"www.acfun.cn"},
			VideoShareUrlParser: acFun{requestClient: client},
			VideoIdParser:       acFun{requestClient: client},
		},
		SourceWeiBo: {
			VideoShareUrlDomain: []string{"weibo.com"},
			VideoShareUrlParser: weiBo{requestClient: client},
			VideoIdParser:       weiBo{requestClient: client},
		},
		SourceLvZhou: {
			VideoShareUrlDomain: []string{"weibo.cn"},
			VideoShareUrlParser: lvZhou{requestClient: client},
			VideoIdParser:       lvZhou{requestClient: client},
		},
		SourceMeiPai: {
			VideoShareUrlDomain: []string{"meipai.com"},
			VideoShareUrlParser: meiPai{requestClient: client},
			VideoIdParser:       meiPai{requestClient: client},
		},
		SourceDouPai: {
			VideoShareUrlDomain: []string{"doupai.cc"},
			VideoShareUrlParser: douPai{requestClient: client},
			VideoIdParser:       douPai{requestClient: client},
		},
		SourceQuanMinKGe: {
			VideoShareUrlDomain: []string{"kg.qq.com"},
			VideoShareUrlParser: quanMinKGe{requestClient: client},
			VideoIdParser:       quanMinKGe{requestClient: client},
		},
		SourceSixRoom: {
			VideoShareUrlDomain: []string{"6.cn"},
			VideoShareUrlParser: sixRoom{requestClient: client},
			VideoIdParser:       sixRoom{requestClient: client},
		},
		SourceXinPianChang: {
			VideoShareUrlDomain: []string{"xinpianchang.com"},
			VideoShareUrlParser: xinPianChang{requestClient: client},
		},
		SourceHaoKan: {
			VideoShareUrlDomain: []string{
				"haokan.baidu.com",
				"haokan.hao123.com",
			},
			VideoShareUrlParser: haoKan{requestClient: client},
			VideoIdParser:       haoKan{requestClient: client},
		},
		SourceRedBook: {
			VideoShareUrlDomain: []string{
				"www.xiaohongshu.com",
				"xhslink.com",
				"xhslink.cn",
			},
			VideoShareUrlParser: redBook{requestClient: client},
		},
		SourceBiliBili: {
			VideoShareUrlDomain: []string{
				"bilibili.com",
				"b23.tv",
			},
			VideoShareUrlParser: biliBili{requestClient: client},
		},
		SourceTwitter: {
			VideoShareUrlDomain: []string{
				"x.com",
				"twitter.com",
				"t.co",
			},
			VideoShareUrlParser: twitter{requestClient: client},
			VideoIdParser:       twitter{requestClient: client},
		},
		SourceQQVideo: {
			VideoShareUrlDomain: []string{
				"v.qq.com",
			},
			VideoShareUrlParser: qqVideo{requestClient: client},
			VideoIdParser:       qqVideo{requestClient: client},
		},
		SourceSohu: {
			VideoShareUrlDomain: []string{
				"tv.sohu.com",
				"my.tv.sohu.com",
			},
			VideoShareUrlParser: sohuVideo{requestClient: client},
			VideoIdParser:       sohuVideo{requestClient: client},
		},
		SourceCCTV: {
			VideoShareUrlDomain: []string{
				"tv.cctv.cn",
				"tv.cctv.com",
			},
			VideoShareUrlParser: cctvVideo{requestClient: client},
			VideoIdParser:       cctvVideo{requestClient: client},
		},
	}

}

// VideoSourceInfoMapping 导出的视频渠道映射信息，供外部使用
var VideoSourceInfoMapping = videoSourceInfoMapping
