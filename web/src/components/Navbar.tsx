import React from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { Search, ShoppingCart, User } from 'lucide-react';
import { useCart } from '../context/CartContext';

const CATEGORIES = ['All', 'Clothing', 'Footwear', 'Accessories'];

export function Navbar() {
  const [searchParams, setSearchParams] = useSearchParams();
  const currentCategory = searchParams.get('category') || 'All';
  const { state: cart } = useCart();

  const cartItemCount = cart.items.reduce((total, item) => total + item.quantity, 0);

  const handleCategoryChange = (category: string) => {
    if (category === 'All') {
      searchParams.delete('category');
    } else {
      searchParams.set('category', category);
    }
    searchParams.delete('page');
    setSearchParams(searchParams);
  };

  return (
    <nav className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col">
          <div className="flex justify-between h-16 items-center">
            <Link to="/" className="flex items-center">
              <span className="text-2xl font-bold text-gray-900">Store</span>
            </Link>
            
            <div className="hidden sm:block flex-1 max-w-lg mx-8">
              <div className="relative">
                <input
                  type="text"
                  placeholder="Search products..."
                  className="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <Search className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
              </div>
            </div>

            <div className="flex items-center space-x-4">
              <Link to="/cart" className="text-gray-700 hover:text-gray-900">
                <div className="relative">
                  <ShoppingCart className="h-6 w-6" />
                  {cartItemCount > 0 && (
                    <span className="absolute -top-2 -right-2 bg-blue-600 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                      {cartItemCount}
                    </span>
                  )}
                </div>
              </Link>
              <Link to="/login" className="text-gray-700 hover:text-gray-900">
                <User className="h-6 w-6" />
              </Link>
            </div>
          </div>
          
          <div className="flex space-x-8 h-12 items-center -mb-px">
            {CATEGORIES.map((category) => (
              <button
                key={category}
                onClick={() => handleCategoryChange(category)}
                className={`
                  text-sm font-medium px-1 py-2 border-b-2 transition-colors
                  ${currentCategory === category
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                `}
              >
                {category}
              </button>
            ))}
          </div>
        </div>
      </div>
    </nav>
  );
}