import { NavLink as RouterNavLink, useLocation } from 'react-router-dom';
import { Home, ShoppingBag, Package, User } from 'lucide-react';
import { useCart } from '@/context/CartContext';
import { cn } from '@/lib/utils';

const tabs = [
  { path: '/', icon: Home, label: 'Home' },
  { path: '/cart', icon: ShoppingBag, label: 'Cart' },
  { path: '/orders', icon: Package, label: 'Orders' },
  { path: '/profile', icon: User, label: 'Profile' },
];

export const BottomNav = () => {
  const location = useLocation();
  const { itemCount } = useCart();

  // Hide BottomNav on product detail page
  const isProductDetailPage = location.pathname.startsWith('/product/');
  if (isProductDetailPage) {
    return null;
  }

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 bg-background border-t border-border safe-bottom">
      <div className="flex items-center justify-around h-14">
        {tabs.map((tab) => {
          const isCart = tab.path === '/cart';
          const isActive = location.pathname === tab.path;
          const Icon = tab.icon;

          return (
            <RouterNavLink
              key={tab.path}
              to={tab.path}
              className={cn(
                "flex flex-col items-center justify-center w-full h-full btn-press relative",
                "transition-colors",
                isActive ? "text-primary" : "text-muted-foreground"
              )}
            >
              <div className="relative">
                <Icon className="w-5 h-5 stroke-[1.5]" />
                {isCart && itemCount > 0 && (
                  <span className="absolute -top-1.5 -right-2 min-w-[18px] h-[18px] flex items-center justify-center bg-primary text-primary-foreground text-[10px] font-medium rounded-full px-1">
                    {itemCount}
                  </span>
                )}
              </div>
              <span className="text-[10px] mt-0.5 font-medium">{tab.label}</span>
            </RouterNavLink>
          );
        })}
      </div>
    </nav>
  );
};
