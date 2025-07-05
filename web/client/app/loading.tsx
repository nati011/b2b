import { PiSpinner } from "react-icons/pi";
import Image from 'next/image'


export default function Loading() {
  return (
    <div className="min-h-screen min-w-screen flex justify-center items-center">
      <div className="grid-cols-1 align-center">
         <Image src="/logo.svg" width={100} height={100} alt="logo" className="mr-4" />
        <PiSpinner className="animate-spin text-primary text-6xl" />
      </div>

    </div>
  );
}
