"use client";
import { useState } from 'react';
import { MdOutlineDashboard, MdOutlineStorefront, MdOutlineSupervisorAccount, MdOutlineRealEstateAgent, MdOutlineAddShoppingCart } from "react-icons/md";
import { LiaFileInvoiceDollarSolid, LiaCartArrowDownSolid } from "react-icons/lia";
import { IoIosLogOut } from "react-icons/io";
import { CiBank, CiMenuFries } from "react-icons/ci";
import Link from "next/link";
import { NavLink } from "@/app/components/navlink";
import { useRouter } from 'next/navigation';
import { AiOutlineProduct } from "react-icons/ai";
import { PiUsersThreeLight } from "react-icons/pi";
import { CiSettings } from 'react-icons/ci';

// import {
//   Accordion,
//   AccordionContent,
//   AccordionItem,
//   AccordionTrigger,
// } from "@/components/ui/accordion"


const SideBar = () => {
  const [open, setOpen] = useState(false);
  const router = useRouter();
  const navigationLinks = [
    { heading: "Dashboard", link: "", icon: MdOutlineDashboard },
    { heading: "Properties", link: "properties", icon: MdOutlineStorefront },
    { heading: "Property Loan", link: "home-loan", icon: LiaFileInvoiceDollarSolid },
    { heading: "Loaners", link: "loaners", icon: CiBank },
    { heading: "Auctions", link: "auctions", icon: LiaCartArrowDownSolid },
    { heading: "Customers", link: 'customers', icon: MdOutlineSupervisorAccount },
    { heading: "Tours", link: "requested-tours", icon: MdOutlineSupervisorAccount },
    { heading: "Real Estates", link: "homeowners", icon: MdOutlineRealEstateAgent }
  ];

  const SidebarContent = () => (
    <>
      <Link
        href="/"
        className="flex flex-col items-center justify-center w-full font-bold text-lg my-8 sm:mb-4 sm:p-2"
      >
        <p className="hidden lg:block dark:text-white font-semibold lg:text-2xl">
          Efoyeta Market
        </p>
      </Link>

      <div className="list mt-4 grid grid-cols-1 gap-3">
        <NavLink
          heading={"Dashboard"}
          link={""}
          icon={MdOutlineDashboard}
        />
        {/* <Accordion type="single" collapsible>
          <AccordionItem value="Product">
            <AccordionTrigger>Products</AccordionTrigger>
            <AccordionContent> */}
        <NavLink
          heading={"Product"}
          link={"products"}
          icon={AiOutlineProduct}
        />
        {/* </AccordionContent>
          </AccordionItem>
        </Accordion> */}
        <NavLink
          heading={"Retailer"}
          link={"retailers"}
          icon={PiUsersThreeLight}
        />
        <NavLink
          heading={"Distributor"}
          link={"/distributor"}
          icon={PiUsersThreeLight}
        />
        <NavLink
          heading={"Orders"}
          link={"/orders"}
          icon={MdOutlineAddShoppingCart}
        />

        <div
          className="flex items-center group pl-8 py-3 gap-2 cursor-pointer"
          onClick={() => router.push('/login')}
        >
          <IoIosLogOut className="group-hover:text-red-500 text-2xl sm:text-md" />
          <span className="group-hover:text-black text-gray-500 dark:group-hover:text-white font-semibold">
            Logout
          </span>
        </div>
      </div>
    </>
  );

  return (
    <div className="print:hidden">
      {/* Mobile Menu Button */}
      <div
        className='block lg:hidden absolute left-1 z-50 cursor-pointer bg-transparent dark:text-white p-4 rounded-full'
        onClick={() => setOpen(!open)}
      >
        <CiMenuFries className='text-xl' />
      </div>

      {/* Mobile Sidebar */}
      {open && (
        <>
          <div
            className="block lg:hidden fixed top-0 left-0 h-screen w-screen bg-black bg-opacity-50 z-40"
            onClick={() => setOpen(false)}
          />
          <aside className="block lg:hidden fixed top-0 left-0 bg-white dark:bg-black dark:text-white h-screen w-64 z-50">
            <SidebarContent />
          </aside>
        </>
      )}

      {/* Desktop Sidebar */}
      <aside className="hidden lg:block h-screen overflow-y-auto print:hidden border-r-[1.5px] border-gray-200 dark:bg-black dark:text-white bg-white">
        <SidebarContent />
      </aside>
    </div>
  );
};

export default SideBar;

