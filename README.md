# XrayHelper
XrayHelper for Android, some scripts in [Xray4Magisk](https://github.com/Asterisk4Magisk/Xray4Magisk) rewritten with golang, provide arm64 and amd64 binary.

## Control Core Service
`xrayhelper service start`, start core service  
`xrayhelper service stop`, stop core service  
`xrayhelper service restart`, restart core service

## Control System Proyx
Support application package proxy list run with blacklist and whitelist, bypass specific network interface, and proxy ap interface, should configure **proxy**
```yaml
proxy:
    method: tproxy
    tproxyPort: 65535
    enableIPv6: false
    mode: whitelist
    pkgList:
        - com.kiwibrowser.browser
        - com.termux
    apList:
        - wlan2
        - rndis0
    ignoreList:
        - ignore
```
`xrayhelper proxy enable`, enable system proxy  
`xrayhelper proxy disable`, disable system proxy    
`xrayhelper proxy refresh`, refresh system proxy  

## Update Components
- update core  
  `xrayhelper update core`, should configure **xrayHelper.coreType** first, support xray, v2fly, sagernet  
- update geodata  
  `xrayhelper update geodata`, update geodata from [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)  
- update subscribe nodes  
  `xrayhelper update subscribe`, update your subscribe, should configure **xrayHelper.subList** first, compatible with [v2rayNg](https://github.com/2dust/v2rayNG)'s subscription link standard

## Switch Proxy Node(Currently only compatible with xray-core)  
`xrayhelper swtich`, should configure **xrayHelper.proxyTag** and update subscribe first, **warning: it will replace your outbounds configuration which has the same proxy tag**

## Credits
- [@Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- [@Asterisk4Magisk/Xray4Magisk](https://github.com/Asterisk4Magisk/Xray4Magisk)
- [@2dust/v2rayNG](https://github.com/2dust/v2rayNG)
