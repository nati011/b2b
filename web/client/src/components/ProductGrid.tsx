import { useState, useEffect } from "react";
import { Filter, ShoppingCart, Search, User } from "lucide-react";
import { Link } from "react-router-dom";
import { MdOutlineKeyboardDoubleArrowRight } from "react-icons/md";
import { fetchCategories } from "@/api/CategoryApi";
import { Skeleton } from "@/components/ui/skeleton"
import { useCart } from "@/contexts/CartContext";
import { fetchProducts } from "@/api/ProductApi";

export const ProductGrid = () => {
    const [categories, setCategories] = useState([{ id: 0, name: "All" }]);
    const [loading, setLoading] = useState(false)
    const [products, setProducts] = useState([]);
    const [selectedCategory, setSelectedCategory] = useState(0);
    const [searchQuery, setSearchQuery] = useState("");
    const [currentPage, setCurrentPage] = useState(1);
    const [showSearch, setShowSearch] = useState(false);
    const { getTotalItems } = useCart();
    const productsPerPage = 9;

    useEffect(() => {
        const loadCategories = async () => {
            try {
                const fetchedCategories = await fetchCategories();
                setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
            } catch (error) {
                console.error("Error fetching categories:", error);
            }
        };

        loadCategories();
    }, []);

    useEffect(() => {
        const loadProducts = async () => {
            setLoading(true)
            try {
                const productList = await fetchProducts();
                setProducts(productList);
            } catch (error) {
                console.error("Failed to load products:", error);
            } finally {
                setLoading(false)
            }
        };

        loadProducts();
    }, []);

    const getPriceRange = (product) => {
        if (!product.configurables || product.configurables.length === 0) {
            return `$${product.price}`;
        }

        const prices = product.configurables.map((config) => config.price);
        const minPrice = Math.min(...prices);
        const maxPrice = Math.max(...prices);

        if (minPrice === maxPrice) {
            return `$${minPrice}`;
        } else {
            return `$${minPrice} - $${maxPrice}`;
        }
    };

    const startIndex = (currentPage - 1) * productsPerPage;

    // Filter products based on selected category and search query
    const filteredProducts = products.filter((product) => {
        const matchesCategory =
            selectedCategory === 0 || // "All" category
            product.configurables.some((configurable) =>
                configurable.categories.includes(selectedCategory)
            );

        const matchesSearchQuery = product.name
            .toLowerCase()
            .includes(searchQuery.toLowerCase());

        return matchesCategory && matchesSearchQuery;
    });

    const totalPages = Math.ceil(filteredProducts.length / productsPerPage);
    const displayedProducts = filteredProducts.slice(
        startIndex,
        startIndex + 4
    );
    const isAllOutOfStock = (product) => {
        return product.configurables.every(
            (configurable) => configurable.stock === 0
        );
    };


    return (
        <section id="products" className="container mx-auto px-4 py-16">
            <div className="flex flex-col items-center mb-12">
                <h2 className="text-3xl font-medium mb-4">Featured Products</h2>
                <div className="h-1 w-20 bg-primary mb-8"></div>
                <div className='flex items-center justify-between mb-8 w-full'>
                    <div className="flex gap-4 items-center">
                        <Filter className='w-5 h-5 shrink-0' />
                        <div className='flex gap-4 overflow-x-auto no-scrollbar w-full'>
                            {categories.map((category) => (
                                <button
                                    key={category.id}
                                    onClick={() => setSelectedCategory(category.id)}
                                    className={`shrink-0 px-4 py-2 rounded-full text-sm transition-colors whitespace-nowrap ${selectedCategory === category.id
                                        ? "bg-primary text-white"
                                        : "bg-secondary text-primary hover:bg-opacity-80"
                                        }`}
                                >
                                    {category.name}
                                </button>
                            ))}
                        </div>
                    </div>
                    <div className="text-blue-900 font-semibold flex items-center gap-2">
                        <Link to="/product">
                            Load All
                        </Link>
                        <MdOutlineKeyboardDoubleArrowRight />
                    </div>


                </div>

            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-6">
                {
                    (loading || displayedProducts.length == 0) ? (
                        <>
                            <div className="flex flex-col space-y-3">
                                <Skeleton className="h-[300px] w-[300px]" />
                                <div className="space-y-2">
                                    <Skeleton className="h-4 w-[250px]" />
                                    <Skeleton className="h-4 w-[200px]" />
                                </div>
                            </div>
                            <div className="flex flex-col space-y-3">
                                <Skeleton className="h-[300px] w-[300px]" />
                                <div className="space-y-2">
                                    <Skeleton className="h-4 w-[250px]" />
                                    <Skeleton className="h-4 w-[200px]" />
                                </div>
                            </div>
                            <div className="flex flex-col space-y-3">
                                <Skeleton className="h-[300px] w-[300px]" />
                                <div className="space-y-2">
                                    <Skeleton className="h-4 w-[250px]" />
                                    <Skeleton className="h-4 w-[200px]" />
                                </div>
                            </div>
                            <div className="flex flex-col space-y-3">
                                <Skeleton className="h-[300px] w-[300px]" />
                                <div className="space-y-2">
                                    <Skeleton className="h-4 w-[250px]" />
                                    <Skeleton className="h-4 w-[200px]" />
                                </div>
                            </div>

                        </>

                    ) : (
                        <>
                            {displayedProducts.map((product, index) => (
                                <Link
                                    key={product.id}
                                    to={`/product/${product.id}`}
                                    state={{ product }}
                                    className='group animate-fade-in'
                                >
                                    <div className='aspect-square overflow-hidden rounded-lg bg-secondary mb-4'>
                                        <img

                                            src={product.images[0].ImageUrl}
                                            alt={product.name}
                                            className='w-full h-full object-cover transform transition-transform group-hover:scale-105'
                                        />
                                    </div>
                                    <h3 className='text-lg font-medium text-primary mb-2'>
                                        {product.name}
                                    </h3>
                                    <p className='text-sm text-primary'>
                                        {getPriceRange(product)}
                                        {isAllOutOfStock(product) && (
                                            <span className='text-red-500 ml-2'>(Out of Stock)</span>
                                        )}
                                    </p>
                                </Link>
                            ))}
                        </>
                    )
                }

            </div>

        </section>
    );
};
