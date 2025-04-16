import { useState } from "react";
import { Link, useParams, useLocation } from "react-router-dom";
import { ArrowLeft, ShoppingCart, Plus, Minus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCart } from "@/contexts/CartContext";
import { toast } from "sonner";

// const mockProducts = [
//   {
//     id: 1,
//     name: "Wireless Headphones",
//     price: 199.99,
//     description: "Premium wireless headphones with active noise cancellation and up to 30 hours of battery life.",
//     images: ["https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500&q=80"],
//     type: "electronics",
//     colors: ["Black", "White", "Blue"],
//   },
//   {
//     id: 2,
//     name: "Smart Watch",
//     price: 299.99,
//     description: "Advanced smartwatch with health tracking features and AMOLED display.",
//     images: ["https://images.unsplash.com/photo-1546868871-7041f2a55e12?w=500&q=80"],
//     type: "electronics",
//     colors: ["Black", "Silver", "Gold"],
//   },
//   {
//     id: 3,
//     name: "Cotton T-Shirt",
//     price: 29.99,
//     description: "Comfortable 100% cotton t-shirt for everyday wear.",
//     images: ["https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500&q=80"],
//     type: "clothing",
//     sizes: ["XS", "S", "M", "L", "XL"],
//     colors: ["White", "Black", "Gray", "Navy"],
//   },
//   {
//     id: 4,
//     name: "Leather Wallet",
//     price: 49.99,
//     description: "Genuine leather wallet with multiple card slots and coin pocket.",
//     images: ["https://images.unsplash.com/photo-1627123424574-724758594e93?w=500&q=80"],
//     type: "accessories",
//     colors: ["Brown", "Black"],
//   },
//   {
//     id: 5,
//     name: "Running Shoes",
//     price: 89.99,
//     description: "Lightweight running shoes with responsive cushioning.",
//     images: ["https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=500&q=80"],
//     type: "footwear",
//     sizes: ["7", "8", "9", "10", "11", "12"],
//     colors: ["Black/Red", "Blue/White", "Gray/Yellow"],
//   },
//   {
//     id: 6,
//     name: "Wireless Earbuds",
//     price: 159.99,
//     description: "True wireless earbuds with premium sound quality.",
//     images: ["https://images.unsplash.com/photo-1572569511254-d8f925fe2cbb?w=500&q=80"],
//     type: "electronics",
//     colors: ["White", "Black"],
//   },
//   {
//     id: 7,
//     name: "Denim Jeans",
//     price: 79.99,
//     description: "Classic fit denim jeans with stretch comfort.",
//     images: ["https://images.unsplash.com/photo-1542272604-787c3835535d?w=500&q=80"],
//     type: "clothing",
//     sizes: ["30x30", "32x32", "34x32", "36x32"],
//     colors: ["Blue", "Black", "Gray"],
//   },
//   {
//     id: 8,
//     name: "Backpack",
//     price: 69.99,
//     description: "Durable backpack with laptop compartment and multiple pockets.",
//     images: ["https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=500&q=80"],
//     type: "accessories",
//     colors: ["Black", "Navy", "Gray"],
//   },
//   {
//     id: 9,
//     name: "Smart Speaker",
//     price: 129.99,
//     description: "Voice-controlled smart speaker with premium sound.",
//     images: ["https://images.unsplash.com/photo-1543512214-318c7553f230?w=500&q=80"],
//     type: "electronics",
//     colors: ["Black", "White"],
//   },
//   {
//     id: 10,
//     name: "Summer Dress",
//     price: 59.99,
//     description: "Lightweight summer dress with floral pattern.",
//     images: ["https://images.unsplash.com/photo-1585487000160-6ebcfceb0d03?w=500&q=80"],
//     type: "clothing",
//     sizes: ["XS", "S", "M", "L"],
//     colors: ["Blue Floral", "Pink Floral", "White"],
//   },
//   {
//     id: 11,
//     name: "Sunglasses",
//     price: 149.99,
//     description: "Polarized sunglasses with UV protection.",
//     images: ["https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=500&q=80"],
//     type: "accessories",
//     colors: ["Black/Gold", "Tortoise/Brown"],
//   },
//   {
//     id: 12,
//     name: "Fitness Tracker",
//     price: 89.99,
//     description: "Water-resistant fitness tracker with heart rate monitoring.",
//     images: ["https://images.unsplash.com/photo-1575311373937-040b8e1fd5b6?w=500&q=80"],
//     type: "electronics",
//     colors: ["Black", "Blue", "Pink"],
//   },
// ];
interface ConfigurableAttribute {
  product_id: number;
  attribute_value: string;
}

interface Configurable {
  id: number;
  name: string;
  desc: string;
  price: number;
  stock: number;
  external_id: string;
  attributes: {
    color: string;
    size: string;
  };
  images: string[];
  distributor_id: number;
  categories: number[];
  is_active: boolean;
}

interface Product {
  id: number;
  name: string;
  price: number;
  desc: string;
  is_active: boolean;
  images: string[];
  configurable_attributes: {
    color: ConfigurableAttribute[];
    size: ConfigurableAttribute[];
  };
  configurables: Configurable[];
}

const ProductDetail = () => {
  const { id } = useParams();
  const [selectedSize, setSelectedSize] = useState("");
  const [selectedColor, setSelectedColor] = useState("");
  const [quantity, setQuantity] = useState(1);
  const { addToCart, getTotalItems } = useCart();
  const location = useLocation();
  const { product } = location.state as { product: Product };

  if (!product) {
    return (
      <div className='min-h-screen bg-white flex items-center justify-center'>
        <div className='text-center'>
          <h1 className='text-2xl font-semibold text-primary mb-4'>
            Product Not Found
          </h1>
          <Link to='/'>
            <Button>Back to Store</Button>
          </Link>
        </div>
      </div>
    );
  }

  const handleQuantityChange = (increment) => {
    setQuantity((prev) => Math.max(1, increment ? prev + 1 : prev - 1));
  };

  const handleAddToCart = () => {
    // Check if both attributes are required and selected
    const errors = [];
    if (
      product.configurable_attributes.size &&
      product.configurable_attributes.size.length > 0 &&
      !selectedSize
    ) {
      errors.push("size");
    }
    if (
      product.configurable_attributes.color &&
      product.configurable_attributes.color.length > 0 &&
      !selectedColor
    ) {
      errors.push("color");
    }

    if (errors.length > 0) {
      const errorMessage = `Please select ${errors.join(
        " and "
      )} before adding to cart.`;
      toast.error(errorMessage);
      return;
    }

    // Find the selected product configuration based on selected size and color
    const selectedProduct = product.configurables.find((configurable) => {
      const matchesSize =
        !selectedSize || configurable.attributes.size === selectedSize;
      const matchesColor =
        !selectedColor || configurable.attributes.color === selectedColor;
      return matchesSize && matchesColor;
    });

    if (!selectedProduct) {
      toast.error("Invalid combination of size and color selected.");
      return;
    }

    addToCart(
      {
        id: selectedProduct.id,
        name: selectedProduct.name, // Construct name dynamically
        price: selectedProduct.price,
        image: selectedProduct.images[0],
      },
      quantity
    );
    console.log(selectedProduct);
    toast.success("Added to cart!");
  };

  const uniqueColors = product.configurable_attributes.color
    ? Array.from(
        new Set(
          product.configurable_attributes.color.map((c) => c.attribute_value)
        )
      )
    : [];

  const uniqueSizes = product.configurable_attributes.size
    ? Array.from(
        new Set(
          product.configurable_attributes.size.map((s) => s.attribute_value)
        )
      )
    : [];

  return (
    <div className='min-h-screen bg-white'>
      <header className='fixed top-0 left-0 right-0 bg-white z-50 border-b border-gray-100'>
        <nav className='container mx-auto px-4 py-4 flex justify-between items-center'>
          <Link
            to='/'
            className='flex items-center gap-2 text-primary hover:opacity-80'
          >
            <ArrowLeft className='w-5 h-5' />
            <span>Back</span>
          </Link>
          <Link
            to='/cart'
            className='relative p-2 hover:bg-secondary rounded-full transition-colors'
          >
            <ShoppingCart className='w-5 h-5' />
            {getTotalItems() > 0 && (
              <span className='absolute -top-1 -right-1 bg-primary text-white text-xs w-5 h-5 rounded-full flex items-center justify-center'>
                {getTotalItems()}
              </span>
            )}
          </Link>
        </nav>
      </header>

      <main className='container mx-auto px-4 pt-24 pb-16'>
        <div className='grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12'>
          <div className='aspect-square bg-secondary rounded-lg overflow-hidden'>
            <img
              src={product.images[0]}
              alt={product.name}
              className='w-full h-full object-cover'
            />
          </div>

          <div className='flex flex-col'>
            <h1 className='text-2xl md:text-3xl font-semibold text-primary mb-4'>
              {product.name}
            </h1>
            <p className='text-xl text-primary mb-8'>${product.price}</p>
            {/* Size Selection */}
            {uniqueSizes.length > 0 && (
              <div className='mb-8'>
                <h3 className='text-sm font-medium text-primary mb-4'>Size</h3>
                <div className='flex flex-wrap gap-4'>
                  {uniqueSizes.map((size) => (
                    <button
                      key={size}
                      onClick={() => setSelectedSize(size)}
                      className={`w-12 h-12 rounded-full flex items-center justify-center text-sm transition-colors ${
                        selectedSize === size
                          ? "bg-primary text-white"
                          : "bg-secondary text-primary hover:bg-opacity-80"
                      }`}
                    >
                      {size}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Color Selection */}
            {uniqueColors.length > 0 && (
              <div className='mb-8'>
                <h3 className='text-sm font-medium text-primary mb-4'>Color</h3>
                <div className='flex flex-wrap gap-4'>
                  {uniqueColors.map((color) => (
                    <button
                      key={color}
                      onClick={() => setSelectedColor(color)}
                      className={`px-4 py-2 rounded-full text-sm transition-colors ${
                        selectedColor === color
                          ? "bg-primary text-white"
                          : "bg-secondary text-primary hover:bg-opacity-80"
                      }`}
                    >
                      {color}
                    </button>
                  ))}
                </div>
              </div>
            )}

            <div className='mb-8'>
              <h3 className='text-sm font-medium text-primary mb-4'>
                Quantity
              </h3>
              <div className='flex items-center gap-4'>
                <Button
                  variant='outline'
                  size='icon'
                  onClick={() => handleQuantityChange(false)}
                  disabled={quantity <= 1}
                >
                  <Minus className='w-4 h-4' />
                </Button>
                <span className='w-12 text-center text-lg'>{quantity}</span>
                <Button
                  variant='outline'
                  size='icon'
                  onClick={() => handleQuantityChange(true)}
                >
                  <Plus className='w-4 h-4' />
                </Button>
              </div>
            </div>

            <p className='text-sm text-gray-600 mb-8'>{product.desc}</p>

            <Button
              onClick={handleAddToCart}
              className='w-full bg-primary text-white hover:bg-primary/90'
              size='lg'
            >
              Add to Cart
            </Button>
          </div>
        </div>
      </main>
    </div>
  );
};

export default ProductDetail;
