export interface Peer {
    id: string
    publicKey: string
    userFName?: string
    userLName?: string
    hostname: string
    allowedIp: string
    endpoint?: string
    latestHandshake?: number
    transferRxBytes?: number
    transferTxBytes?: number
}

export interface NewPeer {
    userFName?: string
    userLName?: string
    hostname: string
    publicKey: string
    allowedIp: string
}
