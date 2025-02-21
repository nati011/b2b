export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  image: string;
  category: string;
}

export const SAMPLE_PRODUCTS: Product[] = [
  {
    id: '1',
    name: 'Basic White T-Shirt',
    description: 'Classic cotton t-shirt',
    price: 29.99,
    image: 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Clothing'
  },
  {
    id: '2',
    name: 'Leather Backpack',
    description: 'Durable everyday backpack',
    price: 89.99,
    image: 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Accessories'
  },
  {
    id: '3',
    name: 'Running Shoes',
    description: 'Comfortable athletic shoes',
    price: 119.99,
    image: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Footwear'
  },
  {
    id: '4',
    name: 'Denim Jeans',
    description: 'Classic blue jeans',
    price: 79.99,
    image: 'https://images.unsplash.com/photo-1542272604-787c3835535d?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Clothing'
  },
  {
    id: '5',
    name: 'Canvas Sneakers',
    description: 'Casual everyday sneakers',
    price: 59.99,
    image: 'https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Footwear'
  },
  {
    id: '6',
    name: 'Wool Sweater',
    description: 'Warm winter sweater',
    price: 89.99,
    image: 'https://images.unsplash.com/photo-1576566588028-4147f3842f27?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Clothing'
  },
  {
    id: '7',
    name: 'Leather Wallet',
    description: 'Classic leather wallet',
    price: 49.99,
    image: 'https://images.unsplash.com/photo-1627123424574-724758594e93?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Accessories'
  },
  {
    id: '8',
    name: 'Dress Shoes',
    description: 'Formal leather shoes',
    price: 129.99,
    image: 'https://images.unsplash.com/photo-1533867617858-e7b97e060509?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Footwear'
  },
  {
    id: '9',
    name: 'Denim Jacket',
    description: 'Classic denim jacket',
    price: 99.99,
    image: 'https://images.unsplash.com/photo-1576871337632-b9aef4c17ab9?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Clothing'
  },
  {
    id: '10',
    name: 'Travel Bag',
    description: 'Spacious travel bag',
    price: 149.99,
    image: 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
    category: 'Accessories'
  }
];