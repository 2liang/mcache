# VCutter

## 关于vender

1. 安装vender

`go get -u -v github.com/kardianos/govendor`

2. 初始化,生成vendor文件夹

`govendor init`

3. 执行命令将当前应用必须的文件包含进来

`govendor add +external`

4. 更新和删除

`update和remove`