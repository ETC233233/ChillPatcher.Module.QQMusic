package main

import (
	"C"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 歌曲搜索 SearchSong(keyword string, page int)
type SearchResult struct {
	Code int `json:"code"`
	Data struct {
		Song struct {
			List []struct{
				SongName string   `json:"songname"`
				SongMid  string   `json:"songmid"`
				AlbumMid string   `json:"albummid"`
				Singer   []struct{ Name string `json:"name"` } `json:"singer"`
			} `json:"list"`
		} `json:"song"`
	} `json:"data"`
}

//export QQMusicSearchSong
func QQMusicSearchSong(keyword *C.char, page C.int) *C.char {
	k := C.GoString(keyword)
	searchURL := fmt.Sprintf("https://c.y.qq.com/soso/fcgi-bin/client_search_cp?format=json&w=%s&p=%d", url.QueryEscape(k), int(page))
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Add("Referer", "https://y.qq.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return C.CString("{\"code\":-1}") }
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return C.CString(string(body))
}

//export QQMusicGetLyric
func QQMusicGetLyric(songmid *C.char) *C.char {
	mid := C.GoString(songmid)
	lyricURL := fmt.Sprintf("https://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg?songmid=%s&format=json&nobase64=1", mid)
	req, _ := http.NewRequest("GET", lyricURL, nil)
	req.Header.Add("Referer", "https://y.qq.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return C.CString("{\"code\":-1}") }
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return C.CString(string(body))
}

//export QQMusicGetSongURL
func QQMusicGetSongURL(songmid *C.char, quality *C.char) *C.char {
	mid := C.GoString(songmid)
	fmtq := C.GoString(quality)
	// [TODO] 目前此处未实现签名和key抓取，只返回样例URL，建议后续用更全API完善
	urlLike := fmt.Sprintf("https://api.qqmusicapi.com/song/url?id=%s&type=%s", mid, fmtq)
	return C.CString(fmt.Sprintf(`{"url":"%s"}`, urlLike))
}
