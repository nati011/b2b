import React from "react"
import Link from 'next/link'

import { IoHomeOutline } from "react-icons/io5"
import { Breadcrumb } from "antd"

interface Props {
    page: any[],
    heading?: string
    subheading?: string
}

const Heading: React.FC<Props> = ({ page, heading, subheading }) => {
    return (
        <div className="mb-4 print:hidden flex justify-between items-center">
            <div className="">
                {
                    heading && (
                        <div className="flex">
                            <p className="text-xl font-semibold text-black">{heading}</p>
                        </div>
                    )
                }
                {subheading && (
                    <p className="text-md font-medium text-gray-700">{subheading}</p>
                )
                }
            </div>
            <div className="flex items-center gap-2 text-cyan-900">
                <div className="breadcrumb flex items-center text-sm sm:mb-2 my-4 ">
                    <Link href="/">
                        <IoHomeOutline />
                    </Link>
                    {page.map((p, index) => (
                        <div key={index}>
                            <span className="mx-2">/</span>
                            <Link className='' href={`${p.href}`}>{p.title}</Link>
                        </div>
                    ))}
                </div>

            </div>
        </div>

    )
}

export default Heading