import { PiSpinner } from "react-icons/pi";

export default function Loading() {
  return (
    <div className="min-h-screen min-w-screen flex justify-center items-center">
      <PiSpinner className="animate-spin text-primary text-6xl" />
    </div>
  );
}
