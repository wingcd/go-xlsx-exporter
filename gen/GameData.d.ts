// DO NOT EDIT! This is a generated file. Edit the JSDoc in src/*.js instead and run 'npm run types'.

/** Namespace GameData. */
export namespace GameData {

    /** Corpus enum. */
    enum Corpus {
        UNIVERSAL = 0,
        WEB = 1,
        IMAGES = 2,
        LOCAL = 3,
        NEWS = 4,
        PRODUCTS = 5,
        VIDEO = 6
    }

    /** Properties of a Chain. */
    interface IChain {

        /** Chain VarDouble1 */
        VarDouble1?: (number|null);

        /** Chain VarBool */
        VarBool?: (boolean|null);

        /** Chain VarInt32 */
        VarInt32?: (number|null);

        /** Chain VarUint32 */
        VarUint32?: (number|null);

        /** Chain VarInt64 */
        VarInt64?: (number|Long|null);

        /** Chain VarUint64 */
        VarUint64?: (number|Long|null);

        /** Chain VarFloat */
        VarFloat?: (number|null);

        /** Chain VarDouble */
        VarDouble?: (number|null);

        /** Chain VarString */
        VarString?: (string|null);

        /** Chain VarBools */
        VarBools?: (boolean[]|null);

        /** Chain VarInt32s */
        VarInt32s?: (number[]|null);

        /** Chain VarUint32s */
        VarUint32s?: (number[]|null);

        /** Chain VarInt64s */
        VarInt64s?: ((number|Long)[]|null);

        /** Chain VarUint64s */
        VarUint64s?: ((number|Long)[]|null);

        /** Chain VarFloats */
        VarFloats?: (number[]|null);

        /** Chain VarDoubles */
        VarDoubles?: (number[]|null);

        /** Chain VarStrings */
        VarStrings?: (string[]|null);

        /** Chain corpus */
        corpus?: (GameData.Corpus|null);

        /** Chain corpuss */
        corpuss?: (GameData.Corpus[]|null);
    }

    /** Represents a Chain. */
    class Chain implements IChain {

        /**
         * Constructs a new Chain.
         * @param [properties] Properties to set
         */
        constructor(properties?: GameData.IChain);

        /** Chain VarDouble1. */
        public VarDouble1: number;

        /** Chain VarBool. */
        public VarBool: boolean;

        /** Chain VarInt32. */
        public VarInt32: number;

        /** Chain VarUint32. */
        public VarUint32: number;

        /** Chain VarInt64. */
        public VarInt64: (number|Long);

        /** Chain VarUint64. */
        public VarUint64: (number|Long);

        /** Chain VarFloat. */
        public VarFloat: number;

        /** Chain VarDouble. */
        public VarDouble: number;

        /** Chain VarString. */
        public VarString: string;

        /** Chain VarBools. */
        public VarBools: boolean[];

        /** Chain VarInt32s. */
        public VarInt32s: number[];

        /** Chain VarUint32s. */
        public VarUint32s: number[];

        /** Chain VarInt64s. */
        public VarInt64s: (number|Long)[];

        /** Chain VarUint64s. */
        public VarUint64s: (number|Long)[];

        /** Chain VarFloats. */
        public VarFloats: number[];

        /** Chain VarDoubles. */
        public VarDoubles: number[];

        /** Chain VarStrings. */
        public VarStrings: string[];

        /** Chain corpus. */
        public corpus: GameData.Corpus;

        /** Chain corpuss. */
        public corpuss: GameData.Corpus[];

        /**
         * Creates a new Chain instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Chain instance
         */
        public static create(properties?: GameData.IChain): GameData.Chain;

        /**
         * Encodes the specified Chain message. Does not implicitly {@link GameData.Chain.verify|verify} messages.
         * @param message Chain message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: GameData.IChain, writer?: protobuf.Writer): protobuf.Writer;

        /**
         * Encodes the specified Chain message, length delimited. Does not implicitly {@link GameData.Chain.verify|verify} messages.
         * @param message Chain message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: GameData.IChain, writer?: protobuf.Writer): protobuf.Writer;

        /**
         * Decodes a Chain message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Chain
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: (protobuf.Reader|Uint8Array), length?: number): GameData.Chain;

        /**
         * Decodes a Chain message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Chain
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: (protobuf.Reader|Uint8Array)): GameData.Chain;

        /**
         * Verifies a Chain message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Chain message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Chain
         */
        public static fromObject(object: { [k: string]: any }): GameData.Chain;

        /**
         * Creates a plain object from a Chain message. Also converts values to other types if specified.
         * @param message Chain
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: GameData.Chain, options?: protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Chain to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };
    }

    /** Properties of a Chain_ARRAY. */
    interface IChain_ARRAY {

        /** Chain_ARRAY Items */
        Items?: (GameData.IChain[]|null);
    }

    /** Represents a Chain_ARRAY. */
    class Chain_ARRAY implements IChain_ARRAY {

        /**
         * Constructs a new Chain_ARRAY.
         * @param [properties] Properties to set
         */
        constructor(properties?: GameData.IChain_ARRAY);

        /** Chain_ARRAY Items. */
        public Items: IChain[];

        /**
         * Creates a new Chain_ARRAY instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Chain_ARRAY instance
         */
        public static create(properties?: GameData.IChain_ARRAY): GameData.Chain_ARRAY;

        /**
         * Encodes the specified Chain_ARRAY message. Does not implicitly {@link GameData.Chain_ARRAY.verify|verify} messages.
         * @param message Chain_ARRAY message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: GameData.IChain_ARRAY, writer?: protobuf.Writer): protobuf.Writer;

        /**
         * Encodes the specified Chain_ARRAY message, length delimited. Does not implicitly {@link GameData.Chain_ARRAY.verify|verify} messages.
         * @param message Chain_ARRAY message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: GameData.IChain_ARRAY, writer?: protobuf.Writer): protobuf.Writer;

        /**
         * Decodes a Chain_ARRAY message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Chain_ARRAY
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: (protobuf.Reader|Uint8Array), length?: number): GameData.Chain_ARRAY;

        /**
         * Decodes a Chain_ARRAY message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Chain_ARRAY
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: (protobuf.Reader|Uint8Array)): GameData.Chain_ARRAY;

        /**
         * Verifies a Chain_ARRAY message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Chain_ARRAY message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Chain_ARRAY
         */
        public static fromObject(object: { [k: string]: any }): GameData.Chain_ARRAY;

        /**
         * Creates a plain object from a Chain_ARRAY message. Also converts values to other types if specified.
         * @param message Chain_ARRAY
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: GameData.Chain_ARRAY, options?: protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Chain_ARRAY to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };
    }
}
