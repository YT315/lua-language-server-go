function a10()
    return 10,true,"string"
end
--当函数后面有表达式时,函数返回值只取第一个
local a10a,a10b,a10c="string",a10()
print(a10a,a10b,a10c)