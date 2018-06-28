# CGO编程初识

## 1. why cgo
> GO语言有自己的擅长的领域 [web后端,分布式] ,但许多传统领域仍是C的主场
> 
> 通过CGO可以继承C语言近半个世纪的软件积累
> 
> CGO是GO语言直接和其他语言通讯的桥梁,C=>GO ,GO=>C
> 
> 可以直接用纯粹的GO语言解决的问题不用CGO,能让其他组提供RPC调用的话,不使用CGO [本来有一个问题,现在有一堆问题]

## 2. 快速入门


## 2. 类型转换,复杂类型[string]不可以
    
## 3. C=>GO GO=>C

## 4. 内存管理

## 5. 线程模型

## 6. revover

## 7. 库引用 与 部署

## 8. 项目介绍
### 8.1 连麦视频录制
### 8.2 截图服务
### 8.2.1 pipeline
### 8.2.2 channel监控

## 9. References

[1]  [cgo](https://golang.org/cmd/cgo/)

[2]  [深入CGO编程](https://github.com/chai2010/gopherchina2018-cgo-talk)

[3]  [cgo is not go](https://dave.cheney.net/2016/01/18/cgo-is-not-go)