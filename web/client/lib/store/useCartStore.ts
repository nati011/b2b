import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { CartItem } from '@/lib/types';

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
                const existingItemIndex = state.cartItems.findIndex(
                    (item) => item.id === product.id
                );

                if (existingItemIndex >= 0) {
                    const updatedItems = [...state.cartItems];
                    updatedItems[existingItemIndex] = {
                        ...updatedItems[existingItemIndex],
                        quantity: updatedItems[existingItemIndex].quantity + quantity
                    };
                    return {
                        cartItems: updatedItems,
                        totalItems: updatedItems.reduce((total, item) => total + item.quantity, 0),
                        totalPrice: updatedItems.reduce((total, item) => total + item.price * item.quantity, 0),
                        success: 'Item quantity updated in cart'
                    };
                } else {
                    const newItem = { ...product, quantity };
                    const cartItems = [...state.cartItems, newItem];
                    return {
                        cartItems,
                        totalItems: cartItems.reduce((total, item) => total + item.quantity, 0),
                        totalPrice: cartItems.reduce((total, item) => total + item.price * item.quantity, 0),
                        success: 'Item added to cart'
                    };
                }
            }),
            removeCartItems: (cartItem) => set((state) => {
                const updatedCartItems = state.cartItems.filter((item) => item.id !== cartItem.id);
                return {
                    cartItems: updatedCartItems,
                    totalItems: updatedCartItems.reduce((total, item) => total + item.quantity, 0),
                    totalPrice: updatedCartItems.reduce((total, item) => total + item.price * item.quantity, 0),
                    success: 'Item removed from cart'
                };
            }),
        }),
        {
            name: 'cart-storage',
        }
    )
);

export default useCartStore;