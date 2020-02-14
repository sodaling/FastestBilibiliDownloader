package model

import "sync"

type VideoAidInfo struct {
	Aid       int64
	cidMap    map[int64]*VideoCidInfo
	totalPage int64
	Title     string
	Quality   int64
	cidLock   sync.RWMutex
	pageLock  sync.RWMutex
}

func (videoAid *VideoAidInfo) AddCid(videoCid *VideoCidInfo) {
	videoAid.cidLock.Lock()
	defer videoAid.cidLock.Unlock()
	videoAid.cidMap[videoCid.Cid] = videoCid
}

func (videoAid *VideoAidInfo) GetCid(cid int64) *VideoCidInfo {
	videoAid.cidLock.RLock()
	defer videoAid.cidLock.RUnlock()
	return videoAid.cidMap[cid]
}
func (videoAid *VideoAidInfo) SetPage(num int64) {
	videoAid.pageLock.Lock()
	defer videoAid.pageLock.Unlock()
	videoAid.totalPage = num
}

func (videoAid *VideoAidInfo) GetPage() int64 {
	videoAid.pageLock.RLock()
	defer videoAid.pageLock.RUnlock()
	return videoAid.totalPage
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
