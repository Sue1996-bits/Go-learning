https://www.liwenzhou.com/posts/Go/struct/

type关键字来定义自定义类型。
1.基于 内置的基本类型--int、string 定义：
```
//将MyInt定义为int类型
type MyInt int
```
2.通过struct 定义
`让一个自定义的数据类型封装多个基本数据类型。—— 实现面向对象。`
```
type person struct {  //person 类型名，同一个包不能重复
	name string         //结构体中的字段名必须唯一
	city string
	age  int8
}
```
内置的基础数据类型是用来描述一个值的，而结构体是用来描述一组值。

1.2结构体的构造函数：

“ 因为struct是值类型，如果结构体比较复杂的话，值拷贝性能开销会比较大，所以该构造函数返回的是结构体指针类型。”
```
func newPerson(name, city string, age int8) *person {
	return &person{
		name: name,
		city: city,
		age:  age,
	}
}
```
1.3 方法（区别`函数`）

方法（Method）是一种作用于`特定类型变量——接收者（Receiver）——类似this或者 self。`的函数。
```
//Dream Person做梦的方法
func (p Person) Dream() {
	fmt.Printf("%s的梦想是学好Go语言！\n", p.name)
}

func main() {
	p1 := NewPerson("小王子", 25)
	p1.Dream()
}
```
方法名、参数列表、返回参数：具体格式与函数定义相同

接收者变量：

·命名为接收者类型名称首字母的小写

·接收者类型和参数类似，可以是指针类型和非指针类型。

(1)若为指针：在方法结束后，修改都是有效的。
```
// 方法SetAge 设置p的年龄
// 使用指针接收者
func (p *Person) SetAge(newAge int8) {
	p.age = newAge
}

func main() {
	p1 := NewPerson("小王子", 25)
	fmt.Println(p1.age) // 25
	p1.SetAge(30)
	fmt.Println(p1.age) // 30
}
```
什么时候应该使用指针类型接收者?
需要修改接收者中的值.
接收者是拷贝代价比较大的大对象.
保证一致性，如果有某个方法使用了指针接收者，那么其他的方法也应该使用指针接收者。

(2)值类型:代码运行时将接收者的值复制一份。修改操作只是针对副本，无法修改接收者变量本身。
```
// SetAge2 设置p的年龄
// 使用值接收者
func (p Person) SetAge2(newAge int8) {
	p.age = newAge
}
func main() {
	p1 := NewPerson("小王子", 25)
	p1.Dream()
	fmt.Println(p1.age) // 25
	p1.SetAge2(30) // (*p1).SetAge2(30)
	fmt.Println(p1.age) // 25
}
```
`方法与函数的区别是，函数不属于任何类型，方法属于特定的类型。`
但`接收者的类型可以是任何类型，不仅仅是结构体，任何类型都可以拥有方法。`

1.4结构体内可以嵌套子结构体/结构体指针

--继承：
```
//Dog 狗
type Dog struct {
	Feet    int8
	*Animal //通过嵌套匿名结构体实现继承
}
```

1.5 tag
```
type WatermarkLog struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	ImageID     string    `json:"image_id" gorm:"not null;index"`
	WatermarkID string    `json:"watermark_id" gorm:"not null;index;size:64"`
	RequestIP   string    `json:"request_ip" gorm:"size:45"`
	UserAgent   string    `json:"user_agent" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time `json:"updated_at"`
//Tag在结构体字段的后方定义，由一对反引号包裹起来
//格式`key1:"value1" key2:"value2"` 如：`json:"id" gorm:"primarykey"`
//`json:"id"`--通过指定tag实现json序列化该字段时的key（若未指定，则json序列化时默认使用字段名作为key）
}
```
