mytable = {}                          -- 普通表
mymetatable = {} 
mymetatable.aaa=100
mymetatable.__index=mymetatable
nn=setmetatable(mytable,mymetatable)     -- 把 mymetatable 设为 mytable 的元表 
bb=mytable
mymetatable=500
print(bb,nn,mytable,mymetatable)
print(bb.aaa)