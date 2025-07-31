在这里补充需要的背景知识。

### JWT 是什么--Web安全：JWT（JSON Web Token）认证

p1：认证授权基础概念详解
------


认证 (Authentication)： 你是谁 + 授权 (Authorization)： 你有权限干什么。

最常采用的访问控制模型就是 RBAC 模型--基于角色的权限访问控制（Role-Based Access Control）

用户与角色、角色与权限=n:m:z关系。

Cookie 存放在客户端(用户本地)，一般用来保存用户信息（通常经过加密）。

使用 Cookie 保存 SessionId 或者 Token ，向后端发送请求的时候带上 Cookie，这样后端就能取到 Session 或者 Token 了。这样就能记录用户当前的状态了，因为 HTTP 协议是无状态的。

also：记录和分析用户行为--如果服务器想要获取你在某个页面的停留状态或者看了哪些商品，一种常用的实现方式就是将这些信息存放在 Cookie

Session 的主要作用就是通过服务端记录用户的状态--服务端给特定的用户创建特定的 Session 之后就可以标识这个用户并且跟踪这个用户了（因为 HTTP 协议无状态）--数据保存在服务器端。

Session-Cookie 身份验证：
1.用户成功登陆系统，然后返回给客户端具有 SessionID 的 Cookie 。
2.当用户向后端发起请求的时候会把 SessionID 带上，这样后端就知道你的身份状态了。

！如果别人通过 Cookie 拿到了 SessionId 后就可以代替你的身份访问系统了。
different with `Token`:获得 Token 之后，一般会选择存放在 localStorage （浏览器本地存储）中。然后我们在前端通过某些方式会给每个发到后端的请求加上这个 Token

p2：认证授权-JWT
------
JWT （JSON Web Token） 是目前最流行的跨域认证解决方案，是一种基于 Token 的认证授权机制--一种规范化之后的 JSON 结构的 Token。

JWT 自身包含了身份验证所需要的所有信息，因此，我们的服务器不需要存储 Session 信息。

------
著作权归JavaGuide(javaguide.cn)所有
基于MIT协议
原文链接：https://javaguide.cn/system-design/security/basis-of-authority-certification.html
