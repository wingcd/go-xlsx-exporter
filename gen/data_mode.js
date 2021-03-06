// DO NOT EDIT!
// This code is auto generated by go-xlsx-exporter
// VERSION 1.2
// go-protobuf v1.27.1

/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

var DataConverter = $root.DataConverter = {
    convertHandler: null,
};

var DataModel = $root.DataModel = function() {
    this._converted = {};

    this.getConvertData = function(fieldName, value)
    {
        if(this._converted[fieldName])
        {
            this._converted[fieldName];
        }

        if(DataConverter.convertHandler == null)
        {
            throw `convert field ${fieldName} value need a convetor`;
        }

        var data = DataConverter.convertHandler(this, fieldName, value);
        this._converted[fieldName] = data;
        return data;
    } 
};

$root.GameData = (function() {
    /**
     * Namespace GameData.
     * @exports GameData
     * @namespace
     */
    var GameData = {};    

    var ALLTYPES = GameData.ALLTYPES = {};

    
    // Defined in table: 
    var EMsgType = GameData.EMsgType = (function() {
        var valuesById = {}, values = Object.create(valuesById);
    
        
        values[valuesById[0] = "UNKNOWN"] = 0;
    
        
        values[valuesById[1] = "JSON"] = 1;
    
        
        values[valuesById[2] = "XML"] = 2;
    
        return values;
    })();
    
    
    // Defined in table: data/message.xml
    var MessageWrapper = GameData.MessageWrapper = (function() {
        this.prototype = Object.create(DataModel.prototype);

        MessageWrapper.__type_name__ = "MessageWrapper";   
        MessageWrapper.__array_type_name__ = "MessageWrapper_ARRAY";
        MessageWrapper.getArrayType() = function() {
            return ALLTYPES[ALLTYPES.__type_name__];
        }

        function MessageWrapper(properties) {
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        
        MessageWrapper.prototype.id =  0;
                
        
        MessageWrapper.prototype.data =  $util.newBuffer([]);
                
         

        MessageWrapper.create = function create(properties) {
            return new MessageWrapper(properties);
        };

        MessageWrapper.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
                
             
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.id); 
             
            if (message.data != null && Object.hasOwnProperty.call(message, "data"))
                writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.data); 
            return writer;
        };

        MessageWrapper.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        MessageWrapper.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new MessageWrapper();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {          
                case 1:
                    message.id = reader.int32();  
                    break;           
                case 2:
                    message.data = reader.bytes();  
                    break;  
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        MessageWrapper.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        MessageWrapper.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
                if (message.id != null && message.hasOwnProperty("id")) {
                    if (!$util.isInteger(message.id))
                        return "id: integer expected"; 
                } 
                if (message.data != null && message.hasOwnProperty("data")) { 
                    if (!(message.data && typeof message.data.length === "number" || $util.isString(message.data)))
                        return "data: buffer expected"; 
                }  
            return null;
        };

        MessageWrapper.fromObject = function fromObject(object) {
            if (object instanceof MessageWrapper)
                return object;
            var message = new MessageWrapper();                
            if (object.id != null)
                message.id = object.id | 0;                 
            if (object.data != null)
                if (typeof object.data === "string")
                    $util.base64.decode(object.data, message.data = $util.newBuffer($util.base64.length(object.data)), 0);
                else if (object.data.length)
                    message.data = object.data;  
            return message;
        };

        MessageWrapper.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults) {  
            }

            if (options.defaults) {
                object.id = 0; 
                object.data = new Uint8Array(0);  
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id; 
            if (message.data != null && message.hasOwnProperty("data"))
                object.data = options.bytes === String ? $util.base64.encode(message.data, 0, message.data.length) : options.bytes === Array ? Array.prototype.slice.call(message.data) : message.data; 
            return object;
        };

        MessageWrapper.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return MessageWrapper;
    })(); 
    ALLTYPES["MessageWrapper"] = MessageWrapper;
    
    // Defined in table: data/message.xml
    var Item = GameData.Item = (function() {
        this.prototype = Object.create(DataModel.prototype);

        Item.__type_name__ = "Item";   
        Item.__array_type_name__ = "Item_ARRAY";
        Item.getArrayType() = function() {
            return ALLTYPES[ALLTYPES.__type_name__];
        }

        function Item(properties) {
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        
        Item.prototype.id =  0;
                
        
        Item.prototype.name =  "";
                
         

        Item.create = function create(properties) {
            return new Item(properties);
        };

        Item.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
                
             
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.id); 
             
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.name); 
            return writer;
        };

        Item.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        Item.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new Item();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {          
                case 1:
                    message.id = reader.int32();  
                    break;           
                case 2:
                    message.name = reader.string();  
                    break;  
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        Item.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        Item.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
                if (message.id != null && message.hasOwnProperty("id")) {
                    if (!$util.isInteger(message.id))
                        return "id: integer expected"; 
                } 
                if (message.name != null && message.hasOwnProperty("name")) {
                    if (!$util.isString(message.name))
                        return "name: string expected"; 
                }  
            return null;
        };

        Item.fromObject = function fromObject(object) {
            if (object instanceof Item)
                return object;
            var message = new Item();                
            if (object.id != null)
                message.id = object.id | 0;                 
            if (object.name != null)
                message.name = String(object.name);  
            return message;
        };

        Item.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults) {  
            }

            if (options.defaults) {
                object.id = 0; 
                object.name = "";  
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id; 
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name; 
            return object;
        };

        Item.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return Item;
    })(); 
    ALLTYPES["Item"] = Item;
    
    // Defined in table: data/message.xml
    var C2S_GetPlayerInfo = GameData.C2S_GetPlayerInfo = (function() {
        this.prototype = Object.create(DataModel.prototype);

        C2S_GetPlayerInfo.__type_name__ = "C2S_GetPlayerInfo";   
        C2S_GetPlayerInfo.__array_type_name__ = "C2S_GetPlayerInfo_ARRAY";
        C2S_GetPlayerInfo.getArrayType() = function() {
            return ALLTYPES[ALLTYPES.__type_name__];
        }

        function C2S_GetPlayerInfo(properties) {
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        
        C2S_GetPlayerInfo.prototype.name =  "";
                
         

        C2S_GetPlayerInfo.create = function create(properties) {
            return new C2S_GetPlayerInfo(properties);
        };

        C2S_GetPlayerInfo.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
                
             
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.name); 
            return writer;
        };

        C2S_GetPlayerInfo.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        C2S_GetPlayerInfo.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new C2S_GetPlayerInfo();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {          
                case 1:
                    message.name = reader.string();  
                    break;  
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        C2S_GetPlayerInfo.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        C2S_GetPlayerInfo.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
                if (message.name != null && message.hasOwnProperty("name")) {
                    if (!$util.isString(message.name))
                        return "name: string expected"; 
                }  
            return null;
        };

        C2S_GetPlayerInfo.fromObject = function fromObject(object) {
            if (object instanceof C2S_GetPlayerInfo)
                return object;
            var message = new C2S_GetPlayerInfo();                
            if (object.name != null)
                message.name = String(object.name);  
            return message;
        };

        C2S_GetPlayerInfo.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults) { 
            }

            if (options.defaults) {
                object.name = "";  
            }
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name; 
            return object;
        };

        C2S_GetPlayerInfo.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C2S_GetPlayerInfo;
    })(); 
    ALLTYPES["C2S_GetPlayerInfo"] = C2S_GetPlayerInfo;
    
    // Defined in table: data/message.xml
    var S2C_GetPlayerInfo = GameData.S2C_GetPlayerInfo = (function() {
        this.prototype = Object.create(DataModel.prototype);

        S2C_GetPlayerInfo.__type_name__ = "S2C_GetPlayerInfo";   
        S2C_GetPlayerInfo.__array_type_name__ = "S2C_GetPlayerInfo_ARRAY";
        S2C_GetPlayerInfo.getArrayType() = function() {
            return ALLTYPES[ALLTYPES.__type_name__];
        }

        function S2C_GetPlayerInfo(properties) {
            
            this.items = [];
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        
        S2C_GetPlayerInfo.prototype.name =  "";
                
        
        S2C_GetPlayerInfo.prototype.type =  EMsgType.UNKNOWN;
                
        
        S2C_GetPlayerInfo.prototype.items =  $util.emptyArray;
         

        S2C_GetPlayerInfo.create = function create(properties) {
            return new S2C_GetPlayerInfo(properties);
        };

        S2C_GetPlayerInfo.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
                
             
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.name); 
             
            if (message.type != null && Object.hasOwnProperty.call(message, "type"))
                writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.type); 
            
            if (message.items != null && message.items.length) {
                writer.uint32(/* id 3, wireType 2 =*/26).fork();
                for (var i = 0; i < message.items.length; ++i)
                    writer.Item(message.items[i]);
                writer.ldelim();
            } 
            return writer;
        };

        S2C_GetPlayerInfo.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        S2C_GetPlayerInfo.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new S2C_GetPlayerInfo();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {          
                case 1:
                    message.name = reader.string();  
                    break;           
                case 2:
                    message.type = reader.uint32();  
                    break;           
                case 3:                    
                    if (!(message.items && message.items.length))
                        message.items = [];
                    if ((tag & 7) === 2) {
                        var end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2)
                            message.items.push(reader.Item());
                    } else
                        message.items.push(reader.Item());   
                    break;  
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        S2C_GetPlayerInfo.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        S2C_GetPlayerInfo.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
                if (message.name != null && message.hasOwnProperty("name")) {
                    if (!$util.isString(message.name))
                        return "name: string expected"; 
                } 
                if (message.type != null && message.hasOwnProperty("type")) {
                    switch (message.type) {
                        default:
                            return "type: enum value expected";
                        case 0:
                        case 1:
                        case 2:
                            break;
                    } 
                } 
                if (message.items != null && message.hasOwnProperty("items")) {
                    if (!Array.isArray(message.items))
                        return "items: array expected";
                        "error type Item items"; 
                }  
            return null;
        };

        S2C_GetPlayerInfo.fromObject = function fromObject(object) {
            if (object instanceof S2C_GetPlayerInfo)
                return object;
            var message = new S2C_GetPlayerInfo();                
            if (object.name != null)
                message.name = String(object.name); 
            switch (object.type) {
                default:
                case "UNKNOWN":
                case 0:
                    message.type = 0;
                    break;
                case "JSON":
                case 1:
                    message.type = 1;
                    break;
                case "XML":
                case 2:
                    message.type = 2;
                    break;    
                } 
            if(object.items) {
                if (!Array.isArray(object.items))
                    throw TypeError("S2C_GetPlayerInfo.items: array expected");            
                message.items = [];
                for (var i = 0; i < object.items.length; ++i)
            }  
            return message;
        };

        S2C_GetPlayerInfo.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults) {  
                object.items = []; 
            }

            if (options.defaults) {
                object.name = ""; 
                
                object.type = options.enums === String ? "UNKNOWN" : 0; 
                object.items = null;  
            }
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name; 
            if (message.type != null && message.hasOwnProperty("type"))
                    object.type = options.enums === String ? EMsgType[message.type] : message.type; 
                
            if (message.items && message.items.length) {
                object.items = [];
                for (var j = 0; j < message.items.length; ++j)
                    object.items[j] = message.items[j];
            } 
            return object;
        };

        S2C_GetPlayerInfo.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return S2C_GetPlayerInfo;
    })(); 
    ALLTYPES["S2C_GetPlayerInfo"] = S2C_GetPlayerInfo; 

    return GameData;
})(); 

module.exports = $root;