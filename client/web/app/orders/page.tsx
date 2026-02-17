"use client";
import { useEffect, useState } from "react";
import { Loader2, Calendar, ChevronDown, ChevronUp, ExternalLink } from "lucide-react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Skeleton } from "@/components/ui/skeleton";
import { getCustomerByEmail } from "@/app/actions/customer";
import { ListCustomerOrders, type OrderResponse, type OrderItemResponse } from "@/app/actions/orders";
import { GetProduct, type ProductResponse } from "@/app/actions/product";
import { useAuth } from "@/context/AuthContext";
import { toast } from "sonner";

const statusColors: Record<string, string> = {
  PENDING: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20',
  CONFIRMED: 'bg-blue-500/10 text-blue-600 border-blue-500/20',
  PROCESSING: 'bg-purple-500/10 text-purple-600 border-purple-500/20',
  SHIPPED: 'bg-indigo-500/10 text-indigo-600 border-indigo-500/20',
  DELIVERED: 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20',
  CANCELLED: 'bg-red-500/10 text-red-600 border-red-500/20',
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
};

const formatCurrency = (amount?: number) => {
  if (!amount) return 'N/A';
  return new Intl.NumberFormat('en-ET', {
    style: 'currency',
    currency: 'ETB',
    minimumFractionDigits: 0,
  }).format(amount);
};

function getFirstImageUrl(attributes: any): string | null {
  if (!attributes) return null;
  let attrs = typeof attributes === 'string' ? (() => { try { return JSON.parse(attributes); } catch { return null; } })() : attributes;
  if (!attrs) return null;
  if (typeof attrs.image_url === 'string') return attrs.image_url;
  if (Array.isArray(attrs.images) && attrs.images.length > 0) {
    const first = attrs.images[0];
    if (typeof first === 'string') return first;
    return first?.ImageUrl ?? first?.url ?? first?.imageUrl ?? null;
  }
  if (typeof attrs.images === 'string') return attrs.images;
  if (attrs.images?.ImageUrl || attrs.images?.url || attrs.images?.imageUrl) {
    return attrs.images.ImageUrl || attrs.images.url || attrs.images.imageUrl;
  }
  return null;
}

const Orders = () => {
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedStatus, setSelectedStatus] = useState<string>('ALL');
  const [customerId, setCustomerId] = useState<number | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [expandedOrderItemKey, setExpandedOrderItemKey] = useState<string | null>(null);
  const [expandedProduct, setExpandedProduct] = useState<ProductResponse | null>(null);
  const [expandedProductLoading, setExpandedProductLoading] = useState(false);

  // Check if user is logged in based on localStorage
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const userEmail = localStorage.getItem('user_email');
      const savedPassword = localStorage.getItem('saved_password');
      setIsLoggedIn(!!(userEmail && savedPassword));
    }
  }, []);

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!isAuthenticated && !isLoggedIn) {
      router.push('/login');
    }
  }, [isAuthenticated, isLoggedIn, router]);

  // Fetch customer ID from API using logged-in user's email
  useEffect(() => {
    if (!isAuthenticated && !isLoggedIn) return;
    if (!user?.email) {
      const userEmail = typeof window !== 'undefined' 
        ? (localStorage.getItem('user_email') || localStorage.getItem('saved_email'))
        : null;
      if (!userEmail) {
        setIsLoading(false);
        return;
      }
      // Fetch customer by email
      const fetchCustomer = async () => {
        try {
          const customer = await getCustomerByEmail(userEmail);
          if (customer) {
            setCustomerId(customer.id);
          } else {
            setIsLoading(false);
          }
        } catch (error) {
          console.error('Error fetching customer:', error);
          setIsLoading(false);
        }
      };
      fetchCustomer();
    } else {
      const fetchCustomer = async () => {
        try {
          const customer = await getCustomerByEmail(user.email);
          if (customer) {
            setCustomerId(customer.id);
          } else {
            setIsLoading(false);
          }
        } catch (error) {
          console.error('Error fetching customer:', error);
          setIsLoading(false);
        }
      };
      fetchCustomer();
    }
  }, [isAuthenticated, isLoggedIn, user]);

  // Fetch orders when customer ID is available
  useEffect(() => {
    if (!customerId) return;

    const fetchOrders = async () => {
      setIsLoading(true);
      try {
        const params: {
          customer_id: number;
          status?: string;
          limit?: number;
          offset?: number;
        } = {
          customer_id: customerId,
          limit: 100,
          offset: 0,
        };

        if (selectedStatus !== 'ALL') {
          params.status = selectedStatus;
        }

        const response = await ListCustomerOrders(params);
        setOrders(response.orders || []);
      } catch (error: any) {
        console.error('Error fetching orders:', error);
        toast.error('Failed to load orders', {
          description: error?.message || 'An error occurred while fetching your orders.',
        });
        setOrders([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchOrders();
  }, [customerId, selectedStatus]);

  // Fetch product details when an order item is expanded
  useEffect(() => {
    if (!expandedOrderItemKey) {
      setExpandedProduct(null);
      return;
    }
    const [orderIdStr, itemIndexStr] = expandedOrderItemKey.split('-');
    const orderId = Number(orderIdStr);
    const itemIndex = Number(itemIndexStr);
    const order = orders.find((o) => o.id === orderId);
    if (!order?.cart_snapshot) {
      setExpandedProduct(null);
      return;
    }
    let displayItems: OrderItemResponse[] = [];
    try {
      let cartSnapshot: any = order.cart_snapshot;
      if (typeof cartSnapshot === 'string') {
        const trimmed = cartSnapshot.trim();
        if (trimmed && trimmed !== 'null' && trimmed !== '') cartSnapshot = JSON.parse(trimmed);
        else cartSnapshot = null;
      }
      if (Array.isArray(cartSnapshot) && cartSnapshot.length > 0) {
        displayItems = cartSnapshot
          .filter((item: any) => item != null && item !== undefined)
          .map((item: any) => ({
            product_id: item.product_id || item.productId || item.ProductID || 0,
            quantity: item.quantity || item.Quantity || 0,
            price: item.price || item.Price || undefined,
          }))
          .filter((item: OrderItemResponse) => item.product_id > 0 && item.quantity > 0);
      }
    } catch {
      setExpandedProduct(null);
      return;
    }
    const item = displayItems[itemIndex];
    if (!item?.product_id) {
      setExpandedProduct(null);
      return;
    }
    setExpandedProductLoading(true);
    setExpandedProduct(null);
    GetProduct(item.product_id)
      .then(setExpandedProduct)
      .catch(() => setExpandedProduct(null))
      .finally(() => setExpandedProductLoading(false));
  }, [expandedOrderItemKey, orders]);

  const statusFilters = [
    { value: 'ALL', label: 'All Orders' },
    { value: 'PENDING', label: 'Pending' },
    { value: 'DELIVERED', label: 'Delivered' },
    { value: 'CANCELLED', label: 'Cancelled' },
  ];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading orders...</p>
        </div>
      </div>
    );
  }

  const filteredOrders = orders.filter(order => {
    if (selectedStatus === 'ALL') return true;
    return order.status?.toUpperCase() === selectedStatus.toUpperCase();
  });

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Status Filters - Fixed Header (always visible when not loading) */}
      <div className="sticky top-0 z-10 bg-white px-4 py-4 border-b border-gray-200 overflow-x-auto">
        <div className="max-w-6xl mx-auto flex gap-2 min-w-max justify-center">
          {statusFilters.map((filter) => (
            <button
              key={filter.value}
              onClick={() => setSelectedStatus(filter.value)}
              className={`shrink-0 px-4 py-2 rounded-full text-xs font-medium transition-all ${
                selectedStatus === filter.value
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-muted text-muted-foreground hover:bg-muted/80'
              }`}
            >
              {filter.label}
            </button>
          ))}
        </div>
      </div>

      {/* Orders List */}
      <div className="max-w-6xl mx-auto px-4 py-6">
        {filteredOrders.length > 0 && (
          <div className="space-y-4">
            {filteredOrders.map((order) => (
              <div
                key={order.id}
                className="bg-white border border-gray-200 rounded-md overflow-hidden"
              >
                {/* Order Header */}
                <div className="p-4 border-b border-gray-200">
                  <div className="flex items-start justify-between mb-2">
                    <div>
                      <h3 className="font-semibold text-sm mb-1">Order #{order.id}</h3>
                      <div className="flex items-center gap-2 text-xs text-muted-foreground">
                        <Calendar className="w-3 h-3" />
                        <span>{formatDate(order.created_at)}</span>
                      </div>
                    </div>
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium border ${
                        statusColors[order.status?.toUpperCase() || 'PENDING'] ||
                        statusColors.PENDING
                      }`}
                    >
                      {order.status || 'PENDING'}
                    </span>
                  </div>
                </div>

                {/* Order Details */}
                <div className="p-4 space-y-3">
                  {/* Order Items from cart_snapshot */}
                  {(() => {
                    let displayItems: OrderItemResponse[] = [];
                    
                    if (order.cart_snapshot != null && order.cart_snapshot !== '' && order.cart_snapshot !== 'null') {
                      try {
                        let cartSnapshot = order.cart_snapshot;
                        
                        if (typeof cartSnapshot === 'string') {
                          const trimmed = cartSnapshot.trim();
                          if (trimmed && trimmed !== 'null' && trimmed !== '') {
                            cartSnapshot = JSON.parse(trimmed);
                          } else {
                            cartSnapshot = null;
                          }
                        }
                        
                        if (Array.isArray(cartSnapshot) && cartSnapshot.length > 0) {
                          displayItems = cartSnapshot
                            .filter((item: any) => item != null && item !== undefined)
                            .map((item: any) => ({
                              product_id: item.product_id || item.productId || item.ProductID || 0,
                              quantity: item.quantity || item.Quantity || 0,
                              price: item.price || item.Price || undefined,
                              selected_attributes: item.selected_attributes || item.SelectedAttributes,
                            }))
                            .filter((item: OrderItemResponse) => item.product_id > 0 && item.quantity > 0);
                        } else if (cartSnapshot && typeof cartSnapshot === 'object' && !Array.isArray(cartSnapshot)) {
                          const singleItem = {
                            product_id: cartSnapshot.product_id || cartSnapshot.productId || cartSnapshot.ProductID || 0,
                            quantity: cartSnapshot.quantity || cartSnapshot.Quantity || 0,
                            price: cartSnapshot.price || cartSnapshot.Price || undefined,
                            selected_attributes: cartSnapshot.selected_attributes || cartSnapshot.SelectedAttributes,
                          };
                          if (singleItem.product_id > 0 && singleItem.quantity > 0) {
                            displayItems = [singleItem];
                          }
                        }
                      } catch (error) {
                        console.error('Error parsing cart_snapshot for order', order.id, ':', error);
                      }
                    }

                    return displayItems.length > 0 ? (
                      <div className="space-y-3 pb-3 border-b border-gray-200">
                        <p className="text-xs font-semibold text-muted-foreground uppercase mb-1">
                          Order Items ({displayItems.length})
                        </p>
                        <div className="space-y-2">
                          {displayItems.map((item, itemIndex) => {
                            const itemKey = `${order.id}-${itemIndex}`;
                            const isExpanded = expandedOrderItemKey === itemKey;
                            return (
                              <div
                                key={itemIndex}
                                className="rounded-md overflow-hidden border border-transparent bg-muted/30"
                              >
                                <button
                                  type="button"
                                  onClick={() => setExpandedOrderItemKey((k) => (k === itemKey ? null : itemKey))}
                                  className="w-full flex items-start justify-between gap-3 p-2 hover:bg-muted/50 active:bg-muted/60 transition-colors text-left"
                                >
                                  <div className="flex-1 min-w-0">
                                    <p className="font-medium text-sm">Product #{item.product_id}</p>
                                    {item.selected_attributes && Object.keys(item.selected_attributes).length > 0 && (
                                      <p className="text-xs text-muted-foreground mt-0.5">
                                        {Object.entries(item.selected_attributes).map(([k, v]) => `${k}: ${v}`).join(", ")}
                                      </p>
                                    )}
                                    <p className="text-xs text-muted-foreground mt-0.5">
                                      Quantity: {item.quantity}
                                    </p>
                                  </div>
                                  <div className="flex items-center gap-1">
                                    {item.price && (
                                      <div className="text-right">
                                        <p className="font-semibold text-sm">{formatCurrency(item.price * item.quantity)}</p>
                                        <p className="text-xs text-muted-foreground">
                                          {formatCurrency(item.price)} × {item.quantity}
                                        </p>
                                      </div>
                                    )}
                                    {isExpanded ? (
                                      <ChevronUp className="w-4 h-4 shrink-0 text-muted-foreground" />
                                    ) : (
                                      <ChevronDown className="w-4 h-4 shrink-0 text-muted-foreground" />
                                    )}
                                  </div>
                                </button>
                                {isExpanded && (
                                  <div className="border-t border-gray-200 bg-muted/20 p-3 space-y-3">
                                    {expandedProductLoading ? (
                                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                        <Loader2 className="w-4 h-4 animate-spin shrink-0" />
                                        <span>Loading details...</span>
                                      </div>
                                    ) : expandedProduct ? (
                                      <>
                                        <div className="flex gap-3">
                                          {getFirstImageUrl(expandedProduct.attributes) ? (
                                            <div className="shrink-0 w-20 h-20 rounded-md overflow-hidden bg-gray-100 border border-gray-200">
                                              <img
                                                src={getFirstImageUrl(expandedProduct.attributes)!}
                                                alt={expandedProduct.name}
                                                className="w-full h-full object-cover"
                                              />
                                            </div>
                                          ) : null}
                                          <div className="min-w-0 flex-1 space-y-1">
                                            <p className="font-medium text-sm">{expandedProduct.name}</p>
                                            {expandedProduct.description && (
                                              <p className="text-xs text-muted-foreground">
                                                {expandedProduct.description}
                                              </p>
                                            )}
                                          </div>
                                        </div>
                                        {item.selected_attributes && typeof item.selected_attributes === 'object' && Object.keys(item.selected_attributes).length > 0 && (
                                          <div className="space-y-1">
                                            <p className="text-xs font-semibold text-muted-foreground uppercase">Selected options</p>
                                            <ul className="text-xs text-muted-foreground space-y-0.5">
                                              {Object.entries(item.selected_attributes).map(([key, value]) => (
                                                <li key={key} className="flex gap-2">
                                                  <span className="capitalize font-medium">{String(key).replace(/_/g, ' ')}:</span>
                                                  <span>{String(value)}</span>
                                                </li>
                                              ))}
                                            </ul>
                                          </div>
                                        )}
                                        <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground border-t border-gray-200 pt-2">
                                          <span>Quantity: {item.quantity}</span>
                                          {item.price != null && (
                                            <>
                                              <span>Unit price: {formatCurrency(item.price)}</span>
                                              <span>Line total: {formatCurrency(item.price * item.quantity)}</span>
                                            </>
                                          )}
                                        </div>
                                        <Link
                                          href={`/product/${item.product_id}`}
                                          className="inline-flex items-center gap-1.5 text-xs font-medium text-primary hover:underline"
                                        >
                                          View full details
                                          <ExternalLink className="w-3.5 h-3.5" />
                                        </Link>
                                      </>
                                    ) : (
                                      <Link
                                        href={`/product/${item.product_id}`}
                                        className="inline-flex items-center gap-1.5 text-xs font-medium text-primary hover:underline"
                                      >
                                        View product #{item.product_id}
                                        <ExternalLink className="w-3.5 h-3.5" />
                                      </Link>
                                    )}
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    ) : (
                      <div className="text-xs text-muted-foreground pb-3 border-b border-gray-200">
                        {order.cart_snapshot ? 'Unable to parse order items' : 'No items available'}
                      </div>
                    );
                  })()}

                  {/* Order Total */}
                  <div className="flex items-center justify-between pt-2">
                    <div className="text-sm font-semibold">
                      <span>Total</span>
                    </div>
                    <span className="font-semibold text-lg">
                      {formatCurrency(order.total)}
                    </span>
                  </div>

                  {/* Delivery Status */}
                  {order.delivery_status && (
                    <div className="flex items-center justify-between text-xs pt-2 border-t border-gray-200">
                      <span className="text-muted-foreground">Delivery Status</span>
                      <span className="font-medium capitalize">{order.delivery_status.toLowerCase()}</span>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default Orders;
