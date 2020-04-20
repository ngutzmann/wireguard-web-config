import 'antd/dist/antd.css'
import React from 'react'
import { Layout, Menu } from 'antd'
import { grey } from '@ant-design/colors'
import ApolloClient from 'apollo-boost'
import { ApolloProvider } from '@apollo/react-hooks'
import { HashRouter as Router, Route, Switch } from 'react-router-dom'

import { AddPeerForm } from './AddPeerForm'
import { PeersTable } from './PeersTable'

const { Header, Content, Footer } = Layout

export const apolloClient = new ApolloClient({
  uri: 'http://localhost:8771/query',
})

const GREY_IDX = 0

export const Main: React.FC = () => {
  const serverName = process.env.SERVER_NAME || 'Wireguard Server'
  return (
    <ApolloProvider client={apolloClient}>
      <Router>
        <Layout className="layout">
          <Header>
            <h1 style={{ color: grey[GREY_IDX] }}>{serverName}</h1>
            <Menu theme="dark" mode="horizontal"></Menu>
          </Header>
          <Content style={{ height: 'calc(100vh - 55px)' }}>
            <Switch>
              <Route path="/add-peer">
                <AddPeerForm />
              </Route>
              <Route path="/">
                <PeersTable />
              </Route>
            </Switch>
          </Content>
          <Footer
            style={{
              position: 'absolute',
              bottom: 0,
              textAlign: 'center',
              width: '100%',
            }}
          >
            Wireguard Web Configuration
          </Footer>
        </Layout>
      </Router>
    </ApolloProvider>
  )
}
