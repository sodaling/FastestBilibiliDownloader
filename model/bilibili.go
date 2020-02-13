package model

import "sync"

type VideoAidInfo struct {
	Aid       int64
	cidMap    map[int64]*VideoCidInfo
	TotalPage int64
	Title     string
	Quality   int64
	sync.RWMutex
}

func (videoAid *VideoAidInfo) AddCid(videoCid *VideoCidInfo) {
	videoAid.Lock()
	defer videoAid.Unlock()
	videoAid.cidMap[videoCid.Cid] = videoCid
}

func (videoAid *VideoAidInfo) GetCid(cid int64) *VideoCidInfo {
	videoAid.RLock()
	defer videoAid.RUnlock()
	return videoAid.cidMap[cid]
}

func NewVideoAidInfo(aid int64, title string) *VideoAidInfo {
	return &VideoAidInfo{Aid: aid, Title: title, cidMap: make(map[int64]*VideoCidInfo)}
}

type VideoCidInfo struct {
	Cid    int64
	ParAid *VideoAidInfo
	Page   int64
}

func NewVideoCidInfo(cid int64, parAid *VideoAidInfo, page int64) *VideoCidInfo {
	return &VideoCidInfo{Cid: cid, ParAid: parAid, Page: page}
}
