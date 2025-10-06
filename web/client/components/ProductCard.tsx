import useCatalogueStore from "@/lib/store/useCatalogueStore";
import { Catalogue } from "@/lib/types";
import Image from "next/image";
import Link from "next/link";

interface Props {
  product: Catalogue;
}

export const ProductCard: React.FC<Props> = ({ product }) => {
  const setProduct = useCatalogueStore((state) => state.setProduct);

  const getPriceRange = (product: Catalogue) => {
    if (!product.configurables || product.configurables.length === 0) {
      return `$${product.price}`;
    }

    const prices = product.configurables.map((config) => config.price);
    const minPrice = Math.min(...prices);
    const maxPrice = Math.max(...prices);

    return minPrice === maxPrice
      ? `${minPrice.toLocaleString()} ETB`
      : `${minPrice.toLocaleString()} ETB - ${maxPrice.toLocaleString()} ETB`;
  };

  const isAllOutOfStock = (product: Catalogue) => {
    return product.configurables.every(
      (configurable) => configurable.stock === 0
    );
  };

  return (
    <Link
      href={`/product/${product.name}`}
      className="group animate-fade-in"
      onClick={() => setProduct(product)}
    >
      <div className="aspect-square overflow-hidden rounded-lg bg-secondary mb-4">
        <Image
          width={400}
          height={400}
          src={product?.images?.[0]?.ImageUrl ?? ""}
          alt={product.name}
          className="w-full h-full object-cover transform transition-transform group-hover:scale-105"
          loading="lazy"
        />
      </div>
      <h3 className="text-lg font-semibold mb-2">{product.name}</h3>
      <p className="text-md text-gray-800">
        {getPriceRange(product)}
        {isAllOutOfStock(product) && (
          <span className="text-red-500 ml-2">(Out of Stock)</span>
        )}
      </p>
    </Link>
  );
};
