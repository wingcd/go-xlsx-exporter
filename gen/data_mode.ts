// DO NOT EDIT!
// This code is auto generated by go-xlsx-exporter
// VERSION 1.2
// go-protobuf v1.27.1

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
    static convertHandler: (data: DataModel, fieldName:string, value: string)=>any = null;
}

export class DataModel {
    private _converted = {};

    protected getConvertData(fieldName: string, value: any): any
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
}

export namespace GameData {
    var ALLTYPES: {[key: string]: any} = {};

    
    // Defined in table: 
    export enum EMsgType {
        UNKNOWN = 0,
        
        JSON = 1,
        
        XML = 2,
        
    }
    

    
    // Defined in table: data/message.xml
    /** Properties of a MessageWrapper. */
    export interface IMessageWrapper {
                       
        id?: (number|null);
             
                       
        data?: (Uint8Array|null);
             
         
    }

     /** Represents a MessageWrapper. */
    export class MessageWrapper extends DataModel implements IMessageWrapper { 
        private static __type_name__ = "MessageWrapper";
        private static __array_type_name__ = "MessageWrapper_ARRAY";
        static getArrayType(): any {
            return ALLTYPES[MessageWrapper.__array_type_name__];
        }

        
        id?: (number|null) = 0;
                
        
        data =  $util.newBuffer([]);
                
        

        /**
         * Constructs a new MessageWrapper.
         * @param [properties] Properties to set
         */
        constructor(properties?: IMessageWrapper) {
            super();
            
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        static create(properties?: MessageWrapper): MessageWrapper {
            return new MessageWrapper(properties);
        }

        static encode(message: IMessageWrapper, writer?: $protobuf.Writer): $protobuf.Writer {
            if (!writer)
                writer = $Writer.create();
                
            
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.id); 
            
            if (message.data != null && Object.hasOwnProperty.call(message, "data"))
                writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.data); 
            

            return writer;
        }

        static encodeDelimited(message: IMessageWrapper, writer?: $protobuf.Writer): $protobuf.Writer {
            return this.encode(message, writer).ldelim();
        }

        static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): MessageWrapper {
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
        }

        static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): MessageWrapper {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        }

        static verify(message: { [k: string]: any }): (string|null) {
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
        }

        static fromObject(object: { [k: string]: any }): MessageWrapper {
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
        }

        static toObject(message: MessageWrapper, options?: $protobuf.IConversionOptions): { [k: string]: any } {
            if (!options)
                options = {};
            var object: any = {};
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
        }

        toJSON(): { [k: string]: any } {
            return MessageWrapper.toObject(this, $protobuf.util.toJSONOptions);
        }
    } 
    ALLTYPES["MessageWrapper"] = MessageWrapper;

    
    // Defined in table: data/message.xml
    /** Properties of a Item. */
    export interface IItem {
                       
        id?: (number|null);
             
                       
        name?: (string|null);
             
         
    }

     /** Represents a Item. */
    export class Item extends DataModel implements IItem { 
        private static __type_name__ = "Item";
        private static __array_type_name__ = "Item_ARRAY";
        static getArrayType(): any {
            return ALLTYPES[Item.__array_type_name__];
        }

        
        id?: (number|null) = 0;
                
        
        name?: (string|null) = "";
                
        

        /**
         * Constructs a new Item.
         * @param [properties] Properties to set
         */
        constructor(properties?: IItem) {
            super();
            
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        static create(properties?: Item): Item {
            return new Item(properties);
        }

        static encode(message: IItem, writer?: $protobuf.Writer): $protobuf.Writer {
            if (!writer)
                writer = $Writer.create();
                
            
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.id); 
            
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.name); 
            

            return writer;
        }

        static encodeDelimited(message: IItem, writer?: $protobuf.Writer): $protobuf.Writer {
            return this.encode(message, writer).ldelim();
        }

        static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Item {
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
        }

        static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Item {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        }

        static verify(message: { [k: string]: any }): (string|null) {
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
        }

        static fromObject(object: { [k: string]: any }): Item {
            if (object instanceof Item)
                return object;
            var message = new Item();                
            if (object.id != null)
                message.id = object.id | 0;                 
            if (object.name != null)
                message.name = String(object.name);  
            return message;
        }

        static toObject(message: Item, options?: $protobuf.IConversionOptions): { [k: string]: any } {
            if (!options)
                options = {};
            var object: any = {};
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
        }

        toJSON(): { [k: string]: any } {
            return Item.toObject(this, $protobuf.util.toJSONOptions);
        }
    } 
    ALLTYPES["Item"] = Item;

    
    // Defined in table: data/message.xml
    /** Properties of a C2S_GetPlayerInfo. */
    export interface IC2S_GetPlayerInfo {
                       
        name?: (string|null);
             
         
    }

     /** Represents a C2S_GetPlayerInfo. */
    export class C2S_GetPlayerInfo extends DataModel implements IC2S_GetPlayerInfo { 
        private static __type_name__ = "C2S_GetPlayerInfo";
        private static __array_type_name__ = "C2S_GetPlayerInfo_ARRAY";
        static getArrayType(): any {
            return ALLTYPES[C2S_GetPlayerInfo.__array_type_name__];
        }

        
        name?: (string|null) = "";
                
        

        /**
         * Constructs a new C2S_GetPlayerInfo.
         * @param [properties] Properties to set
         */
        constructor(properties?: IC2S_GetPlayerInfo) {
            super();
            
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        static create(properties?: C2S_GetPlayerInfo): C2S_GetPlayerInfo {
            return new C2S_GetPlayerInfo(properties);
        }

        static encode(message: IC2S_GetPlayerInfo, writer?: $protobuf.Writer): $protobuf.Writer {
            if (!writer)
                writer = $Writer.create();
                
            
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.name); 
            

            return writer;
        }

        static encodeDelimited(message: IC2S_GetPlayerInfo, writer?: $protobuf.Writer): $protobuf.Writer {
            return this.encode(message, writer).ldelim();
        }

        static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): C2S_GetPlayerInfo {
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
        }

        static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): C2S_GetPlayerInfo {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        }

        static verify(message: { [k: string]: any }): (string|null) {
            if (typeof message !== "object" || message === null)
                return "object expected";     
                if (message.name != null && message.hasOwnProperty("name")) {
                    if (!$util.isString(message.name))
                        return "name: string expected"; 
                }  
            return null;
        }

        static fromObject(object: { [k: string]: any }): C2S_GetPlayerInfo {
            if (object instanceof C2S_GetPlayerInfo)
                return object;
            var message = new C2S_GetPlayerInfo();                
            if (object.name != null)
                message.name = String(object.name);  
            return message;
        }

        static toObject(message: C2S_GetPlayerInfo, options?: $protobuf.IConversionOptions): { [k: string]: any } {
            if (!options)
                options = {};
            var object: any = {};
            if (options.arrays || options.defaults) { 
            }

            if (options.defaults) {
                object.name = "";  
            }
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name; 
            return object;
        }

        toJSON(): { [k: string]: any } {
            return C2S_GetPlayerInfo.toObject(this, $protobuf.util.toJSONOptions);
        }
    } 
    ALLTYPES["C2S_GetPlayerInfo"] = C2S_GetPlayerInfo;

    
    // Defined in table: data/message.xml
    /** Properties of a S2C_GetPlayerInfo. */
    export interface IS2C_GetPlayerInfo {
                       
        name?: (string|null);
             
                       
        type?: (EMsgType|null);
             
                       
        items?: (Item[]|null);
             
         
    }

     /** Represents a S2C_GetPlayerInfo. */
    export class S2C_GetPlayerInfo extends DataModel implements IS2C_GetPlayerInfo { 
        private static __type_name__ = "S2C_GetPlayerInfo";
        private static __array_type_name__ = "S2C_GetPlayerInfo_ARRAY";
        static getArrayType(): any {
            return ALLTYPES[S2C_GetPlayerInfo.__array_type_name__];
        }

        
        name?: (string|null) = "";
                
        
        type?: (EMsgType|null) = EMsgType.UNKNOWN;
                
        
        items =  $util.emptyArray;
        

        /**
         * Constructs a new S2C_GetPlayerInfo.
         * @param [properties] Properties to set
         */
        constructor(properties?: IS2C_GetPlayerInfo) {
            super();
            
            
            this.items = [];
                    
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        static create(properties?: S2C_GetPlayerInfo): S2C_GetPlayerInfo {
            return new S2C_GetPlayerInfo(properties);
        }

        static encode(message: IS2C_GetPlayerInfo, writer?: $protobuf.Writer): $protobuf.Writer {
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
        }

        static encodeDelimited(message: IS2C_GetPlayerInfo, writer?: $protobuf.Writer): $protobuf.Writer {
            return this.encode(message, writer).ldelim();
        }

        static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): S2C_GetPlayerInfo {
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
        }

        static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): S2C_GetPlayerInfo {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        }

        static verify(message: { [k: string]: any }): (string|null) {
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
        }

        static fromObject(object: { [k: string]: any }): S2C_GetPlayerInfo {
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
        }

        static toObject(message: S2C_GetPlayerInfo, options?: $protobuf.IConversionOptions): { [k: string]: any } {
            if (!options)
                options = {};
            var object: any = {};
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
        }

        toJSON(): { [k: string]: any } {
            return S2C_GetPlayerInfo.toObject(this, $protobuf.util.toJSONOptions);
        }
    } 
    ALLTYPES["S2C_GetPlayerInfo"] = S2C_GetPlayerInfo; 
}