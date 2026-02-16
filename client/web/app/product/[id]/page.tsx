"use client";
import { useEffect, useState } from "react";
import { Plus, Minus, ShoppingCart, Check, Truck, Shield, RotateCcw, ArrowLeft } from "lucide-react";
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
    catalogue?.configurables?.[0]?.price || catalogue?.price || 0
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

  // Filter out category_ids, quantity fields, and other non-attribute fields from configurable_attributes
  const attributeTypes = catalogue.configurable_attributes
    ? Object.keys(catalogue.configurable_attributes).filter(
        (key) => 
          key !== 'category_ids' && 
          key !== 'images' &&
          key !== 'total_quantity' &&
          key !== 'TOTAL_QUANTITY' &&
          key !== 'reserved_quantity' &&
          key !== 'RESERVED_QUANTITY' &&
          key !== 'available_quantity' &&
          key !== 'AVAILABLE_QUANTITY'
      )
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

  // Auto-select first available option for each attribute
  useEffect(() => {
    if (!catalogue || attributeTypes.length === 0) return;

    // Check if we already have all attributes selected
    const allSelected = attributeTypes.every(
      (attr) => selectedAttributes[attr]
    );
    if (allSelected) return;

    // Build selections progressively, considering previously selected attributes
    const newSelections: Record<string, string> = {};
    let currentSelections = { ...selectedAttributes };

    attributeTypes.forEach((attributeName) => {
      // Skip if already selected
      if (currentSelections[attributeName]) return;

      // Get all options for this attribute
      const options = getAllOptions(attributeName);
      
      // Find the first selectable option based on current selections
      for (const option of options) {
        // Check if this option is selectable given current selections
        let isSelectable = false;
        
        if (!catalogue.configurables || catalogue.configurables.length === 0) {
          isSelectable = catalogue.is_active;
        } else {
          const filteredConfigs = catalogue.configurables.filter((configurable) => {
            return Object.entries(currentSelections).every(
              ([key, selectedValue]) =>
                key === attributeName ||
                configurable.attributes?.[key] === selectedValue
            );
          });

          isSelectable = filteredConfigs.some(
            (configurable) =>
              configurable.attributes?.[attributeName] === option &&
              configurable.stock > 0
          );
        }

        if (isSelectable) {
          newSelections[attributeName] = option;
          currentSelections[attributeName] = option; // Update for next iteration
          break;
        }
      }
    });

    // Only update if we found new selections
    if (Object.keys(newSelections).length > 0) {
      setSelectedAttributes((prev) => ({
        ...prev,
        ...newSelections,
      }));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [catalogue, attributeTypes.join(',')]);
  
  const isAllOutOfStock = () => {
    // If no configurables, check if product is active and has stock
    if (!catalogue.configurables || catalogue.configurables.length === 0) {
      return !catalogue.is_active;
    }
    return catalogue.configurables.every(
      (configurable) => configurable.stock === 0
    );
  };

  const getAllOptions = (attributeName: string) => {
    // If no configurables, get options from configurable_attributes
    if (!catalogue.configurables || catalogue.configurables.length === 0) {
      const attrValue = catalogue.configurable_attributes?.[attributeName];
      if (!attrValue) return [];
      return Array.isArray(attrValue) ? attrValue : [attrValue];
    }
    
    return Array.from(
      new Set(
        catalogue.configurables
          .map((configurable) => configurable.attributes?.[attributeName])
          .filter((val) => val !== undefined && val !== null)
      )
    );
  };

  const isOptionSelectable = (attributeName: string, value: string) => {
    // If no configurables, all options are selectable if product is active
    if (!catalogue.configurables || catalogue.configurables.length === 0) {
      return catalogue.is_active;
    }
    
    const filteredConfigs = catalogue.configurables.filter((configurable) => {
      return Object.entries(selectedAttributes).every(
        ([key, selectedValue]) =>
          key === attributeName ||
          configurable.attributes?.[key] === selectedValue
      );
    });

    return filteredConfigs.some(
      (configurable) =>
        configurable.attributes?.[attributeName] === value &&
        configurable.stock > 0
    );
  };

  useEffect(() => {
    // If no configurables, use base product price
    if (!catalogue?.configurables || catalogue.configurables.length === 0) {
      setCurrentPrice(catalogue?.price || 0);
      return;
    }

    const selectedProduct = catalogue.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes?.[attr] === selectedAttributes[attr]
      );
    });

    if (selectedProduct) {
      setCurrentPrice(selectedProduct.price);
    } else if (catalogue.configurables[0]) {
      setCurrentPrice(catalogue.configurables[0].price);
    }
  }, [selectedAttributes, catalogue, attributeTypes]);

  const getSelectedProduct = () => {
    // If no configurables, return a product from the base catalogue
    if (!catalogue.configurables || catalogue.configurables.length === 0) {
      // Get available quantity from configurable_attributes if stored there
      const availableQuantity = catalogue.configurable_attributes?.available_quantity;
      // Use actual available_quantity if available, otherwise check is_active
      const stock = availableQuantity !== undefined && availableQuantity !== null
        ? availableQuantity 
        : (catalogue.is_active ? 999 : 0);
      
      return {
        id: catalogue.id,
        name: catalogue.name,
        price: catalogue.price || 0,
        stock: stock,
        images: catalogue.images || [],
        attributes: catalogue.configurable_attributes || {}
      };
    }
    
    return catalogue.configurables.find((configurable) => {
      return attributeTypes.every(
        (attr) => configurable.attributes?.[attr] === selectedAttributes[attr]
      );
    }) || catalogue.configurables[0]; // Fallback to first configurable if none matches
  };
  
  // Get stock information for display
  const getStockInfo = () => {
    const product = getSelectedProduct();
    if (!product) {
      return {
        stock: 0,
        isAvailable: false,
        message: "Product not available"
      };
    }
    
    const stock = product.stock || 0;
    return {
      stock: stock,
      isAvailable: stock > 0 && catalogue.is_active,
      message: stock > 0 
        ? `${stock} ${stock === 1 ? 'unit' : 'units'} available`
        : "Currently out of stock"
    };
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

    // Ensure we have a valid id
    const productId = selectedProduct.id ?? catalogue?.id;
    if (typeof productId !== 'number') {
      toast.error("Product ID is missing. Please try again.");
      return;
    }

    addCartItems(
      {
        id: productId,
        name: selectedProduct.name,
        price: selectedProduct.price,
        image: selectedProduct.images?.[0]?.ImageUrl || catalogue.images?.[0]?.ImageUrl || '',
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
        <Link
          href="/product"
          className="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 mb-6 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to products
        </Link>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-16">
          {/* Image Gallery Section */}
          <div className="space-y-4">
            {/* Main Image */}
            <div className="relative aspect-square w-full overflow-hidden rounded-2xl bg-white border border-gray-200 shadow-lg group">
              {catalogue?.images && catalogue.images.length > 0 ? (
                <Image
                  width={800}
                  height={800}
                  loading="lazy"
                  src={selectedImage?.ImageUrl || catalogue?.images[0]?.ImageUrl || ""}
                  alt={catalogue?.name || "Product"}
                  className="object-contain w-full h-full transition-transform duration-300 group-hover:scale-105"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center bg-gray-100 text-gray-400">
                  No Image Available
                </div>
              )}
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
            <div className="space-y-2">
              <h1 className="text-2xl md:text-3xl font-bold text-gray-900 leading-tight">
                {catalogue.name}
              </h1>
              
              <div className="flex items-baseline gap-2 flex-wrap">
                <span className="text-2xl md:text-3xl font-bold text-primary">
                  {currentPrice.toLocaleString()}
                </span>
                <span className="text-base text-gray-500">ETB</span>
                {stockStatus && (
                  <Badge variant={stockStatus.variant} className="shrink-0 self-center">
                    {stockStatus.text}
                  </Badge>
                )}
              </div>
            </div>

            {/* Description */}
            <div className="prose max-w-none">
              <p className="text-gray-700 leading-relaxed text-base">
                {catalogue.desc}
              </p>
            </div>

            {/* Product Attributes - e-commerce style grid */}
            <div className="grid gap-4">
              {attributeTypes.map((attributeName) => {
                const isColor = isColorAttribute(attributeName);
                const options = getAllOptions(attributeName);

                return (
                  <div
                    key={attributeName}
                    className="grid grid-cols-1 sm:grid-cols-[minmax(0,8rem)_1fr] gap-3 sm:gap-6 py-4 border-b border-gray-100 last:border-b-0"
                  >
                    <div className="flex items-center sm:pt-0.5">
                      <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">
                        {attributeName.charAt(0).toUpperCase() + attributeName.slice(1)}
                      </span>
                    </div>
                    <div className="grid grid-cols-[repeat(auto-fill,minmax(4.5rem,1fr))] gap-2 sm:gap-2 min-w-0">
                      {options.map((value) => {
                        const isOutOfStock = !isOptionSelectable(attributeName, value);
                        const isSelected = selectedAttributes[attributeName] === value;
                        const colorValue = isColor ? getColorValue(value) : null;

                        if (isColor && colorValue) {
                          return (
                            <button
                              key={value}
                              onClick={() => !isOutOfStock && handleAttributeSelect(attributeName, value)}
                              disabled={isOutOfStock}
                              className={`relative aspect-square max-w-[2.5rem] w-full rounded-full border-2 transition-all duration-200 ${
                                isSelected
                                  ? "ring-2 ring-primary ring-offset-1 border-primary"
                                  : "border-gray-200 hover:border-gray-300"
                              } ${isOutOfStock ? "opacity-40 cursor-not-allowed grayscale" : ""}`}
                              style={{ backgroundColor: colorValue }}
                              title={value}
                              type="button"
                            >
                              {isSelected && (
                                <div className="absolute inset-0 flex items-center justify-center">
                                  <Check className="w-4 h-4 text-white drop-shadow-md" />
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
                            type="button"
                            className={`min-w-0 py-2 px-3 rounded-md text-xs font-medium transition-all duration-200 text-center truncate ${
                              isSelected
                                ? "bg-primary text-white border border-primary"
                                : isOutOfStock
                                ? "bg-gray-50 text-gray-400 cursor-not-allowed border border-gray-100"
                                : "bg-white text-gray-700 border border-gray-200 hover:border-primary hover:text-primary"
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
            </div>

            {/* Stock Info */}
            {(() => {
              const stockInfo = getStockInfo();
              return (
                <div className={`p-4 rounded-lg border ${
                  stockInfo.isAvailable 
                    ? 'bg-emerald-50 border-emerald-200' 
                    : 'bg-red-50 border-red-200'
                }`}>
                  <div className="flex items-center justify-between gap-3">
                    <div className="flex items-center gap-2">
                      <Check className={`w-5 h-5 ${
                        stockInfo.isAvailable ? 'text-emerald-600' : 'text-red-600'
                      }`} />
                      <div>
                        <p className={`text-sm font-semibold ${
                          stockInfo.isAvailable ? 'text-emerald-900' : 'text-red-900'
                        }`}>
                          Stock Availability
                        </p>
                        <p className={`text-sm ${
                          stockInfo.isAvailable ? 'text-emerald-700' : 'text-red-700'
                        }`}>
                          {stockInfo.message}
                        </p>
                      </div>
                    </div>
                    {stockInfo.isAvailable && stockInfo.stock < 10 && (
                      <Badge variant="outline" className="bg-yellow-50 border-yellow-300 text-yellow-800">
                        Low Stock
                      </Badge>
                    )}
                  </div>
                </div>
              );
            })()}

            {/* Quantity Selector */}
            <div className="space-y-2">
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">
                Quantity
              </span>
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
          </div>
        </div>
      </main>
    </div>
  );
};

export default ProductDetail;
