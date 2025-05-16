'use client'
import { Button, Card, Input, Form, message, Select, Breadcrumb } from "antd";
import TextArea from "antd/es/input/TextArea";
import { HomeOutlined } from '@ant-design/icons';

import { useEffect, useState } from "react";
import useDistributorsStore from "@/app/libs/store/useDistributorStore";



export default function AgentForm() {
    const [form] = Form.useForm();
    const [messageApi, contextHolder] = message.useMessage();
    const [success, setSuccesss] = useState(true)
    const [current, setCurrent] = useState(0);

    const {
        distributors,
        loading,
        error,
        fetchDistributors
    } = useDistributorsStore()

    useEffect(() => {
        fetchDistributors();
    }, []);




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
                    Register Distributor Agent
                </p>
                <div className="flex items-center gap-2">
                    <Breadcrumb items={
                        [
                            {
                                title: <HomeOutlined />,
                                href: "/"
                            },
                            {
                                title: "Agents",
                                href: "/agents"
                            },
                            {
                                title: "Form",
                                href: "/agents/form"
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
                    <div className="grid grid-cols-2 gap-2">
                        <Form.Item label="First Name">
                            <Input placeholder="First Name" />
                        </Form.Item>
                        <Form.Item label="Last Name">
                            <Input placeholder="Last Name" />
                        </Form.Item>
                    </div>
                    <Form.Item label="Email"
                        rules={[{ required: true }]}
                    >
                        <Input placeholder="Email" />
                    </Form.Item>
                    <Form.Item label="Phone" rules={[{ required: true }]} >
                        <Input placeholder="Phone" />
                    </Form.Item>
                    <Form.Item label="Username" rules={[{ required: true }]}>
                        <Input placeholder="Username" />
                    </Form.Item>
                    <Form.Item
                        name="Distributor"
                        label="Distributor"
                        rules={[{ required: true }]}
                    >
                        <Select
                            showSearch
                            placeholder="Search to Select"
                            optionFilterProp="Distributor"
                            filterSort={(optionA, optionB) =>
                                (optionA?.name ?? '').toLowerCase().localeCompare((optionB?.name ?? '').toLowerCase())
                            }
                            options={distributors}
                        />
                    </Form.Item>
                </Form>
                <div className="flex w-full justify-end" style={{ marginTop: 24 }}>
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