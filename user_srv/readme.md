### question1
如果设计一个用户服务具备通用性，比如可以让所有的系统都可以使用一套公共代码，但是你同的系统 user 表上可能会有不同的字段，如何设计让系统具备通用性的同时还能具备好的扩展性？
> 思路：
> 
> (1) 基本上所有的系统都需要用户名和密码、登录时间等，这些可以设计成一张通用的表（用户公共表）
> 
> (2) 如何可以扩展表并且不会对现有的表产生影响
>
> (3) 扩展接口，比如将这一套的用户服务完善好，把一整套的用户相关接口都自己实现好

### question2
自己写一个exe文件可以使得生成基本的service微服务脚手架，这个脚本可以在启动的时候让用户输入一些信息，哪写信息可以通过用户输入进行配置？

> (1) 对于service和web端来说，两种代码的目录结构会不一致，所以该命令行可以支持两种类型
> 
> (2) 比如后期可以考虑，是否支持服务注册、日志框架替换、中间件的依赖替换等
> 
> (3) 命令行模式基本是微服务中必备的，go-micro和go-zero等解决方案都支持通过命令行生成模板目录。
### question3
启动service服务虽然比较简单，但是还是不够灵活，我们是否可以考虑自己编写makefile使得微服务启动变得更加简单？
> (1) 每次启动服务前我们都可能设置一定的参数，通过makefile将启动前的准备工作配置到makefile中
> 
> (2) 后期如果想要将每个微服务部署在docker中，这样有了makefile也可以将这个过程一键完成
