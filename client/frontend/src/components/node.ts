export interface IP4Address {
  address: string
  prefix: number
  gateway: string
}

export interface Node {
  uuid: string
  ip4: IP4Address[]
  dhcp: boolean
}
