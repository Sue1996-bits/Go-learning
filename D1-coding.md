 参考书《The Go Programming Language》
 ![屏幕截图 2025-06-30 162648](https://github.com/user-attachments/assets/e36fedc6-9b7c-4a0a-8e50-c4837174e01a)
go是一个编译语言（即通过编译器转为机器码：二进制，即.exe可执行文件，再运行）
区别：解释语言python（源码被解释器解释生成，每次都需解释）

go命令有一系列子命令，调用
$ go run --.go    编译
$ go build --.go    编译并生成.exe

-开头package main，表示该文件属于哪个包
-紧跟着一系列导入（import）的包
-之后是存储在这个文件里的程序语句

fmt 包 ：格式化输出、接收输入的函数  Println
main 包：它定义了一个独立可执行的程序，而不是一个库。main函数是程序的真正入口。
“一个 Go 程序中，只有一个 main 包，且必须包含一个 main() 函数，且只有它可以包含func main（）。即告诉编译器，这不是别人调用的库，而将生成可执行的二进制文件”

处理输入：
命令行参数。
os包：与OS交互
os.Args 变量是一个字符串（string）的 切片（slice）

导师布置问题：
部署两个接口AB，使用go中的gin完成。
A是访问照片，并对照片进行加密处理，生成加密后的文本或者...
B是对已知的加密后信息进行解密，并从数据库中找寻调用原图片资源。
