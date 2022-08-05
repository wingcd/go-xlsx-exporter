import $protobuf from "protobufjs";

var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

export class Language {

    /**
     * Properties of a Language.
     * @memberof GameData
     * @interface ILanguage
     * @property {string|null} [ID] Language ID
     * @property {string|null} [Text] Language Text
     */

    /**
     * Constructs a new Language.
     * @memberof GameData
     * @classdesc Represents a Language.
     * @implements ILanguage
     * @constructor
     * @param {ILanguage=} [properties] Properties to set
     */
    constructor(properties?:any) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Language ID.
     * @member {string} ID
     * @memberof Language
     * @instance
     */
    ID = "";

    /**
     * Language Text.
     * @member {string} Text
     * @memberof Language
     * @instance
     */
    Text = "";

    /**
     * Creates a new Language instance using the specified properties.
     * @function create
     * @memberof Language
     * @static
     * @param {ILanguage=} [properties] Properties to set
     * @returns {Language} Language instance
     */
    static create(properties) {
        return new Language(properties);
    };

    /**
     * Encodes the specified Language message. Does not implicitly {@link Language.verify|verify} messages.
     * @function encode
     * @memberof Language
     * @static
     * @param {ILanguage} message Language message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    static encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.ID != null && Object.hasOwnProperty.call(message, "ID"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.ID);
        if (message.Text != null && Object.hasOwnProperty.call(message, "Text"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.Text);
        return writer;
    };

    /**
     * Encodes the specified Language message, length delimited. Does not implicitly {@link Language.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Language
     * @static
     * @param {ILanguage} message Language message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    static encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Language message from the specified reader or buffer.
     * @function decode
     * @memberof Language
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Language} Language
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    static decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new Language();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.ID = reader.string();
                break;
            case 2:
                message.Text = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Language message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Language
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Language} Language
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    static decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Language message.
     * @function verify
     * @memberof Language
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
     static verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.ID != null && message.hasOwnProperty("ID"))
            if (!$util.isString(message.ID))
                return "ID: string expected";
        if (message.Text != null && message.hasOwnProperty("Text"))
            if (!$util.isString(message.Text))
                return "Text: string expected";
        return null;
    };

    /**
     * Creates a Language message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Language
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Language} Language
     */
     static fromObject(object) {
        if (object instanceof Language)
            return object;
        var message = new Language();
        if (object.ID != null)
            message.ID = String(object.ID);
        if (object.Text != null)
            message.Text = String(object.Text);
        return message;
    };

    /**
     * Creates a plain object from a Language message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Language
     * @static
     * @param {Language} message Language
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    static toObject(message, options) {
        if (!options)
            options = {};
        var object:any = {};
        if (options.defaults) {
            object.ID = "";
            object.Text = "";
        }
        if (message.ID != null && message.hasOwnProperty("ID"))
            object.ID = message.ID;
        if (message.Text != null && message.hasOwnProperty("Text"))
            object.Text = message.Text;
        return object;
    };

    /**
     * Converts this Language to JSON.
     * @function toJSON
     * @memberof Language
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    toJSON() {
        return Language.toObject(this, $protobuf.util.toJSONOptions);
    };
}

export class Language_ARRAY {

    /**
     * Properties of a Language_ARRAY.
     * @memberof GameData
     * @interface ILanguage_ARRAY
     * @property {Array.<ILanguage>|null} [Items] Language_ARRAY Items
     */

    /**
     * Constructs a new Language_ARRAY.
     * @memberof GameData
     * @classdesc Represents a Language_ARRAY.
     * @implements ILanguage_ARRAY
     * @constructor
     * @param {ILanguage_ARRAY=} [properties] Properties to set
     */
    constructor(properties?:any) {
        this.Items = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Language_ARRAY Items.
     * @member {Array.<ILanguage>} Items
     * @memberof Language_ARRAY
     * @instance
     */
    Items = $util.emptyArray;

    /**
     * Creates a new Language_ARRAY instance using the specified properties.
     * @function create
     * @memberof Language_ARRAY
     * @static
     * @param {ILanguage_ARRAY=} [properties] Properties to set
     * @returns {Language_ARRAY} Language_ARRAY instance
     */
     static create(properties) {
        return new Language_ARRAY(properties);
    };

    /**
     * Encodes the specified Language_ARRAY message. Does not implicitly {@link Language_ARRAY.verify|verify} messages.
     * @function encode
     * @memberof Language_ARRAY
     * @static
     * @param {ILanguage_ARRAY} message Language_ARRAY message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
     static encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.Items != null && message.Items.length)
            for (var i = 0; i < message.Items.length; ++i)
                Language.encode(message.Items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Language_ARRAY message, length delimited. Does not implicitly {@link Language_ARRAY.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Language_ARRAY
     * @static
     * @param {ILanguage_ARRAY} message Language_ARRAY message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
     static encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Language_ARRAY message from the specified reader or buffer.
     * @function decode
     * @memberof Language_ARRAY
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Language_ARRAY} Language_ARRAY
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
     static decode(reader, length?) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new Language_ARRAY();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.Items && message.Items.length))
                    message.Items = [];
                message.Items.push(Language.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Language_ARRAY message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Language_ARRAY
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Language_ARRAY} Language_ARRAY
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
     static decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Language_ARRAY message.
     * @function verify
     * @memberof Language_ARRAY
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
     static verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.Items != null && message.hasOwnProperty("Items")) {
            if (!Array.isArray(message.Items))
                return "Items: array expected";
            for (var i = 0; i < message.Items.length; ++i) {
                var error = Language.verify(message.Items[i]);
                if (error)
                    return "Items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Language_ARRAY message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Language_ARRAY
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Language_ARRAY} Language_ARRAY
     */
     static fromObject(object) {
        if (object instanceof Language_ARRAY)
            return object;
        var message = new Language_ARRAY();
        if (object.Items) {
            if (!Array.isArray(object.Items))
                throw TypeError(".Language_ARRAY.Items: array expected");
            message.Items = [];
            for (var i = 0; i < object.Items.length; ++i) {
                if (typeof object.Items[i] !== "object")
                    throw TypeError(".Language_ARRAY.Items: object expected");
                message.Items[i] = Language.fromObject(object.Items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Language_ARRAY message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Language_ARRAY
     * @static
     * @param {Language_ARRAY} message Language_ARRAY
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
     static toObject(message, options) {
        if (!options)
            options = {};
        var object:any = {};
        if (options.arrays || options.defaults)
            object.Items = [];
        if (message.Items && message.Items.length) {
            object.Items = [];
            for (var j = 0; j < message.Items.length; ++j)
                object.Items[j] = Language.toObject(message.Items[j], options);
        }
        return object;
    };

    /**
     * Converts this Language_ARRAY to JSON.
     * @function toJSON
     * @memberof Language_ARRAY
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    toJSON() {
        return Language_ARRAY.toObject(this, $protobuf.util.toJSONOptions);
    };
}