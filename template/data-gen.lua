ValueTypes = {
    ["int"] = "i",
    ["int32"] = "i",
    ["uint"] = "i",
    ["uint32"] = "i",
    ["float"] = "f",
    ["float32"] = "f",
    ["string"] = "s",
    ["bool"] = "b",
    ["int[]"] = "ia",
    ["uint[]"] = "ia",
    ["float[]"] = "fa",
    ["float32[]"] = "fa",
    ["string[]"] = "sa",
    ["bool[]"] = "ba",
}

function genTableDesc(t)    
    line = t.TypeName
    for i,v in pairs(t.Headers) do
        structInfo = v.StructInfo
        if nil == ValueTypes[structInfo.RawValueType] then
            print("Error: " .. structInfo.RawValueType)
            goto continue            
        end

        line = line .. ":" .. structInfo.FieldName .. "-" .. ValueTypes[structInfo.RawValueType]
        ::continue::
    end
    return line .. "\n"
end

function genTableData(t)
    ret = "##########\n"
    ret = ret .. t.TypeName .. "\n"

    if t.TypeName == "novice" then
        -- print(GXE.json_encode(t.Data))
    end

    for i,v in pairs(t.Data) do
        line = ""
        idx = 0
        for _,v2 in pairs(v) do
            if idx == 0 then
                line = string.format("%s", v2)
            else
                line = line .. "," .. string.format("%s", v2)
            end
            idx = idx + 1
        end
        ret = ret .. line .. "\n"
    end
    return ret
end

function generate() 
    ret = ""
    for i,v in pairs(GXE.fileDesc.Tables) do
        if not v.IsArray then
            ret = ret .. genTableDesc(v)
        end
    end

    for i,v in pairs(GXE.fileDesc.Tables) do
        if not v.IsArray then
            ret = ret .. genTableData(v)
        end
    end
    return ret
end
