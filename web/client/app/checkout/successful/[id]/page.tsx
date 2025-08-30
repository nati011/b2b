"use client";
import { useEffect } from "react";
import Image from "next/image";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CheckCircle } from "lucide-react";
import { useParams } from "next/navigation";
import Link from "next/link";
import useInvoiceStore from "@/lib/store/useInvoiceStore";
import ThankYouSkeleton from "@/components/ThankYouSkeleton";
import { InvoiceCard } from "@/components/invoice/invoiceCard";

const ThankYou = () => {
  const routeParam = useParams<{ id: string }>();
  const { invoice, loading, verified, fetchInvoice, paymentVerification } =
    useInvoiceStore();

  useEffect(() => {
    const verify = async () => {
      try {
        await paymentVerification();
      } catch (error) {}
    };
    verify();
  }, []);

  useEffect(() => {
    if (verified) {
      fetchInvoice(parseInt(routeParam.id));
    }
  }, [verified]);

  const handlePrint = () => {
    window.print();
  };

  if (loading) {
    return <ThankYouSkeleton />;
  }

  if (!invoice) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Order not found
          </h1>
          <Link href="/">
            <Button>Return to Store</Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <main className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          {verified ? (
            <>
              {/* Success Message */}
              <Card className="mb-8 border-green-200 bg-green-50 shadow-none rounded-sm print:hidden">
                <CardContent className="p-6 text-center">
                  <CheckCircle className="w-16 h-16 text-green-600 mx-auto mb-4" />
                  <h1 className="text-3xl font-bold text-green-800 mb-2">
                    Thank You for Your Purchase!
                  </h1>
                  <p className="text-green-700">
                    Your order has been confirmed and will be processed shortly.
                  </p>
                </CardContent>
              </Card>
              {
                invoice && (
                  <>
 {/* Invoice */}
 <InvoiceCard loading={loading} invoice={invoice} />

 {/* Additional Actions */}
 <div className="mt-8 text-center print:hidden">
   <div className="flex flex-col sm:flex-row gap-4 justify-center">
     <Link href="/">
       <Button variant="outline" className="w-full sm:w-auto">
         Continue Shopping
       </Button>
     </Link>
     <Button
       className="w-full sm:w-auto"
       onClick={() => {
         handlePrint();
       }}
     >
       Print Invoice
     </Button>
   </div>
 </div>
</>
                )
              }
              </>
             
          ) : (
            <>
              <div className="w-full flex flex-col gap-6 items-center justify-center text-center py-20">
                <Image
                  src={"/failed-payment.png"}
                  width={300}
                  height={400}
                  alt="failed transaction"
                />
                <div>
                  <h4 className="text-2xl font-semibold text-gray-900 mb-2">
                    Failed to complete payment.
                  </h4>
                </div>
                <div className="mt-8 text-center print:hidden">
                  <div className="flex flex-col sm:flex-row gap-4 justify-center">
                    <Link href="/">
                      <Button variant="outline" className="w-full sm:w-auto">
                        Continue Shopping
                      </Button>
                    </Link>
                    <Link href="/orders">
                      <Button className="w-full sm:w-auto">View Orders</Button>
                    </Link>
                  </div>
                </div>
              </div>
            </>
          )}
        </div>
      </main>
    </div>
  );
};

export default ThankYou;
