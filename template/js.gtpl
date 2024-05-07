// DO NOT EDIT!
// This code is auto generated by go-xlsx-exporter
// VERSION {{.Version}}
// go-protobuf {{.GoProtoVersion}}

{{- $G := .}}
{{- $NS := .Namespace}}

/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

var DataConverter = $root.DataConverter = {
    convertHandler: null,    

    getConvertData: function(target, fieldName, value, alias, cachable) {
        target._converted = target._converted || {};
        if(target._converted[fieldName]) {
            return target._converted[fieldName];
        }

        if(DataConverter.convertHandler == null) {
            throw `convert field ${fieldName} value need a convetor`;
        }

        var data = DataConverter.convertHandler(target, fieldName, value, alias);
        if(cachable) {
            target._converted[fieldName] = data;
        }
        return data;
    },
};

var DataModel = $root.DataModel = function() {
    this._converted = {};

    this.getConvertData = function(fieldName, value, alias, cachable)
    {
        if(this._converted[fieldName])
        {
            this._converted[fieldName];
        }

        if(DataConverter.convertHandler == null)
        {
            throw `convert field ${fieldName} value need a convetor`;
        }

        var data = DataConverter.convertHandler(this, fieldName, value, alias);
        if(cachable) {
            this._converted[fieldName] = data;
        }
        return data;
    };
};

$root.{{$NS}} = (function() {
    /**
     * Namespace {{$NS}}.
     * @exports {{$NS}}
     * @namespace
     */
    var {{$NS}} = {};    

    var ALLTYPES = {{$NS}}.ALLTYPES = {};

    {{/*生成枚举类型*/}}
    {{- range .Enums}}
    // Defined in table: {{.DefinedTable}}
    var {{.TypeName}} = {{$NS}}.{{.TypeName}} = (function() {
        var valuesById = {}, values = Object.create(valuesById);
    {{range .Items}}
        {{if ne .Desc ""}} /** {{.Desc}} */{{end}}
        values[valuesById[{{.Index}}] = "{{.FieldName}}"] = {{.Value}};
    {{end}}
        return values;
    })();
    {{end}}

    {{- /*生成配置类类型*/}}
    {{- range .Consts}}    
    {{$TypeName := .TypeName}}
    // Defined in table: {{.DefinedTable}}
    var {{.TypeName}} = {{$NS}}.{{.TypeName}} = (function() {
        var valuesById = {}, values = Object.create(valuesById);
    {{range .Items}}
        {{if ne .Desc ""}} /** {{.Desc}} */{{end}}
    {{- if not .IsVoid }}
        values.{{.FieldName}} = {{value_format .Value .}};
    {{- end}}    
        {{- if .Convertable}}
        values.get{{upperF .FieldName}} = function() {
            return DataConverter.getConvertData({{$TypeName}}, '{{.FieldName}}', {{$TypeName}}.{{.FieldName}}, '{{get_alias .Alias}}', {{.Cachable}});
        };
        {{- end}}
    {{end}}
        return values;
    })();
    {{end}}

    {{- /*生成类类型*/}}
    {{- range .Tables}}
    {{$TypeName := .TypeName}}
    // Defined in table: {{.DefinedTable}}
    var {{$TypeName}} = {{$NS}}.{{$TypeName}} = (function() {
        this.prototype = Object.create(DataModel.prototype);

        {{$TypeName}}.__type_name__ = "{{$TypeName}}";   
        {{- if not .IsArray}}   
        {{$TypeName}}.__array_type_name__ = "{{$TypeName}}_ARRAY";
        {{$TypeName}}.getArrayType() = function() {
            return ALLTYPES[ALLTYPES.__type_name__];
        }
        {{- end}}

        function {{$TypeName}}(properties) {
            {{range .Headers}}
        {{- if not .IsVoid }}
            {{- if .IsArray}}
            this.{{.FieldName}} = [];
            {{end -}}
        {{end -}}     
        {{end}}        
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }{{/*end function */}}

        {{range .Headers}}
        {{if ne .Desc ""}} /** {{.Desc}} */{{end}}
            {{- if not .IsVoid }}
                {{- if .IsArray}}
        {{$TypeName}}.prototype.{{.FieldName}} =  $util.emptyArray;
                {{- else if eq .StandardValueType "bytes"}}
        {{$TypeName}}.prototype.{{.FieldName}} =  $util.newBuffer([]);
                {{else}}
        {{$TypeName}}.prototype.{{.FieldName}} =  {{default .}};
                {{end -}}   
            {{end -}} 
            {{- if .Convertable}}
        {{$TypeName}}.prototype.get{{upperF .FieldName}} = function() {
            return this.getConvertData("{{.FieldName}}", {{if .IsVoid}}null{{else}}this.{{.FieldName}}{{end}}, '{{get_alias .Alias}}', {{.Cachable}});
        };
            {{- end}}
        {{end}} 

        {{$TypeName}}.create = function create(properties) {
            return new {{$TypeName}}(properties);
        };

        {{$TypeName}}.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
                
            {{range .Headers}}        
                {{- $wireType := get_wire_type .}}
                {{- $count := calc_wire_offset .}}           
                
            {{if ne .Desc ""}} /** {{.Desc}} */ {{end}}
                {{- if not .IsVoid }}
                    {{- if .IsArray}}
                            {{- if .IsMessage}}
            if (message.{{.FieldName}} != null && message.{{.FieldName}}.length)
                for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                    {{.ValueType}}.encode(message.{{.FieldName}}[i], writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).fork()).ldelim();
                            {{- else}}
            if (message.{{.FieldName}} != null && message.{{.FieldName}}.length) {
                {{- if out .StandardValueType "string" "bytes"}}
                writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).fork();
                {{- end}}
                for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                    {{- if eq "string" .StandardValueType}}
                    writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).string(message.{{.FieldName}}[i]);
                    {{- else if eq "bytes" .StandardValueType}}
                    writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).bytes(message.{{.FieldName}}[i]);
                    {{- else}}
                    {{- $pbType := get_pb_type .StandardValueType}}
                    writer.{{$pbType}}(message.{{.FieldName}}[i]);
                    {{- end}}
                {{- if out .StandardValueType "string" "bytes"}}
                writer.ldelim();
                {{- end}}
            }
                            {{- end}}{{/*end message*/}}
                    {{- else}} {{/*not array */}}
            if (message.{{.FieldName}} != null && Object.hasOwnProperty.call(message, "{{.FieldName}}"))
                            {{- if .IsMessage}} 
                {{.ValueType}}.encode(message.{{.FieldName}}[i], writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).fork()).ldelim();
                            {{- else}}
                    {{- $pbType := get_pb_type .StandardValueType}}
                writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).{{$pbType}}(message.{{.FieldName}});
                            {{- end}}{{/*end message*/}}
                    {{- end -}} {{/*end if*/}}   
                {{- end}} 
            {{end -}} 

            return writer;
        };

        {{$TypeName}}.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        {{$TypeName}}.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new {{$TypeName}}();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                {{- range .Headers}} 
                {{- if not .IsVoid}}          
                case {{.Index}}:
                    {{- if .IsArray}}                    
                    if (!(message.{{.FieldName}} && message.{{.FieldName}}.length))
                        message.{{.FieldName}} = [];

                        {{- if .IsMessage}}
                    message.{{.FieldName}}.push({{.ValueType}}.decode(reader, reader.uint32()));                    
                        {{- else}}
                            {{- if out .ValueType "string" "bytes"}}
                    if ((tag & 7) === 2) {
                        var end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2)
                            message.{{.FieldName}}.push(reader.{{.PBValueType}}());
                    } else
                        message.{{.FieldName}}.push(reader.{{.PBValueType}}());
                            {{- else}}
                    message.{{.FieldName}}.push(reader.{{.PBValueType}}());
                            {{- end}} {{/*end if string */}}
                        {{- end}} {{/*end if message*/}}
                    {{- else}}
                        {{- if .IsMessage}}
                    {{.ValueType}}.decode(reader, reader.uint32());
                        {{- else}}
                    message.{{.FieldName}} = reader.{{.PBValueType}}();
                        {{- end}} {{/*end message*/}}
                    {{- end}} 
                    break;
                {{- end}} {{/**end .IsVoid */}}
                {{- end}} {{/**end case */}}
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        {{$TypeName}}.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        {{$TypeName}}.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            {{- range .Headers}}      
            {{- if not .IsVoid }}
                if (message.{{.FieldName}} != null && message.hasOwnProperty("{{.FieldName}}")) {
                {{- if .IsArray}}
                    if (!Array.isArray(message.{{.FieldName}}))
                        return "{{.FieldName}}: array expected";
                    {{- if .IsEnum}}                
                    {{- $enums := get_enum_values .StandardValueType}}
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                        switch (message.{{.FieldName}}[i]) {
                            default:
                                return "{{.FieldName}}: enum value[] expected";
                            {{- range $enums}}
                            case {{.}}:
                            {{- end}}
                                break;
                        }
                    {{- else if is_interger .StandardValueType}}
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                        if (!$util.isInteger(message.{{.FieldName}}[i]))
                            return "{{.FieldName}}: integer[] expected";
                    {{- else if is_float .StandardValueType}}
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                        if (typeof message.{{.FieldName}}[i] !== "number")
                            return "{{.FieldName}}: number[] expected";
                    {{- else if is_bool .StandardValueType}}
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                        if (typeof message.{{.FieldName}}[i] !== "boolean")
                            return "{{.FieldName}}: boolean[] expected";
                    {{- else if is_string .StandardValueType}}
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                        if (!$util.isString(message.{{.FieldName}}[i]))
                            return "{{.FieldName}}: string[] expected";
                    {{- else if is_bytes .StandardValueType}}
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                        if (!(message.{{.FieldName}}[i] && typeof message.{{.FieldName}}[i].length === "number" || $util.isString(message.{{.FieldName}}[i])))
                            return "{{.FieldName}}: buffer[] expected";
                    {{- else if .IsMessage}}     
                    for (var i = 0; i < message.{{.FieldName}}.length; ++i) {
                        var error = {{.ValueType}}.verify(message.{{.FieldName}}[i]);
                        if (error)
                            return "{{.FieldName}}." + error;
                    }       
                    {{- else}}
                        "error type {{.ValueType}} {{.FieldName}}";
                    {{- end}}
                {{- else}}
                    {{- if .IsEnum}}
                    {{- $enums := get_enum_values .StandardValueType}}
                    switch (message.{{.FieldName}}) {
                        default:
                            return "{{.FieldName}}: enum value expected";
                        {{- range $enums}}
                        case {{.}}:
                        {{- end}}
                            break;
                    }
                    {{- else if is_interger .StandardValueType}}
                    if (!$util.isInteger(message.{{.FieldName}}))
                        return "{{.FieldName}}: integer expected";
                    {{- else if is_float .StandardValueType}}
                    if (typeof message.{{.FieldName}} !== "number")
                        return "{{.FieldName}}: number expected";
                    {{- else if is_bool .StandardValueType}}
                    if (typeof message.{{.FieldName}} !== "boolean")
                        return "{{.FieldName}}: boolean expected";
                    {{- else if is_string .StandardValueType}}
                    if (!$util.isString(message.{{.FieldName}}))
                        return "{{.FieldName}}: string expected";
                    {{- else if is_bytes .StandardValueType}} 
                    if (!(message.{{.FieldName}} && typeof message.{{.FieldName}}.length === "number" || $util.isString(message.{{.FieldName}})))
                        return "{{.FieldName}}: buffer expected";
                    {{- else if .IsMessage}}
                        return {{.ValueType}}.verify(message.{{.FieldName}});
                    {{- else}}
                        "error type {{.ValueType}} {{.FieldName}}";
                    {{- end}}
                {{- end}} {{/* end if IsArray */}}
                }
            {{- end}} {{/* end if not IsVoid */}}
            {{- end}} {{/**end range Headers */}}
            return null;
        };

        {{$TypeName}}.fromObject = function fromObject(object) {
            if (object instanceof {{$TypeName}})
                return object;
            var message = new {{$TypeName}}();
            {{- range .Headers}}  
            {{- if not .IsVoid }}
                {{- if .IsArray}}
            if(object.{{.FieldName}}) {
                if (!Array.isArray(object.{{.FieldName}}))
                    throw TypeError("{{$TypeName}}.{{.FieldName}}: array expected");            
                message.{{.FieldName}} = [];
                for (var i = 0; i < object.{{.FieldName}}.length; ++i) {
                    {{- if .IsEnum}} 
                        {{- $Ofn := .FieldName}}
                        {{- $enum := get_enum .StandardValueType}}
                    switch (object.{{.FieldName}}[i]) {
                        default:
                            {{- range $enum.Items}}
                        case "{{.FieldName}}":
                        case {{.Value}}:
                            message.{{$Ofn}}[i] = {{.Value}};
                            break;
                        {{- end}}    
                    }                
                    {{- else if is_long .StandardValueType}}
                    if ($util.Long)
                        (message.{{.FieldName}}[i] = $util.Long.fromValue(object.{{.FieldName}})).unsigned = {{if eq .StandardValueType "uint64"}}true{{else}}false{{end}};
                    else if (typeof object.{{.FieldName}} === "string")
                        message.{{.FieldName}}[i] = parseInt(object.{{.FieldName}}, 10);
                    else if (typeof object.{{.FieldName}} === "number")
                        message.{{.FieldName}}[i] = object.{{.FieldName}};
                    else if (typeof object.{{.FieldName}} === "object")
                        message.{{.FieldName}}[i] = new $util.LongBits(object.{{.FieldName}}[i].low >>> 0, object.{{.FieldName}}[i].high >>> 0).toNumber();
                    {{- else if eq .StandardValueType "int"}}
                    message.{{.FieldName}}[i] = object.{{.FieldName}}[i] | 0;
                    {{- else if eq .StandardValueType "uint"}}
                    message.{{.FieldName}}[i] = object.{{.FieldName}}[i] >>> 0;
                    {{- else if eq .StandardValueType "bool"}}
                    message.{{.FieldName}}[i] = Boolean(object.{{.FieldName}}[i]);
                    {{- else if is_float .StandardValueType}}
                    message.{{.FieldName}}[i] = Number(object.{{.FieldName}}[i]);                
                    {{- else if eq .StandardValueType "string"}}
                    message.{{.FieldName}}[i] = String(object.{{.FieldName}}[i]);
                    {{- else if eq .StandardValueType "bytes"}}
                    if (typeof object.{{.FieldName}}[i] === "string")
                        $util.base64.decode(object.{{.FieldName}}[i], message.{{.FieldName}}[i] = $util.newBuffer($util.base64.length(object.{{.FieldName}}[i])), 0);
                    else if (object.{{.FieldName}}[i].length)
                        message.{{.FieldName}}[i] = object.{{.FieldName}}[i];
                    {{- else if .IsMessage}}
                    if (typeof object.{{.FieldName}}[i] !== "object")
                        throw TypeError("{{$TypeName}}.{{.FieldName}}: object expected");
                    message.{{.FieldName}}[i] = {{$TypeName}}.fromObject(object.{{.FieldName}}[i]);
                    {{- else}}
                        throw TypeError("{{$TypeName}}.{{.FieldName}}: object expected");
                    {{- end}}
                }
            }
                {{- else}}{{/**end if Array */}}
                    {{- if not .IsEnum}}                
            if (object.{{.FieldName}} != null)
                    {{- end}}
                    {{- if .IsEnum}} 
                    {{- $Ofn := .FieldName}}
                    {{- $enum := get_enum .StandardValueType}}
            switch (object.{{.FieldName}}) {
                default:
                    {{- range $enum.Items}}
                case "{{.FieldName}}":
                case {{.Value}}:
                    message.{{$Ofn}} = {{.Value}};
                    break;
                    {{- end}}    
                }               
                    {{- else if is_long .StandardValueType}}
                if ($util.Long)
                    (message.{{.FieldName}} = $util.Long.fromValue(object.{{.FieldName}})).unsigned = {{if eq .StandardValueType "uint64"}}true{{else}}false{{end}};
                else if (typeof object.{{.FieldName}} === "string")
                    message.{{.FieldName}} = parseInt(object.{{.FieldName}}, 10);
                else if (typeof object.{{.FieldName}} === "number")
                    message.{{.FieldName}} = object.{{.FieldName}};
                else if (typeof object.{{.FieldName}} === "object")
                    message.{{.FieldName}} = new $util.LongBits(object.{{.FieldName}}.low >>> 0, object.{{.FieldName}}.high >>> 0).toNumber();
                    {{- else if eq .StandardValueType "int"}}
                message.{{.FieldName}} = object.{{.FieldName}} | 0;
                    {{- else if eq .StandardValueType "uint"}}
                message.{{.FieldName}} = object.{{.FieldName}} >>> 0;
                    {{- else if eq .StandardValueType "bool"}}
                message.{{.FieldName}} = Boolean(object.{{.FieldName}});
                    {{- else if is_float .StandardValueType}}
                message.{{.FieldName}} = Number(object.{{.FieldName}});
                    {{- else if eq .StandardValueType "string"}}
                message.{{.FieldName}} = String(object.{{.FieldName}});
                    {{- else if eq .StandardValueType "bytes"}}
                if (typeof object.{{.FieldName}} === "string")
                    $util.base64.decode(object.{{.FieldName}}, message.{{.FieldName}} = $util.newBuffer($util.base64.length(object.{{.FieldName}})), 0);
                else if (object.{{.FieldName}}.length)
                    message.{{.FieldName}} = object.{{.FieldName}};
                    {{- else if .IsMessage}}
                if (typeof object.{{.FieldName}}[i] !== "object")
                    throw TypeError("{{$TypeName}}.{{.FieldName}}: object expected");
                message.{{.FieldName}} = {{$TypeName}}.fromObject(object.{{.FieldName}});
                    {{- end}}
                {{- end}}
            {{- end}} {{/* end if not IsVoid */}}
            {{- end}} {{/* end range Headers */}}
            return message;
        };

        {{$TypeName}}.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults) {
            {{- range .Headers}}  
            {{- if not .IsVoid }}
                {{- if .IsArray}}
                object.{{.FieldName}} = [];
                {{- end}}
            {{- end}} {{/* end if not IsVoid */}}
            {{- end}}
            }

            if (options.defaults) {
            {{- range .Headers}}
            {{- if not .IsVoid }}
                {{- if is_long .StandardValueType}}
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.{{.FieldName}} = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.{{.FieldName}} = options.longs === String ? "0" : 0;
                {{- else if .IsEnum}}
                {{$enumDefault := get_enum_default .StandardValueType}}
                object.{{.FieldName}} = options.enums === String ? "{{$enumDefault.FieldName}}" : {{$enumDefault.Value}};
                {{- else}}
                object.{{.FieldName}} = {{default .}};
                {{- end}}
            {{- end}} {{/* end if not IsVoid */}}
            {{- end}} {{/* end range Headers */}}
            }

            {{- range .Headers}}  
            {{- if not .IsVoid }}
                {{- if .IsArray}}
                
            if (message.{{.FieldName}} && message.{{.FieldName}}.length) {
                object.{{.FieldName}} = [];
                for (var j = 0; j < message.{{.FieldName}}.length; ++j)
                    {{- if is_float .StandardValueType}}
                    object.{{.FieldName}}[j] = options.json && !isFinite(message.{{.FieldName}}[j]) ? String(message.{{.FieldName}}[j]) : message.{{.FieldName}}[j];
                    {{- else if is_long .StandardValueType}}
                    if (typeof message.{{.FieldName}}[j] === "number")
                        object.{{.FieldName}}[j] = options.longs === String ? String(message.{{.FieldName}}[j]) : message.{{.FieldName}}[j];
                    else
                        object.{{.FieldName}}[j] = options.longs === String ? $util.Long.prototype.toString.call(message.{{.FieldName}}[j]) : options.longs === Number ? new $util.LongBits(message.{{.FieldName}}[j].low >>> 0, message.{{.FieldName}}[j].high >>> 0).toNumber() : message.{{.FieldName}}[j];
                    {{- else if .IsEnum}}
                    object.{{.FieldName}}[j] = options.enums === String ? {{.ValueType}}[message.{{.FieldName}}[j]] : message.{{.FieldName}}[j];
                    {{- else if eq .StandardValueType "bytes"}}
                    object.{{.FieldName}}[j] = options.bytes === String ? $util.base64.encode(message.{{.FieldName}}[j], 0, message.{{.FieldName}}[j].length) : options.bytes === Array ? Array.prototype.slice.call(message.{{.FieldName}}[j]) : message.{{.FieldName}}[j];
                    {{- else if .IsMessage}}
                    object.{{.FieldName}}[j] = {{.ValueType}}.toObject(message.{{.FieldName}}[j], options);
                    {{- else}}
                    object.{{.FieldName}}[j] = message.{{.FieldName}}[j];
                    {{- end}}
            }
                {{- else}}
            if (message.{{.FieldName}} != null && message.hasOwnProperty("{{.FieldName}}"))
                    {{- if is_float .StandardValueType}}
                object.{{.FieldName}} = options.json && !isFinite(message.{{.FieldName}}) ? String(message.{{.FieldName}}) : message.{{.FieldName}};
                    {{- else if is_long .StandardValueType}}
                if (typeof message.{{.FieldName}} === "number")
                    object.{{.FieldName}} = options.longs === String ? String(message.{{.FieldName}}) : message.{{.FieldName}};
                else
                    object.{{.FieldName}} = options.longs === String ? $util.Long.prototype.toString.call(message.{{.FieldName}}) : options.longs === Number ? new $util.LongBits(message.{{.FieldName}}.low >>> 0, message.{{.FieldName}}.high >>> 0).toNumber() : message.{{.FieldName}};
                    {{- else if .IsEnum}}
                    object.{{.FieldName}} = options.enums === String ? {{.ValueType}}[message.{{.FieldName}}] : message.{{.FieldName}};
                    {{- else if eq .StandardValueType "bytes"}}
                object.{{.FieldName}} = options.bytes === String ? $util.base64.encode(message.{{.FieldName}}, 0, message.{{.FieldName}}.length) : options.bytes === Array ? Array.prototype.slice.call(message.{{.FieldName}}) : message.{{.FieldName}};
                    {{- else if .IsMessage}}
                object.{{.FieldName}} = {{.ValueType}}.toObject(message.{{.FieldName}}, options);
                    {{- else}}
                object.{{.FieldName}} = message.{{.FieldName}};
                    {{- end}}
                {{- end}}
            {{- end}} {{/* end if not IsVoid */}}
            {{- end}}{{/*end range headers */}}
            return object;
        };{{/*end toObject function*/}}

        {{$TypeName}}.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return {{$TypeName}};
    })(); {{/*end class */}}
    ALLTYPES["{{$TypeName}}"] = {{$TypeName}};
    
        {{- end}} {{/*end tables */}}

    return {{$NS}};
})(); {{/*end all */}}

module.exports = $root;