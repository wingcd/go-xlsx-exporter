ValueTypes = {
    ["int"] = "number",
    ["int32"] = "number",
    ["uint"] = "number",
    ["uint32"] = "number",
    ["float"] = "number",
    ["float32"] = "number",
    ["string"] = "string",
    ["bool"] = "boolean",
    ["int[]"] = "number[]",
    ["uint[]"] = "number[]",
    ["float[]"] = "number[]",
    ["float32[]"] = "number[]",
    ["string[]"] = "string[]",
    ["bool[]"] = "boolean[]",
}

function genItem(t)    
    line = "    interface " .. t.TypeName .. "Item {\n"
    for i,v in pairs(t.Headers) do
        structInfo = v.StructInfo
        if nil == ValueTypes[structInfo.RawValueType] then
            print("Error: " .. structInfo.RawValueType)
            goto continue            
        end

        line = line .. "        /** " .. structInfo.Desc .. " */\n"
        line = line .. "        readonly " .. structInfo.FieldName .. ": " .. ValueTypes[structInfo.RawValueType] .. ";\n"
        ::continue::
    end
    return line .. "    }\n\n"
end

function genTable(t)   
    line = "    let " .. t.TypeName .. ": {[key: number]: " .. t.TypeName .. "Item };\n"
    return line
end

function generate() 
    ret = "//由工具自动生成的代码，请勿手动修改！\n"
    ret = ret .. "declare namespace Configs {\n"

    for i,v in pairs(GXE.fileDesc.Tables) do
        if not v.IsArray then
            ret = ret .. genItem(v)
        end
    end

    for i,v in pairs(GXE.fileDesc.Tables) do
        if not v.IsArray then
            ret = ret .. genTable(v)
        end
    end

    ret = ret .. "}\n"
    return ret
end
