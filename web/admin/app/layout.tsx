"use client"
import { Poppins, DM_Sans } from "next/font/google";
import "./globals.css";
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  TranslationOutlined,
  UserOutlined
} from '@ant-design/icons';
import { Avatar, Button, ConfigProvider, Layout, Menu, theme } from 'antd';

import { useState } from "react";
import { Sidebar } from "@/app/components/sidebar";

const { Header, Sider, Content } = Layout;

const poppins = DM_Sans({
  weight: ['100', '300', '700'],
  subsets: ['latin']
});


export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [collapsed, setCollapsed] = useState(false);
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#1A1F2C',
          colorPrimaryBg: '#ebeef7',
        },

      }}
    >
      <html lang="en">
        <body
          className={`${poppins} antialiased`}
        >
          <Layout>
            <Sider width={250} style={{

              background: colorBgContainer, overflow: 'auto',
              height: '100vh',
              position: 'sticky',
              insetInlineStart: 0,
              top: 0,
              bottom: 0,
              scrollbarWidth: 'thin',
              scrollbarGutter: 'stable',
            }} trigger={null} collapsible collapsed={collapsed}>
              <div className="px-16 py-10">
                <h2 className="font-semibold text-2xl text-blue-900">
                  Efoyta
                </h2>
              </div>

              <Menu
                style={{ fontSize: 14 }}
                mode="inline"

                theme="light"
                items={Sidebar}

                defaultSelectedKeys={['1']} />

            </Sider>
            <Layout style={{
              background: colorBgContainer
            }}>
              <Header style={{ padding: 10, background: colorBgContainer, display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <Button
                  type="text"
                  icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                  onClick={() => setCollapsed(!collapsed)}
                  style={{
                    fontSize: '16px',

                  }}
                />
                <div className="flex items-center ">
                  <Button
                    type="text"
                    icon={<TranslationOutlined />}
                    onClick={() => setCollapsed(!collapsed)}
                    style={{
                      fontSize: '20px',

                    }}
                  />
                  <Avatar
                    style={{
                      backgroundColor: '#1A1F2C',
                    }}
                    size={40}
                    icon={<UserOutlined />}
                  />
                </div>


              </Header>
              <div className="bg-[#F4F5F9] py-4 px-2 md:p-10 min-h-screen">
                {children}
              </div>


            </Layout>
          </Layout>
        </body>
      </html>
    </ConfigProvider>
  );
}
