--类型中转,语义分析顺序
function change()
    a={}
    --a,此处获取a.类型
end

a=100
change()
 --a,此处获取a.类型