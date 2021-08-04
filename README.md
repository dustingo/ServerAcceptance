##### 服务器验收
USAGE:
```shell
	./ServerAcceptance --precheck --config=xxx
	./ServerAcceptance --printjson
	./ServerAcceptance --lastcheck --config=xxx
```
> --precheck 服务器基础信息验收

  基础信息验收【precheck】主要是对于基础硬件方面的验收，主要服务器交付时的： 
  系统信息  
  CPU信息  
  内存信息  
  硬盘信息  
  网卡信息  
  输出会直接显示配置文件和实际参数是否相符

> --printjson 服务器基础信息打印

  将会以json格式打印出服务器基础信息，便于以此配置配置文件

> --lastcheck 最终结果校验

  检查的点有： 
  系统服务是否开启或关闭  
  yum、pip、perl包是否安装  
  目录是否存在以及权限是否正常  
  ulimit信息是否正确  
  内核参数是否正确  
  dns信息是否正确


> NOTE

  所有的工作只是校验，不进行操作

 
  


> precheck和lastcheck 模板如下

  precheck 模板：https://raw.githubusercontent.com/dustingo/configtoml/main/config.toml  
  lastcheck 模板: https://raw.githubusercontent.com/dustingo/configtoml/main/lastcheck.toml
