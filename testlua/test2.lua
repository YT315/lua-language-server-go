--闭包
function getfunc()
    local a=200
    return function()
        a=a+100
        --此处获取a类型
        print(a)
        return a
    end
end

getfunc()()
func = getfunc()
return func()

