"use client";
import { useEffect, useState } from "react";
import { Plus, Minus } from "lucide-react";
import { Button } from "@/components/ui/button";
// import { useCart } from "@/contexts/CartContext";
import Image from 'next/image'
import { toast } from "sonner";
import Link from "next/dist/client/link";
import useCatalogueStore from "@/lib/store/useCatalogueStore";
import useCartStore from "@/lib/store/useCartStore";

const ProductDetail = () => {
  const { catalogue } = useCatalogueStore();
  const [quantity, setQuantity] = useState(1);
  const [currentPrice, setCurrentPrice] = useState(
    catalogue!.configurables[0].price
  );
  const [selectedImage, setSelectedImage] = useState(catalogue?.images[0]);

  const { addCartItems, removeCartItems } = useCartStore();

  const handleQuantityChange = (increment: boolean) => {
    setQuantity((prev) => Math.max(1, increment ? prev + 1 : prev - 1));
  };

  if (!catalogue) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-semibold text-primary mb-4">
            Product Not Found
          </h1>
          <Link href="/">
            <Button>Back to Store</Button>
          </Link>
        </div>
      </div>
    );
  }

  const attributeTypes = catalogue.configurable_attributes
    ? Object.keys(catalogue.configurable_attributes)
    : [];

  const [selectedAttributes, setSelectedAttributes] = useState<
    Record<string, string>
  >({});

  const handleAttributeSelect = (attributeName: string, value: string) => {
    setSelectedAttributes((prev) => ({
      ...prev,
      [attributeName]: value,
    }));
  };
  const isAllOutOfStock = () => {
    return catalogue.configurables.every(
      (configurable) => configurable.stock === 0
    );
  };

  const getAllOptions = (attributeName: string) => {
    return Array.from(
      new Set(
        catalogue.configurables.map(
          (configurable) => configurable.attributes[attributeName]
        )
      )
    );
  };

  const isOptionSelectable = (attributeName: string, value: string) => {
    const filteredConfigs = catalogue.configurables.filter((configurable) => {
      return Object.entries(selectedAttributes).every(
        ([key, selectedValue]) =>
          key === attributeName ||
          configurable.attributes[key] === selectedValue
      );
    });

    return filteredConfigs.some(
      (configurable) =>
        configurable.attributes[attributeName] === value &&
        configurable.stock > 0
    );
  };

  useEffect(() => {
    const selectedProduct = catalogue.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes[attr] === selectedAttributes[attr]
      );
    });

    if (selectedProduct) {
      setCurrentPrice(selectedProduct.price);
    } else {
      setCurrentPrice(catalogue!.configurables[0].price);
    }
  }, [selectedAttributes, catalogue, attributeTypes]);

  const handleAddToCart = () => {
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

    const selectedProduct = catalogue.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes[attr] === selectedAttributes[attr]
      );
    });

    if (!selectedProduct || selectedProduct.stock === 0) {
      toast.error("Selected product is out of stock.");
      return;
    }

    addCartItems(
      {
        id: selectedProduct.id,
        name: selectedProduct.name,
        price: selectedProduct.price,
        image: selectedProduct.images[0].ImageUrl,
      },
      quantity
    );
    toast.success("Added to cart!");
  };

  const [api, setApi] = useState<any>(null);
  const [current, setCurrent] = useState(0);

  useEffect(() => {
    if (!api) return;

    const interval = setInterval(() => {
      api.scrollNext();
    }, 5000);

    api.on("select", () => {
      setCurrent(api.selectedScrollSnap());
    });

    return () => {
      clearInterval(interval);
      api.off("select");
    };
  }, [api]);

  return (
    <div className="min-h-screen bg-white">
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

      <main className="container mx-auto px-4 pt-24 pb-16">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12">
          <div className="grid grid-cols-1 gap-4">
            <div className="">
              <Image
                width={400}
                height={300}
                loading="lazy"
                src={selectedImage?.ImageUrl}
                alt={catalogue?.name}
                className="w-full h-full object-cover rounded-md"
              />
            </div>
          </div>

          <div className="flex flex-col mt-32">
            <h1 className="text-2xl md:text-3xl font-semibold text-primary mb-2">
              {catalogue.name}
            </h1>
            <p className="text-xl mb-4">${currentPrice}</p>
            <p className="text-md text-gray-600 mb-8">{catalogue.desc}</p>

            {attributeTypes.map((attributeName) => (
              <div key={attributeName} className="mb-8">
                <h3 className="text-sm font-medium mb-4">
                  {attributeName.charAt(0).toUpperCase() +
                    attributeName.slice(1)}
                </h3>
                <div className="flex flex-wrap gap-4">
                  {getAllOptions(attributeName).map((value) => {
                    const isOutOfStock = !isOptionSelectable(
                      attributeName,
                      value
                    );

                    return (
                      <button
                        key={value}
                        onClick={() =>
                          !isOutOfStock &&
                          handleAttributeSelect(attributeName, value)
                        }
                        className={`px-4 py-2 rounded-lg border-primary border uppercase text-sm transition-colors ${
                          selectedAttributes[attributeName] === value
                            ? "bg-primary text-white"
                            : isOutOfStock
                            ? "bg-gray-300 text-gray-600 cursor-not-allowed"
                            : "bg-secondary text-primary hover:bg-opacity-80"
                        }`}
                        disabled={isOutOfStock}
                      >
                        {value}
                      </button>
                    );
                  })}
                </div>
              </div>
            ))}

            <div className="mb-8">
              <h3 className="text-sm font-medium mb-4">Quantity</h3>
              <div className="flex items-center gap-4">
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleQuantityChange(false)}
                  disabled={quantity <= 1}
                >
                  <Minus className="w-4 h-4" />
                </Button>
                <input
                  type="number"
                  min="1"
                  value={quantity}
                  onChange={(e) => {
                    const value = parseInt(e.target.value);
                    if (!isNaN(value) && value >= 1) {
                      setQuantity(value);
                    }
                  }}
                  className="w-12 text-center text-lg border rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleQuantityChange(true)}
                >
                  <Plus className="w-4 h-4" />
                </Button>
              </div>
            </div>

            <Button
              onClick={handleAddToCart}
              className={`w-full ${
                isAllOutOfStock()
                  ? "bg-gray-300 text-gray-600 cursor-not-allowed"
                  : "bg-primary text-white hover:bg-primary/90"
              }`}
              size="lg"
              disabled={isAllOutOfStock()}
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
