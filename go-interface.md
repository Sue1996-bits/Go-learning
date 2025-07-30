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

### 1.2 接口：值接收者和指针接收者
```
// Move 使用值接收者定义Move方法实现Mover接口
func (d Dog) Move() {
	fmt.Println("狗会动")
}
```
使用值接收者实现接口之后，不管是结构体类型还是对应的结构体指针类型的变量都可以赋值给该接口变量。
```
// Move 使用指针接收者定义Move方法实现Mover接口
func (c *Cat) Move() {
	fmt.Println("猫会动")
}
```
此时只能将 *Cat类型的变量 赋值给 接口类型变量x（var x Mover）

### 1.3 多种类型实现同一接口
一个接口的所有方法，不一定需要由一个类型完全实现，接口的方法可以通过在类型中嵌入其他类型或者结构体来实现。
```
tpye X interface{
  A()
  B()
}//接口要求实现两个方法


type X1 struct{}
func (x1 X1)A(){
  fmt.Println("实现了A方法")
}


type X1-tpye202 struct{
  X1  //嵌入X1，继承A方法
}
func (hyf X1-tpye202) B() {
	fmt.Println("实现了B方法")
}
//X1-tpye202通过嵌入X1的办法实现X1可包含的方法，即是功能件和整体的关系（eg：显示器和电脑）
```
### 1.4 接口嵌套
```
// ReadWriter 是组合Reader接口和Writer接口形成的新接口类型
type ReadWriter interface {
	Reader
	Writer
}
```
通过在结构体中嵌入一个接口类型，从而让该结构体类型实现了该接口类型，并且还可以改写该接口的方法。
### 空接口
通常我们在使用空接口类型时不必使用type关键字声明，可以像下面的代码一样直接使用interface{}。

var x interface{}  // 声明一个空接口类型变量x
1.实现可以接收任意类型的函数参数。
2.使用空接口实现可以保存任意值的字典。
