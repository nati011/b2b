'use client'
import useProfileStore from "@/app/libs/store/useProfileStore"
import { useEffect, useState } from "react";
import { Button, Card, Form, Input } from "antd";
import { VscAccount } from "react-icons/vsc";
import { DatePicker, Space } from 'antd';
import type { DatePickerProps, GetProps } from 'antd';

export default function Settings() {
    const {
        profile,
        loading,
        error,
        fetchProfile
    } = useProfileStore()
    const [form] = Form.useForm();
    const [success, setSuccesss] = useState(true)


    type RangePickerProps = GetProps<typeof DatePicker.RangePicker>;

    const { RangePicker } = DatePicker;

    const onOk = (value: DatePickerProps['value'] | RangePickerProps['value']) => {
        console.log('onOk: ', value);
    };


    useEffect(() => {
        fetchProfile();
    }, []);

    return (
        <Card title="Account Settings">
            <Form
                layout={'vertical'}
                form={form}
            >
                <div className="flex gap-8 w-full">
                    <div className="">
                        <VscAccount className="text-9xl text-gray-700" />
                    </div>
                    <div className="grid grid-cols-1 w-full">
                        <Form.Item label="First Name">
                            <Input placeholder="First Name" />
                        </Form.Item>
                        <Form.Item label="Last Name">
                            <Input placeholder="Last Name" />
                        </Form.Item>
                    </div>
                </div>
                <Form.Item label="Email">
                    <Input placeholder="Email" />
                </Form.Item>
                <Form.Item label="Phone">
                    <Input placeholder="Phone" />
                </Form.Item>
                <Form.Item label="Username">
                    <Input placeholder="Username" />
                </Form.Item>
                <Form.Item label="BirthDate">
                    <DatePicker
                        className="w-full"
                        onChange={(value, dateString) => {
                            console.log('Selected Time: ', value);
                            console.log('Formatted Selected Time: ', dateString);
                        }}
                        onOk={onOk}
                    />
                </Form.Item>
            </Form>
            <div style={{ marginTop: 24 }} className="flex w-full justify-end">
                <Button style={{ margin: '0 8px' }}>
                    Cancel
                </Button>

                <Button type="primary" onClick={() => setSuccesss(true)} className="px-10 border border-[#451f87]">
                    Submit
                </Button>
            </div>
        </Card>

    );
}
