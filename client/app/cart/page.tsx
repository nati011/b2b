"use client";
import { useEffect, useCallback } from "react";
import { Minus, Plus, X, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { CheckoutRequest } from "@/lib/types";
import useCartStore from "@/lib/store/useCartStore";
import useOrderStore from "@/lib/store/useOrderStore";
import { useRouter } from "next/navigation";

const Cart = () => {
  const {
    cartItems,
    totalItems,
    totalPrice,
    addCartItems,
    clearCart,
    removeCartItems,
  } = useCartStore();

  const { checkout, error, loading: checkoutLoading } = useOrderStore();

  const router = useRouter();

  const handleClearCart = () => {
    clearCart();
    toast.success("Cart cleared successfully");
  };

  const handleCheckout = async () => {
    if (cartItems.length === 0) {
      toast.error("Your cart is empty");
      return;
    }
    const request: CheckoutRequest = {
      customer_id: 25,
      items: cartItems.map((cartItem) => ({
        id: cartItem.id,
        quantity: cartItem.quantity,
      })),
    };
    try {
      await checkout(request);
      clearCart();
      toast.success("Order placed successfully");
      router.push("/orders");
    } catch (checkoutError: any) {
      toast.error(checkoutError?.message || "Failed to place order");
    }
  };

  const handleBackToShopping = () => {
    router.push("/product");
  };

  const handleQuantityDecrease = useCallback((item: any) => {
    addCartItems(item, -1);
  }, [addCartItems]);

  const handleQuantityIncrease = useCallback((item: any) => {
    addCartItems(item, 1);
  }, [addCartItems]);

  useEffect(() => {
    if (error != null) {
      toast.error(error);
    }
  }, [error]);

  return (
    <div className="min-h-screen bg-gray-50">
        <main className="container px-4 py-8 mx-auto">
          <Card className="max-w-6xl px-4 mx-auto shadow-none">
            {cartItems.length === 0 ? (
              <div className="text-center py-16 bg-white rounded-lg border border-gray-100">
                <div className="w-24 h-24 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
                  <svg
                    className="w-12 h-12 text-gray-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M3 3h2l.4 2M7 13h10l4-8H5.4m0 0L7 13m0 0l-2.5 5M7 13l-2.5 5m4.5 0h9"
                    />
                  </svg>
                </div>
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  Your cart is empty
                </h3>
                <p className="text-gray-500 mb-6">
                  Start shopping to add items to your cart
                </p>
                <Button onClick={handleBackToShopping} className="px-6">
                  Continue Shopping
                </Button>
              </div>
            ) : (
              <div className="grid lg:grid-cols-6 gap-8">
                <div className="lg:col-span-4 space-y-4">
                  <div className="flex justify-between items-center">
                    <h2 className="text-xl font-semibold text-gray-900">
                      Cart Items
                    </h2>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleClearCart}
                      className="text-red-600 border-red-200 hover:bg-red-50 hover:border-red-300"
                    >
                      <Trash2 className="w-4 h-4 mr-2" />
                      Clear Cart
                    </Button>
                  </div>

                  <div className="space-y-3">
                    {cartItems.map((item) => {
                      const itemTotal = item.price * item.quantity;
                      return (
                      <div
                        key={item.id}
                        className="bg-white rounded-lg  border p-4 border-gray-100 transition-shadow"
                      >
                        <div className="flex gap-4">
                          <div className="w-20 h-20 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0">
                            <img
                              src={item.image}
                              alt={item.name}
                              className="w-full h-full object-cover"
                            />
                          </div>

                          <div className="flex-1 min-w-0">
                            <div className="flex justify-between items-start">
                              <div>
                                <h3 className="font-medium text-gray-900 line-clamp-2">
                                  {item.name}
                                </h3>
                                <p className="text-sm text-gray-500 mt-1">
                                  {item.price.toLocaleString()} ETB each
                                </p>
                              </div>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => removeCartItems(item)}
                                className="text-gray-400 hover:text-red-500 hover:bg-red-50"
                              >
                                <X className="w-4 h-4" />
                              </Button>
                            </div>

                            <div className="flex items-center justify-between mt-4">
                              <div className="flex items-center gap-2">
                                <Button
                                  variant="outline"
                                  size="icon"
                                  onClick={() => handleQuantityDecrease(item)}
                                  className="h-8 w-8"
                                >
                                  <Minus className="w-3 h-3" />
                                </Button>
                                <span className="w-12 text-center font-medium">
                                  {item.quantity}
                                </span>
                                <Button
                                  variant="outline"
                                  size="icon"
                                  onClick={() => handleQuantityIncrease(item)}
                                  className="h-8 w-8"
                                >
                                  <Plus className="w-3 h-3" />
                                </Button>
                              </div>
                              <div className="text-right">
                                <p className="font-semibold text-gray-900">
                                  {itemTotal.toLocaleString()} ETB
                                </p>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      );
                    })}
                  </div>
                </div>

                <div className="lg:col-span-2">
                  <div className="bg-white rounded-lg border border-gray-100 p-6 sticky top-8">
                    <h3 className="text-lg font-medium text-gray-900 mb-4">
                      Order Summary
                    </h3>

                    <div className="space-y-3 mb-4">
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">
                          Subtotal ({totalItems} items)
                        </span>
                        <span className="font-medium">{totalPrice.toLocaleString()} ETB</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="font-medium">
                          Delivery Estimation:
                        </span>
                        <span className="text-gray-600">30 minutes to 1 hour</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="font-medium">
                          Delivery Fee:
                        </span>
                        <span className="text-gray-600">Depends on the location</span>
                      </div>
                    </div>

                    <Separator className="my-4" />

                    <div className="flex justify-between items-center mb-6">
                      <span className="text-lg font-semibold text-gray-900">
                        Total
                      </span>
                      <span className="text-lg font-semibold text-gray-900">
                        {totalPrice.toLocaleString()} ETB
                      </span>
                    </div>

                    <Button
                      className="w-full"
                      size="lg"
                      onClick={handleCheckout}
                      disabled={cartItems.length === 0 || checkoutLoading}
                    >
                      {checkoutLoading ? "Placing Order..." : "Place Order"}
                    </Button>
                  </div>
                </div>
              </div>
            )}
          </Card>
        </main>
      </div>
  );
};

export default Cart;
