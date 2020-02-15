# FastestBibiliDownloader
提供两个视频下载方案:
1. 视频aid下载单个视频
2. up主的id下载他所有视频

## 特点:
彻底利用了golang的goroutine,贼快:
- 下载视频的数量越多越快.比如当单个视频分了几个part时候,或者下载up主的视频有一大堆时候
- 当下载视频分了若干个part,下载完要合并时候.
