'use client'
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { Distributor, DistributorRequest } from '@/app/libs/types';
import { Button, Card, Checkbox, Form, Input, message } from "antd";
import dynamic from "next/dynamic";
import { useState } from "react";


export default function DistributorsForm() {
  const {
    distributors,
    loading,
    error,
    fetchDistributors
  } = useDistributorsStore()

  const [form] = Form.useForm();
  const [markerPosition, setMarkerPosition] = useState<[number, number]>([9, 38])
  const [useCurrentLocation, setUseCurrentLocation] = useState(true)
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

  const [formData, setFormData] = useState<DistributorRequest>({
    name: "",
    tin: "",
    latitude: "",
    longitude: "",
    general_zone: "",
    region: "",
    woreda: "",
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    username: "",
    dob: "",
    external_id: "",

  })


  const Map = dynamic(
    () => import('@/app/components/map'),
    { ssr: false }
  )

  return (
    <div className="grid grid-cols-1 gap-4">
      <p className="font-semibold text-3xl text-black">
        Register Distributor
      </p>
      <Form
        layout={'vertical'}
        form={form}
        className="grid grid-cols-1 gap-4"
      >
        <Card title="Profile Information">
          <div className="grid grid-cols-2 gap-2">
            <Form.Item label="First Name">
              <Input placeholder="First Name" />
            </Form.Item>
            <Form.Item label="Last Name">
              <Input placeholder="Last Name" />
            </Form.Item>
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
        </Card>

        <Card title="Business Information">
          <Form.Item label="Tin">
            <Input placeholder="Tin Number" minLength={13} maxLength={13} />
          </Form.Item>
          <Form.Item label="General Zone">
            <Input placeholder="General Zone" />
          </Form.Item>
          <Form.Item label="Region">
            <Input placeholder="Region" />
          </Form.Item>
          <Form.Item label="Woreda">
            <Input placeholder="Woreda" />
          </Form.Item>
          <div className="flex items-center space-x-2 my-4">
            <Checkbox
              id="useLocation"
              checked={useCurrentLocation}
              onChange={(checked: any) => setUseCurrentLocation(checked as boolean)}
            />
            <label htmlFor="useLocation">Use my current location</label>
          </div>

          <div className="h-[400px] rounded-lg overflow-hidden">
            <Map
              markerPosition={markerPosition}
              useCurrentLocation={useCurrentLocation}
              setMarkerPosition={setMarkerPosition}
              setFormData={setFormData}
            />
          </div>
        </Card>

        <div style={{ marginTop: 24 }} className="flex w-full justify-end">
          <Button size='large' style={{ margin: '0 8px' }} className="px-8">
            Cancel
          </Button>

          <Button size='large' type="primary" onClick={() => setSuccesss(true)} className="px-10 border border-[#451f87]">
            Submit
          </Button>
        </div>
      </Form>
    </div>
  );
}
