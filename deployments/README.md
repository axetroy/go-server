### 这里把 Docker 分为 3 部分

应该按照顺序依次启动，因为他们有依赖管理

1. database

数据库部分

2. message_queue

消息队列

3. api

接口

4. caddy
   HTTP 服务器，类似于 nginx。用于反向代理，压缩和全自动 HTTPS
