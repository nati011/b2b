import { motion } from 'framer-motion';
import { User, FileText, MapPin, HelpCircle, LogOut, ChevronRight } from 'lucide-react';

const menuItems = [
  { icon: FileText, label: 'Invoices', description: 'Request quotes, view invoices' },
  { icon: MapPin, label: 'Shipping Locations', description: 'Warehouse addresses' },
  { icon: HelpCircle, label: 'Support', description: 'Contact your account manager' },
];

const Profile = () => {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Header */}
      <div className="p-6 pt-8 bg-card border-b border-border">
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 bg-primary rounded-md flex items-center justify-center">
            <User className="w-6 h-6 text-primary-foreground" />
          </div>
          <div>
            <h1 className="font-display text-lg font-semibold">Business Account</h1>
          </div>
        </div>
      </div>

      {/* Menu Items */}
      <div className="px-4">
        <div className="border border-border rounded-sm overflow-hidden">
          {menuItems.map((item, index) => {
            const Icon = item.icon;
            return (
              <button
                key={item.label}
                className={`w-full flex items-center justify-between p-4 btn-press hover:bg-secondary transition-colors ${
                  index !== menuItems.length - 1 ? 'border-b border-border' : ''
                }`}
              >
                <div className="flex items-center gap-3">
                  <Icon className="w-5 h-5 text-muted-foreground" />
                  <div className="text-left">
                    <p className="text-sm font-medium">{item.label}</p>
                    <p className="text-xs text-muted-foreground">{item.description}</p>
                  </div>
                </div>
                <ChevronRight className="w-4 h-4 text-muted-foreground" />
              </button>
            );
          })}
        </div>
      </div>

      {/* Sign Out */}
      <div className="px-4 mt-6">
        <button className="w-full flex items-center justify-center gap-2 p-4 text-sm text-muted-foreground hover:text-foreground border border-border rounded-sm btn-press transition-colors">
          <LogOut className="w-4 h-4" />
          Sign Out
        </button>
      </div>

      {/* App Info */}
      <div className="text-center mt-8 text-xs text-muted-foreground">
        <p>Version 1.0.0</p>
      </div>
    </motion.div>
  );
};

export default Profile;
