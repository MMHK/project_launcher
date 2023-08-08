# MM Project Launcher

基于 `docker` 容器的 windows 10 / MacOS 开发环境


### How to use

1. 到 [Github 发布页](https://github.com/MMHK/project_launcher/releases)，下载最新版本zip文档，解压 exe 到项目 `public` 目录下
2. 执行 exe 文件，根据提示进行操作
3. 等候 exe 执行完毕后，重新刷新打开的网页，就会见到已经跑起来的网站了。

### 常见问题

1. 工具提示 `请先安装 Docker Desktop`, 请先安装 [docker desktop](https://docs.docker.com/docker-for-windows/install/)
2. 工具提示 `DockerDesktop 还未运行`, 表示没有启动 `docker Desktop` 引擎, 工具会自行找到已经安装的 `Docker Desktop` 并启动。
如果还未有问题，请重启引擎 
![image01](./doc/image01.png)
3. 如果关闭正在运行的 `docker` 项目, 打开 docker desktop, 找到 `Containers / Apps`， 找到对应的项目干掉。
![image02](./doc/image02.png)
   

### 自己编译

- 需要 Golang >= 1.16 (使用最新的 `embed` 特性)
- 需要安装 [rsrc](https://github.com/akavel/rsrc) 用于引入 windows10 的管理员权限
- 执行 根目录下的 编译脚本 `build.cmd` （windows10 版本）
- 执行 根目录下的 编译脚本 `build-darwin.sh` （MacOS 版本）

#### macOS 编译说明

- 请使用[platypus](https://sveinbjorn.org/platypus) 工具打包成 bundle app，启动脚本在 `bin/darwin/start.sh`
