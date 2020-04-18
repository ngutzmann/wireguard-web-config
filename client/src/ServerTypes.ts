export interface Peer {
  readonly id: string
  readonly publicKey: string
  readonly userFName?: string
  readonly userLName?: string
  readonly hostname: string
  readonly allowedIp: string
  readonly endpoint?: string
  readonly latestHandshake?: number
  readonly transferRxBytes?: number
  readonly transferTxBytes?: number
}

export interface NewPeer {
  readonly userFName?: string
  readonly userLName?: string
  readonly hostname: string
  readonly publicKey: string
  readonly allowedIp: string
}
