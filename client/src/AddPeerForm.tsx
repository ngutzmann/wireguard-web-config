import * as React from 'react'
import { Row, Col, Form, Input, Button } from 'antd'
import { Link, useHistory } from 'react-router-dom'
import gql from 'graphql-tag'
import { useMutation } from '@apollo/react-hooks'

const CREATE_PEER = gql`
    mutation CreatePeer($input: NewPeer!) {
        createPeer(input: $input )  {
            id
        }
    }`

export function AddPeerForm() {

    const [createPeer] = useMutation(CREATE_PEER)
    const history = useHistory()

    async function onFinish(values) {
        const { hostname, publicKey, ipAddress } = values
        try {
            await createPeer({ variables: { input: { name: hostname, publicKey, allowedIp: ipAddress } } })
            history.push('/')
        }
        catch (e) {
            console.error("Error creating peer: ", e)
        }
    }

    return (<>
        <Row style={{ marginTop: 20 }}>
            <Col span={4} offset={4}><h1>Add Peer:</h1></Col>
            <Col span={2} offset={10}>
                <Link to="/">
                    <Button>
                        Back
                        </Button>
                </Link>
            </Col>
        </Row>
        <Row>
            <Col span={8} offset={2}>
                <Form
                    labelCol={{ span: 8 }}
                    wrapperCol={{ span: 16 }}
                    name="basic"
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                >
                    <Form.Item
                        label="Hostname"
                        name="hostname"
                        rules={[{ required: true, message: 'Please input the system hostname!' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="Public Key"
                        name="publicKey"
                        rules={[{ required: true, message: 'Please input the system public key!' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="IP Address"
                        name="ipAddress"
                        rules={[{ required: true, message: 'Please input the system IP!' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
                        <Button type="primary" htmlType="submit">
                            Submit
                    </Button>
                    </Form.Item>
                </Form>
            </Col>
        </Row>
    </>)
}
