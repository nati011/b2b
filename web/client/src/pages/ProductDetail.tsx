import { useEffect, useState } from "react";
import { Link, useParams, useLocation } from "react-router-dom";
import { ArrowLeft, ShoppingCart, Plus, Minus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCart } from "@/contexts/CartContext";
import { toast } from "sonner";

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
    [key: string]: ConfigurableAttribute[];
  };
  configurables: Configurable[];
}

const ProductDetail = () => {
  const { id } = useParams();
  const [quantity, setQuantity] = useState(1);
  const { addToCart, getTotalItems } = useCart();
  const location = useLocation();
  const { product } = location.state as { product: Product };
  const [currentPrice, setCurrentPrice] = useState(product.price);

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
    // Check if all required attributes are selected
    const missingAttributes = attributeTypes.filter(
      (attr) => !selectedAttributes[attr]
    );

    if (missingAttributes.length > 0) {
      const errorMessage = `Please select ${missingAttributes.join(
        " and "
      )} before adding to cart.`;
      toast.error(errorMessage);
      return;
    }

    // Find the selected product configuration
    const selectedProduct = product.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes[attr] === selectedAttributes[attr]
      );
    });

    if (!selectedProduct) {
      toast.error("Invalid combination of attributes selected.");
      return;
    }

    addToCart(
      {
        id: selectedProduct.id,
        name: selectedProduct.name,
        price: selectedProduct.price,
        image: selectedProduct.images[0],
      },
      quantity
    );
    toast.success("Added to cart!");
  };
  const attributeTypes = product.configurable_attributes
    ? Object.keys(product.configurable_attributes)
    : [];

  // Function to get unique values for an attribute
  const getUniqueAttributeValues = (attributeName: string) => {
    if (!product.configurable_attributes[attributeName]) return [];
    return Array.from(
      new Set(
        product.configurable_attributes[attributeName].map(
          (attr) => attr.attribute_value
        )
      )
    );
  };

  // State for selected attributes
  const [selectedAttributes, setSelectedAttributes] = useState<
    Record<string, string>
  >({});

  const handleAttributeSelect = (attributeName: string, value: string) => {
    setSelectedAttributes((prev) => ({
      ...prev,
      [attributeName]: value,
    }));
  };
  useEffect(() => {
    // Only update price if all required attributes are selected
    if (
      attributeTypes.length > 0 &&
      attributeTypes.every((attr) => selectedAttributes[attr])
    ) {
      const selectedProduct = product.configurables.find((configurable) => {
        return attributeTypes.every(
          (attr) => configurable.attributes[attr] === selectedAttributes[attr]
        );
      });

      if (selectedProduct) {
        setCurrentPrice(selectedProduct.price);
      }
    } else {
      // Show base price if not all attributes are selected
      setCurrentPrice(product.price);
    }
  }, [selectedAttributes, product, attributeTypes]);

  return (
    <div className='min-h-screen bg-white'>
      <style>
        {`
          input[type="number"]::-webkit-inner-spin-button,
          input[type="number"]::-webkit-outer-spin-button {
            display: none;
          }

          input[type="number"] {
            -moz-appearance: textfield; /* Firefox */
          }
        `}
      </style>
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
            <p className='text-xl text-primary mb-8'>${currentPrice}</p>
            {attributeTypes.map((attributeName) => (
              <div key={attributeName} className='mb-8'>
                <h3 className='text-sm font-medium text-primary mb-4'>
                  {attributeName.charAt(0).toUpperCase() +
                    attributeName.slice(1)}
                </h3>
                <div className='flex flex-wrap gap-4'>
                  {getUniqueAttributeValues(attributeName).map((value) => (
                    <button
                      key={value}
                      onClick={() =>
                        handleAttributeSelect(attributeName, value)
                      }
                      className={`px-4 py-2 rounded-full text-sm transition-colors ${
                        selectedAttributes[attributeName] === value
                          ? "bg-primary text-white"
                          : "bg-secondary text-primary hover:bg-opacity-80"
                      }`}
                    >
                      {value}
                    </button>
                  ))}
                </div>
              </div>
            ))}

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
                <input
                  type='number'
                  min='1'
                  value={quantity}
                  onChange={(e) => {
                    const value = parseInt(e.target.value);
                    if (!isNaN(value) && value >= 1) {
                      setQuantity(value);
                    }
                  }}
                  className='w-12 text-center text-lg border rounded-md focus:outline-none focus:ring-1 focus:ring-primary'
                />
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
