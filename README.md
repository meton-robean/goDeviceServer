# goDeviceServer
Golang-based Device-Server 
基于golang的轻量级物联网网关集群管理后台微服务 

项目参与的一个基于golang的物联网后台项目. 该后台用于对接响其他 前继的java服务/webserver 发来的设备控制请求、以及管理控制lora物联网网关设备集群.基于这两个明确的需求，该系统也清晰的分为WechatAPI和DeviceServer两部分，以实现维护上的功能解耦.  
 
![系统架构](/design/DeviceServer.png)  

本微服务在项目中的位置:  
```
java后台 <-------------> **goDeviceServer** <-----> 物联网网关集群设备 <----> 大量物联网终端设备
                 |   |
网站webserver ----    |
                     |
手机微信端/测试APP端 ---

```

设计细节见 design/ 目录下文档.
