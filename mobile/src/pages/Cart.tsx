import { motion } from 'framer-motion';
import { X, Minus, Plus, ShoppingBag } from 'lucide-react';
import { useCart } from '@/context/CartContext';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

const Cart = () => {
  const { items, updateQuantity, removeItem, subtotal } = useCart();

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Items */}
      <div className="min-h-[calc(100vh-200px)]">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center text-center py-16 px-4">
            <div className="bg-muted rounded-full p-8 mb-4">
              <ShoppingBag className="w-12 h-12 text-muted-foreground" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Your cart is empty</h3>
            <p className="text-sm text-muted-foreground mb-6 max-w-sm">
              Add items to your cart to see them here. Start shopping to fill your cart.
            </p>
            <Link
              to="/"
              className="px-6 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium btn-press"
            >
              Start Shopping
            </Link>
          </div>
        ) : (
          <>
            <div className="p-4 space-y-4">
              {items.map((item) => (
                <motion.div
                  key={`${item.product.id}-${item.size}-${item.color}`}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
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
                          {((item.product.price * item.quantity).toLocaleString())} ETB
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
            </div>

            {/* Footer */}
            <div className="sticky bottom-0 bg-background border-t border-border p-4 space-y-4 mt-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Subtotal</span>
                <span className="text-lg font-medium">{subtotal.toLocaleString()} ETB</span>
              </div>
              <p className="text-xs text-muted-foreground">
                Shipping and taxes calculated at checkout
              </p>
              <Link to="/checkout">
                <Button className="w-full h-12 text-sm tracking-wide btn-press rounded-sm">
                  Checkout
                </Button>
              </Link>
            </div>
          </>
        )}
      </div>
    </motion.div>
  );
};

export default Cart;

