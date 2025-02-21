import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { SAMPLE_PRODUCTS } from '../types/product';
import { Minus, Plus, X, ArrowLeft } from 'lucide-react';

export function Cart() {
  const navigate = useNavigate();
  const { state: cart, dispatch } = useCart();

  const cartItems = cart.items.map(item => ({
    ...item,
    product: SAMPLE_PRODUCTS.find(p => p.id === item.productId)!
  }));

  const subtotal = cartItems.reduce(
    (total, item) => total + item.product.price * item.quantity,
    0
  );

  const handleQuantityChange = (productId: string, delta: number) => {
    const item = cart.items.find(item => item.productId === productId);
    if (!item) return;

    const newQuantity = item.quantity + delta;
    if (newQuantity < 1) {
      dispatch({ type: 'REMOVE_FROM_CART', payload: { productId } });
    } else {
      dispatch({
        type: 'UPDATE_QUANTITY',
        payload: { productId, quantity: newQuantity }
      });
    }
  };

  const handleRemove = (productId: string) => {
    dispatch({ type: 'REMOVE_FROM_CART', payload: { productId } });
  };

  if (cartItems.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900">Your cart is empty</h2>
          <button
            onClick={() => navigate('/')}
            className="mt-4 inline-flex items-center text-blue-600 hover:text-blue-500"
          >
            <ArrowLeft className="h-5 w-5 mr-2" />
            Continue shopping
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Shopping Cart</h1>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div className="lg:col-span-8">
          <div className="space-y-4">
            {cartItems.map(({ product, quantity }) => (
              <div
                key={product.id}
                className="flex items-center gap-4 bg-white p-4 rounded-lg shadow-sm"
              >
                <img
                  src={product.image}
                  alt={product.name}
                  className="h-24 w-24 object-cover rounded"
                />
                
                <div className="flex-grow">
                  <h3 className="text-lg font-medium text-gray-900">{product.name}</h3>
                  <p className="text-gray-600">${product.price}</p>
                </div>

                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => handleQuantityChange(product.id, -1)}
                    className="p-1 text-gray-400 hover:text-gray-500"
                  >
                    <Minus className="h-4 w-4" />
                  </button>
                  <span className="text-gray-900 w-8 text-center">{quantity}</span>
                  <button
                    onClick={() => handleQuantityChange(product.id, 1)}
                    className="p-1 text-gray-400 hover:text-gray-500"
                  >
                    <Plus className="h-4 w-4" />
                  </button>
                </div>

                <div className="text-right">
                  <p className="text-lg font-medium text-gray-900">
                    ${(product.price * quantity).toFixed(2)}
                  </p>
                  <button
                    onClick={() => handleRemove(product.id)}
                    className="text-gray-400 hover:text-gray-500"
                  >
                    <X className="h-5 w-5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="lg:col-span-4">
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-lg font-medium text-gray-900 mb-4">Order Summary</h2>
            
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-gray-600">Subtotal</span>
                <span className="text-gray-900">${subtotal.toFixed(2)}</span>
              </div>
              
              <div className="flex justify-between">
                <span className="text-gray-600">Shipping</span>
                <span className="text-gray-900">Free</span>
              </div>
              
              <div className="border-t pt-4">
                <div className="flex justify-between">
                  <span className="text-lg font-medium text-gray-900">Total</span>
                  <span className="text-lg font-medium text-gray-900">
                    ${subtotal.toFixed(2)}
                  </span>
                </div>
              </div>

              <button className="w-full bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 transition-colors">
                Checkout
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}