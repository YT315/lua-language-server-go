--函数级闭包
local a=100
function change()
    a=300
end
change()
print(a)