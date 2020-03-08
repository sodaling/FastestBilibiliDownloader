# FastestBibiliDownloader

#### 项目地址：**[ FastestBilibiliDownloader](https://github.com/sodaling/FastestBilibiliDownloader)**

**东半球第二快的Bilibili.com（B站）视频下载器！**

如果你想下载b站某个up主的所有视频，而且要飞快的那种，那么你可以试试这个项目-.-

目前提供两个视频下载方案:

1. 通过视频的aid,下载单个视频.
2. 通过up主的upid(b站叫mid),下载这个up主所投稿的所有视频.



> 特性:
>
> Github上下载b站视频代码已经有很多了.那么本下载器的特点是啥呢?
>
> 因为这是用Golang写的,当然了,也就利用了Golang的特性:goroutine.
>
> 简单来说,特点就是:
>
> **快!贼快!下载的视频越多越快！**
>
> * 当单个aid视频分了若干个part时候,或者当你选了下载up主下所有视频时候.多个视频将会同时并行下载,跑满你的网速绝对不是问题.
> * 下载与合并视频并行处理.如果视频分了多个part,下载完成的同时就会立即合并.该视频合并处理和其他与其他下载和合并同时进行且互不影响.

## 运行

下载的视频会存放在运行路径下的**download**文件夹下，每个视频（aid）一个文件夹。

```shell
Please enter your id type(`aid` or `upid`) #aid是视频id，upid是up主id，可以在up主主页地址栏获得。
Please enter your id #输入id
```



1. 如果你有Golang环境的话，程序入口在**cmd/start-concurrent-engine.go**。自己build或者run就好。但是注意，合并视频需要ffmeg的支持。不然只会下载并不会自动合并。

2. 如果你没有Golang编译环境，或者没有ffmeg环境。那么推荐用docker方式运行。已经写好了dockefile和makefile。你只需要：

   ```shell
   $ cd FastestBilibiliDownloader
   $ make build #下载镜像
   $ make run #运行镜像
   ```

   

3. 后续有空会打包bin文件到release的。

## 感谢

1. engine部分的框架参考**ccmouse**的思路，后面自己调整了整体架构部分，非常感谢。
2. [bilibili-downloader](https://github.com/stevenjoezhang/bilibili-downloader)：b站请求视频的API等等都是从这位的代码获得，本身的py代码注释也非常清晰，非常感谢。
3. 我亲爱的娜娜姐～
4. @sshwy帮忙抓虫纠错

>欢迎各位提pr或者fork或者什么都行，能帮助到你的话欢迎star！疫情无聊在家磨时间的产物，粗糙了一点，欢迎各位完善～

