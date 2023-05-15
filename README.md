# 使用go写一个简易版的解释器

最近花了几天时间看了一本书head first go，了解一些go语法，感觉很有意思。很简洁，而且在有些地方觉得和ts语法有些类似。

刚好最近也在看一本解释器的书，于是想试试能否跟着书用go写出来写一个简单的解释器，顺便我也想用node再写一遍。
书中把要实现的`语言`叫做`Monkey`。这里也就延续这个叫法吧。
这里想实现：

1. 一种类似C的语法
2. 变量绑定
3. 整形和布尔值
4. 算数表达式
5. 内置函数
6. 头等函数和高阶函数（函数第一公民）
7. 闭包
8. 字符串数据结构
9. 数组数据结构
10. 哈希数据结构

## 解释器

解释器应该包含如下部分：

1. 词法分析器
2. 语法分析器
3. 抽象语法树AST
4. 内部对象语法
5. 求值器

## 词法分析

为了解释源代码，我们需要将其转换为其他易于处理的形式。具体来说就是在对最终代码求值之前，需要两次转换源代码的表示形式。
由源码先转为词法单元，也就是所谓的Token，然后将token转换为抽象语法树AST。

### token

实现词法分析肯定离不开token.我们将token相关的东西放在了token包下。
具体的代码可以到该文件夹下看了。

### lexer

lexer包下主要是负责真正的词法解析的。将源代码解析生成一个个token，去除无用的空白符。

## 语法分析

词法分析只是解析为一个个的token，并不会检测语法错误。向上面我们的测试代码，`lexer/lexer_test.go`文件中的：

```go
 !-*/5;
 5 < 10 > 5;
```

这很明显就是错误的语法，但是我们依旧可以进行词法分析。其实这也是单一职责的一种体现。
代码中的错误之所以可以报告出来，都是靠语法分析来做的。
语法分析将输入的内容转换成对应的数据结构。听起来很抽象，但是我们可以结合一个下面这个js的例子！

```js
const input = '{ "name":"zs", "age": 22 }'
const output = JSON.parse(input)
console.log(output.name, output.age)
```

输入input虽然是一个JSON字符串，但是通过parse解析（其实就是一个语法分析器），就可以得到一个js语言中的对象，获取对象中的两个属性name和age.其实，parse解析和我们这里说的语法解析器并没有本质的区别。
经过语法分析后生成的数据结构，在多数语言中都称为"语法树"或者说"抽象语法树"(Abstract Syntax Tree AST)。在抽象语法树中会省略一些源代码中可见的某些细节。
比如说我们会省略分号，空格，注释，花括号，方括号，括号等信息，不让其出现在AST中。

### 上下文无关法

**上下文无关法（context-free grammer, CFG）**：CFG是一组规则，描述了如何根据一种语言的语法狗证正确的语句。CFG最常用的符号格式是Backus-Naur形式（BNF）或Extended Backus-Naur形式（EBNF）。

### 编写语法分析器

编写策略： 自上而下的分析或之下而上的分析。每种策略都有很多变体。
例如：

1. 递归下降分析
2. Earley分析
3. 预测分析
这些都是智商而下分析的变体。
那么我们这里也采用递归下降的语法分析器。具体的说，它是基于自上而下的运算符优先级分析法的语法分析器。
这里编写的语法分析器其实局限性很大，比如它可能不是很快，也没用对其正确性和错误恢复过程进行形式化的证明，错误语法的检测也不是无懈可击的。如果不深入研究语法分析相关的理论，那么很难真的解决最后一个问题。

### 第一步，解析let语句

```js
let x = 5;
let y = 10;
let foobar = add(5, 5)
let barfoo = 5 * 5 / 10 + 18 - add(5, 5) + multiply(124);
let anotherName = barfoo;
```

上面这个是一个复杂的例子，是我们完成语法分析后才能解决的。我们可以先解析不带表达式的let语句。比如下面这个例子：

```js
let x = 10;
let y = 15;
let add = fn(a, b){
  return a + b;
}
```

我们使用let语句实现了3个变量绑定。let的形式如下：

```js
let <标识符> = <表达式>;
```

对于语句，不会产生值，但是表达式会。
就像`let x = 5;`不会产生值，而`5`会产生值（也就是`5`）。`return 5;`不会产生值，但是`add(5,5)`会产生值.

### 解析return语句

例如：

```js
return 5;
return add(10);
```

**return语句的结构：**

```js
return <表达式>;
```

return语句仅由关键字return和表达式组成，因此：我们可以定义returnStatement结构来表达。

### 解析表达式

我们只有let和return两种语句。那么接下来我们需要解析表达式了。
解析表达式应该算是语法分析中最难的部分吧。解析语句的过程中我们是从左到右处理词法单元，然后期望或拒绝下一个词法单元，如果一切正常，最后就返回一个对应的AST节点。
但是表达式不一样，比如：我们可能遇到的第一个难点就是运算符优先级。这算是一个挑战了。
举例如下：

```js
5 * 5 + 10
```

对于这个小例子，其实我们应该是先计算`5*5`，在计算`25+10`.也就是说`5*5`应该是更深一个层级的ast，因为它是优先于加法运算求值的。为了生成目标AST，语法分析器必须知道`*`的优先级是高度`+`的。
但是对于下面这个小例子：

```js
5 * (5 + 10)
```

这个优先级又不一样了。因为括号提升了`+`的优先级，我们应该先计算`+`，在去计算乘法。
表达式的种类也有很多种：
比如：
前缀表达式：

```js
-5
!true
+10
```

中缀表达式：

```js
5 + 10
5 - 5
10 / 2
10 * 2
```

比较运算符表达式：

```js
10 > 5
true  == true
foo == bar
foo > bar
```

分组表达式：

```js
(5 + 3) * (10 - 5)
add(5, add(1, (10 + 2)))
```

标识符也是表达式的一种：

```js
foo * bar / foobar + barfoo
add(foo, bar)
```

函数字面量也是表达式，可以绑定在一个变量上，也可以直接使用：

```js
let foo = fn(x, y) {return x + y; }
(fn(x, y) { return x + y; })(10, 20)
```

if表达式：

```js
let result = if (10 > 5) { true } else { false }
```

### 自上而下的运算符优先级分析

一种基于上下文无关文法和Backus-Naur-Form语法分析器的替代方法。

### 标识符

标识符是最简单的表达式类型。标识符是产生值的。

```js
foo;
```

### 整数字面量

和标识符基本一样，解析方式大同小异。

```js
10;
5
```

### 前缀运算符

这里支持`-x`和`!x`两个前缀运算符。

```js
-5;
!foobar;
!true;
5 + -10;
```

用法结构如下：

```js
<前缀运算符><表达式>
```

### 中缀表达式

```js
5 + 5;
5 * 5;
5 > 10;
5 == 10;
```

中缀运算符的左右可以使用任何表达式，不仅仅是数字。

```js
<表达式> <中缀运算符> <表达式>
```

可以看出中缀表达式有三部分。左右两侧都是操作数，或者说表达式。因此我们也可以称为二元表达式。而前缀表达式则称为一元表达式。
现在先支持两侧表达式都是数字的中缀表达式。

### 布尔字面量

能使用表达式的地方都可以使用布尔字面量。

```js
true;
false;
let foobar = true;
let foobar = false;
```

### 分组表达式

就是使用分组括号可以提升优先级。

```js
(5 + 5) * 10;
```

### if表达式

和其他编程语言一样，都可以使用if和else。

```js
if (10 > 5) {

} else {

}
```

else是可以省略的。

在我们这里，if-else条件语句是表达式。这意味着其中的语句会产生值，对于if表达式而言，是最后求值的代码产生值。因此这里不需要return语句。

```js
let foobar = if (x > y) { x } else { y }
```

那么if-else条件句的结构如下：

```js
if (条件) <结果> else <结果>
```

注意if语句的小括号是不可少的.

### 函数字面量

函数字面量是定义函数的方式，其中包括函数的参数及其作用。函数字面量如下所示：

```js
fn(x, y) {
  return x + y;
}
```

函数字面量一般以关键字fn开头，后跟一个参数列表，再后面跟一个块语句。块语句是函数的主体，调用函数时会执行块语句。
函数字面量的抽象结构如下所示：

```js
fn <参数列表> <块语句>
```

参数标识符列表：

```js
(<参数1>, <参数2>, <参数3>, ...)
```

当然，参数可以留空：

```js
fn() {
  return 10;
}
```

### 调用表达式

了解如何解析函数字面量了，下一步是来看解析函数的调用。即调用表达式。
结构如下所示：

```js
<表达式>(<以逗号分隔的表达式列表>)
```

没错，就是这么简单。就像这样：

```js
add(10, 20)
```

标识符add是一个表达式，经过替换可能是这样：

```js
fn(x, y) { return x + y; }(10, 20)
```

而且，函数字面量还可以当做参数：

```js
callFn(2, 10, fn(x, y){ return x + y; })
```

那么上面的结构如下：

```js
<表达式>(<逗号分隔的表达式列表>)
```

### repl

当前的repl类似read-lex-print Loop（读取词法分析词打印循环）.
现在我们可以替换词法分析（Lex）为语法分析（Parse）来构建新的repl。

## 求值

### 为符号赋予含义

来到了真正开始求值的环节。求值就是evaluation这一步，在repl中的e就是表示求值。是解释器处理源码过程中最后一步。代码求值后才会变得更有意义。如果不进行求值，那么类似 `1 + 2`的表达式转换后也只是一组字符串，一组词法单元构成的树结构罢了。没有真正的含义的。

解释器的求值过程决定了语言的解释方式。

```js
let num = 5;
if (num) {
  return a;
} else {
  return b;
}
```

如上面的例子，解释器的求值过程决定了上面的代码是返回a还是返回b。我们需要分析num变量是否是真值。

### 求值策略

在解释器中，求值都是差异最大的。源代码求值可以有不同的策略。

1. 最直观的方式就是解释执行AST。也就是说遍历AST，访问每个节点并执行该节点的语义，像是打印字符串，添加两个数字，执行函数的主体等。这些都是实时进行的。所以这种方式也被称为树遍历解释器，最经典的解释器之一。当然，有时候在求值步骤之前，解释器会进行少量优化，包括重写AST（删除未使用的变量绑定）或者将其转换为更适合递归和重复求值的另一个中间表示（IR）。
2. 其他类型的解释器也是遍历AST，不过不是解释AST本身，而是先转为字节码再解释。

字节码算是AST的另一种中间表示，信息密度比AST高。不同解释器的字节码格式及其操作码（构成字节码的指令）各不相同，这取决于解释器本身及其所实现的宿主编程语言。通常操作码与大多数汇编语言的助记符很相似。和第一种方式来比，这种方式的性能更好。
3. 最后一种策略是不涉及AST。因为语法分析器不构建AST，直接生成字节码了。
在上面最后一种策略中，会显得解释器和编译器的解析不明显，这到底是解释器还是编译器？生成并解释（或执行）字节码不就是一种形式的编译吗。之前说了，编译和解释之间的界限非常模糊。有些情况更加模糊。就像某些语言的实现会解析源代码，构建AST来生成字节码。在执行之前，虚拟机会即使将字节码编译为机器码，而不是直接在虚拟机中执行字节码指定的操作，这就是所谓的JIT解释器/编译器。

### 树遍历解释器

使用语法分析构建的AST，直接结束它，不经过任何预处理或编译步骤。
虽然从上面到现在感觉求值的过程很复杂，实际上就是一个eval函数，其工作就是对AST求值。

### 表示对象

对象系统也可以称为值系统或者对象表示方法。这里的重点是，需要为eval函数所返回的内容添加一个定义。也就是说，我们需要一个系统用来表示AST的值或者表示在内存中对AST求值时生成的值。
可以阅读`wren`解释器项目的源代码。

### 对象系统集成

这里解释器的性能没有什么要求，因此选择了简单的方法，也就是在对monkey源代码求值时，每个遇到的值都会用Object表示。
这个Object是我们为monkey语言设计的接口。也就是说所有的值都会封装到一个符合Object接口的结构体中。

#### 整数

支持Integer类型实现Object接口。
每当源代码遇到整数字面量的时候，都需要先转为`ast.IntegerLiteral`。然后对该AST节点求值时，将其转换为`object.Integer`，以便将整数的值保存在结构体中并传递这个结构体的引用。

#### 布尔值

布尔类型和整数一样，都是很简单的。同上。

#### 空值

语言如果有了null类型，那么遇到使用null的地方，都应该三思而后行。如果语言不存在null，那么使用起来会更加安全。

### 求值表达式

eval函数将Node作为输入并返回一个Object。提醒一下，之前在ast包的结构体都实现了Node接口。

#### 整数字面量

我们有一个表达式，其中仅仅包含一个整数字面量，我们将其作为输入进行求值，希望返回的结果是整数本身。

#### 布尔字面量

实现原理类似整数字面量。

#### 空值

和上面的布尔值一样，布尔值对应的对象只有两个，空值也应该是全局唯一的。

### 前缀表达式

前缀表达式也是一元运算符表达式，它是monkey支持的最简单的运算符表达式形式。
我们需要先对前缀表达式的右侧表达式求值，然后再和运算符求值。

### 中缀表达式

```js
5 + 5;
5 - 5;
5 * 5;
5 / 5;

5 > 5;
5 < 5;
5 == 5;
5 != 5;
```

上面四组表达式产生的结果不是布尔值，而下面表达式结果都是布尔值。

### 条件语句

在条件语句中，我们唯一的难点就是如何决定在何时对那一部分内容求值，这就是条件语句的全部要点：根据条件决定求值内容！

```js
if (x > 10) {
  puts("ok")
} else {
  puts("fail")
}
```

条件正确，就对该分支求值，但是不能对else分支求值。
这里的条件正确，指的就是条件表达式代表的值是真值。不需要真值必须是true或者false。

```js
if (false) {
  return 1
}
```

对于上面这个例子，我们条件是false，且没有else分支，那么条件语句返回的结果就是NULL。

### return语句

和其他语言一样，monkey也具备return语句。可以用在函数体中，也可以作为顶层语句来使用。在哪里使用return语句并不重要，其工作方式都是相同的：停止对一系列语句的求值，同时保存其中表达式求值的结果~。

```js
return 10;
5 * 5;
```

像上面的例子中，不会对最后的`5*5`求值。

### 错误处理

真正的错误处理不是指用户定义的异常，而是指处理程序内部的错误处理流程。用于处理错误的运算符，不支持的操作以及执行期间可能发生的其他用户错误或内部错误。

### 绑定与环境

接下来添加对let语句的支持，为解释器增加绑定功能。除了支持let语句，还需要支持对标识符求值。假设对以下代码求值：

```js
let x = 5 * 5;
```

仅实现对这个语句的求值还不够，还需要确保求值后x的值是25.

### 函数和函数调用

实现函数和函数调用的支持~。不仅仅是这些，我们应该还需要实现传递函数，高阶函数和闭包功能。
