import { useLocation } from "react-router-dom";
import { useEffect } from "react";
import error from "@/public/404.svg"

const NotFound = () => {
  const location = useLocation();

  useEffect(() => {
    console.error(
      "404 Error: User attempted to access non-existent route:",
      location.pathname
    );
  }, [location.pathname]);

  return (
    <div className="w-screen h-screen items-center flex justify-center align-middle">
      <div className="grid grid-cols-1 gap-4 w-full align-middle justify-center place-items-center">
        <img src="/404.svg" alt="404" className="max-w-lg" />
        <a href={"/"} className="w-fit border-2 border-blue-900 text-blue-900 text-center rounded-md font-semibold px-10 py-2 hover:bg-blue-900/[5%]" >Back to home</a>

      </div>
    </div>
  );
};

export default NotFound;
