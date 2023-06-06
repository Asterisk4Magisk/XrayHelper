[English](README.md) | 简体中文

# XrayHelper  
一个安卓专属的Xray助手，使用golang实现[Xray4Magisk](https://github.com/Asterisk4Magisk/Xray4Magisk)的部分脚本，提供arm64和amd64二进制文件.

## 配置文件  
XrayHelper使用yml格式的配置文件，默认使用`/data/adb/xray/xrayhelper.yml`，当然你可以使用`-c`选项自定义配置文件路径
- xrayHelper  
    - `coreType`默认值`xray`，指定所使用的核心类型，可选`xray`、`sing-box`
    - `corePath`必填，指定核心路径
    - `coreConfig`必填，指定核心配置文件，可指向文件或目录，影响核心的启动命令（`-c`或`-confdir`）
    - `dataDir`必填，指定XrayHelper的数据目录，用于存储GEO数据文件、自定义节点和订阅节点信息等
    - `runDir`必填，用于存储运行时所产生的文件，例如核心的pid值，核心日志，inotify的监控日志等
    - `proxyTag`默认值`proxy`，使用XrayHelper进行节点切换时，将进行替换的出站代理Tag
    - `subList`可选，数组，与 [v2rayNg](https://github.com/2dust/v2rayNG) 兼容的节点订阅链接
- proxy  
    - `method`默认值`tproxy`，代理模式，可选`tproxy`、`tun`，使用tun模式时，需要提前下载 tun2socks 二进制文件（`xrayhelper update tun2socks`）
    - `tproxyPort`默认值`65535`，透明代理端口，该值需要与核心的`Dokodemo-Door`入站代理端口相对应，`tproxy`模式需要
    - `socksPort`默认值`65534`，socks5代理端口，该值需要与核心的`socks`入站代理端口相对应，`tun`模式需要
    - `enableIPv6`默认值`false`，是否启用ipv6代理，需要代理节点支持
    - `mode`默认值`blacklist`，代理应用名单模式，可选`whitelist`、`blacklist`
    - `pkgList`，可选，数组，代理应用名单，apk包名
    - `apList`，可选，数组，需代理的ap接口名，例如`wlan+`可代理wlan热点，`rndis+`可代理usb网络共享
    - `ignoreList`，可选，数组，需要忽略的接口名，例如可以实现连上wifi不走代理
    - `intraList`，可选，数组，CIDR，默认情况下，内网地址不会被标记，若需要将部分内网地址标记，可配置此项
## 命令
- service
    - `start`启动核心服务
    - `stop`停止核心服务
    - `restart`重启核心服务
    - `status`检查核心服务状态
- proxy
    - `enable`启用系统代理规则
    - `disable`停用系统代理规则
    - `refresh`刷新系统代理规则
- update
    - `core`更新核心，需要指定coreType
    - `geodata`从 [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) 更新geo数据库
    - `subscribe`更新订阅节点
    - `tun2socks`更新 tun2socks
- switch（仅支持xray核心）
    - 不带任何参数时，从订阅`${xrayHelper.dataDir}/sub.txt`获取节点信息并选择
    - `custom`从`${xrayHelper.dataDir}/custom.txt`获取节点信息并选择，因此，可将自定义节点的分享链接放置于此方便选择

## 鸣谢
- [@Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- [@Asterisk4Magisk/Xray4Magisk](https://github.com/Asterisk4Magisk/Xray4Magisk)
- [@2dust/v2rayNG](https://github.com/2dust/v2rayNG)
- [@heiher/hev-socks5-tunnel](https://github.com/heiher/hev-socks5-tunnel)
