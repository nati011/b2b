"use client";
import { useEffect, useCallback, useState } from "react";
import { Minus, Plus, X, Trash2, Copy, CheckCircle2, ArrowLeft, Check } from "lucide-react";
import { PiTelegramLogo } from "react-icons/pi";
import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import useCartStore from "@/lib/store/useCartStore";
import { useRouter } from "next/navigation";
import { getCustomerByEmail } from "@/app/actions/customer";
import { Login, getAccessToken } from "@/app/actions/auth";
import { createOrder, CreateOrderRequest } from "@/app/actions/orders";
import { cn } from "@/lib/utils";

type CheckoutStep = 'payment' | 'review';

const steps: { id: CheckoutStep; label: string }[] = [
  { id: 'review', label: 'Review' },
  { id: 'payment', label: 'Payment' },
];

const Cart = () => {
  const {
    cartItems,
    totalItems,
    totalPrice,
    addCartItems,
    clearCart,
    removeCartItems,
  } = useCartStore();

  const [customerId, setCustomerId] = useState<number | null>(null);
  const [customerData, setCustomerData] = useState<any>(null);
  const [isLoadingCustomer, setIsLoadingCustomer] = useState(false);
  const [isPlacingOrder, setIsPlacingOrder] = useState(false);
  const [currentStep, setCurrentStep] = useState<CheckoutStep>('review');
  const [copiedField, setCopiedField] = useState<string | null>(null);

  const router = useRouter();
  
  const currentStepIndex = steps.findIndex(s => s.id === currentStep);

  // Prefetch routes on mount
  useEffect(() => {
    router.prefetch("/orders");
    router.prefetch("/product");
  }, [router]);

  // Load customer ID from localStorage or fetch by email
  useEffect(() => {
    const loadCustomerId = async () => {
      if (typeof window === 'undefined') return;

      try {
        // First, try to get customer_id from localStorage
        const storedCustomerId = localStorage.getItem('customer_id');
        const storedCustomerData = localStorage.getItem('customer_data');
        const userEmail = localStorage.getItem('user_email') || localStorage.getItem('saved_email');

        if (storedCustomerId && storedCustomerData) {
          try {
            const customerIdNum = parseInt(storedCustomerId, 10);
            const customerDataParsed = JSON.parse(storedCustomerData);
            
            // Verify the stored customer email matches the logged-in user
            if (userEmail && customerDataParsed.email?.toLowerCase() === userEmail.toLowerCase() && !isNaN(customerIdNum)) {
              setCustomerId(customerIdNum);
              setCustomerData(customerDataParsed);
              return;
            } else {
              // Clear mismatched data
              localStorage.removeItem('customer_id');
              localStorage.removeItem('customer_data');
            }
          } catch (parseError) {
            console.error('Error parsing stored customer data:', parseError);
            localStorage.removeItem('customer_id');
            localStorage.removeItem('customer_data');
          }
        }

        // If customer_id not in localStorage, fetch by email
        if (userEmail && !customerId) {
          setIsLoadingCustomer(true);
          try {
            const customer = await getCustomerByEmail(userEmail);
            if (customer) {
              setCustomerId(customer.id);
              setCustomerData(customer);
              // Store in localStorage for future use
              localStorage.setItem('customer_id', customer.id.toString());
              localStorage.setItem('customer_data', JSON.stringify(customer));
            } else {
              // Customer not found - this is OK, user can complete profile later
              console.log('Customer profile not found for email:', userEmail);
            }
          } catch (error) {
            console.error('Error fetching customer:', error);
            // Don't show error toast here - let checkout handle it
          } finally {
            setIsLoadingCustomer(false);
          }
        }
      } catch (error) {
        console.error('Error loading customer ID:', error);
        setIsLoadingCustomer(false);
      }
    };

    loadCustomerId();
  }, [customerId]);

  const handleClearCart = () => {
    clearCart();
    toast.success("Cart cleared successfully");
  };

  const handleNext = async () => {
    const nextIndex = currentStepIndex + 1;
    if (nextIndex < steps.length) {
      setCurrentStep(steps[nextIndex].id);
    } else {
      // Place order
      await handlePlaceOrder();
    }
  };

  const handleBack = () => {
    const prevIndex = currentStepIndex - 1;
    if (prevIndex >= 0) {
      setCurrentStep(steps[prevIndex].id);
    } else {
      router.push("/product");
    }
  };

  const handlePlaceOrder = async () => {
    if (cartItems.length === 0) {
      toast.error("Your cart is empty");
      return;
    }

    // Ensure we have a valid access token using stored credentials
    if (typeof window !== 'undefined') {
      const accessToken = getAccessToken();
      const userEmail = localStorage.getItem('user_email') || localStorage.getItem('saved_email');
      const savedPassword = localStorage.getItem('saved_password');

      // If no token or credentials, redirect to login
      if (!accessToken && (!userEmail || !savedPassword)) {
        toast.error("Please log in to place an order");
        router.push("/login");
        return;
      }

      // If no token but we have credentials, refresh the token
      if (!accessToken && userEmail && savedPassword) {
        try {
          await Login({ email: userEmail, password: savedPassword });
          toast.success("Session refreshed");
        } catch (error: any) {
          console.error('Error refreshing token:', error);
          toast.error("Failed to authenticate. Please log in again.");
          router.push("/login");
          return;
        }
      }
    }

    // Get customer data from state or localStorage
    let finalCustomerId = customerId;
    let finalCustomerData = customerData;

    // If not in state, try localStorage first
    if (!finalCustomerId || !finalCustomerData) {
      const storedCustomerId = typeof window !== 'undefined' ? localStorage.getItem('customer_id') : null;
      const storedCustomerData = typeof window !== 'undefined' ? localStorage.getItem('customer_data') : null;
      
      if (storedCustomerId && storedCustomerData) {
        try {
          const parsedId = parseInt(storedCustomerId, 10);
          const parsedData = JSON.parse(storedCustomerData);
          
          // Verify the data is valid
          if (!isNaN(parsedId) && parsedData && parsedData.id) {
            finalCustomerId = parsedId;
            finalCustomerData = parsedData;
            // Update state
            setCustomerId(parsedId);
            setCustomerData(parsedData);
          }
        } catch (e) {
          console.error('Error parsing stored customer data:', e);
          // Clear invalid data
          if (typeof window !== 'undefined') {
            localStorage.removeItem('customer_id');
            localStorage.removeItem('customer_data');
          }
        }
      }
    }

    // If still no customer_id, try to fetch by email
    if (!finalCustomerId || !finalCustomerData) {
      const userEmail = typeof window !== 'undefined' 
        ? (localStorage.getItem('user_email') || localStorage.getItem('saved_email'))
        : null;
      
      if (!userEmail) {
        toast.error("Please log in to place an order");
        router.push("/login");
        return;
      }

      // Try to fetch customer one more time
      setIsLoadingCustomer(true);
      try {
        const customer = await getCustomerByEmail(userEmail);
        if (customer) {
          finalCustomerId = customer.id;
          finalCustomerData = customer;
          setCustomerId(customer.id);
          setCustomerData(customer);
          localStorage.setItem('customer_id', customer.id.toString());
          localStorage.setItem('customer_data', JSON.stringify(customer));
        } else {
          // Customer profile doesn't exist - redirect to profile page
          toast.error("Please complete your profile before placing an order", {
            description: "We need your profile information to process your order.",
          });
          setIsLoadingCustomer(false);
          router.push("/profile");
          return;
        }
      } catch (error: any) {
        console.error('Error fetching customer:', error);
        toast.error("Unable to load your profile information", {
          description: "Please try again or complete your profile.",
        });
        setIsLoadingCustomer(false);
        router.push("/profile");
        return;
      } finally {
        setIsLoadingCustomer(false);
      }
    }

    // Final check - ensure we have valid customer data
    if (!finalCustomerId || !finalCustomerData || !finalCustomerData.id) {
      toast.error("Customer information is incomplete", {
        description: "Please complete your profile to continue.",
      });
      router.push("/profile");
      return;
    }

    // Get user info for snapshots
    const userEmail = typeof window !== 'undefined' 
      ? (localStorage.getItem('user_email') || localStorage.getItem('saved_email'))
      : null;
    const userName = typeof window !== 'undefined' 
      ? localStorage.getItem('user_name')
      : null;

    // Prepare order items (exactly matching mobile format)
    const orderItems = cartItems.map(item => ({
      product_id: item.id,
      quantity: item.quantity,
      price: item.price,
    }));

    // Prepare customer snapshot (exactly matching mobile format)
    const customerSnapshot = JSON.stringify({
      id: finalCustomerData.id,
      full_name: finalCustomerData.full_name || userName || '',
      email: finalCustomerData.email || userEmail || '',
      phone_number: finalCustomerData.phone_number || '',
      city: finalCustomerData.city || '',
      region: finalCustomerData.region || '',
      woreda: finalCustomerData.woreda || '',
      status: finalCustomerData.status || 'active',
    });

    // Prepare shipping address snapshot (exactly matching mobile format)
    const shippingAddressSnapshot = JSON.stringify({
      full_name: finalCustomerData.full_name || userName || '',
      email: finalCustomerData.email || userEmail || '',
      phone_number: finalCustomerData.phone_number || '',
      street_address: finalCustomerData.woreda || '',
      city: finalCustomerData.city || '',
      region: finalCustomerData.region || '',
    });

    // Create order request (exactly matching mobile format)
    const orderRequest: CreateOrderRequest = {
      customer_id: finalCustomerId,
      status: 'PENDING',
      payment_status: 'PENDING',
      delivery_status: 'PENDING',
      confirmation_status: 'PENDING',
      total: totalPrice,
      customer_snapshot: customerSnapshot as any,
      shipping_address_snapshot: shippingAddressSnapshot as any,
      items: orderItems,
    };

    // Log the request to verify it matches mobile format
    console.log('📤 Order Request (matching mobile):', JSON.stringify(orderRequest, null, 2));

    setIsPlacingOrder(true);
    
    try {
      const orderResponse = await createOrder(orderRequest);
      
      toast.success("Order placed successfully!", {
        description: `Order #${orderResponse.id} has been created.`,
      });
      
      clearCart();
      // Redirect immediately to orders page
      router.push("/orders");
    } catch (orderError: any) {
      // If 401 error, try to refresh token and retry
      if (orderError?.response?.status === 401 || orderError?.message?.includes('401')) {
        const userEmail = typeof window !== 'undefined' 
          ? (localStorage.getItem('user_email') || localStorage.getItem('saved_email'))
          : null;
        const savedPassword = typeof window !== 'undefined' 
          ? localStorage.getItem('saved_password')
          : null;

        if (userEmail && savedPassword) {
          try {
            await Login({ email: userEmail, password: savedPassword });
            // Retry the order
            const retryResponse = await createOrder(orderRequest);
            toast.success("Order placed successfully!", {
              description: `Order #${retryResponse.id} has been created.`,
            });
            clearCart();
            // Redirect immediately to orders page
            router.push("/orders");
            return;
          } catch (retryError: any) {
            toast.error("Session expired. Please log in again.");
            router.push("/login");
            return;
          }
        }
      }
      toast.error(orderError?.message || "Failed to place order");
    } finally {
      setIsPlacingOrder(false);
    }
  };

  const handleBackToShopping = () => {
    router.push("/product");
  };

  const copyToClipboard = (text: string, field: string) => {
    navigator.clipboard.writeText(text);
    setCopiedField(field);
    toast.success("Copied to clipboard");
    setTimeout(() => setCopiedField(null), 2000);
  };

  // Show completion screen

  const handleQuantityDecrease = useCallback((item: any) => {
    addCartItems(item, -1);
  }, [addCartItems]);

  const handleQuantityIncrease = useCallback((item: any) => {
    addCartItems(item, 1);
  }, [addCartItems]);


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
              <div className="space-y-6">
                {/* Header with Back Button */}
                <div className="flex items-center gap-4 pt-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={handleBack}
                    className="h-9 w-9"
                  >
                    <ArrowLeft className="h-5 w-5" />
                  </Button>
                  <h1 className="text-2xl font-semibold">Checkout</h1>
                </div>

                {/* Step Indicator */}
                <div className="flex items-center justify-center gap-4 pb-4">
                  {steps.map((step, index) => {
                    const isActive = currentStepIndex === index;
                    const isCompleted = currentStepIndex > index;
                    const canNavigate = isCompleted || isActive;

                    return (
                      <div key={step.id} className="flex items-center gap-2">
                        <button
                          onClick={() => {
                            if (canNavigate) {
                              setCurrentStep(step.id);
                            }
                          }}
                          disabled={!canNavigate}
                          className={cn(
                            "flex items-center gap-2 transition-all rounded-lg px-2 py-1",
                            canNavigate && "cursor-pointer hover:bg-accent hover:opacity-90",
                            !canNavigate && "cursor-not-allowed opacity-50"
                          )}
                        >
                          <div
                            className={cn(
                              "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-all",
                              isActive && "bg-primary text-primary-foreground ring-2 ring-primary ring-offset-2",
                              isCompleted && "bg-primary text-primary-foreground",
                              !isActive && !isCompleted && "bg-secondary text-muted-foreground"
                            )}
                          >
                            {isCompleted ? <Check className="w-4 h-4" /> : index + 1}
                          </div>
                          <span
                            className={cn(
                              "text-sm font-medium transition-colors",
                              isActive ? "text-foreground" : "text-muted-foreground",
                              canNavigate && "hover:text-foreground"
                            )}
                          >
                            {step.label}
                          </span>
                        </button>
                        {index < steps.length - 1 && (
                          <div className={cn(
                            "w-12 h-0.5 transition-colors",
                            isCompleted ? "bg-primary" : "bg-border"
                          )} />
                        )}
                      </div>
                    );
                  })}
                </div>

                {/* Step Content */}
                <div className="bg-white rounded-lg border border-gray-100 p-6">
                  {currentStep === 'payment' && (
                    <div className="space-y-6">
                      <h2 className="text-xl font-semibold mb-4">Payment Details</h2>

                      {/* Payment Instructions */}
                      <div className="bg-primary/5 border border-primary/20 rounded-lg p-4 space-y-2">
                        <div className="flex items-start gap-2">
                          <CheckCircle2 className="w-4 h-4 text-primary mt-0.5 shrink-0" />
                          <div className="space-y-1">
                            <p className="text-xs font-medium text-foreground">Payment Instructions</p>
                            <ul className="text-xs text-muted-foreground space-y-1 list-disc list-inside">
                              <li>Use either CBE or TELEBIRR to pay. Include your order reference in the payment description.</li>
                              <li>Payment should be made within 3 business days</li>
                              <li>Your order will be processed once payment is confirmed</li>
                            </ul>
                          </div>
                        </div>
                      </div>

                      {/* Option 1: CBE */}
                      <div className="bg-card border border-border rounded-lg p-6 space-y-4">
                        <div className="flex items-center gap-2 mb-4">
                          <Image src="/cbe-logo.png" alt="CBE" width={32} height={32} className="w-8 h-8" />
                          <h3 className="font-semibold text-base">CBE</h3>
                        </div>
                        <div className="space-y-4">
                          <div>
                            <label className="text-xs text-muted-foreground mb-2 block">Account Number</label>
                            <div className="flex items-center gap-2">
                              <div className="flex-1 flex items-center gap-2 p-3 bg-secondary rounded-md">
                                <span className="text-sm font-mono font-medium">1000024909364</span>
                              </div>
                              <button
                                onClick={() => copyToClipboard('1000024909364', 'cbe-account')}
                                className="p-3 bg-secondary hover:bg-secondary/80 rounded-md transition-colors"
                                aria-label="Copy CBE account number"
                              >
                                {copiedField === 'cbe-account' ? (
                                  <CheckCircle2 className="w-4 h-4 text-primary" />
                                ) : (
                                  <Copy className="w-4 h-4 text-muted-foreground" />
                                )}
                              </button>
                            </div>
                          </div>
                          <div>
                            <label className="text-xs text-muted-foreground mb-2 block">Name</label>
                            <div className="p-3 bg-secondary rounded-md">
                              <span className="text-sm font-medium">ARAGAW MELAK</span>
                            </div>
                          </div>
                        </div>
                      </div>

                      {/* Option 2: TELEBIRR */}
                      <div className="bg-card border border-border rounded-lg p-6 space-y-4">
                        <div className="flex items-center gap-2 mb-4">
                          <div className="relative h-8 w-20 shrink-0">
                            <Image src="/telebirr-logo.png" alt="Telebirr" fill className="object-contain object-left" sizes="80px" />
                          </div>
                          <h3 className="font-semibold text-base">TELEBIRR</h3>
                        </div>
                        <div className="space-y-4">
                          <div>
                            <label className="text-xs text-muted-foreground mb-2 block">Phone Number</label>
                            <div className="flex items-center gap-2">
                              <div className="flex-1 flex items-center gap-2 p-3 bg-secondary rounded-md">
                                <span className="text-sm font-mono font-medium">+251 93 763 9608</span>
                              </div>
                              <button
                                onClick={() => copyToClipboard('+251937639608', 'telebirr-phone')}
                                className="p-3 bg-secondary hover:bg-secondary/80 rounded-md transition-colors"
                                aria-label="Copy TELEBIRR phone"
                              >
                                {copiedField === 'telebirr-phone' ? (
                                  <CheckCircle2 className="w-4 h-4 text-primary" />
                                ) : (
                                  <Copy className="w-4 h-4 text-muted-foreground" />
                                )}
                              </button>
                            </div>
                          </div>
                          <div>
                            <label className="text-xs text-muted-foreground mb-2 block">Name</label>
                            <div className="p-3 bg-secondary rounded-md">
                              <span className="text-sm font-medium">ARAGAW MELAK</span>
                            </div>
                          </div>
                        </div>
                      </div>

                      {/* Share payment details on Telegram */}
                      <div className="bg-gradient-to-br from-blue-50/50 to-cyan-50/30 border-2 border-blue-200/50 rounded-lg p-5 shadow-sm">
                        <div className="flex items-start gap-3 mb-4">
                          <div className="mt-0.5">
                            <PiTelegramLogo className="w-6 h-6 text-[#0088cc]" />
                          </div>
                          <p className="text-sm font-medium text-gray-700 leading-relaxed">
                            After paying, share your payment details (e.g. transaction ID, amount, date) with us on Telegram so we can confirm and process your order.
                          </p>
                        </div>
                        <div className="flex justify-end">
                          <Link
                            href="https://t.me/efoyetastore"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center justify-center gap-2 rounded-md bg-[#0088cc] hover:bg-[#0077b5] text-white px-5 py-2.5 text-sm font-semibold transition-colors shadow-md hover:shadow-lg"
                          >
                            <PiTelegramLogo className="w-5 h-5" />
                            Share payment details on Telegram
                          </Link>
                        </div>
                      </div>
                    </div>
                  )}

                  {currentStep === 'review' && (
                    <div className="space-y-6">
                      <h2 className="text-xl font-semibold mb-4">Order Review</h2>
                      
                      {/* Order Items */}
                      <div className="space-y-4">
                        {cartItems.map((item) => {
                          const itemTotal = item.price * item.quantity;
                          return (
                            <div key={item.id} className="flex gap-4">
                              <div className="w-20 h-20 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0">
                                {item.image ? (
                                  <img
                                    src={item.image}
                                    alt={item.name}
                                    className="w-full h-full object-cover"
                                  />
                                ) : (
                                  <div className="w-full h-full flex items-center justify-center text-gray-400">
                                    <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                    </svg>
                                  </div>
                                )}
                              </div>
                              <div className="flex-1">
                                <p className="text-sm font-medium">{item.name}</p>
                                <p className="text-xs text-muted-foreground">
                                  Quantity: {item.quantity}
                                </p>
                                <p className="text-sm mt-1 font-medium">{itemTotal.toLocaleString()} ETB</p>
                              </div>
                            </div>
                          );
                        })}
                      </div>

                      {/* Order Summary */}
                      <div className="border-t border-border pt-4 space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">Subtotal</span>
                          <span>{totalPrice.toLocaleString()} ETB</span>
                        </div>
                        <div className="flex justify-between text-base font-semibold pt-2 border-t border-border">
                          <span>Total</span>
                          <span>{totalPrice.toLocaleString()} ETB</span>
                        </div>
                      </div>
                    </div>
                  )}
                </div>

                {/* Action Button */}
                <div className="flex gap-4">
                  {currentStepIndex > 0 && (
                    <Button
                      variant="outline"
                      onClick={handleBack}
                      className="flex-1"
                    >
                      Back
                    </Button>
                  )}
                  <Button
                    className="flex-1"
                    size="lg"
                    onClick={handleNext}
                    disabled={cartItems.length === 0 || isPlacingOrder || isLoadingCustomer}
                  >
                    {isLoadingCustomer 
                      ? "Loading customer info..." 
                      : isPlacingOrder 
                      ? "Placing Order..." 
                      : currentStep === 'payment' 
                      ? "Place Order" 
                      : "Continue"}
                  </Button>
                </div>
              </div>
            )}
          </Card>
        </main>
      </div>
  );
};

export default Cart;
