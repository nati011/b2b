"use client";
import { useEffect, useState } from "react";
import { Plus, Minus, ShoppingCart, Check, Star, Truck, Shield, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
// import { useCart } from "@/contexts/CartContext";
import Image from "next/image";
import { toast } from "sonner";
import Link from "next/dist/client/link";
import useCatalogueStore from "@/lib/store/useCatalogueStore";
import useCartStore from "@/lib/store/useCartStore";

const ProductDetail = () => {
  const { catalogue } = useCatalogueStore();
  const [quantity, setQuantity] = useState(1);
  const [currentPrice, setCurrentPrice] = useState(
    catalogue?.configurables?.[0]?.price || 0
  );
  const [selectedImage, setSelectedImage] = useState(catalogue?.images?.[0]);

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
    if (!catalogue?.configurables?.length) return;

    const selectedProduct = catalogue.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes[attr] === selectedAttributes[attr]
      );
    });

    if (selectedProduct) {
      setCurrentPrice(selectedProduct.price);
    } else if (catalogue.configurables[0]) {
      setCurrentPrice(catalogue.configurables[0].price);
    }
  }, [selectedAttributes, catalogue, attributeTypes]);

  const getSelectedProduct = () => {
    return catalogue.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes[attr] === selectedAttributes[attr]
      );
    });
  };

  const getStockStatus = () => {
    const product = getSelectedProduct();
    if (!product) return null;
    if (product.stock === 0) return { status: "out", text: "Out of Stock", variant: "destructive" as const };
    if (product.stock < 10) return { status: "low", text: "Low Stock", variant: "outline" as const };
    return { status: "in", text: "In Stock", variant: "success" as const };
  };

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

    const selectedProduct = getSelectedProduct();

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

  const isColorAttribute = (attrName: string) => {
    return attrName.toLowerCase().includes("color") || attrName.toLowerCase().includes("colour");
  };

  const getColorValue = (colorName: string) => {
    const colorMap: Record<string, string> = {
      black: "#000000",
      white: "#FFFFFF",
      red: "#EF4444",
      blue: "#3B82F6",
      green: "#10B981",
      yellow: "#F59E0B",
      purple: "#8B5CF6",
      pink: "#EC4899",
      gray: "#6B7280",
      grey: "#6B7280",
      silver: "#9CA3AF",
      gold: "#F59E0B",
    };
    return colorMap[colorName.toLowerCase()] || null;
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

  const stockStatus = getStockStatus();
  const selectedProduct = getSelectedProduct();

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
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

      <main className="container mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-16">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-16">
          {/* Image Gallery Section */}
          <div className="space-y-4">
            {/* Main Image */}
            <div className="relative aspect-square w-full overflow-hidden rounded-2xl bg-white border border-gray-200 shadow-lg group">
              <Image
                width={800}
                height={800}
                loading="lazy"
                src={selectedImage?.ImageUrl || catalogue?.images[0]?.ImageUrl || ""}
                alt={catalogue?.name || "Product"}
                className="object-contain w-full h-full transition-transform duration-300 group-hover:scale-105"
              />
            </div>

            {/* Thumbnail Gallery */}
            {catalogue?.images && catalogue.images.length > 1 && (
              <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
                {catalogue.images.map((image, index) => (
                  <button
                    key={index}
                    onClick={() => setSelectedImage(image)}
                    className={`relative flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden border-2 transition-all duration-200 ${
                      selectedImage?.ImageUrl === image.ImageUrl
                        ? "border-primary ring-2 ring-primary ring-offset-2 scale-105"
                        : "border-gray-200 hover:border-gray-300"
                    }`}
                  >
                    <Image
                      width={80}
                      height={80}
                      loading="lazy"
                      src={image.ImageUrl || ""}
                      alt={`${catalogue.name} - View ${index + 1}`}
                      className="object-cover w-full h-full"
                    />
                    {selectedImage?.ImageUrl === image.ImageUrl && (
                      <div className="absolute inset-0 bg-primary/10 flex items-center justify-center">
                        <Check className="w-5 h-5 text-primary" />
                      </div>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product Info Section */}
          <div className="flex flex-col space-y-6">
            {/* Product Title & Price */}
            <div className="space-y-3">
              <div className="flex items-start justify-between gap-4">
                <h1 className="text-3xl md:text-4xl font-bold text-gray-900 leading-tight">
                  {catalogue.name}
                </h1>
                {stockStatus && (
                  <Badge variant={stockStatus.variant} className="shrink-0">
                    {stockStatus.text}
                  </Badge>
                )}
              </div>
              
              <div className="flex items-baseline gap-3">
                <span className="text-4xl md:text-5xl font-bold text-primary">
                  {currentPrice.toLocaleString()}
                </span>
                <span className="text-xl text-gray-500">ETB</span>
              </div>

              {/* Rating & Reviews Placeholder */}
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <div className="flex items-center gap-1">
                  {[...Array(5)].map((_, i) => (
                    <Star key={i} className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                  ))}
                </div>
                <span>(4.5)</span>
                <span className="text-gray-400">•</span>
                <span>24 reviews</span>
              </div>
            </div>

            {/* Description */}
            <div className="prose max-w-none">
              <p className="text-gray-700 leading-relaxed text-base">
                {catalogue.desc}
              </p>
            </div>

            {/* Product Attributes */}
            {attributeTypes.map((attributeName) => {
              const isColor = isColorAttribute(attributeName);
              
              return (
                <div key={attributeName} className="space-y-3">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide">
                      {attributeName.charAt(0).toUpperCase() + attributeName.slice(1)}
                    </h3>
                    {selectedAttributes[attributeName] && (
                      <span className="text-xs text-gray-500">
                        Selected: <span className="font-medium">{selectedAttributes[attributeName]}</span>
                      </span>
                    )}
                  </div>
                  
                  <div className={`flex flex-wrap gap-3 ${isColor ? 'gap-2' : ''}`}>
                    {getAllOptions(attributeName).map((value) => {
                      const isOutOfStock = !isOptionSelectable(attributeName, value);
                      const isSelected = selectedAttributes[attributeName] === value;
                      const colorValue = isColor ? getColorValue(value) : null;

                      if (isColor && colorValue) {
                        return (
                          <button
                            key={value}
                            onClick={() => !isOutOfStock && handleAttributeSelect(attributeName, value)}
                            disabled={isOutOfStock}
                            className={`relative w-12 h-12 rounded-full border-2 transition-all duration-200 ${
                              isSelected
                                ? "ring-2 ring-primary ring-offset-2 scale-110 border-primary"
                                : "border-gray-300 hover:border-gray-400"
                            } ${isOutOfStock ? "opacity-40 cursor-not-allowed grayscale" : "hover:scale-105"}`}
                            style={{ backgroundColor: colorValue }}
                            title={value}
                          >
                            {isSelected && (
                              <div className="absolute inset-0 flex items-center justify-center">
                                <Check className="w-5 h-5 text-white drop-shadow-md" />
                              </div>
                            )}
                          </button>
                        );
                      }

                      return (
                        <button
                          key={value}
                          onClick={() => !isOutOfStock && handleAttributeSelect(attributeName, value)}
                          disabled={isOutOfStock}
                          className={`px-5 py-2.5 rounded-lg font-medium text-sm transition-all duration-200 ${
                            isSelected
                              ? "bg-primary text-white shadow-md scale-105"
                              : isOutOfStock
                              ? "bg-gray-100 text-gray-400 cursor-not-allowed border border-gray-200"
                              : "bg-white text-gray-700 border-2 border-gray-300 hover:border-primary hover:text-primary hover:shadow-sm"
                          }`}
                        >
                          {value}
                        </button>
                      );
                    })}
                  </div>
                </div>
              );
            })}

            {/* Stock Info */}
            {selectedProduct && (
              <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <p className="text-sm text-blue-900">
                  <span className="font-semibold">Availability:</span>{" "}
                  {selectedProduct.stock > 0
                    ? `${selectedProduct.stock} units available`
                    : "Currently out of stock"}
                </p>
              </div>
            )}

            {/* Quantity Selector */}
            <div className="space-y-3">
              <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide">
                Quantity
              </h3>
              <div className="flex items-center gap-4">
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleQuantityChange(false)}
                  disabled={quantity <= 1}
                  className="rounded-lg border-2 hover:bg-gray-50"
                >
                  <Minus className="w-4 h-4" />
                </Button>
                <input
                  type="number"
                  min="1"
                  max={selectedProduct?.stock || 999}
                  value={quantity}
                  onChange={(e) => {
                    const value = parseInt(e.target.value);
                    const maxStock = selectedProduct?.stock || 999;
                    if (!isNaN(value) && value >= 1 && value <= maxStock) {
                      setQuantity(value);
                    }
                  }}
                  className="w-16 text-center text-lg font-semibold border-2 border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary py-2"
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleQuantityChange(true)}
                  disabled={selectedProduct ? quantity >= selectedProduct.stock : false}
                  className="rounded-lg border-2 hover:bg-gray-50"
                >
                  <Plus className="w-4 h-4" />
                </Button>
              </div>
            </div>

            {/* Add to Cart Button */}
            <Button
              onClick={handleAddToCart}
              disabled={isAllOutOfStock() || !selectedProduct || selectedProduct.stock === 0}
              className="w-full h-14 text-lg font-semibold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              size="lg"
            >
              <ShoppingCart className="w-5 h-5 mr-2" />
              {isAllOutOfStock() || !selectedProduct || selectedProduct.stock === 0
                ? "Out of Stock"
                : "Add to Cart"}
            </Button>

            {/* Trust Badges */}
            <div className="grid grid-cols-3 gap-4 pt-4 border-t border-gray-200">
              <div className="flex flex-col items-center text-center">
                <Truck className="w-6 h-6 text-primary mb-2" />
                <span className="text-xs font-medium text-gray-700">Free Shipping</span>
              </div>
              <div className="flex flex-col items-center text-center">
                <RotateCcw className="w-6 h-6 text-primary mb-2" />
                <span className="text-xs font-medium text-gray-700">Easy Returns</span>
              </div>
              <div className="flex flex-col items-center text-center">
                <Shield className="w-6 h-6 text-primary mb-2" />
                <span className="text-xs font-medium text-gray-700">Secure Ordering</span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default ProductDetail;
