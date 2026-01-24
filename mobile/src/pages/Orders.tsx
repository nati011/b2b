import { motion } from 'framer-motion';
import { Package } from 'lucide-react';
import { Link } from 'react-router-dom';

const Orders = () => {
  // TODO: Fetch orders from API
  const orders: any[] = [];

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Orders List */}
      {orders.length === 0 ? (
        <div className="flex flex-col items-center justify-center text-center py-16 px-4">
          <div className="bg-muted rounded-full p-8 mb-4">
            <Package className="w-12 h-12 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold mb-2">No orders yet</h3>
          <p className="text-sm text-muted-foreground mb-6 max-w-sm">
            When you place your first order, it will appear here. Start shopping to see your order history.
          </p>
          <Link
            to="/"
            className="px-6 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium btn-press"
          >
            Start Shopping
          </Link>
        </div>
      ) : (
        <div className="p-4 space-y-4">
          {orders.map((order, index) => (
            <motion.div
              key={order.id || index}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="bg-card border border-border rounded-md p-4"
            >
              {/* Order content will go here */}
            </motion.div>
          ))}
        </div>
      )}
    </motion.div>
  );
};

export default Orders;

