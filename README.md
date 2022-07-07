# jenkinscli

## 简介

jenkinscli 是一个简单的运维工具，大体上分为两个功能，远程触发Jenkins模块和部署模块。

远程触发Jenkins模块：你可以远程触发jenkins同时还能下载工件，能够发送邮件显示构建结果信息。

部署相关的模块主要集成了sftp相关功能，你可以同时向多台主机上传文件或向多台主机下载文件到本地。你也可以使用ssh子命令向多个主机远程发送指令，或者将本地的脚本批量发送到远程主机执行

**注意：本版本将使用yaml格式的配置文件**

**3.0.4更新：添加了远程执行脚本时为脚本输入位置参数。更新了对于ubuntu系统远程执行脚本时指定sudo密码(-s或--sudo)**

**3.0.3更新：添加了tmpl功能**

**更新说明：支持远程触发参数化构建类型，换句话说，可以通过远程触发构建enable launch 第二个参数起到任意个参数数量为参数构建选择值，格式必须为name:value 形式。name为参数名，在Jenkins中配置的参数名，value就是设置参数的值。经过测试，不仅可以用于jenkins自带的参数类型，也可以用于扩展参数插件（Extended Choice Parameter）中的类型，比如单选radio类型**

示例：

```shell
./jenkinscli launch  pc-system xhh:123456789 "xhhstring:remote string" "xhhtext:remote3 text" "xhhradio:三"
```

建议注入参数时使用引号，特别是有空格的参数名或值必须使用引号。

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

```yaml
[root@localhost jenkinscli]# vi ~/.config/config.yaml 
Server: http://xxx.xxx.xxx.xxx:8080/
JenkinsUser: admin
Token: 113a8cxxxxxxxxxxxxxxxe8b8
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
  download    download related commands(手动下载工件,指定3个参数分别是:流水线名,构建号,工件保存的路径)
  enable      Enable a resource in Jenkins(启动Jenkins流水线,相对于launch功能少很多且实现方法不一样)
  get         Get a resource Jenkins(用于获取远程Jenkins的资源信息)
  help        Help about any command
  launch      Start a Jenkins resource and you can trigger an artifact download(启动Jenkins用于java可以下载工件)
  mail        Ability to send custom content emails(发送邮件模块)
  sftp        Upload or download files or folders to a remote host(集成了sftp相关功能)
  ssh         remote exec bash(集成ssh相关功能)

Flags:
      --config string   Path to config file(指定配置文件路径)
  -h, --help            help for jenkinscli
  -I, --immunity        Used to prevent the jenkins server from exiting with an initialized error(免疫Jenkins初始化保存导致的程序退出)
  -v, --version         version for jenkinscli

Use "jenkinscli [command] --help" for more information about a command.

```

#### 1.3.2 远程构建Jenkins并下载工件到指定路径

（以java程序为例：需要jenkins安装插件Pipeline Maven Integration Plugin来提供下载jar包 -p集成了自动下载工件功能，目前适用于java）

```shell
#-p 指定下载路径
./jenkinscli launch -p /tmp/jar 你的流水线名称
```

## 2.配置邮箱

### 2.1 编写配置文件

同样在~/.config/jenkinscli/config 文件里

```yaml
[root@localhost jenkinscli]# cat config.yaml 
Server: http://xxx.xxx.xxx.xxx:8080/
JenkinsUser: admin
Token: 113a8cxxxxxxxxxxxxxxxe8b8
MailSmpt: smtp.qq.com
MailPort: 25
MailUser: xxxxxxxxx7@qq.com
MailToken: lrdxxxxxxxxxgfa
MailFrom: xxxxxxxxx@qq.com
MailTo:
- xxxxxxxxx@qq.com
MailSub: pc-system Test!!!

```

> MsilSmpt , MailPort, MailUser, MailToken分别指定: 服务器地址，端口，邮箱用户名，邮箱token	
>
> 邮箱token 需要你登录到邮箱，在设置里可以拿到自己的token（这里以qq邮箱为例）
>
> MailFrom 设置发件人邮箱用户
>
> MailTo 设置收件人邮箱，可以设置多个，比如：
>
> ```yaml
> MailTo:
> - xxxxxxxxx@qq.com
> - yyyyyyyyy@qq.com
> ```
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

```yaml
MailCc: 
- xxxxxxx.qq.com
- yyyyyyy.qq.com
MailBcc: 
- xxxxxxxxx@qq.com
- yyyyyyyyy.qq.com
MailAttach: /root.1.jpg

```

> MailCc 配置抄送邮箱，可以配置多个
>
> MailBcc 配置密送邮箱，可以配置多个
>
> MailAttach 配置默认附件，同样命令行优先级最高

效果图：

![image-20220513140116132](C:\Users\EDZ\AppData\Roaming\Typora\typora-user-images\image-20220513140116132.png)

## 3.配置sftp和ssh

### 3.1编写ssh配置文件

```yaml
root@localhost ~]# cat .config/jenkinscli/config.yaml 
Sshs: 
- {User: root,Host: 192.168.100.34}
- {User: root,Password: admin12345,Host: 192.168.100.31,Port: 22}
- {User: root,Password: admin12345,Host: 192.168.100.32,Port: 22}
```

> Sshs:是一个主机数组，每一个{}都是一个单独主机 所有字段如下
>
> User: 主机的用户名 （必须设置）
>
> Host: 主机地址	（必须设置）
>
> Password: 主机密码	(若没有设置免密登录就必须设置明文密码)
>
> Port: ssh端口 默认22	
>
> Cmd: 远程执行的bash指令，优先级高于命令行，适用于设置某台主机独立与命令行的bash指令 
>
> Disbash 布尔类型 默认false 禁用某台主机执行如：jenkinscli ssh bash "date" 若字段为true这台主机将不会执行，不会影响其它主机执行
>
> Privatekey 设置私钥路径，用于设置免密登录的主机,不设置时默认为$HOME/.ssh/id_rsa

**提示：当设置密码时使用密码登录，没有密码时会获取密钥来免密登录主机**

### 3.2 使用示例

配置文件：

```yaml
[root@localhost ~]# cat .config/jenkinscli/config.yaml 
Server: http://192.168.xxx.xxx:8080/
JenkinsUser: admin
Token: 11xxxxxxxxxxxxxxxxxxxb8
MailSmpt: smtp.qq.com
MailPort: 25
MailUser: tengxun@qq.com
MailToken: lrxxxxxxxxxxfa
MailFrom: tengxun@qq.com
MailTo: 
- xxxxxxxxx@qq.com
MailSub: pc-system Test!!!
Sshs: 
- {User: root,Host: 192.168.100.34}
- {User: root,Password: a12345,Host: 192.168.100.31,Port: 22}
- {User: root,Password: a12345,Host: 192.168.100.32,Port: 22}
```

其中192.168.100.34 设置了免密登录，可以不用设置密码。

#### 3.2.1 sftp上传文件

jenkinscli sftp upfile

```shell
[root@localhost jenkinscli]# ./jenkinscli sftp upfile -t /root/test /root/1.jpg
2022/05/27 17:50:25 ------------------------------remote host:192.168.100.34:22 path:/root/1.jpg copy file to remote server finished!------------------------------
2022/05/27 17:50:25 ------------------------------remote host:192.168.100.32:22 path:/root/1.jpg copy file to remote server finished!------------------------------
2022/05/27 17:50:25 ------------------------------remote host:192.168.100.31:22 path:/root/1.jpg copy file to remote server finished!----------------
```

-t :--target 用于指定远程目录，会将文件上传到远程目录下

#### 3.2.2 sftp上传文件夹

jenkinscli sftp updir

```shell
[root@localhost jenkinscli]# ./jenkinscli sftp updir -t /root/test /root/hello
```

只会将/root/hello 文件夹里的内容上传到/root/test目录下，指定的目录都必须存在，不会自动创建

#### 3.2.3 sftp下载文件

jenkinscli sftp downfile

```
[root@localhost jenkinscli]# ./jenkinscli sftp downfile -t /root/test/1.jpg /root/aaa
主机数量： 3
2022/05/27 18:04:46 ------------------------------remote host:192.168.100.34:22 path:/root/test/1.jpg copy file to remote server finished!------------------------------
2022/05/27 18:04:46 ------------------------------remote host:192.168.100.32:22 path:/root/test/1.jpg copy file to remote server finished!------------------------------
2022/05/27 18:04:46 ------------------------------remote host:192.168.100.31:22 path:/root/test/1.jpg copy file to remote server finished!------------------------------
```

下载的文件都保存/root/aaa目录下，通过ip文件夹来区分不通机器上的文件：

```shell
[root@localhost aaa]# tree
.
├── 192.168.100.31
│   └── 1.jpg
├── 192.168.100.32
│   └── 1.jpg
└── 192.168.100.34
    └── 1.jpg
```

#### 3.2.4 sftp下载文件夹

jenkinscli sftp downdir

```shell
[root@localhost jenkinscli]# ./jenkinscli sftp downdir -t /root/test/ /root/aaa
```

同下载文件一样，也是通过创建ip目录来存放不通主机上的数据

#### 3.2.5 sftp上传多个文件

> sftp上传多个文件使用go语言的正则表达式来匹配文件，使用RE2语法书写表达式
>
> 不确定自己写的正则表达式是否正确匹配到文件，你可以通过jenkinscli sftp regexp 工具来校验是否匹配到自己想要的文件，示例如下：

jenkinscli sftp regexp

```
[root@localhost jenkinscli]# ./jenkinscli sftp regexp -R ".*txt$" /root
匹配成功的文件： 123.txt
匹配成功的文件： a.txt
匹配成功的文件： t.txt
匹配成功的文件： test.txt
```

> -R 是指定正则表达式，这里可以找出本地/root目录下以txt结尾的所有文件，你也可以将“txt”改为“jar”来匹配jar包文件，实现批量上传

通过上面的命令校验后，我们可以将这些文件一次全部上传到远程主机目录下

例如：将/root目录下通过正则匹配到的文件上传到远程主机/root/test目录下

```shell
[root@localhost jenkinscli]# ./jenkinscli sftp upfilereg -R ".*txt$" -t /root/test /root
```

-t 指定远程目录

#### 3.2.6 ssh发送bash指令

jenkinscli ssh bash

```shell
[root@localhost jenkinscli]# ./jenkinscli ssh bash "date"
2022/05/27 18:25:29 ------------------------------remote host:192.168.100.34:22 exec bash remote server finished!------------------------------
Fri May 27 18:25:29 CST 2022

2022/05/27 18:25:29 ------------------------------remote host:192.168.100.32:22 exec bash remote server finished!------------------------------
Fri May 27 18:25:29 CST 2022

2022/05/27 18:25:29 ------------------------------remote host:192.168.100.31:22 exec bash remote server finished!------------------------------
Fri May 27 18:25:29 CST 2022
```

> 多个命令时需要用双引号括起来，比如./jenkinscli ssh bash "ps -aux | grep java"

我们也可以通过配置文件（Cmd）修改某个主机的执行命令

比如

```yaml
Sshs: 
- {User: root,Host: 192.168.100.34,Cmd: "ps aux |grep nginx"}
```

再执行上面的命令时，192.168.100.34会执行配置文件中加载的命令，其它主机不变，还是执行命令行指定的指令。

通过在位置文件设置Disbash： true 这台主机将不会执行bash指令

#### 3.2.7 ssh远程执行脚本

jenkinscli ssh task

```shell
[root@localhost jenkinscli]# ./jenkinscli ssh task -t /root/test /root/test.sh

2022/05/30 09:48:25 ------------------------------remote host:192.168.100.34:22 path:/root/test.sh copy file to remote server finished!------------------------------
2022/05/30 09:48:25 ------------------------------remote host:192.168.100.34:22 exec bash remote server finished!------------------------------
2022/05/30 09:48:25 ------------------------------remote host:192.168.100.31:22 path:/root/test.sh copy file to remote server finished!------------------------------
2022/05/30 09:48:25 ------------------------------remote host:192.168.100.32:22 path:/root/test.sh copy file to remote server finished!------------------------------
2022/05/30 09:48:25 ------------------------------remote host:192.168.100.32:22 exec bash remote server finished!------------------------------
2022/05/30 09:48:25 ------------------------------remote host:192.168.100.31:22 exec bash remote server finished!------------------------------
```

-t 指定远程执行脚本的工作目录 

-s ubuntu 系统执行sudo密码

-c 为脚本输入位置参数，字符串类型多个参数用空格隔开


### 3.3使用Extend管理主机群

Extend是原来Sshs主机群的一个加强版，你可以配置多个主机群，在命令行指定使用哪个主机群，适用于sftp和ssh模块，如果不知道这个选项，默认使用Sshs主机群

#### 3.3.1 编写配置文件

```yaml
Extend:
- centos:
  - {User: root,Host: 192.168.100.34}
  - {User: root,Password: admin12345,Host: 192.168.100.31,Port: 22}
- openeul:
  - {User: root,Password: admin12345,Host: 192.168.100.31,Port: 22}
  - {User: root,Password: admin12345,Host: 192.168.100.32,Port: 22}

```

字段为Extend

你可以自定义多个主机群，在需要时使用它，比如，这里定义了两个主机群，centos和openeul (名字可以随意取)

使用时使用参数 -H centos 表示作用域centos下的主机群,不指定-H 也就时使用默认的Sshs字段下的主机群。

#### 3.3.2 示例

使用自定义centos主机群:

```shell
[root@localhost jenkinscli]# ./jenkinscli ssh -H centos bash "date"
2022/05/30 16:00:45  ------------------- remote host:192.168.100.34:22 exec bash remote server finished!
Mon May 30 16:00:45 CST 2022

2022/05/30 16:00:45  ------------------- remote host:192.168.100.31:22 exec bash remote server finished!
Mon May 30 16:00:45 CST 2022
```

同时sftp和ssh模块都是一样的原理，不指定-H就表示使用默认的Sshs字段的主机群。

#### 3.3.3 完整配置文件示例参考

```yaml
Server: http://127.0.0.1:8080/
JenkinsUser: admin
Token: 113a8b8
MailSmpt: smtp.qq.com
MailPort: 25
MailUser: 123456789@qq.com
MailToken: lrfa
MailFrom: 123456789@qq.com
MailTo:
- 987654321@qq.com
MailSub: pc-system Test!!!
Sshs:
- {User: root,Host: 192.168.100.34}
- {User: root,Password: admin12345,Host: 192.168.100.31,Port: 22}
- {User: root,Password: admin12345,Host: 192.168.100.32,Port: 22}
Extend:
- centos:
  - {User: root,Host: 192.168.100.34}
  - {User: root,Password: admin12345,Host: 192.168.100.31,Port: 22}
- openeul:
  - {User: root,Password: admin12345,Host: 192.168.100.31,Port: 22}
  - {User: root,Password: admin12345,Host: 192.168.100.32,Port: 22}
```

## 4.使用动态生成文件tmpl功能

tmpl动态生成文件使用的go语言template包，需要提供一个模板文件，模板定义示例：tmpl.yaml 文件名

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: {{.number}}
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: {{.image}}
        ports:
        - containerPort: 80
```

{{.name}} 来定义变量的

### 4.1 tmpl local 用法

这个是用于本机动态生成文件，可以解析模板文件后输出到终端和文件里

上面的例子中，我们可以动态生成文件：

```bash
 ./jenkinscli -I tmpl local -l /root/tmpl.yaml "/root/test.yaml" "image:redis" "number:4"
```

-I :表示忽略jenkins 初始化错误，如果你没有配置Jenkins登录信息程序会退出

-l :指定模板文件

第一个参数 :"/root/test.yaml" 指定输出位置

第二个到无限个开始为自定义变量，用于替换模板文件中的变量。格式为name:value  。其中name为你模板中变量的名一致。

**注意：如果没有为模板文件中的变量在命令行设置替换值name:value ,那么生成的新文件的变量位置就会为空。**

如果第一个参数指定为nil 就会将解析后的模板文件输出到终端：

```bash
./jenkinscli -I tmpl local -l /root/tmpl.yaml "nil" "image:redis"
```

### 4.2 tmpl remote 用法

可以将模板文件解析后发送到远端机器的目录下,操作类似于sftp上传文件：

```bash
 ./jenkinscli tmpl remote -I -l -H openeul /root/tmpl.yaml "/root/k8stest.yaml" "image:mongo" "number:3"
```

-I :表示忽略jenkins 初始化错误，如果你没有配置Jenkins登录信息程序会退出

-l :指定模板文件

-H : 指定主机组名称，如果不指定默认就是默认主机组Sshs。主机组就是定义一个主机群，详细见sftp功能的使用说明

**注意：remote 相比local 是不能将解析的模板文件输出到终端的，也就是不支持"nil"**

### 4.3 tmpl 模板高级用法

和go语言的模板用法一致，这里列举几个常见用法：

1.当我们在命令行中没有为模板的变量设置值时，我们或许需要为它们设置一个默认值，可以使用模板的条件语句，用法如下(还是已上方yaml文件为例)

```bash
{{if pipeline}} T1 {{else}} T0 {{end}}
    如果pipeline的值为empty，输出T0执行结果，否则输出T1执行结果。不改变dot的值。
```

修改如下：

```yaml
 replicas: {{if .number}}{{.number}}{{else}}1{{end}}
```

这样当没有指定replicas值时就会设置为1 

我们也可以用其它的语法：

```yaml
 replicas: {{with .number}}{{.}}{{else}}1{{end}}
```

与上方的效果是一样的

```bash
{{with pipeline}} T1 {{else}} T0 {{end}}
    如果pipeline为empty，不改变dot并执行T0，否则dot设为pipeline的值并执行T1。
```

with 表示会取出 .number 的值再赋值给 '.'  因此直接使用{{.}}就是它的值。如果不是with是if, '.'表示一个map类型。我们需要'.number'来获取值，也就是说要指定一个字段名

## 5.常见问题

## 4.常见问题
1.如果在配置文件没有指定Jenkins登录的信息会导致jengkins初始化失败从而退出程序。通常在单独使用某个子命令时遇到，比如使用sftp 、ssh、email等。**我们可以通过-I来忽略Jenkins初始化报错导致的程序退出。这样你就可以单独使用其它子命令了。**
