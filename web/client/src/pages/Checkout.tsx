import { Card } from "@/components/ui/card";
import { CiCircleCheck } from "react-icons/ci";
import { Navbar } from "./Navbar";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { Footer } from "@/components/Footer";

const Checkout = () => {
  return (
    <div className="min-h-screen">
      <Navbar />
      <div className="mx-auto  flex justify-center text-center px-20 my-[10rem]">
        <div className="grid grid-cols-1 gap-4">
          <CiCircleCheck className="mx-auto text-emerald-900 text-9xl" />
          <p className="text-emerald-900 font-semibold text-3xl">
            THANK YOU FOR YOUR PURCHASE
          </p>
          <Link to="/product">
            <Button>
              Continue Shopping
            </Button>
          </Link>
        </div>


      </div>

      <Footer />
    </div>

  );
};

export default Checkout;
