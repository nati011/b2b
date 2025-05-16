
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ArrowLeft, Minus, Plus, X, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCart } from "@/contexts/CartContext";
import { toast } from "sonner";
import usePartnerStore from "@/lib/store/usePaymentStore";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { CheckoutRequest } from "@/lib/types";
import { PiSpinner } from "react-icons/pi";

const Cart = () => {
  const { items, removeFromCart, updateQuantity, getTotalPrice, clearCart } = useCart();
  const [partnerDialog, setPartnerDialog] = useState(false)
  const {
    loading,
    partners,
    fetchPaymentPartners,
    checkout
  } = usePartnerStore()

  const handleUpdateQuantity = (productId: number, currentQuantity: number, increment: boolean) => {
    const newQuantity = increment ? currentQuantity + 1 : currentQuantity - 1;
    if (newQuantity < 1) {
      removeFromCart(productId);
      toast.success("Product removed from cart");
      return;
    }
    updateQuantity(productId, newQuantity);
  };

  const handleClearCart = () => {
    clearCart();
    toast.success("Cart cleared successfully");
  };
  const [selectedPartner, setSelectedPartner] = useState({
    id: 0,
    name: ""
  })

  const viewPartnerDialog = () => {
    if (items.length === 0) {
      toast.error("Your cart is empty");
      return;
    }
    setPartnerDialog(true)
    fetchPaymentPartners()
  };

  const handleCheckout = () => {
    const request: CheckoutRequest = {
      retailer_id: 25,
      items: items,
      payment_partner_id: selectedPartner.id
    }
    checkout(request)
  }

  return (
    <>
      <div className="min-h-screen bg-white flex flex-col">
        <header className="fixed top-0 left-0 right-0 bg-white z-50 border-b border-gray-100">
          <nav className="container mx-auto px-4 py-4">
            <Link to="/" className="flex items-center gap-2 text-primary hover:opacity-80">
              <ArrowLeft className="w-5 h-5" />
              <span>Back to Store</span>
            </Link>
          </nav>
        </header>

        <main className="container mx-auto px-4 pt-24 pb-16 flex-1">
          <div className="max-w-4xl mx-auto">
            <div className="flex justify-between items-center mb-8">
              <h1 className="text-2xl font-semibold text-primary">Shopping Cart</h1>
              {items.length > 0 && (
                <Button
                  variant="destructive"
                  onClick={handleClearCart}
                  className="flex items-center gap-2"
                >
                  <Trash2 className="w-4 h-4" />
                  Clear Cart
                </Button>
              )}
            </div>

            {items.length === 0 ? (
              <div className="text-center py-16">
                <p className="text-gray-500 mb-4">Your cart is empty</p>
                <Link to="/">
                  <Button>Continue Shopping</Button>
                </Link>
              </div>
            ) : (
              <div className="space-y-8">
                <div className="space-y-4">
                  {items.map((item) => (
                    <div
                      key={item.id}
                      className="flex items-center gap-4 p-4 bg-white border border-gray-200 rounded-lg"
                    >
                      <div className="w-24 h-24 bg-secondary rounded-md overflow-hidden shrink-0">
                        <img
                          src={item.image}
                          alt={item.name}
                          className="w-full h-full object-cover"
                        />
                      </div>

                      <div className="flex-1 min-w-0">
                        <div className="flex justify-between items-start gap-4">
                          <div>
                            <h3 className="font-medium text-primary">{item.name}</h3>
                            <p className="text-sm text-gray-500">
                              ${item.price.toFixed(2)}
                            </p>
                          </div>
                          <button
                            onClick={() => removeFromCart(item.id)}
                            className="text-gray-400 hover:text-gray-500"
                          >
                            <X className="w-5 h-5" />
                          </button>
                        </div>

                        <div className="flex items-center gap-2 mt-4">
                          <Button
                            variant="outline"
                            size="icon"
                            onClick={() => handleUpdateQuantity(item.id, item.quantity, false)}
                          >
                            <Minus className="w-4 h-4" />
                          </Button>
                          <span className="w-12 text-center">{item.quantity}</span>
                          <Button
                            variant="outline"
                            size="icon"
                            onClick={() => handleUpdateQuantity(item.id, item.quantity, true)}
                          >
                            <Plus className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="font-medium text-primary">
                          ${(item.price * item.quantity).toFixed(2)}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="border-t border-gray-200 pt-4">
                  <div className="flex justify-between items-center py-2">
                    <span className="text-gray-600">Subtotal</span>
                    <span className="font-medium text-primary">
                      ${getTotalPrice().toFixed(2)}
                    </span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                    <span className="text-gray-600">Shipping</span>
                    <span className="text-gray-600">Calculated at checkout</span>
                  </div>
                  <div className="flex justify-between items-center py-4 border-t border-gray-200 mt-2">
                    <span className="text-lg font-medium text-primary">Total</span>
                    <span className="text-lg font-medium text-primary">
                      ${getTotalPrice().toFixed(2)}
                    </span>
                  </div>

                  <Button className="w-full mt-4" size="lg" onClick={viewPartnerDialog}>
                    Proceed to Checkout
                  </Button>
                </div>
              </div>
            )}
          </div>
        </main>
      </div>
      <Dialog open={partnerDialog} onOpenChange={setPartnerDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Select Payment Partner</DialogTitle>
          </DialogHeader>
          <Separator />
          <div className="grid md:grid-cols-2 gap-8">
            {
              (loading && partners.length == 0) ? (
                <>
                  <div className="flex flex-col space-y-3">
                    <Skeleton className="h-[200px] w-[200px]" />
                    <div className="space-y-2">
                      <Skeleton className="h-4 w-[50px]" />
                    </div>
                  </div>
                  <div className="flex flex-col space-y-3">
                    <Skeleton className="h-[200px] w-[200px]" />
                    <div className="space-y-2">
                      <Skeleton className="h-4 w-[50px]" />
                    </div>
                  </div>
                </>
              ) : (
                <>
                  {partners.map((partner, index) => (
                    <Card key={index} className="shadow-none  hover:shadow-md hover:shadow-gray-100 flex flex-col items-center text-center cursor-pointer" onClick={() => setSelectedPartner({ id: partner.id, name: partner.name })}>
                      <CardContent>
                        <img src={partner.icon} className="rounded-full" />
                        <p className="font-semibold text-lg">
                          {partner.name}
                        </p>
                      </CardContent>
                    </Card>
                  ))}

                </>
              )
            }

          </div>
          <Button className="w-full" onClick={handleCheckout} disabled={selectedPartner.id == 0}>
            {
              loading ? (
                <>
                  <PiSpinner className="animate-spin text-white" />
                  Loading
                </>
              ) :
                <>
                  Continue {
                    selectedPartner.name != "" && (
                      `with ${selectedPartner.name}`
                    )
                  }
                </>
            }

          </Button>
        </DialogContent>
      </Dialog>

    </>

  );
};

export default Cart;
