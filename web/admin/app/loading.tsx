import { PiPlantThin, PiSpinner } from "react-icons/pi";

export default function Loading() {
  return (
    <div className="min-h-screen min-w-screen flex justify-center items-center">
      <div className="grid-cols-1 align-middle">
        <PiSpinner className="animate-spin text-[#1A1F2C] text-6xl" />
      </div>

    </div>
  );
}
