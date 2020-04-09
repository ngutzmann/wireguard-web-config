import * as React from 'react'
import { Table, Input, Button, Row, Col } from 'antd';
import Highlighter from 'react-highlight-words';
import { SearchOutlined } from '@ant-design/icons';
import { useQuery } from '@apollo/react-hooks'
import gql from "graphql-tag"

import { Peer } from "./ServerTypes"

const initialState = {
    searchText: ""
}

interface PeersData {
    peers: Peer[]
}

const GET_PEERS = gql`
{
    peers {
        name,
        latestHandshake,
        id,
        publicKey,
        allowedIp,
        endpoint,
        transferRxBytes,
        transferTxBytes
    }
}
`

type State = Readonly<typeof initialState>

function formatBytes(bytes: number, decimals = 2) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

function formatDate(seconds: number) {
    if (seconds === 0) return "Disconnected"

    const d = new Date(seconds * 1000)
    return d.toLocaleString()
}
export function PeersTable() {

    let children: React.ReactChild = null
    const { loading, data, error } = useQuery<PeersData>(GET_PEERS)
    if (loading) children = <p>Loading... </p>
    else if (error) children = <p>Error: {error}</p>

    else {
        const { peers } = data
        const columns = [
            { title: "Name", dataIndex: "name", key: "name" },
            { title: "Public Key", dataIndex: "publicKey", key: "publicKey" },
            { title: "IP Address", dataIndex: "allowedIp", key: "allowedIp" },
            { title: "Endpoint", dataIndex: "endpoint", key: "endpoint" },
            { title: "Latest Handshake", dataIndex: "latestHandshake", key: "latestHandshake", render: (seconds: number) => formatDate(seconds) },
            { title: "Received Bytes", dataIndex: "transferRxBytes", key: "transferRxBytes", render: (bytes: number) => formatBytes(bytes) },
            { title: "Transmitted Bytes", dataIndex: "transferTxBytes", key: "transferTxBytes", render: (bytes: number) => formatBytes(bytes) },
        ]

        children = <Table dataSource={peers} columns={columns} />
    }
    return (
        <Row style={{ marginTop: 20 }}>
            <Col span={20} offset={2}>
                <h1>Peers: </h1>
                {children}
            </Col>
        </Row>
    )
}
