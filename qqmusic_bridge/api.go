package qqmusic_bridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// QQMusicAPIResponse 定义标准返回格式
// QQ音乐通用API响应结构体，用于统一返回结果
// 业务数据一般存放在 Data 字段
// Code 表示响应码，Message 返回提示
// 非零 Code 需特殊处理
//
type QQMusicAPIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Playlist 歌单基本信息
// 用于 QQMusicGetPlaylist / QQMusicGetPlaylistDetail 返回
// 歌曲列表 SongList 需单独查询
//
type Playlist struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	CoverURL  string   `json:"cover_url"`
	Desc      string   `json:"desc"`
	SongCount int      `json:"song_count"`
	SongList  []Song   `json:"song_list,omitempty"`
	Creator   string   `json:"creator"`
	Tags      []string `json:"tags"`
}

// Song 歌曲信息
// 用于通用歌曲返回
//
type Song struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Album     Album    `json:"album"`
	Artists   []Singer `json:"artists"`
	Duration  int      `json:"duration"`
	CoverURL  string   `json:"cover_url"`
	URL       string   `json:"url,omitempty"`
	Copyright string   `json:"copyright,omitempty"`
}

// Album 专辑基本信息
//
type Album struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Cover string `json:"cover"`
}

// MV MV结构体
//
type MV struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Singer 歌手信息
//
type Singer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ========================= 主要对外接口实现 =========================

// QQMusicLogin_QRGetKey 获取QQ音乐扫码登录二维码key及图片
// 业务流程：向 QQ音乐API 发送请求，生成可用于扫码登录的 key，并生成二维码（base64/png）
// 参数说明：无（部分实现可能需cookie，仅匿名尝试即可）
// 返回：二维码 key、二维码图片（base64/png）、过期时间等错误
func QQMusicLogin_QRGetKey() (key string, qrImageBase64 string, expireAt int64, err error) {
	// TODO: 如果官方API稳定，可实现改为更高可靠请求。
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://ssl.hwui.top/api/qqmusic/qrkey") // 示例api，实际需替换生产可用
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, err
	}
	type qrResp struct {
		Code int    `json:"code"`
		Key  string `json:"key"`
		Img  string `json:"img"`
		Exp  int64  `json:"expire_at"`
	}
	var qr qrResp
	err = json.Unmarshal(body, &qr)
	if err != nil || qr.Code != 0 {
		return "", "", 0, errors.New("二维码获取失败")
	}
	return qr.Key, qr.Img, qr.Exp, nil
}

// QQMusicLogin_QRCheckStatus 轮询二维码扫码状态
// 业务流程：周期请求对应key状态，返回扫码结果，如uin、token、cookie、状态码
// 参数说明：key string，二维码对应值
// 返回：uin、cookie、token、状态码、错误提示等
func QQMusicLogin_QRCheckStatus(key string) (uin string, cookie string, token string, status int, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	postbody := map[string]string{"key": key}
	bs, _ := json.Marshal(postbody)
	resp, err := client.Post("https://ssl.hwui.top/api/qqmusic/qrstatus", "application/json", bytes.NewReader(bs))
	if err != nil {
		return "", "", "", -1, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", -1, err
	}
	type statusResp struct {
		Code   int    `json:"code"`
		Uin    string `json:"uin"`
		Token  string `json:"token"`
		Cookie string `json:"cookie"`
		Stat   int    `json:"status"` // 0未扫码/1已扫码/2授权成功/3二维码失效 etc
	}
	var sr statusResp
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return "", "", "", -1, err
	}
	if sr.Code != 0 {
		return "", "", "", sr.Stat, errors.New("登录状态异常")
	}
	return sr.Uin, sr.Cookie, sr.Token, sr.Stat, nil
}

// QQMusicSearchSong 搜索单曲（标准json返回，支持cookie与完整API请求）
// 业务流程：拼接参数，带cookie请求 QQ音乐搜索接口，解析返回歌曲列表结构
// 参数说明：keyword 搜索关键字, page 页码, pageSize 每页数量, cookie 可选用户cookie
// 返回：API标准结构体/错误
func QQMusicSearchSong(keyword string, page int, pageSize int, cookie string) (QQMusicAPIResponse, error) {
	url := fmt.Sprintf("https://c.y.qq.com/soso/fcgi-bin/client_search_cp?p=%d&n=%d&w=%s", page, pageSize, keyword)
	req, _ := http.NewRequest("GET", url, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return QQMusicAPIResponse{Code: 500, Message: err.Error()}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return QQMusicAPIResponse{Code: 500, Message: err.Error()}, err
	}
	// TODO: 按需解析更多字段
	return QQMusicAPIResponse{Code: 0, Message: "", Data: string(body)}, nil
}

// QQMusicGetLyric 获取歌词接口（支持cookie/VIP歌词）
// 业务流程：带cookie拉取歌词接口，无用户则拉公开API。若VIP歌词遇到特殊授权失败应给出提示。
// 参数说明：songmid 歌曲MID, cookie 用户cookie(可选)
// 返回：API标准结构体/错误
func QQMusicGetLyric(songmid string, cookie string) (QQMusicAPIResponse, error) {
	url := fmt.Sprintf("https://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg?songmid=%s&format=json", songmid)
	req, _ := http.NewRequest("GET", url, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return QQMusicAPIResponse{Code: 500, Message: err.Error()}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return QQMusicAPIResponse{Code: 500, Message: err.Error()}, err
	}
	return QQMusicAPIResponse{Code: 0, Message: "", Data: string(body)}, nil
}

// QQMusicGetSongURL 获取歌曲播放URL（带vkey高音质，失败Fallback开放接口）
// 业务流程：生成with vkey的高音质播放链接，若不可用则请求开放接口fallback。
// 参数说明：songmid 歌曲MID, cookie 用户cookie(可选)
// 返回：播放url、错误
func QQMusicGetSongURL(songmid string, cookie string) (url string, err error) {
	// 官方协议需要uin与vkey，简易方案尝试第三方服务。
	realurl := ""
	// 首选推荐vkey 1080p接口
	api1 := fmt.Sprintf("https://ssl.hwui.top/api/qqmusic/songurl?songmid=%s", songmid)
	client := http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", api1, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		type urlResp struct {
			Code int    `json:"code"`
			URL  string `json:"url"`
		}
		var ur urlResp
		if json.Unmarshal(body, &ur) == nil && ur.Code == 0 && ur.URL != "" {
			realurl = ur.URL
		}
	}
	if realurl == "" {
		// fallback开放接口（如网易云/其余镜像）
		realurl = fmt.Sprintf("https://music.jsososo.com/api/qq/track/url?id=%s&type=320", songmid)
	}
	return realurl, nil
}

// QQMusicGetPlaylist 获取歌单基础信息
// 业务流程：拉QQ音乐主流歌单接口并解析，主流歌单下歌曲需用GetPlaylistDetail二次请求
// 参数说明：playlistID 歌单ID, cookie 用户cookie(可选)
// 返回：Playlist结构体/错误
func QQMusicGetPlaylist(playlistID string, cookie string) (Playlist, error) {
	url := fmt.Sprintf("https://c.y.qq.com/qzone/fcg-bin/fcg_ucc_getcdinfo_byids_cp.fcg?type=1&disstid=%s&format=json", playlistID)
	req, _ := http.NewRequest("GET", url, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Playlist{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Playlist{}, err
	}
	// TODO: 解析更多歌单字段
	var respmap map[string]interface{}
	err = json.Unmarshal(body, &respmap)
	if err != nil {
		return Playlist{}, err
	}
	// 仅解析基础
	cdlist, ok := respmap["cdlist"].([]interface{})
	if !ok || len(cdlist) == 0 {
		return Playlist{}, errors.New("暂未查到歌单信息")
	}
	main := cdlist[0].(map[string]interface{})
	playlist := Playlist{
		ID:        fmt.Sprintf("%v", main["disstid"]),
		Title:     fmt.Sprintf("%v", main["dissname"]),
		CoverURL:  fmt.Sprintf("%v", main["logo"]),
		Desc:      fmt.Sprintf("%v", main["desc"]),
		SongCount: int(main["songnum"].(float64)),
		Creator:   fmt.Sprintf("%v", main["creator"]),
	}
	return playlist, nil
}

// QQMusicGetPlaylistDetail 获取歌单详情(含歌曲)
// 业务流程：拉取歌单基本+歌曲列表，支持cookie
// 参数说明：playlistID 歌单ID, cookie 用户cookie(可选)
// 返回：Playlist完整结构体/错误
func QQMusicGetPlaylistDetail(playlistID string, cookie string) (Playlist, error) {
	basic, err := QQMusicGetPlaylist(playlistID, cookie)
	if err != nil {
		return Playlist{}, err
	}
	url := fmt.Sprintf("https://c.y.qq.com/qzone/fcg-bin/fcg_ucc_getcdinfo_byids_cp.fcg?type=1&disstid=%s&format=json", playlistID)
	req, _ := http.NewRequest("GET", url, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Playlist{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Playlist{}, err
	}
	var respmap map[string]interface{}
	err = json.Unmarshal(body, &respmap)
	if err != nil {
		return Playlist{}, err
	}
	cdlist, ok := respmap["cdlist"].([]interface{})
	if !ok || len(cdlist) == 0 {
		return Playlist{}, errors.New("未查到歌单详情")
	}
	main := cdlist[0].(map[string]interface{})
	songs := []Song{}
	if songlist, ok := main["songlist"].([]interface{}); ok {
		for _, s := range songlist {
			ss := s.(map[string]interface{})
			// 解析歌曲基本字段
			id := fmt.Sprintf("%v", ss["songmid"])
			name := fmt.Sprintf("%v", ss["songname"])
			album := Album{
				ID:    fmt.Sprintf("%v", ss["albummid"]),
				Name:  fmt.Sprintf("%v", ss["albumname"]),
				Cover: fmt.Sprintf("https://y.gtimg.cn/music/photo_new/T002R300x300M000%v.jpg", ss["albummid"]),
			}
			artists := []Singer{}
			if ss["singer"] != nil {
				for _, sinfo := range ss["singer"].([]interface{}) {
					sm := sinfo.(map[string]interface{})
					artists = append(artists, Singer{
						ID:   fmt.Sprintf("%v", sm["mid"]),
						Name: fmt.Sprintf("%v", sm["name"]),
					})
				}
			}
			duration := 0
			if dur, ok := ss["interval"]; ok { duration = int(dur.(float64)) }
			cover := album.Cover
			songs = append(songs, Song{
				ID:       id,
				Name:     name,
				Album:    album,
				Artists:  artists,
				Duration: duration,
				CoverURL: cover,
			})
		}
	}
	basic.SongList = songs
	return basic, nil
}
