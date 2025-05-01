import Link from 'next/link';
import { PiUsersThreeLight } from "react-icons/pi";
import { AiOutlineProduct } from "react-icons/ai";;
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
    icon: <PiUsersThreeLight />,
    label: (
      <Link href='/retailers'>
        Retailers
      </Link>
    )
  },

  {
    key: '4',
    icon: <PiUsersThreeLight />,
    label: (
      <p>
        Distributors
      </p>
    ),
    children: [
      {
        icon: <PiUsersThreeLight />,
        label: (
          <Link href='/distributors'>
            Distributors
          </Link>
        ),
      },
      {
        icon: <PiUsersThreeLight />,
        label: (
          <Link href='/distributors/agents'>
            Agents
          </Link>
        ),
      },
    ]
  },


  {
    key: '5',
    icon: <MdOutlineAddShoppingCart />,
    label: (
      <Link href='/orders'>
        Orders
      </Link>
    ),
  },
  {
    key: '6',
    icon: <CiSettings />,
    label: (
      <Link href='/settings'>
        Account Settings
      </Link>
    ),
  },
]
