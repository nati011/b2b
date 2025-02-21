import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { SAMPLE_PRODUCTS } from '../types/product';
import { useCart } from '../context/CartContext';
import { Minus, Plus, ArrowLeft } from 'lucide-react';

export function ProductDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { dispatch } = useCart();
  const [quantity, setQuantity] = useState(1);

  const product = SAMPLE_PRODUCTS.find(p => p.id === id);

  if (!product) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900">Product not found</h2>
          <button
            onClick={() => navigate('/')}
            className="mt-4 inline-flex items-center text-blue-600 hover:text-blue-500"
          >
            <ArrowLeft className="h-5 w-5 mr-2" />
            Back to products
          </button>
        </div>
      </div>
    );
  }

  const handleQuantityChange = (delta: number) => {
    setQuantity(prev => Math.max(1, prev + delta));
  };

  const handleAddToCart = () => {
    dispatch({
      type: 'ADD_TO_CART',
      payload: { productId: product.id, quantity },
    });
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <button
        onClick={() => navigate('/')}
        className="inline-flex items-center text-blue-600 hover:text-blue-500 mb-8"
      >
        <ArrowLeft className="h-5 w-5 mr-2" />
        Back to products
      </button>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
        <div className="aspect-square overflow-hidden rounded-lg bg-gray-100">
          <img
            src={product.image}
            alt={product.name}
            className="h-full w-full object-cover object-center"
          />
        </div>

        <div>
          <h1 className="text-3xl font-bold text-gray-900">{product.name}</h1>
          <p className="mt-4 text-xl text-gray-900">${product.price}</p>
          <p className="mt-4 text-gray-600">{product.description}</p>

          <div className="mt-8">
            <label htmlFor="quantity" className="block text-sm font-medium text-gray-700">
              Quantity
            </label>
            <div className="mt-2 flex items-center space-x-4">
              <button
                onClick={() => handleQuantityChange(-1)}
                className="rounded-md p-2 text-gray-400 hover:text-gray-500"
              >
                <Minus className="h-5 w-5" />
              </button>
              <span className="text-gray-900 w-8 text-center">{quantity}</span>
              <button
                onClick={() => handleQuantityChange(1)}
                className="rounded-md p-2 text-gray-400 hover:text-gray-500"
              >
                <Plus className="h-5 w-5" />
              </button>
            </div>
          </div>

          <button
            onClick={handleAddToCart}
            className="mt-8 w-full bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 transition-colors"
          >
            Add to Cart
          </button>
        </div>
      </div>
    </div>
  );
}