# 简介
这是一个测试ollama的项目

# install
1. 安装[ollama](https://github.com/ollama/ollama)
2. 下载 llama3
3. 安装redis
4. 安装mysql, 创建名为 ollama的数据库 create database ollama utfbmb4
5. 把script 下的sql脚本写入数据库中
6. 把修改pkg/config/cfg.toml配置文件，并把它移动把 /etc/ollama-hertz目录下
7. 启动项目访问 http:127.0.0.1:你配置的端口/index.html
