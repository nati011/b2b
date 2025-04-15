
import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Users,
  Package,
  Menu,
  X,
  LogIn,
  ShoppingBag
} from "lucide-react";

const MainLayout = ({ children }: { children: React.ReactNode }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const location = useLocation();
  const navigate = useNavigate();

  const menuItems = [
    { icon: LayoutDashboard, label: "Dashboard", path: "/" },
    { icon: Users, label: "Customers", path: "/customers" },
    { icon: Users, label: "Users", path: "/users" },
    { icon: Package, label: "Products", path: "/products" },
    { icon: ShoppingBag, label: "Orders", path: "/orders" },
  ];

  const handleProfileClick = () => {
    navigate("/login");
  };

  return (
    <div className="min-h-screen bg-background flex">
      {/* Sidebar */}
      <aside
        className={cn(
          "fixed left-0 top-0 z-40 h-screen transition-transform duration-300 ease-in-out",
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <div className="h-full px-3 py-4 overflow-y-auto w-64 bg-gradient-to-b from-[#1A1F2C] to-[#403E43] text-white">
          <div className="flex items-center justify-between mb-8">
            <h1 className="text-2xl font-semibold">Efoyta</h1>
            <button
              onClick={() => setIsSidebarOpen(false)}
              className="p-1 rounded-lg hover:bg-[#555555]/30 transition-colors"
            >
              <X className="h-6 w-6" />
            </button>
          </div>
          <nav className="space-y-2">
            {menuItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                className={cn(
                  "flex items-center p-3 rounded-lg transition-colors",
                  location.pathname === item.path
                    ? "bg-[#555555]/50 text-white"
                    : "hover:bg-[#555555]/30"
                )}
              >
                <item.icon className="h-5 w-5 mr-3" />
                <span>{item.label}</span>
              </Link>
            ))}
          </nav>
        </div>
      </aside>

      {/* Main content */}
      <div
        className={cn(
          "flex-1 transition-all duration-300 ease-in-out",
          isSidebarOpen ? "ml-64" : "ml-0"
        )}
      >
        <header className="bg-white border-b h-16 fixed top-0 right-0 left-0 z-30 flex items-center justify-between px-4">
          {/* Navigation toggle button (when sidebar is closed) */}
          {!isSidebarOpen && (
            <button
              onClick={() => setIsSidebarOpen(true)}
              className="p-2 rounded-lg hover:bg-gray-100 text-gray-700 shadow-sm border border-gray-200 transition-all duration-200 hover:shadow-md"
            >
              <Menu className="h-6 w-6" />
            </button>
          )}
          
          {/* Empty div for spacing when sidebar is open */}
          {isSidebarOpen && <div></div>}
          
          {/* Login button */}
          <button 
            className="p-2 rounded-full hover:bg-gray-100 text-blue-600 flex items-center gap-2 border border-gray-200 shadow-sm"
            onClick={handleProfileClick}
          >
            <LogIn className="h-5 w-5" />
            <span className="text-sm font-medium">Login</span>
          </button>
        </header>
        <main className="pt-20 px-4 pb-4">{children}</main>
      </div>
    </div>
  );
};

export default MainLayout;
