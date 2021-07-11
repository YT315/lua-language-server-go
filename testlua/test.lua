a=100

local ff=function ()

goto tag


print("hello")

::tag::
print("world")    

do
    ::tag::
    goto tag
    ::tag::
    print("world")    
end


end
ff()

