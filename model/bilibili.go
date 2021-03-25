package model

import "sync"

type VideoAid struct {
	Aid       int64
	cidMap    map[int64]*VideoCid
	totalPage int64
	Title     string
	Quality   int64
	cidLock   sync.RWMutex
	pageLock  sync.RWMutex
}

func (videoAid *VideoAid) AddCid(videoCid *VideoCid) {
	videoAid.cidLock.Lock()
	defer videoAid.cidLock.Unlock()
	videoAid.cidMap[videoCid.Cid] = videoCid
}

func (videoAid *VideoAid) GetCid(cid int64) *VideoCid {
	videoAid.cidLock.RLock()
	defer videoAid.cidLock.RUnlock()
	return videoAid.cidMap[cid]
}
func (videoAid *VideoAid) SetPage(num int64) {
	videoAid.pageLock.Lock()
	defer videoAid.pageLock.Unlock()
	videoAid.totalPage = num
}

func (videoAid *VideoAid) GetPage() int64 {
	videoAid.pageLock.RLock()
	defer videoAid.pageLock.RUnlock()
	return videoAid.totalPage
}

func NewVideoAidInfo(aid int64, title string) *VideoAid {
	return &VideoAid{Aid: aid, Title: title, cidMap: make(map[int64]*VideoCid)}
}

type VideoCid struct {
	Cid      int64
	ParAid   *VideoAid
	Page     int64
	Part     string
	AllOrder int64
}

type Video struct {
	Order  int64
	ParCid *VideoCid
}

func NewVideoCidInfo(cid int64, parAid *VideoAid, page int64, part string) *VideoCid {
	return &VideoCid{Cid: cid, ParAid: parAid, Page: page, Part: part}
}
