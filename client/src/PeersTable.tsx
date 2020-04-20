import React from 'react'
import { Table, Button, Row, Col } from 'antd'
import { useQuery } from '@apollo/react-hooks'
import gql from 'graphql-tag'
import { Link } from 'react-router-dom'
import cloneDeep from 'lodash.clonedeep'
import { Peer } from './ServerTypes'

interface PeersData {
  readonly peers: readonly Peer[]
}

const NO_BYTES = 0
const DEFAULT_DECIMALS = 2
const NO_DECIMALS = 0
const NO_TIME = 0
const S_TO_MS = 1000
const GET_PEERS = gql`
  {
    peers {
      hostname
      userFName
      userLName
      latestHandshake
      id
      publicKey
      allowedIp
      endpoint
      transferRxBytes
      transferTxBytes
    }
  }
`

function formatBytes(bytes: number, decimals = DEFAULT_DECIMALS): string {
  if (bytes === NO_BYTES) return '0 Bytes'

  const k = 1024
  const dm = decimals < NO_DECIMALS ? NO_DECIMALS : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

function formatDate(seconds: number): string {
  if (seconds === NO_TIME) return 'Disconnected'

  const d = new Date(seconds * S_TO_MS)
  return d.toLocaleString()
}

function formatUserName(firstName?: string, lastName?: string): string {
  if (firstName === null && lastName === null) return ''
  if (firstName === null) return lastName
  if (lastName === null) return firstName
  return `${firstName} ${lastName}`
}

export const PeersTable: React.FC = () => {
  let children = null
  const { loading, data, error } = useQuery<PeersData>(GET_PEERS)
  if (loading) children = <p>Loading... </p>
  else if (error) children = <p>Error: {error}</p>
  else {
    const { peers } = data
    const columns = [
      { title: 'Hostname', dataIndex: 'hostname', key: 'hostname' },
      {
        title: 'User Name',
        dataIndex: 'username',
        key: 'username',
        render: (_, record: Peer): string =>
          formatUserName(record.userFName, record.userLName),
      },
      { title: 'Public Key', dataIndex: 'publicKey', key: 'publicKey' },
      { title: 'IP Address', dataIndex: 'allowedIp', key: 'allowedIp' },
      { title: 'Endpoint', dataIndex: 'endpoint', key: 'endpoint' },
      {
        title: 'Latest Handshake',
        dataIndex: 'latestHandshake',
        key: 'latestHandshake',
        render: (seconds: number): string => formatDate(seconds),
      },
      {
        title: 'Received Bytes',
        dataIndex: 'transferRxBytes',
        key: 'transferRxBytes',
        render: (bytes: number): string => formatBytes(bytes),
      },
      {
        title: 'Transmitted Bytes',
        dataIndex: 'transferTxBytes',
        key: 'transferTxBytes',
        render: (bytes: number): string => formatBytes(bytes),
      },
    ]

    children = (
      <>
        <Row justify="end">
          <Col span={1}>
            <h1>Peers: </h1>
          </Col>
          <Col span={2} offset={21}>
            <Link to="/add-peer">
              <Button>Add Peer</Button>
            </Link>
          </Col>
        </Row>
        <Table
          dataSource={cloneDeep(peers)}
          columns={columns}
          pagination={false}
          rowKey="publicKey"
        />
      </>
    )
  }
  return (
    <Row style={{ marginTop: 20 }}>
      <Col span={20} offset={2}>
        {children}
      </Col>
    </Row>
  )
}
