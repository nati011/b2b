import { motion, AnimatePresence } from 'framer-motion';
import { X, Minus, Plus, ShoppingBag, ShoppingCart, ArrowRight } from 'lucide-react';
import { useCart } from '@/context/CartContext';
import { Button } from '@/components/ui/button';
import { Link, useNavigate } from 'react-router-dom';

export const CartDrawer = () => {
  const { items, isOpen, closeCart, updateQuantity, removeItem, subtotal } = useCart();
  const navigate = useNavigate();
  
  const handlePurchase = () => {
    closeCart();
    navigate('/checkout');
  };
  
  const handleCheckout = () => {
    closeCart();
    navigate('/checkout');
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Overlay */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={closeCart}
            className="fixed inset-0 bg-black/40 z-50"
          />

          {/* Drawer */}
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 30, stiffness: 300 }}
            className="fixed right-0 top-0 bottom-0 w-full max-w-md bg-background z-50 flex flex-col shadow-2xl"
          >
            {/* Header */}
            <div className="flex items-center justify-between p-4 border-b border-border">
              <h2 className="font-display text-lg">Cart</h2>
              <button
                onClick={closeCart}
                className="p-2 hover:bg-secondary rounded-sm btn-press transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {/* Items */}
            <div className="flex-1 overflow-y-auto">
              {items.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-center p-8">
                  <ShoppingBag className="w-12 h-12 text-muted-foreground mb-4 stroke-[1]" />
                  <p className="text-muted-foreground mb-4">Your cart is empty</p>
                  <Button onClick={closeCart} variant="outline" size="sm">
                    Continue Shopping
                  </Button>
                </div>
              ) : (
                <div className="p-4 space-y-4">
                  {items.map((item) => (
                    <motion.div
                      key={`${item.product.id}-${item.size}-${item.color}`}
                      layout
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, x: 100 }}
                      className="flex gap-4"
                    >
                      <div className="w-24 h-32 img-soft rounded-sm overflow-hidden flex-shrink-0">
                        <img
                          src={item.product.attributes?.images?.[0] || '/placeholder.png'}
                          alt={item.product.name}
                          className="w-full h-full object-cover"
                        />
                      </div>
                      <div className="flex-1 flex flex-col justify-between py-1">
                        <div>
                          {item.product.attributes?.brand && (
                            <p className="text-xs text-muted-foreground uppercase tracking-wide">
                              {item.product.attributes.brand}
                            </p>
                          )}
                          <p className="text-sm font-medium leading-snug">
                            {item.product.name}
                          </p>
                          <p className="text-xs text-muted-foreground mt-1">
                            {item.size} / {item.color}
                          </p>
                        </div>
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <button
                              onClick={() => updateQuantity(item.product.id, item.size, item.color, item.quantity - 1)}
                              className="p-1 hover:bg-secondary rounded-sm btn-press"
                            >
                              <Minus className="w-4 h-4" />
                            </button>
                            <span className="text-sm font-medium w-6 text-center">
                              {item.quantity}
                            </span>
                            <button
                              onClick={() => updateQuantity(item.product.id, item.size, item.color, item.quantity + 1)}
                              className="p-1 hover:bg-secondary rounded-sm btn-press"
                            >
                              <Plus className="w-4 h-4" />
                            </button>
                          </div>
                          <div className="flex items-center gap-3">
                            <span className="text-sm font-medium">
                              {(item.product.price * item.quantity).toLocaleString()} ETB
                            </span>
                            <button
                              onClick={() => removeItem(item.product.id, item.size, item.color)}
                              className="text-muted-foreground hover:text-foreground btn-press"
                            >
                              <X className="w-4 h-4" />
                            </button>
                          </div>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                  
                  {/* Checkout Button in Items Section */}
                  <div className="pt-4 pb-2">
                    <Button 
                      onClick={handleCheckout}
                      className="w-full h-12 text-sm font-semibold tracking-wide btn-press rounded-sm"
                      size="lg"
                    >
                      <ArrowRight className="w-4 h-4 mr-2" />
                      Proceed to Checkout
                    </Button>
                  </div>
                </div>
              )}
            </div>

            {/* Footer */}
            {items.length > 0 && (
              <div className="p-4 border-t border-border space-y-4 safe-bottom bg-background">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Subtotal</span>
                  <span className="text-lg font-medium">{subtotal.toLocaleString()} ETB</span>
                </div>
                <p className="text-xs text-muted-foreground">
                  Shipping and taxes calculated at checkout
                </p>
                <div className="flex gap-3">
                  <Button 
                    onClick={handleCheckout}
                    variant="outline"
                    className="flex-1 h-12 text-sm font-medium tracking-wide btn-press rounded-sm"
                  >
                    Checkout
                  </Button>
                  <Button 
                    onClick={handlePurchase}
                    className="flex-1 h-12 text-base font-semibold tracking-wide btn-press rounded-sm shadow-lg"
                    size="lg"
                  >
                    <ShoppingCart className="w-5 h-5 mr-2" />
                    Purchase
                    <ArrowRight className="w-4 h-4 ml-2" />
                  </Button>
                </div>
              </div>
            )}
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
};
