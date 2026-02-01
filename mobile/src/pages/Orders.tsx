import { motion } from 'framer-motion';
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Loader2, Calendar, DollarSign } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { ListCustomerOrders, type OrderResponse, type OrderItemResponse } from '@/lib/api/order';
import { GetCustomerByEmail } from '@/lib/api/customer';
import { toast } from 'sonner';

const statusColors: Record<string, string> = {
  PENDING: 'bg-yellow-500/10 text-yellow-600 border-yellow-500/20',
  CONFIRMED: 'bg-blue-500/10 text-blue-600 border-blue-500/20',
  PROCESSING: 'bg-purple-500/10 text-purple-600 border-purple-500/20',
  SHIPPED: 'bg-indigo-500/10 text-indigo-600 border-indigo-500/20',
  DELIVERED: 'bg-green-500/10 text-green-600 border-green-500/20',
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

const Orders = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedStatus, setSelectedStatus] = useState<string>('ALL');
  const [customerId, setCustomerId] = useState<number | null>(null);
  const [customerFetchComplete, setCustomerFetchComplete] = useState(false);

  // Fetch customer ID from API using logged-in user's email
  useEffect(() => {
    if (authLoading) {
      return; // Wait for auth to finish loading
    }

    if (!isAuthenticated || !user?.email) {
      setIsLoading(false);
      setCustomerFetchComplete(true);
      return;
    }

    const fetchCustomerId = async () => {
      try {
        const customer = await GetCustomerByEmail(user.email);
        
        if (customer) {
          setCustomerId(customer.id);
          setCustomerFetchComplete(true);
          // Don't set loading to false here - let the orders fetch handle it
        } else {
          console.warn('Customer profile not found for email:', user.email);
          setIsLoading(false); // No customer found, stop loading
          setCustomerFetchComplete(true);
        }
      } catch (error) {
        console.error('Error fetching customer:', error);
        setIsLoading(false); // Error fetching customer, stop loading
        setCustomerFetchComplete(true);
      }
    };

    fetchCustomerId();
  }, [isAuthenticated, user, authLoading]);

  // Fetch orders when customer ID is available
  useEffect(() => {
    if (!customerFetchComplete) {
      return; // Wait for customer fetch to complete
    }

    if (!customerId) {
      // Customer fetch completed but no customer ID found
      setIsLoading(false);
      return;
    }

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
          description: error || 'An error occurred while fetching your orders.',
        });
        setOrders([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchOrders();
  }, [customerId, selectedStatus, customerFetchComplete]);

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, authLoading, navigate]);

  const statusFilters = [
    { value: 'ALL', label: 'All Orders' },
    { value: 'PENDING', label: 'Pending' },
    { value: 'DELIVERED', label: 'Delivered' },
    { value: 'CANCELLED', label: 'Cancelled' },
  ];

  if (authLoading || isLoading) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="page-transition pb-20 flex items-center justify-center min-h-screen"
      >
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading orders...</p>
        </div>
      </motion.div>
    );
  }

  const filteredOrders = orders.filter(order => {
    if (selectedStatus === 'ALL') return true;
    return order.status?.toUpperCase() === selectedStatus.toUpperCase();
  });

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Status Filters - Fixed Header */}
      {orders.length > 0 && (
        <div className="sticky top-0 z-10 bg-background px-4 py-4 border-b border-border overflow-x-auto">
          <div className="flex gap-2 min-w-max justify-center">
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
      )}

      {/* Orders List */}
      {filteredOrders.length === 0 ? (
        <div className="flex flex-col items-center justify-center text-center py-16 px-4">
          <div className="bg-muted rounded-full p-8 mb-4">
            <Package className="w-12 h-12 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold mb-2">
            {selectedStatus === 'ALL' ? 'No orders yet' : `No ${selectedStatus.toLowerCase()} orders`}
          </h3>
          <p className="text-sm text-muted-foreground mb-6 max-w-sm">
            {selectedStatus === 'ALL'
              ? 'When you place your first order, it will appear here. Start shopping to see your order history.'
              : `You don't have any ${selectedStatus.toLowerCase()} orders at the moment.`}
          </p>
          {selectedStatus === 'ALL' && (
            <Link
              to="/"
              className="px-6 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium btn-press"
            >
              Start Shopping
            </Link>
          )}
        </div>
      ) : (
        <div className="p-4 space-y-4">
          {filteredOrders.map((order, index) => (
            <motion.div
              key={order.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              className="bg-card border border-border rounded-md overflow-hidden"
            >
              {/* Order Header */}
              <div className="p-4 border-b border-border">
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
                  // Always parse and display cart_snapshot
                  let displayItems: OrderItemResponse[] = [];
                  
                  // Check if cart_snapshot exists and is not null/empty
                  if (order.cart_snapshot != null && order.cart_snapshot !== '' && order.cart_snapshot !== 'null') {
                    try {
                      let cartSnapshot = order.cart_snapshot;
                      
                      // Handle string format (JSON string)
                      if (typeof cartSnapshot === 'string') {
                        const trimmed = cartSnapshot.trim();
                        if (trimmed && trimmed !== 'null' && trimmed !== '') {
                          cartSnapshot = JSON.parse(trimmed);
                        } else {
                          cartSnapshot = null;
                        }
                      }
                      
                      // Handle array format
                      if (Array.isArray(cartSnapshot) && cartSnapshot.length > 0) {
                        displayItems = cartSnapshot
                          .filter((item: any) => item != null && item !== undefined) // Filter out null/undefined items
                          .map((item: any) => ({
                            product_id: item.product_id || item.productId || item.ProductID || 0,
                            quantity: item.quantity || item.Quantity || 0,
                            price: item.price || item.Price || undefined,
                          }))
                          .filter((item: OrderItemResponse) => item.product_id > 0 && item.quantity > 0); // Filter valid items
                      } else if (cartSnapshot && typeof cartSnapshot === 'object' && !Array.isArray(cartSnapshot)) {
                        // Handle single object (wrap in array)
                        const singleItem = {
                          product_id: cartSnapshot.product_id || cartSnapshot.productId || cartSnapshot.ProductID || 0,
                          quantity: cartSnapshot.quantity || cartSnapshot.Quantity || 0,
                          price: cartSnapshot.price || cartSnapshot.Price || undefined,
                        };
                        if (singleItem.product_id > 0 && singleItem.quantity > 0) {
                          displayItems = [singleItem];
                        }
                      }
                    } catch (error) {
                      console.error('Error parsing cart_snapshot for order', order.id, ':', error, 'cart_snapshot:', order.cart_snapshot);
                    }
                  }

                  return displayItems.length > 0 ? (
                    <div className="space-y-3 pb-3 border-b border-border">
                      <p className="text-xs font-semibold text-muted-foreground uppercase mb-1">
                        Order Items ({displayItems.length})
                      </p>
                      <div className="space-y-2">
                        {displayItems.map((item, itemIndex) => (
                          <div key={itemIndex} className="flex items-start justify-between gap-3 p-2 bg-muted/30 rounded-md">
                            <div className="flex-1 min-w-0">
                              <p className="font-medium text-sm">Product #{item.product_id}</p>
                              <p className="text-xs text-muted-foreground mt-0.5">
                                Quantity: {item.quantity}
                              </p>
                            </div>
                            {item.price && (
                              <div className="text-right">
                                <p className="font-semibold text-sm">{formatCurrency(item.price * item.quantity)}</p>
                                <p className="text-xs text-muted-foreground">
                                  {formatCurrency(item.price)} × {item.quantity}
                                </p>
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  ) : (
                    <div className="text-xs text-muted-foreground pb-3 border-b border-border">
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
                  <div className="flex items-center justify-between text-xs pt-2 border-t border-border">
                    <span className="text-muted-foreground">Delivery Status</span>
                    <span className="font-medium capitalize">{order.delivery_status.toLowerCase()}</span>
                  </div>
                )}
              </div>
            </motion.div>
          ))}
        </div>
      )}
    </motion.div>
  );
};

export default Orders;
