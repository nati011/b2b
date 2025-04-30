import { Card, Table, Input, Button, Radio, Divider} from 'antd'
import React, { useState } from "react"
import { CiFilter, CiSearch } from "react-icons/ci"
import type { TableColumnsType, TableProps } from 'antd';


interface Props{
    column: any
    data: any
    heading: string
    subheading?:string
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
  
const TableLayout: React.FC<Props> = ({column, data, heading, subheading}) => {

    return(
        <Card bordered={true}>
            <div className=" mb-2 flex justify-between">
                <p className="font-semibold text-xl w-1/2">
                {heading}
                </p>
                <div className="max-w-sm flex gap-2">
                <Input placeholder="Search  " prefix={<CiSearch/>}/>
                <Button  icon={<CiFilter/>}>
                    Filter
                </Button>

                </div>
                
            </div>
            <div>
     
      <Divider />
      <Table<any>
        className="md:max-w-full max-w-md overflow-x-auto"
        rowSelection={{ type:'checkbox', ...rowSelection }}
        columns={column}
        dataSource={data}
      />
    </div>
           
        </Card>

    )
}

export default TableLayout