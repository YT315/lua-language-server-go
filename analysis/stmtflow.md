```flow
st=>start: 开始
st2nameexpr=>condition: 将st的name是否为nameexpr
nameexprhavetype=>condition: nameexpr是否有类型
errnottype=>operation: 报错不可以有类型
creatnamesymbol=>operation: 将nameexpr创建一个符号symbol
findtaboutside=>condition: 向外层寻找此符号的标签的上下文
symbolctxdefifnil=>condition: 符号上下文中定义数量是否为空
errredef=>operation: 报错错重定义
e=>end: 结束

st->st2nameexpr
st2nameexpr(yes)->nameexprhavetype
st2nameexpr(no)->e
nameexprhavetype(yes)->errnottype->creatnamesymbol
nameexprhavetype(no)->creatnamesymbol->findtaboutside
findtaboutside(yes)->symbolctxdefifnil


```



