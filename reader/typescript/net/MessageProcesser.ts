import { NetworkManager } from "./NetworkManager";

export class MessageProcesser {
    protected networkManger: NetworkManager = null;

    get binaryType(): BinaryType {
        return 'arraybuffer';
    }

    bind(netMgr: NetworkManager) {
        this.networkManger = netMgr;
    }

    parse(data: any): any {
        return null;
    }

    serialize(data: any): any {
        return null;
    }
    onOpen(event: Event) {
    }

    onClose(event: CloseEvent) {
    }

    onMessage(event: MessageEvent): any {
    }

    onError(event: Event) {
    }
}

export class JsonMessageProcesser extends MessageProcesser {
    parse(data: any): any {
        return JSON.parse(data);
    }

    serialize(data: any): any {
        return JSON.stringify(data);
    }
}