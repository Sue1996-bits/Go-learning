在这里补充需要的背景知识。

### JWT 是什么--Web安全：JWT（JSON Web Token）认证

认证 (Authentication)： 你是谁 + 授权 (Authorization)： 你有权限干什么。

最常采用的访问控制模型就是 RBAC 模型--基于角色的权限访问控制（Role-Based Access Control）

用户与角色、角色与权限=n:m:z关系。

Cookie 存放在客户端(用户本地)，一般用来保存用户信息（通常经过加密）。

使用 Cookie 保存 SessionId 或者 Token ，向后端发送请求的时候带上 Cookie，这样后端就能取到 Session 或者 Token 了。这样就能记录用户当前的状态了，因为 HTTP 协议是无状态的。


------
著作权归JavaGuide(javaguide.cn)所有
基于MIT协议
原文链接：https://javaguide.cn/system-design/security/basis-of-authority-certification.html
