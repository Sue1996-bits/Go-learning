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
JWT （JSON Web Token） 是目前最流行的跨域认证解决方案，是一种基于 Token 的认证授权机制--一种规范化之后的 JSON 结构的 Token（令牌）。

JWT 自身包含了身份验证所需要的所有信息，因此，我们的服务器不需要存储 Session 信息。

JWT 本质上就是一组字串，通过（.）切分成三个为 Base64 编码的部分：xxxxx.yyyyy.zzzzz

1.xxxxx--Header（头部） : 描述 JWT 的元数据，定义了生成Signature的算法以及 Token 的类型。
<img width="1324" height="664" alt="image" src="https://github.com/user-attachments/assets/7c604363-d373-4097-b174-5d12f66239cc" />

2.yyyyy--Payload（载荷） : 用来存放实际需要传递的数据，包含声明（Claims），如sub（subject，主题）、jti（JWT ID）。
<img width="1340" height="950" alt="image" src="https://github.com/user-attachments/assets/f21cdb57-3a50-4a91-9fe8-2150daea23a3" />

Claims(声明，包含 JWT 的相关信息)。--三种。
Payload 部分默认是不加密的，一定不要将隐私信息存放在 Payload 当中！！！

3.zzzzz--Signature（签名）：服务器通过 Payload、Header（前两者） 和一个密钥(Secret)使用 Header 里面指定的签名算法（默认是 HMAC SHA256）生成。
<img width="1374" height="830" alt="image" src="https://github.com/user-attachments/assets/0c5360e4-10f4-4827-b463-d3acf7600cba" />

1.2基于 JWT 进行身份验证
<img width="1282" height="754" alt="image" src="https://github.com/user-attachments/assets/641ae3ac-0624-48b4-932e-bffd6ac094cc" />
服务器通过 Payload、Header 和 Secret(密钥)创建 JWT 并将 JWT 发送给客户端。客户端接收到 JWT 之后，会将其保存在 Cookie 或者 localStorage 里面，以后客户端发出的所有请求都会携带这个令牌。

即服务器只需要保存一个统一的Secret(密钥)即可。

------
著作权归JavaGuide(javaguide.cn)所有
基于MIT协议
原文链接：https://javaguide.cn/system-design/security/basis-of-authority-certification.html
