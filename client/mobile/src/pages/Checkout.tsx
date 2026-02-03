import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { ArrowLeft, Check, Building2, CreditCard, Copy, CheckCircle2, Loader2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '@/context/CartContext';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { CreateOrder, type CreateOrderRequest } from '@/lib/api/order';
import { GetCustomerByEmail } from '@/lib/api/customer';
import { toast } from 'sonner';

type CheckoutStep = 'payment' | 'review';

const steps: { id: CheckoutStep; label: string }[] = [
  { id: 'payment', label: 'Payment' },
  { id: 'review', label: 'Review' },
];

const Checkout = () => {
  const navigate = useNavigate();
  const { items, subtotal, clearCart } = useCart();
  const { user, isAuthenticated } = useAuth();
  const [currentStep, setCurrentStep] = useState<CheckoutStep>('payment');
  const [isComplete, setIsComplete] = useState(false);
  const [copiedField, setCopiedField] = useState<string | null>(null);
  const [isPlacingOrder, setIsPlacingOrder] = useState(false);
  const [customerId, setCustomerId] = useState<number | null>(null);
  const [customerData, setCustomerData] = useState<any>(null);
  const [isLoadingCustomer, setIsLoadingCustomer] = useState(true);

  // Load customer data from localStorage or fetch from API
  useEffect(() => {
    const loadCustomerData = async () => {
      if (!isAuthenticated || !user?.email) {
        setIsLoadingCustomer(false);
        return;
      }

      // First, try to load from localStorage
      try {
        const storedCustomerId = localStorage.getItem('customer_id');
        const storedCustomerData = localStorage.getItem('customer_data');
        
        if (storedCustomerId && storedCustomerData) {
          try {
            const customerIdNum = parseInt(storedCustomerId, 10);
            const customerDataParsed = JSON.parse(storedCustomerData);
            
            // Verify the stored customer email matches the logged-in user
            if (customerDataParsed.email?.toLowerCase() === user.email.toLowerCase() && !isNaN(customerIdNum)) {
              setCustomerId(customerIdNum);
              setCustomerData(customerDataParsed);
              setIsLoadingCustomer(false);
              console.log('✅ Loaded customer data from localStorage');
              return;
            } else {
              console.log('⚠️ Stored customer data does not match current user, fetching from API');
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
      } catch (storageError) {
        console.error('Error reading from localStorage:', storageError);
      }

      // API call disabled - only use localStorage data
      // If not in localStorage or data doesn't match, skip API fetch
      // The order placement function will handle fetching customer data if needed
      console.log('⚠️ Customer data not found in localStorage, will fetch during order placement if needed');
      setIsLoadingCustomer(false);
    };

    loadCustomerData();
  }, [isAuthenticated, user, navigate]);

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!isLoadingCustomer && !isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, isLoadingCustomer, navigate]);

  const currentStepIndex = steps.findIndex(s => s.id === currentStep);

  const handleNext = async () => {
    const nextIndex = currentStepIndex + 1;
    if (nextIndex < steps.length) {
      setCurrentStep(steps[nextIndex].id);
    } else {
      // Place order
      await handlePlaceOrder();
    }
  };

  const handlePlaceOrder = async () => {
    if (items.length === 0) {
      toast.error('Your cart is empty');
      return;
    }

    // Validate user is authenticated
    if (!isAuthenticated || !user) {
      toast.error('Authentication required', {
        description: 'Please log in to place an order.',
      });
      navigate('/login');
      return;
    }

    setIsPlacingOrder(true);

    try {
      // Try to get customer ID from localStorage or state
      let finalCustomerId = customerId;
      let finalCustomerData = customerData;

      console.log('🔍 Checking customer data for order:', {
        hasCustomerId: !!customerId,
        hasCustomerData: !!customerData,
        userEmail: user.email,
      });

      // If not available in state, try localStorage
      if (!finalCustomerId || !finalCustomerData) {
        try {
          const storedCustomerId = localStorage.getItem('customer_id');
          const storedCustomerData = localStorage.getItem('customer_data');
          
          console.log('🔍 Checking localStorage:', {
            hasStoredId: !!storedCustomerId,
            hasStoredData: !!storedCustomerData,
          });
          
          if (storedCustomerId && storedCustomerData) {
            try {
              const customerIdNum = parseInt(storedCustomerId, 10);
              const customerDataParsed = JSON.parse(storedCustomerData);
              
              console.log('🔍 Parsed customer data:', {
                customerId: customerIdNum,
                customerEmail: customerDataParsed.email,
                userEmail: user.email,
                emailsMatch: customerDataParsed.email?.toLowerCase() === user.email.toLowerCase(),
              });
              
              // Check if email matches, or if no email in stored data, use it anyway
              // Also allow if customer data has an id (it's valid customer data)
              const emailMatches = !customerDataParsed.email || 
                customerDataParsed.email?.toLowerCase() === user.email.toLowerCase();
              
              // Accept customer data if we have a valid ID and either email matches or customer has an id field
              if (!isNaN(customerIdNum) && (emailMatches || customerDataParsed.id)) {
                finalCustomerId = customerIdNum;
                finalCustomerData = customerDataParsed;
                console.log('✅ Using customer data from localStorage for order');
              } else {
                console.warn('⚠️ Email mismatch or invalid customer ID:', {
                  storedEmail: customerDataParsed.email,
                  userEmail: user.email,
                  customerIdValid: !isNaN(customerIdNum),
                });
              }
            } catch (parseError) {
              console.error('❌ Error parsing stored customer data:', parseError);
            }
          }
        } catch (error) {
          console.error('❌ Error reading customer data from localStorage:', error);
        }
      }

      // If still no customer data, try to create minimal customer data from user context
      if (!finalCustomerId || !finalCustomerData) {
        console.warn('⚠️ No customer data found, attempting to fetch from API...');
        
        // If we have user info but no customer data, we still need customer_id from backend
        // But let's check if we can find it by making a quick API call
        let apiErrorOccurred = false;
        let customerNotFound = false;
        
        try {
          console.log('🔄 Attempting to fetch customer data from API for email:', user.email);
          const customer = await GetCustomerByEmail(user.email);
          
          if (customer) {
            finalCustomerId = customer.id;
            finalCustomerData = customer;
            
            // Store for future use
            localStorage.setItem('customer_id', customer.id.toString());
            localStorage.setItem('customer_data', JSON.stringify(customer));
            console.log('✅ Fetched customer data from API during order placement');
          } else {
            customerNotFound = true;
            console.warn('⚠️ Customer not found in API for email:', user.email);
          }
        } catch (apiError: any) {
          apiErrorOccurred = true;
          console.error('❌ Error fetching customer during order placement:', apiError);
          
          // Log more details about the error
          if (apiError.response) {
            console.error('API Error Details:', {
              status: apiError.response.status,
              data: apiError.response.data,
            });
          } else if (apiError.request) {
            console.error('Network Error: No response received from server');
          }
        }
        
        // If still no customer data after API call, provide specific error message
        if (!finalCustomerId || !finalCustomerData) {
          let errorMessage = 'Customer information not available';
          let errorDescription = 'Please complete your profile before placing an order.';
          
          if (apiErrorOccurred) {
            errorMessage = 'Unable to load customer information';
            errorDescription = 'There was an error connecting to the server. Please check your internet connection and try again.';
          } else if (customerNotFound) {
            errorMessage = 'Customer profile not found';
            errorDescription = `No customer profile found for ${user.email}. Please complete your profile first.`;
          }
          
          console.error('❌ Customer data not available:', {
            finalCustomerId,
            finalCustomerData,
            userEmail: user.email,
            apiErrorOccurred,
            customerNotFound,
          });
          
          toast.error(errorMessage, {
            description: errorDescription,
          });
          navigate('/profile');
          setIsPlacingOrder(false);
          return;
        }
      }
      
      // Prepare order items
      const orderItems = items.map(item => ({
        product_id: item.product.id,
        quantity: item.quantity,
        price: item.product.price,
      }));

      // Prepare customer snapshot using customer data
      const customerSnapshot = JSON.stringify({
        id: finalCustomerData.id,
        full_name: finalCustomerData.full_name || user.name || '',
        email: finalCustomerData.email || user.email || '',
        phone_number: finalCustomerData.phone_number || '',
        city: finalCustomerData.city || '',
        region: finalCustomerData.region || '',
        woreda: finalCustomerData.woreda || '',
        status: finalCustomerData.status || 'active',
      });

      // Prepare shipping address snapshot using customer data
      const shippingAddressSnapshot = JSON.stringify({
        full_name: finalCustomerData.full_name || user.name || '',
        email: finalCustomerData.email || user.email || '',
        phone_number: finalCustomerData.phone_number || '',
        street_address: finalCustomerData.woreda || '',
        city: finalCustomerData.city || '',
        region: finalCustomerData.region || '',
      });

      // Create order request using customer credentials from context
      const orderRequest: CreateOrderRequest = {
        customer_id: finalCustomerId, // From customer API or localStorage using logged-in user's email
        status: 'PENDING',
        payment_status: 'PENDING',
        delivery_status: 'PENDING',
        confirmation_status: 'PENDING',
        total: subtotal,
        customer_snapshot: customerSnapshot as any,
        shipping_address_snapshot: shippingAddressSnapshot as any,
        items: orderItems,
      };

      const orderResponse = await CreateOrder(orderRequest);
      
      toast.success('Order placed successfully!', {
        description: `Order #${orderResponse.id} has been created.`,
      });

      // Clear cart and show success
      clearCart();
      setIsComplete(true);
    } catch (error: any) {
      console.error('Failed to place order:', error);
      toast.error('Failed to place order', {
        description: error || 'An error occurred while placing your order. Please try again.',
      });
    } finally {
      setIsPlacingOrder(false);
    }
  };

  const handleBack = () => {
    const prevIndex = currentStepIndex - 1;
    if (prevIndex >= 0) {
      setCurrentStep(steps[prevIndex].id);
    } else {
      navigate(-1);
    }
  };

  // Show loading state while fetching customer
  if (isLoadingCustomer) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="min-h-screen flex flex-col items-center justify-center p-8"
      >
        <Loader2 className="w-8 h-8 animate-spin text-primary mb-4" />
        <p className="text-muted-foreground">Loading checkout...</p>
      </motion.div>
    );
  }

  if (isComplete) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="min-h-screen flex flex-col items-center justify-center p-8 text-center"
      >
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', damping: 15 }}
          className="w-16 h-16 bg-foreground rounded-full flex items-center justify-center mb-6"
        >
          <Check className="w-8 h-8 text-background" />
        </motion.div>
        <h1 className="font-display text-2xl mb-2">Order Confirmed</h1>
        <p className="text-muted-foreground mb-6">
          Thank you for your purchase. We'll send you a confirmation email shortly.
        </p>
        <Button onClick={() => navigate('/')} variant="outline" className="rounded-sm">
          Continue Shopping
        </Button>
      </motion.div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition min-h-screen pb-32"
    >
      {/* Header */}
      <div className="sticky top-0 bg-background z-10 border-b border-border">
        <div className="flex items-center p-4">
          <button onClick={handleBack} className="p-2 -ml-2 hover:bg-secondary rounded-sm btn-press">
            <ArrowLeft className="w-5 h-5" />
          </button>
          <h1 className="font-display text-lg ml-2">Checkout</h1>
        </div>

        {/* Breadcrumb */}
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
                    "flex items-center gap-2 transition-all",
                    canNavigate && "cursor-pointer hover:opacity-80",
                    !canNavigate && "cursor-not-allowed opacity-60"
                  )}
                >
                  <div
                    className={cn(
                      "w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium transition-colors",
                      isActive && "bg-primary text-primary-foreground",
                      isCompleted && "bg-primary text-primary-foreground",
                      !isActive && !isCompleted && "bg-secondary text-muted-foreground"
                    )}
                  >
                    {isCompleted ? <Check className="w-3 h-3" /> : index + 1}
                  </div>
                  <span
                    className={cn(
                      "text-sm transition-colors",
                      isActive ? "text-foreground font-medium" : "text-muted-foreground"
                    )}
                  >
                    {step.label}
                  </span>
                </button>
                {index < steps.length - 1 && (
                  <div className={cn(
                    "w-8 h-px",
                    isCompleted ? "bg-primary" : "bg-border"
                  )} />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Step Content */}
      <div className="p-4">
        {currentStep === 'payment' && (
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="space-y-4"
          >
            <h2 className="font-display text-lg mb-4">Payment Details</h2>
            
            {/* Payment Information Card */}
            <div className="bg-card border border-border rounded-lg p-4 space-y-4">
              <div className="flex items-center gap-2 mb-3">
                <Building2 className="w-5 h-5 text-primary" />
                <h3 className="font-semibold text-sm">Company Payment Information</h3>
              </div>
              
              <div className="space-y-3">
                {/* Bank Name */}
                <div>
                  <label className="text-xs text-muted-foreground mb-1 block">Bank Name</label>
                  <div className="flex items-center gap-2 p-3 bg-secondary rounded-sm">
                    <CreditCard className="w-4 h-4 text-muted-foreground" />
                    <span className="text-sm font-medium">Commercial Bank of Ethiopia</span>
                  </div>
                </div>

                {/* Account Number */}
                <div>
                  <label className="text-xs text-muted-foreground mb-1 block">Account Number</label>
                  <div className="flex items-center gap-2">
                    <div className="flex-1 flex items-center gap-2 p-3 bg-secondary rounded-sm">
                      <span className="text-sm font-mono font-medium">1000123456789</span>
                    </div>
                    <button
                      onClick={() => {
                        navigator.clipboard.writeText('1000123456789');
                        setCopiedField('account');
                        setTimeout(() => setCopiedField(null), 2000);
                      }}
                      className="p-3 bg-secondary hover:bg-secondary/80 rounded-sm transition-colors"
                      aria-label="Copy account number"
                    >
                      {copiedField === 'account' ? (
                        <CheckCircle2 className="w-4 h-4 text-primary" />
                      ) : (
                        <Copy className="w-4 h-4 text-muted-foreground" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Account Holder */}
                <div>
                  <label className="text-xs text-muted-foreground mb-1 block">Account Holder</label>
                  <div className="p-3 bg-secondary rounded-sm">
                    <span className="text-sm font-medium">Efoyeta Store PLC</span>
                  </div>
                </div>

                {/* SWIFT Code */}
                <div>
                  <label className="text-xs text-muted-foreground mb-1 block">SWIFT Code</label>
                  <div className="flex items-center gap-2">
                    <div className="flex-1 flex items-center gap-2 p-3 bg-secondary rounded-sm">
                      <span className="text-sm font-mono font-medium">CBETETAA</span>
                    </div>
                    <button
                      onClick={() => {
                        navigator.clipboard.writeText('CBETETAA');
                        setCopiedField('swift');
                        setTimeout(() => setCopiedField(null), 2000);
                      }}
                      className="p-3 bg-secondary hover:bg-secondary/80 rounded-sm transition-colors"
                      aria-label="Copy SWIFT code"
                    >
                      {copiedField === 'swift' ? (
                        <CheckCircle2 className="w-4 h-4 text-primary" />
                      ) : (
                        <Copy className="w-4 h-4 text-muted-foreground" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Branch */}
                <div>
                  <label className="text-xs text-muted-foreground mb-1 block">Branch</label>
                  <div className="p-3 bg-secondary rounded-sm">
                    <span className="text-sm">Addis Ababa Main Branch</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Payment Instructions */}
            <div className="bg-primary/5 border border-primary/20 rounded-lg p-4 space-y-2">
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-primary mt-0.5 shrink-0" />
                <div className="space-y-1">
                  <p className="text-xs font-medium text-foreground">Payment Instructions</p>
                  <ul className="text-xs text-muted-foreground space-y-1 list-disc list-inside">
                    <li>Please include your order reference number in the payment description</li>
                    <li>Payment should be made within 3 business days</li>
                    <li>Your order will be processed once payment is confirmed</li>
                  </ul>
                </div>
              </div>
            </div>

            {/* Alternative Payment Methods */}
            <div className="border-t border-border pt-4">
              <p className="text-xs text-muted-foreground text-center">
                For alternative payment methods, please contact our support team
              </p>
            </div>
          </motion.div>
        )}

        {currentStep === 'review' && (
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="space-y-6"
          >
            <h2 className="font-display text-lg mb-4">Order Review</h2>
            
            {/* Order Items */}
            <div className="space-y-3">
              {items.map((item) => (
                <div key={`${item.product.id}-${item.size}-${item.color}`} className="flex gap-3">
                  <div className="w-16 h-20 img-soft rounded-sm overflow-hidden">
                    <img
                      src={item.product.attributes?.images?.[0] || '/placeholder.png'}
                      alt={item.product.name}
                      className="w-full h-full object-cover"
                    />
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium">{item.product.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {item.size} / {item.color} × {item.quantity}
                    </p>
                    <p className="text-sm mt-1">{(item.product.price * item.quantity).toLocaleString()} ETB</p>
                  </div>
                </div>
              ))}
            </div>

            {/* Order Summary */}
            <div className="border-t border-border pt-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Subtotal</span>
                <span>{subtotal.toLocaleString()} ETB</span>
              </div>
              <div className="flex justify-between text-base font-medium pt-2 border-t border-border">
                <span>Total</span>
                <span>{subtotal.toLocaleString()} ETB</span>
              </div>
            </div>
          </motion.div>
        )}
      </div>

      {/* Fixed Footer */}
      <div className="fixed bottom-14 left-0 right-0 p-4 bg-background border-t border-border safe-bottom">
        <Button 
          onClick={handleNext} 
          disabled={isPlacingOrder}
          className="w-full h-12 text-sm tracking-wide btn-press rounded-sm"
        >
          {isPlacingOrder ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Placing Order...
            </>
          ) : (
            currentStep === 'review' ? 'Place Order' : 'Continue'
          )}
        </Button>
      </div>
    </motion.div>
  );
};

export default Checkout;
