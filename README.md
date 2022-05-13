
# jenkinscli

**基于gojenkins cobra viper go-mail 开发的Jenkins远程构建客户端**

**使用指南** 

## 1.配置jenkins登录信息

### 1.1 概览步骤

1.登录你的Jenkins

2.点击右上角你的用户名

3.点击设置

4.点击添加一个tocken，你必须为这个tocken起一个名字

### 1.2 写配置文件

1.2.1 默认配置文件路径在~/.config/jenkinscli

```shell
mkdir -p ~/.config/jenkinscli/
```

```shell
[root@localhost jenkinscli]# vi ~/.config/config.json 
{
    "Server": "https://jenkins.mydomain.com",
    "JenkinsUser": "admin",
    "Token": "113a8xxxxxxxxxxxxxxxxx756fcf4e8b8"
}

```

### 1.3 开始使用远程构建工具

#### 1.3.1 查看使用帮助

```shell
[root@localhost jenkinscli]# ./jenkinscli --help
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
  mail        Ability to send custom content emails

Flags:
      --config string   Path to config file
  -h, --help            help for jenkinscli
  -v, --version         version for jenkinscli

Use "jenkinscli [command] --help" for more information about a command.
```

#### 1.3.2 远程构建Jenkins并下载工件到指定路径

（以java程序为例：需要jenkins安装插件Pipeline Maven Integration Plugin来提供下载jar包）

```shell
#-p 指定下载路径
./jenkinscli launch -p /tmp/jar 你的流水线名称
```

## 2.配置邮箱

### 2.1 编写配置文件

同样在~/.config/jenkinscli/config 文件里

```shell
[root@localhost jenkinscli]# cat config.json 
{
    "Server": "https://jenkins.mydomain.com",
    "JenkinsUser": "admin",
    "Token": "113a8cxxxxxxxxxxxxxxxx756fcf4e8b8",
    "MailSmpt": "smtp.qq.com",
    "MailPort": 25,
    "MailUser": "xxxxxxxxx@qq.com",
    "MailToken": "lrxxxxxxxxxfa",
    "MailFrom": "xxxxxxxx@qq.com",
    "MailTo": ["xxxxxxxxx@qq.com"],
    "MailSub": "pc-system Test!!!"
}
```

> MsilSmpt , MailPort, MailUser, MailToken分别指定: 服务器地址，端口，邮箱用户名，邮箱token	
>
> 邮箱token 需要你登录到邮箱，在设置里可以拿到自己的token（这里以qq邮箱为例）
>
> MailFrom 设置发件人邮箱用户
>
> MailTo 设置收件人邮箱，可以设置多个，比如：["xxxxxxxxx@qq.com"，"yyyyyyyyy@qq.com"]
>
> MailSub 设置邮箱默认标题，在构建时你没有指定标题就会使用这个默认设置，如果在构建的命令参数-s(--subject) ，就会使用命令行的标题，优先级命令行最高

### 2.2 示例

```shell
 #-m 构建后发送邮件 -s 设置邮件标题
 ./jenkinscli launch -m -s "test" pc-system 
```



```shell
 #-m 构建后发送邮件 -s 设置邮件标题 -a 为邮件添加附件
 ./jenkinscli launch -m -s "test" -a /root/1.jpg pc-system 
```

### 2.3 配置抄送，密送

```json
{
    "MailCc": ["xxxxxxx.qq.com"],
    "MailBcc": ["xxxxxxxx@qq.com"],
    "MailAttach": "/root.1.jpg"
}
```

> MailCc 配置抄送邮箱，可以配置多个
>
> MailBcc 配置密送邮箱，可以配置多个
>
> MailAttach 配置默认附件，同样命令行优先级最高

效果图：

![image-20220513140116132](C:\Users\EDZ\AppData\Roaming\Typora\typora-user-images\image-20220513140116132.png)


