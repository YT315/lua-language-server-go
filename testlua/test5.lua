--不确定判断修改变量值
local a={}

a.b=100

if a.b==100 then
    a=200
end

--a,此处获取a类型
print(a)