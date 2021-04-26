local ta={100,200,300,400}
local tb={true,"ffff","ggg",["hell"]="world"}
for i,v in next,tb,nil do
    print(i,v)
end
-- k,v=next(tb)
-- print(k,v)
-- k,v=next(tb,k)
-- print(k,v)
-- k,v=next(tb,k)
-- print(k,v)
-- k,v=next(tb,k)
-- print(k,v)