import { useState } from 'react';
import { motion } from 'framer-motion';
import { ArrowLeft, Check } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '@/context/CartContext';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

type CheckoutStep = 'shipping' | 'payment' | 'review';

const steps: { id: CheckoutStep; label: string }[] = [
  { id: 'shipping', label: 'Shipping' },
  { id: 'payment', label: 'Payment' },
  { id: 'review', label: 'Review' },
];

const Checkout = () => {
  const navigate = useNavigate();
  const { items, subtotal, clearCart } = useCart();
  const [currentStep, setCurrentStep] = useState<CheckoutStep>('shipping');
  const [isComplete, setIsComplete] = useState(false);

  const currentStepIndex = steps.findIndex(s => s.id === currentStep);

  const handleNext = () => {
    const nextIndex = currentStepIndex + 1;
    if (nextIndex < steps.length) {
      setCurrentStep(steps[nextIndex].id);
    } else {
      // Complete order
      setIsComplete(true);
      clearCart();
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

            return (
              <div key={step.id} className="flex items-center gap-2">
                <div
                  className={cn(
                    "w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium transition-colors",
                    isActive && "bg-foreground text-background",
                    isCompleted && "bg-foreground text-background",
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
                {index < steps.length - 1 && (
                  <div className={cn(
                    "w-8 h-px",
                    isCompleted ? "bg-foreground" : "bg-border"
                  )} />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Step Content */}
      <div className="p-4">
        {currentStep === 'shipping' && (
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="space-y-4"
          >
            <h2 className="font-display text-lg mb-4">Shipping Address</h2>
            <div className="space-y-3">
              <input
                type="text"
                placeholder="Full Name"
                className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
              />
              <input
                type="email"
                placeholder="Email Address"
                className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
              />
              <input
                type="text"
                placeholder="Street Address"
                className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
              />
              <div className="grid grid-cols-2 gap-3">
                <input
                  type="text"
                  placeholder="City"
                  className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
                />
                <input
                  type="text"
                  placeholder="ZIP Code"
                  className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
                />
              </div>
            </div>
          </motion.div>
        )}

        {currentStep === 'payment' && (
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="space-y-4"
          >
            <h2 className="font-display text-lg mb-4">Payment Method</h2>
            <div className="space-y-3">
              <input
                type="text"
                placeholder="Card Number"
                className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
              />
              <input
                type="text"
                placeholder="Cardholder Name"
                className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
              />
              <div className="grid grid-cols-2 gap-3">
                <input
                  type="text"
                  placeholder="MM/YY"
                  className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
                />
                <input
                  type="text"
                  placeholder="CVC"
                  className="w-full h-12 px-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
                />
              </div>
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
                      src={item.product.images[0]}
                      alt={item.product.name}
                      className="w-full h-full object-cover"
                    />
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium">{item.product.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {item.size} / {item.color} × {item.quantity}
                    </p>
                    <p className="text-sm mt-1">${(item.product.price * item.quantity).toLocaleString()}</p>
                  </div>
                </div>
              ))}
            </div>

            {/* Order Summary */}
            <div className="border-t border-border pt-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Subtotal</span>
                <span>${subtotal.toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Shipping</span>
                <span>Free</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Tax</span>
                <span>${Math.round(subtotal * 0.08).toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-base font-medium pt-2 border-t border-border">
                <span>Total</span>
                <span>${Math.round(subtotal * 1.08).toLocaleString()}</span>
              </div>
            </div>
          </motion.div>
        )}
      </div>

      {/* Fixed Footer */}
      <div className="fixed bottom-14 left-0 right-0 p-4 bg-background border-t border-border safe-bottom">
        <Button onClick={handleNext} className="w-full h-12 text-sm tracking-wide btn-press rounded-sm">
          {currentStep === 'review' ? 'Place Order' : 'Continue'}
        </Button>
      </div>
    </motion.div>
  );
};

export default Checkout;
