import { VscAccount } from "react-icons/vsc";
import Link from 'next/link';
import { BiTask } from 'react-icons/bi';
import { PiUsersThreeLight } from "react-icons/pi";
import { AiOutlineProduct } from "react-icons/ai";;
import {  TbReportSearch
 } from "react-icons/tb";
import { CiSettings } from 'react-icons/ci';
import { MdOutlineDashboard, MdOutlineAddShoppingCart } from 'react-icons/md';

export const Sidebar = [
  {
    key: '1',
    icon: <MdOutlineDashboard className="text-xl" />,
    label: (
      <Link href='/' >
        Dashboard
      </Link>
    ),
  },
  {
    key: '2',
    icon: <AiOutlineProduct />,
    label: (
      <Link href='/products'>
        Products
      </Link>
    ),
  },
  {
    key: '3',
    icon: <BiTask />,
    label: (
      <Link href='/grns'>
        GRN
      </Link>
    ),
  },
  {
    key: '4',
    icon: <PiUsersThreeLight />,
    label: (
      <Link href='/retailers'>
        Retailers
      </Link>
    )
  },

  {
    key: '5',
    icon: <PiUsersThreeLight />,
    label: (
      <Link href='/distributors'>
        Distributors
      </Link>
    )
  },


  {
    key: '6',
    icon: <MdOutlineAddShoppingCart />,
    label: (
      <Link href='/orders'>
        Orders
      </Link>
    ),
  },
    {
    key: '7',
    icon: <VscAccount />,
    label: (
      <Link href='/agents'>
        Agents
      </Link>
    ),
  },
  {
    key: '8',
    icon: <TbReportSearch/>,
    label:(
      <Link href='/report'>
       Report
      </Link>
    ),
  },
  {
    key: '9',
    icon: <CiSettings/>,
    label: (
      <Link href='/settings'>
       Account Settings
      </Link>
    ),
  },
]
