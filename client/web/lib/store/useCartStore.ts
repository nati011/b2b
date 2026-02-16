import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { CartItem, cartItemKey } from '@/lib/types';

// Helper function to calculate totals efficiently
const calculateTotals = (items: CartItem[]) => {
    let totalItems = 0;
    let totalPrice = 0;
    for (const item of items) {
        totalItems += item.quantity;
        totalPrice += item.price * item.quantity;
    }
    return { totalItems, totalPrice };
};

interface CartStore {
    success: string | null;
    totalItems: number;
    totalPrice: number;
    cartItems: CartItem[];
    addCartItems: (product: Omit<CartItem, "quantity">, quantity: number) => void;
    removeCartItems: (cartItem: CartItem) => void;
    clearCart: () => void;
}

const useCartStore = create<CartStore>()(
    persist(
        (set) => ({
            success: null,
            totalItems: 0,
            totalPrice: 0,
            cartItems: [],
            clearCart: () => {
                set({
                    cartItems: [],
                    totalItems: 0,
                    totalPrice: 0,
                    success: 'Cart cleared'
                });
            },
            addCartItems: (product, quantity = 1) => set((state) => {
                const newItemForKey: CartItem = { ...product, quantity: 0 };
                const key = cartItemKey(newItemForKey);
                const existingItemIndex = state.cartItems.findIndex(
                    (item) => cartItemKey(item) === key
                );

                let updatedItems: CartItem[];
                if (existingItemIndex >= 0) {
                    updatedItems = [...state.cartItems];
                    const newQuantity = updatedItems[existingItemIndex].quantity + quantity;
                    if (newQuantity <= 0) {
                        updatedItems = updatedItems.filter((_, idx) => idx !== existingItemIndex);
                    } else {
                        updatedItems[existingItemIndex] = {
                            ...updatedItems[existingItemIndex],
                            quantity: newQuantity
                        };
                    }
                } else {
                    if (quantity > 0) {
                        const newItem: CartItem = { ...product, quantity };
                        updatedItems = [...state.cartItems, newItem];
                    } else {
                        updatedItems = state.cartItems;
                    }
                }

                const { totalItems, totalPrice } = calculateTotals(updatedItems);
                return {
                    cartItems: updatedItems,
                    totalItems,
                    totalPrice,
                    success: existingItemIndex >= 0 ? 'Item quantity updated in cart' : 'Item added to cart'
                };
            }),
            removeCartItems: (cartItem) => set((state) => {
                const key = cartItemKey(cartItem);
                const updatedCartItems = state.cartItems.filter((item) => cartItemKey(item) !== key);
                const { totalItems, totalPrice } = calculateTotals(updatedCartItems);
                return {
                    cartItems: updatedCartItems,
                    totalItems,
                    totalPrice,
                    success: 'Item removed from cart'
                };
            }),
        }),
        {
            name: 'cart-storage',
            // Optimize persistence - only persist cartItems, recalculate totals on load
            partialize: (state) => ({ cartItems: state.cartItems }),
            // Recalculate totals when loading from storage
            onRehydrateStorage: () => (state) => {
                if (state?.cartItems) {
                    const { totalItems, totalPrice } = calculateTotals(state.cartItems);
                    state.totalItems = totalItems;
                    state.totalPrice = totalPrice;
                }
            },
        }
    )
);

export default useCartStore;