
# Utopia-back

## 亮点介绍

+ [七牛云SDK使用](./doc/七牛云SDK使用.md)
  + 使用七牛云 kodo，采用 uploadToken + 回调 确保密钥安全性
  + 针对未上传封面的视频，异步截取视频首帧，回调替换视频封面
+ [热门视频设计方案](./doc/热门视频设计方案.md)
  + zset存储，异步刷取DB，获取一段时间内点赞量增长最高的视频，保持动态更新
  + version控制缓存版本，避免用户因热门视频突然更新而影响体验
+ [判断用户是否点赞缓存方案](./doc/判断用户是否点赞缓存方案.md)
  + 通过hash类型进行缓存，进行冷热数据分离，极大减少存储成本的同时也不会给DB带来过大负担
  + 根据用户维度存储减少查询成本，TTL合理续期
+ [接口及IP限流](./doc/接口及IP限流.md)
  + 采用令牌桶对可针对不同接口配置不同的限流策略
  + 同时支持各接口对用户IP限流，防止用户恶意攻击

## 项目介绍

项目介绍详见：[项目介绍](./doc/项目介绍.md)

## 使用说明

确保go版本大于等于1.21，通过`go version`命令检查

```shell
# 克隆项目
git clone git@github.com:VideoUtopia/utopia-back.git
cd utopia-back

# 修改配置文件，各字段含义见yaml文件字段备注
cp config/comfig_example.yaml config/congfig.yaml 
vim config/congfig.yaml

# linux 编译运行
go build -o server main.go 
./server

# windows
go build -o server.exe main.go
./server.exe
```