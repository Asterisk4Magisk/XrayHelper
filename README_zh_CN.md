[English](README.md) | 简体中文

# XrayHelper
一个安卓专属的通用代理助手，使用 Golang 实现 [Xray4Magisk](https://github.com/Asterisk4Magisk/Xray4Magisk) 的部分脚本，提供 arm64 和 amd64 二进制文件

## 配置
XrayHelper 使用 yml 格式的配置文件，默认使用`/data/adb/xray/xrayhelper.yml`，当然你可以使用`-c`选项自定义配置文件路径  
[配置示例](config.yml)
- xrayHelper
    - `coreType`默认值`xray`，指定所使用的核心类型，可选`xray`、`v2ray`、`sing-box`、`mihomo(clash.meta)`
    - `corePath`必填，指定核心路径
    - `coreConfig`必填，指定核心配置文件，可指向文件或目录，影响核心的启动命令
    - `dataDir`必填，指定 XrayHelper 的数据目录，用于存储 GEO 数据文件、自定义节点和订阅节点信息等
    - `runDir`必填，用于存储运行时所产生的文件，例如核心的 pid 值，核心日志等
    - `proxyTag`默认值`proxy`，使用 XrayHelper 进行节点切换时，将进行替换的出站代理 Tag
    - `subList`可选，数组，节点订阅链接（SIP002/SSR/v2rayNg/Hysteria），也支持 clash 订阅链接(需要在订阅链接前添加`clash+`前缀)
- proxy
    - `method`默认值`tproxy`，代理模式，可选`tproxy`、`tun`、`tun2socks`，使用 tun 模式时，请确保你的核心支持 tun 并正确配置它；使用 tun2socks 模式时，需要提前下载 tun2socks 二进制文件（可使用命令`xrayhelper update tun2socks`）
    - `tproxyPort`默认值`65535`，透明代理端口，该值需要与核心的 tproxy 入站代理端口相对应，`tproxy`模式需要
    - `socksPort`默认值`65534`，socks5 代理端口，该值需要与核心的 socks5 入站代理端口相对应，`tun2socks`模式需要
    - `tunDevice`默认值`xtun`，核心或 tun2socks 所创建的 tun 设备名
    - `enableIPv6`默认值`false`，是否启用 ipv6 代理，需要代理节点支持
    - `autoDNSStrategy`默认值`true`，是否自动配置核心的 DNS 策略（当未启用 IPv6 代理时，若禁用此特性，请确保你无法从核心的 DNS 解析到任何 AAAA 记录，否则可能导致域名代理策略失效问题）
    - `mode`默认值`blacklist`，代理应用名单模式，可选`whitelist`、`blacklist`，使用白名单模式时，下方应用名单内的应用流量会被标记，其他流量不会被标记（即绕过），反之，黑名单模式则不标记应用名单内的应用流量
    - `pkgList`，可选，数组，代理应用名单，格式为`apk包名:用户`，未指定用户时，默认0，即机主；需要注意当该列表为空时，无论代理名单是什么模式，都会标记所有应用流量
    - `apList`，可选，数组，需代理的 ap 接口名，例如`wlan+`可代理 wlan 热点，`rndis+`可代理 usb 网络共享
    - `ignoreList`，可选，数组，需要忽略的接口名，例如`wlan+`可以实现连上 wifi 不走代理
    - `intraList`，可选，数组，CIDR，默认情况下，内网地址不会被标记，若需要将部分内网地址标记，可配置此项
- clash
  - `dnsPort`默认值`65533`，mihomo(clash.meta) 监听的 dns 端口
  - `template`可选，mihomo(clash.meta) 配置模板，指定配置模板后，该模板会**覆盖（或注入）** mihomo(clash.meta) 配置文件对应内容

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
    - `core`更新核心，需要指定 **xrayHelper.coreType**
    - `geodata`从 [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) 更新 GEO 数据文件
    - `subscribe`更新订阅节点（或 clash 订阅）到`${xrayHelper.dataDir}/sub.txt`（或`${xrayHelper.dataDir}/clashSub#{index}.yaml`），需要指定 **xrayHelper.subList**
    - `tun2socks`从 [hev-socks5-tunnel](https://github.com/heiher/hev-socks5-tunnel) 更新 tun2socks
    - `yacd-meta`更新 [Yacd-meta](https://github.com/MetaCubeX/Yacd-meta) 到`${xrayHelper.dataDir}/Yacd-meta-gh-pages`
### xray、v2ray、sing-box
- switch
    - 不带任何参数时，从订阅`${xrayHelper.dataDir}/sub.txt`获取节点信息并选择
    - `custom`从`${xrayHelper.dataDir}/custom.txt`获取节点信息并选择，因此，可将自定义节点的分享链接放置于此方便选择
### mihomo(clash.meta)
- switch
  - 不带任何参数时，使用`${xrayHelper.dataDir}/clashSub#{index}.yaml`作为配置文件
  - `example.yaml`使用`${xrayHelper.coreConfig}/example.yaml`作为配置文件

**注意：${clash.template} 总是会覆盖（或注入）你所使用的配置文件**

## 许可
[Mozilla Public License Version 2.0 (MPL)](https://raw.githubusercontent.com/Asterisk4Magisk/XrayHelper/master/LICENSE)

## 鸣谢
- [@Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- [@2dust/v2rayNG](https://github.com/2dust/v2rayNG)
- [@heiher/hev-socks5-tunnel](https://github.com/heiher/hev-socks5-tunnel)
- ~~[@haishanh/yacd](https://github.com/haishanh/yacd)~~
- [@MetaCubeX/Yacd-meta](https://github.com/MetaCubeX/Yacd-meta)
