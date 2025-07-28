接口就是方法的集合定义--只定义不实现，由具体的对象来实现细节。

--Go语言中提倡使用面向接口的编程方式实现`解耦`。

```
type 接口类型名 interface{
    方法名1( 参数列表1 ) 返回值列表1
    方法名2( 参数列表2 ) 返回值列表2
    …
}
```
命名：——er ：Writer，closer

当方法名首字母是大写且这个接口类型名首字母也是大写时，
这个方法可以被接口所在的包（package）之外的代码访问。

`一个类型只要实现了接口中规定的所有方法，那么我们就称它实现了这个接口。`

1.定义Sayer接口，有Say()方法
2.然后我们定义一个通用的MakeHungry函数，接收Sayer类型的参数。
3.struct只要实现了Say()方法都能当成Sayer类型的变量来处理。
```
func MakeHungry(s Sayer) {
	s.Say()
}

var c cat
MakeHungry(c)
var d dog
MakeHungry(d)
```
PHP和Java语言中需要显式声明一个类实现了哪些接口，在Go语言中使用隐式声明的方式实现接口
--符合程序开发中抽象的一般规律。

```
// Payer 包含支付方法的接口类型
type Payer interface {
	Pay(int64)
}

// Checkout 结账
func Checkout(obj Payer) {
	// 支付100元
	obj.Pay(100)
}

func main() {
	Checkout(&ZhiFuBao{}) // 之前调用支付宝支付

	Checkout(&WeChat{}) // 现在支持使用微信支付
}
```
一个接口类型的变量能够存储所有实现了该接口的类型变量。
即：var x Sayer 中，x可以被cat、dog..不同的类型变量赋值 

#1.2 接口：值接收者和指针接收者

