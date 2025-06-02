import { Card, Table, Input, Button, Breadcrumb } from 'antd'
import React from "react"
import { CiFilter, CiSearch } from "react-icons/ci"
import type {  TableProps } from 'antd';
import Link from 'next/link';
import { GoPlus } from 'react-icons/go';
import { HomeOutlined } from '@ant-design/icons';

interface Page {
  title: string;
  href: string;
}


interface Props {
  column: any
  data: any
  heading: string
  subheading?: string
  pages: Page[]
  buttonURL?: string
  loading?: boolean
}

const rowSelection: TableProps<any>['rowSelection'] = {
  onChange: (selectedRowKeys: React.Key[], selectedRows: any[]) => {
    console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
  },
  getCheckboxProps: (record: any) => ({
    disabled: record.name === 'Disabled User',
    name: record.name,
  }),
};

const TableLayout: React.FC<Props> = ({ column, data, heading, subheading,
  pages,
  buttonURL,


}) => {

  return (
    <div className="">
      <div className="flex flex-col md:flex-row justify-between mb-4">
        <p className="font-semibold text-lg text-gray-900">
          {heading}
        </p>
        <div className="flex items-center gap-2">
          <Breadcrumb items={
            [{
              title: <HomeOutlined />,
              href:"/"
            },
            // @ts-expect-error
            ...pages
            ]} />

        </div>
      </div>
      <Card>
        <div className="flex flex-col md:flex-row justify-between gap-4 mb-4">
          {
            buttonURL && (
              <Link href={buttonURL}>
              <Button icon={<GoPlus />} >
                <p className="hidden md:block">  Register {heading}</p>
  
              </Button>
            </Link>
  
            )
          }
         
          <div className={`${buttonURL ? 'max-w-sm' :"w-full"} flex gap-1`}>
            <Input placeholder="Search  " prefix={<CiSearch />} />
            <Button icon={<CiFilter />}>
              <p className="hidden md:block">Filter</p>
            </Button>

          </div>
        </div>

        <Table<any>
          bordered
          columns={column}
          dataSource={data}
          className="md:max-w-full max-w-md overflow-x-auto"
        />

      </Card>
    </div>


  )
}

export default TableLayout