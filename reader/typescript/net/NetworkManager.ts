import { EventCenter } from "../core/utils/EventCenter";
import { MessageProcesser } from "./MessageProcesser";

export enum ENetworkState {
    Closed = 0,
    Connecting = 1,
    Open = 2,
}

export class NetworkManager {
    static STATUS_CHANGED = 1;

    private static _instance: NetworkManager = null;
    public static get I(): NetworkManager {
        if (this._instance === null) {
            if(typeof WebSocket === 'undefined') {
                console.error('WebSocket is not supported');
                return null;
            }

            this._instance = new NetworkManager();
        }
        return this._instance;
    }

    private _url: string = '';
    private _socket: WebSocket = null;
    private _processer: MessageProcesser = null;
    private _state: ENetworkState = ENetworkState.Closed;

    private _eventCenter: EventCenter = new EventCenter();

    public init(url: string, processer: MessageProcesser) {
        if (this._socket !== null) {
            this._socket.close();
        }

        this._state = ENetworkState.Closed;
        processer.bind(this);
        this._processer = processer;
        this._url = url;
    }

    public connect(): WebSocket {
        if (this._socket !== null) {
            this._socket.close();
        }

        this._socket = new WebSocket(this._url);
        this._socket.binaryType = this._processer.binaryType;
        this._state = ENetworkState.Connecting;

        this._socket.onopen = this.onOpen.bind(this);
        this._socket.onmessage = this.onMessage.bind(this);
        this._socket.onclose = this.onClose.bind(this);
        this._socket.onerror = this.onError.bind(this);

        return this._socket;
    }

    public send(data: any) {
        if (this._socket !== null) {
            let message = this._processer.serialize(data);
            this._socket.send(message);
        }
    }

    private onOpen(event: Event) {
        this._state = ENetworkState.Open;

        console.log('NetworkManager.onOpen');
        this._processer.onOpen(event);

        this.emit(NetworkManager.STATUS_CHANGED, this._state);
    }

    private onMessage(event: MessageEvent) {
        console.log('NetworkManager.onMessage', event.data);
        this._processer.onMessage(event);
    }

    private onClose(event: CloseEvent) {
        this._state = ENetworkState.Closed;

        console.log('NetworkManager.onClose');
        this._processer.onClose(event);

        this.emit(NetworkManager.STATUS_CHANGED, this._state);
    }

    private onError(event: Event) {
        this._state = ENetworkState.Closed;

        console.log('NetworkManager.onError');
        this._processer.onError(event);

        this.emit(NetworkManager.STATUS_CHANGED, this._state);
    }

    public close() {
        if (this._socket !== null) {
            this._socket.close();
        }

        this._state = ENetworkState.Closed;
    }

    public reconnect() {
        if (this._socket !== null) {
            this._socket.close();
        }

        setTimeout(() => {
            this.connect();
        }, 1000);
    }

    public on(event: number, callback: Function, target: any = null) {
        this._eventCenter.on(event, callback, target);
    }
    
    public off(event: number, callback: Function, target: any = null) {
        this._eventCenter.off(event, callback, target);
    }

    public emit(event: number, ...args: any[]) {
        this._eventCenter.emit(event, ...args);
    }

    public once(event: number, callback: Function, target: any = null) {
        this._eventCenter.once(event, callback, target);
    }
}