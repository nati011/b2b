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
                            <p className="text-2xl font-semibold">{heading}</p>
                        </div>
                    )
                }
                {subheading && (
                    <p className="text-lg font-medium text-gray-700">{subheading}</p>
                )
                }
            </div>
            <div className="flex items-center gap-2">
                <Breadcrumb items={
                    [{
                        title: <IoHomeOutline />,
                        href: "/"
                    },
                    ...page
                    ]} />

            </div>
        </div>

    )
}

export default Heading