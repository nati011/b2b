'use client'
import { Card, Col, Row } from "antd";

import { BiUser } from "react-icons/bi";
import { MdAddShoppingCart } from "react-icons/md";

import dynamic from 'next/dynamic'
// import TableLayout from "@/components/table";
// import { column, data } from "./data";

const Chart = dynamic(() => import('react-apexcharts'), { ssr: false });


export default function Home() {
  const chart1Options =  {
    chart: {
      id: "basic-line",
      height: '4px' ,
    },

    colors:['#FDAD15', '#29C66F'],
    stroke: {
      curve: "smooth" as "smooth",
    },
    xaxis: {
      categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999]
    }
  }
  const chart1Series =  [
      {
        name: "Pending",
        data: [30, 40, 45, 50, 49, 60, 70, 91]
      },
      {
        name: "Completed",
        data: [50, 90, 30, 10, 49, 65, 60, 71]
      }
    ]

    const   dountOptions =  {
      labels: ['A', 'B', 'C','D', 'E'],
      
      theme: {
        monochrome: {
          enabled: true,
          color: '#1C40CA',
          shadeTo: 'light' as "light",
        },
        plotOptions: {
          pie: {
            donut: {
              labels: {
                show: true,
                name: {
                  show: true,
                },
                total: {
                  show: true,
                  label: 'Total',
                  formatter: () => 'Total'
                }
              }
            }
          }
        }
      },
      noData: {
        text: 'No Data Available',
        align: 'center' as "center",
        verticalAlign: 'middle' as "middle",
        offsetX: 0,
        offsetY: 0,
        style: {
          color: "black",
          fontSize: '14px',
        }
      },
    };
    

  return (
    <div className="flex flex-col gap-4">
        <Row gutter={16}>
    <Col xs={24} xl={8} className="md:mb-0 mb-2">
      <Card>
        <div className="flex justify-between">
     
      <div className="">
      <div className="flex gap-1 items-center">
        <h3 className="font-bold text-2xl">
          638
        </h3>
        <p className="text-[#09274b]">+21.01</p>
        </div>
        <p className="text-[#89868D]">Total Retailers</p>
      </div>
      <BiUser className="text-3xl"/>
        </div>
      </Card>
    </Col>
    <Col xs={24} xl={8}  className="md:mb-0 mb-2">
    <Card>
        <div className="flex justify-between">
     
      <div className="">
      <div className="flex gap-1 items-center">
        <h3 className="font-bold text-2xl">
          123
        </h3>
        <p className="text-[#09274b]">+21.01</p>
        </div>
        <p className="text-[#89868D]">Total Distributors</p>
      </div>
      <MdAddShoppingCart className="text-3xl" />
        </div>
      </Card>
    </Col>
    <Col xs={24} xl={8} >
    <Card>
        <div className="flex justify-between">
     
      <div className="">
      <div className="flex gap-1 items-center">
        <h3 className="font-bold text-2xl">
          123
        </h3>
        <p className="text-[#09274b]">+21.01</p>
        </div>
        <p className="text-[#89868D]">Total Orders</p>
      </div>
      <MdAddShoppingCart className="text-3xl" />
        </div>
      </Card>
    </Col>
    
  </Row>
  <Row gutter={16}>
    <Col  xs={24} xl={16} className="md:mb-0 mb-2">
      <Card>
        <p className="text-xl font-semibold">Orders</p>
       <Chart options={chart1Options} series={chart1Series} height={200}/>
      </Card>
    </Col>
    <Col  xs={24} xl={8}>
    <Card>
    <p className="text-xl font-semibold">Yeilds</p>
    <div className="flex justify-center item-center w-full">
    <Chart type="donut" options={dountOptions} series={[10,20,20,60,20]} width={290}/>
    </div>
    
      </Card>
    </Col>
  </Row>
    </div>
  );
}
