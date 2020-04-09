import 'antd/dist/antd.css'
import * as React from 'react';
import { Layout, Menu, Breadcrumb } from 'antd'
import { grey } from '@ant-design/colors'
import ApolloClient from 'apollo-boost'

import { PeersTable } from "./peers"
import { ApolloProvider } from '@apollo/react-hooks';
const { Header, Content, Footer } = Layout


const apolloClient = new ApolloClient({
  uri: 'http://localhost:8771/query',
});


export const Main = () => {

  return (
    <ApolloProvider client={apolloClient}>
      <Layout className="layout">
        <Header>
          <h1 style={{ color: grey[0] }}>Personal AWS Wireguard Server</h1>
          <Menu theme="dark" mode="horizontal">
          </Menu>
        </Header>
        <Content style={{ height: "calc(100vh - 55px)" }}>
          <PeersTable />
        </Content>
        <Footer style={{ position: "absolute", bottom: 0, textAlign: "center", width: "100%" }}>Wireguard Web Configuration</Footer>
      </Layout>
    </ApolloProvider>
  )
}
