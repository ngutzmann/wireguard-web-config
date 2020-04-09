export interface Peer {
    id: string
    publicKey: string
    name: string
    allowedIp: string
    endpoint?: string
    latestHandshake?: number
    transferRxBytes?: number
    transferTxBytes?: number
}

export interface NewPeer {
    name: string
    publicKey: string
    allowedIp: string
}
