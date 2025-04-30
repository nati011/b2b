'use client'
import { Button, Card, Input,  Form, message,Select,Breadcrumb} from "antd";
import TextArea from "antd/es/input/TextArea";
import { HomeOutlined } from '@ant-design/icons';

import { useEffect, useState } from "react";



export default function AgentForm() {
    const [form] = Form.useForm();
    const [messageApi, contextHolder] = message.useMessage();
    const [success, setSuccesss] = useState(true)
    const [current, setCurrent] = useState(0);



    const onFinish = () => {
        console.log(form)
        messageApi.success('Submit success!');
    };

    const onFinishFailed = () => {
        messageApi.error('Submit failed!');
    };

    return (
        <div className="">
        <div className="flex flex-col md:flex-row justify-between mb-4">
          <p className="font-semibold text-lg text-gray-900">
            Agent
          </p>
          <div className="flex items-center gap-2">
            <Breadcrumb items={
              [
                {
                title: <HomeOutlined />,
                href:"/"
              },
              {
                title: "Agents",
                href:"/agents"
              },
              {
                title: "Form",
                href:"/agents/form"
              }
              ]} />
  
          </div>
        </div>
            <Card title="Please fill in the form accordingly">
                <Form

                    form={form}
                    layout="vertical"
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="on"
                >
                    <Form.Item
                        name="Product"
                        label="Product"
                        rules={[{ required: true }]}
                    >
                         <Select
                            showSearch
                            placeholder="Search to Select"
                            optionFilterProp="label"
                            filterSort={(optionA, optionB) =>
                            (optionA?.label ?? '').toLowerCase().localeCompare((optionB?.label ?? '').toLowerCase())
                            }
                            options={[
                            {
                                value: '1',
                                label: 'Crop 1',
                            },
                            {
                                value: '2',
                                label: 'Crop 2',
                            },
                            {
                                value: '3',
                                label: 'Crop 3',
                            },
                            {
                                value: '4',
                                label: 'Crop 4',
                            }
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="Quantity"
                        label="Quantity"
                        rules={[{ required: true }]}
                    >
                        <Input placeholder="0001"type="number"/>
                    </Form.Item>

                    <Form.Item
                        name="Remark"
                        label="Remark"
                        rules={[{ required: true }]}
                    >
                       <TextArea placeholder="Remark"  style={{ height: 120}}/>
                    </Form.Item>
                </Form>
                <div style={{ marginTop: 24 }}>
                    <Button style={{ margin: '0 8px' }} className="px-8">
                        Cancel
                    </Button>

                    <Button type="primary" onClick={() => setSuccesss(true)} className="px-10 border border-[#451f87]">
                        Submit
                    </Button>
                </div>
            </Card>
        </div>
    )
}