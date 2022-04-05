// DO NOT EDIT!
// This code is auto generated by go-xlsx-exporter
// VERSION {{.Version}}
// go-protobuf {{.GoProtoVersion}}

{{- $G := .}}
{{- $NS := .Namespace}}

{{- range .Info.Imports}}
{{.}}
{{- end}}

import $protobuf from "protobufjs";

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {} as any);

export interface Long {
    /** Low bits */
    low: number;

    /** High bits */
    high: number;

    /** Whether unsigned or not */
    unsigned: boolean;
}

export interface Long {
    /** Low bits */
    low: number;

    /** High bits */
    high: number;

    /** Whether unsigned or not */
    unsigned: boolean;
}

export class DataConverter {
    static convertHandler: (typeName: string, fieldName:string, value: string)=>any = null;

    static convertData(typeName: string, fieldName:string, value: string): any {
        if(this.convertHandler) {
            return this.convertHandler(typeName, fieldName, value);
        }
        return null;
    }
}

export namespace {{$NS}} {
    {{/*生成枚举类型*/}}
    {{- range .Enums}}
    // Defined in table: {{.DefinedTable}}
    export enum {{.TypeName}} {
        {{- range .Items}}
        {{- if ne .Desc ""}} //{{.Desc}} {{end}}
        {{.FieldName}} = {{.Value}},
        {{end}}
    }
    {{end}}{{/* enum */}}

    {{- /*生成配置类类型*/}}
    {{- range .Consts}}
    {{$TypeName := .TypeName}}

    // Defined in table: {{.DefinedTable}}
    export var {{.TypeName}}: {
        {{- range .Items}}
            {{- if not .IsVoid }}   
                {{- if ne .Desc ""}} //{{.Desc}} {{end}}                    
        {{.FieldName}}?: {{type_format .StandardValueType .ValueType .IsArray}},
            {{end}}
            {{- if .Convertable}}    
        get{{.FieldName}}(): {{get_alias .Alias}},
            {{end}}
        {{end}} {{/*end .Items */}}
    } = {
        {{- range .Items}}
            {{- if not .IsVoid }}
        {{.FieldName}} : {{value_format .Value .}},
            {{end}}
            {{- if .Convertable}}    
        get{{.FieldName}}(): {{get_alias .Alias}} {
            return DataConverter.convertData("{{$TypeName}}", "{{.FieldName}}", this.{{.FieldName}});
        },
            {{end}}
        {{end}} {{/*end .Items */}}
    }
    {{end}}{{/*end .Consts */}}

    {{- /*生成类类型*/}}
    {{- range .Tables}}

    {{$TypeName := .TypeName}}
    // Defined in table: {{.DefinedTable}}
    /** Properties of a {{$TypeName}}. */
    export interface I{{$TypeName}} {
        {{range .Headers}}
            {{- if not .IsVoid }}               
        {{.FieldName}}?: {{type_format .StandardValueType .ValueType .IsArray}};
            {{end}} {{/*end not Void*/}}
            {{- if .Convertable}}
        get{{.FieldName}}(): {{get_alias .Alias}};
            {{- end}}
        {{end}} {{/*end .Headers */}}
    }

     /** Represents a {{$TypeName}}. */
    export class {{$TypeName}} implements I{{$TypeName}} { 
        private static __type_name__ = "{{$TypeName}}";

        {{range .Headers}}
            {{- if not .IsVoid }}
                {{- if .IsArray}}
                    {{- if ne .Desc ""}} //{{.Desc}} {{end}}
        {{.FieldName}} =  $util.emptyArray;
                {{else}}
        {{.FieldName}}?: {{type_format .StandardValueType .ValueType .IsArray}} = {{default .}};
                {{end -}}   
            {{end -}} 
            {{- if .Convertable}}    
        get{{.FieldName}}(): {{get_alias .Alias}} {
            return DataConverter.convertData("{{$TypeName}}", "{{.FieldName}}", this.{{.FieldName}});
        };
            {{end}}
        {{end}}

        /**
         * Constructs a new {{$TypeName}}.
         * @param [properties] Properties to set
         */
        constructor(properties?: I{{$TypeName}}) {
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

        static create(properties?: {{$TypeName}}): {{$TypeName}} {
            return new {{$TypeName}}(properties);
        }

        static encode(message: I{{$TypeName}}, writer?: $protobuf.Writer): $protobuf.Writer {
            if (!writer)
                writer = $Writer.create();
                
            {{range .Headers}}        
                {{- $wireType := get_wire_type .}}
                {{- $count := calc_wire_offset .}}           
                
                {{- if not .IsVoid }}
                    {{- if .IsArray}}
                        {{- if ne .Desc ""}} //{{.Desc}} {{end}}
                            {{- if .IsMessage}}
            if (message.{{.FieldName}} != null && message.{{.FieldName}}.length)
                for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                    {{.ValueType}}.encode(message.{{.FieldName}}[i], writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).fork()).ldelim();
                            {{- else}}
            if (message.{{.FieldName}} != null && message.{{.FieldName}}.length) {
                writer.uint32(/* id {{.Index}}, wireType {{$wireType}} =*/{{$count}}).fork();
                for (var i = 0; i < message.{{.FieldName}}.length; ++i)
                    writer.uint32(message.{{.FieldName}}[i]);
                writer.ldelim();
            }
                            {{- end}}{{/*end message*/}}
                    {{- else}}
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
        }

        static encodeDelimited(message: I{{$TypeName}}, writer?: $protobuf.Writer): $protobuf.Writer {
            return this.encode(message, writer).ldelim();
        }

        static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): {{$TypeName}} {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new {{$TypeName}}();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                {{- range .Headers}}          
                case {{.Index}}:
                    {{- if .IsArray}}                    
                    if (!(message.{{.FieldName}} && message.{{.FieldName}}.length))
                        message.{{.FieldName}} = [];

                        {{- if .IsMessage}}
                    message.{{.FieldName}}.push({{.ValueType}}.decode(reader, reader.uint32()));                    
                        {{- else}}
                    if ((tag & 7) === 2) {
                        var end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2)
                            message.{{.FieldName}}.push(reader.{{.PBValueType}}());
                    } else
                        message.{{.FieldName}}.push(reader.{{.PBValueType}}());
                        {{- end}}
                    {{- else}}
                        {{- if .IsMessage}}
                    {{.ValueType}}.decode(reader, reader.uint32());
                        {{- else}}
                    message.{{.FieldName}} = reader.{{.PBValueType}}();
                        {{- end}} {{/*end message*/}}
                    {{- end}} 
                    break;
                {{- end}}
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        }

        static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): {{$TypeName}} {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        }

        static verify(message: { [k: string]: any }): (string|null) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            {{- range .Headers}}        
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
                    {{- else if .IsMessage}}
                        return {{.ValueType}}.verify(message.{{.FieldName}});
                    {{- else}}
                        "error type {{.ValueType}} {{.FieldName}}";
                    {{- end}}
                {{- end}} {{/* end if IsArray */}}
                }
            {{- end}} {{/**end range Headers */}}
            return null;
        }

        static fromObject(object: { [k: string]: any }): {{$TypeName}} {
            if (object instanceof {{$TypeName}})
                return object;
            var message = new {{$TypeName}}();
            {{- range .Headers}}  
                {{- if .IsArray}}
            if(object.{{.FieldName}}) {
                if (!Array.isArray(object.{{.FieldName}}))
                    throw TypeError("{{$TypeName}}.{{.FieldName}}: array expected");            
                message.{{.FieldName}} = [];
                for (var i = 0; i < object.{{.FieldName}}.length; ++i)
                    {{- $Ofn := .FieldName}}
                    {{- if .IsEnum}} 
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
                    {{- else if .IsMessage}}
                    if (typeof object.{{.FieldName}}[i] !== "object")
                        throw TypeError("{{$TypeName}}.{{.FieldName}}: object expected");
                    message.{{.FieldName}}[i] = {{$TypeName}}.fromObject(object.{{.FieldName}}[i]);
                    {{- end}}
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
                    {{- else if .IsMessage}}
                if (typeof object.{{.FieldName}}[i] !== "object")
                    throw TypeError("{{$TypeName}}.{{.FieldName}}: object expected");
                message.{{.FieldName}} = {{$TypeName}}.fromObject(object.{{.FieldName}});
                    {{- end}}
                {{- end}}
            {{- end}}
            return message;
        }

        static toObject(message: {{$TypeName}}, options?: $protobuf.IConversionOptions): { [k: string]: any } {
            if (!options)
                options = {};
            var object: any = {};
            if (options.arrays || options.defaults) {
            {{- range .Headers}}  
                {{- if .IsArray}}
                object.{{.FieldName}} = [];
                {{- end}}
            {{- end}}
            }

            if (options.defaults) {
            {{- range .Headers}}  
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
            {{- end}}
            }

            {{- range .Headers}}  
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
                    {{- else if .IsMessage}}
                object.{{.FieldName}} = {{.ValueType}}.toObject(message.{{.FieldName}}, options);
                    {{- else}}
                object.{{.FieldName}} = message.{{.FieldName}};
                    {{- end}}
                {{- end}}
            {{- end}}{{/*end range headers */}}
            return object;
        }{{/*end toObject function*/}}

        toJSON(): { [k: string]: any } {
            return {{$TypeName}}.toObject(this, $protobuf.util.toJSONOptions);
        }
    } {{/*end class */}}
        {{- end}} {{/*end tables */}}
}{{/*end namespace */}}