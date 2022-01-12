# jenkinscli

基于gojenkins cobra viper 开发的Jenkins客户端

#### 使用指南 

#### 第一步

1. 登录你的Jenkins
2. 点击右上角你的用户名
3. 点击设置
4. 点击添加一个tocken，你必须为这个token起一个名字

#### 

#### 第二步



[root@node1 ~]#mkdir -p ~/.config/jenkinscli/

[root@node1 ~]#cat .config/jenkinscli/config.json 

{

 "Server": "https://jenkins.mydomain.com",

 "JenkinsUser": "admin",

 "Token": "1184d0c4d3763a1f119a602f4d4d3a70b4"

}

#### 第三步

git clone https://github.com/xuanhi/jenkinscli.git

cd jenkinscli

go build

#### 第四步

`[root@localhost jenkinscli]# ./jenkinscli 

Client for jenkins, manage resources by the jenkis

Usage:
  jenkinscli [command]

Available Commands:

  completion  Generate the autocompletion script for the specified shell

  download    download related commands

  enable      Enable a resource in Jenkins

  get         Get a resource Jenkins

  help        Help about any command

  launch      Start a Jenkins resource and you can trigger an artifact download

Flags:

​      --config string   Path to config file

  -h, --help            help for jenkinscli

  -v, --version         version for jenkinscli

Use "jenkinscli [command] --help" for more information about a command
