'use clint'
import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react";
import { IconType } from "react-icons";

interface Props {
    heading: string;
    icon: IconType;
    link: string;
}

export const NavLink: React.FC<Props> = ({ heading, icon: Icon, link }) => {
    const pathname = usePathname();
    const isActive = pathname === `/${link}`;

    return (
        <Link
            href={`/${link}`}
            className={`
        flex items-center group pl-8 py-3
        hover:border-r-4 hover:border-blue-900/[0.5] hover:bg-blue-900/[0.1]
        dark:hover:bg-blue-900/[10%]
        transition-all duration-200
        ${isActive ? 'border-r-4 border-blue-900 bg-blue-900/[0.1]' : ''}
      `}
        >
            <Icon
                className={`
          mr-2 text-xl sm:text-md
          ${isActive ? 'text-blue-900' : 'text-gray-500 group-hover:text-blue-900'}
        `}
            />
            <span className={`
        font-semibold
        ${isActive
                    ? 'text-blue-900 dark:text-blue-900'
                    : 'text-gray-500 group-hover:text-black dark:group-hover:text-white'
                }
      `}>
                {heading}
            </span>
        </Link>
    );
};
