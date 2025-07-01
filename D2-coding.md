!使用gin框架（modules）
发现又有新版本需求，打算直接更新最新版本。。。
这里顺便说一下怎么更新:
我直接是控制面板卸载，然后注意一下新版本的环境变量系统变量位置 GOPATH GOROOT
工具VSCODE这里会提示重新更新工具

下载成功结果是：go: added ————

因为我没有接触过用go的基础库的web方向操作，所以这里花了半天时间学习实践了一下。
<img width="630" alt="cgi-bin_mmwebwx-bin_webwxgetmsgimg___MsgID_8608713345148903234_skey__crypt_3761f30_e2a56e0ecbd30496e" src="https://github.com/user-attachments/assets/eb47fe78-ccb7-4f9c-90d5-c7437ab999a0" />

这里顺便复习一下web知识：
 电脑浏览器--输入url--向服务器发送了请求
          --http协议返回--
！一个请求对应一个响应

Gin框架学习：参照博客：https://liwenzhou.com/posts/Go/gin/
传统：接受html文件，服务器http协议读取返回
Gin：只提供json格式数据，前段ajax（前后端分离）

补充RESTful API：软件架构风格--增删改查分化风格
GET用来获取资源
POST用来新建资源
PUT用来更新资源
DELETE用来删除资源。

