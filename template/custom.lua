function genTableData(t)
    ret = ""
    -- print(GXE.json_encode(t.Data))
    for i,v in pairs(t.Data) do
        line = "" .. i
        for _,v2 in pairs(v) do
            line = line .. "," .. string.format("%s", v2)
        end
        ret = ret .. line .. "\n"
    end
    return ret
end

function generate() 
    ret = ""
    -- print(GXE.json_encode(GXE.fileDesc))
    -- print(GXE.fileDesc.Tables)
    for i,v in pairs(GXE.fileDesc.Tables) do
        --  print(v.TypeName)
        --  print(GXE.json_encode(v))
        if not v.IsArray then
            ret = ret .. genTableData(v) .. "\n"
        end
    end
    return ret
end
