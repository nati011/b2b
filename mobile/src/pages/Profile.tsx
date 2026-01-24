import { motion } from 'framer-motion';
import { User, Mail, Phone, Building2, MapPin, FileText, Settings, HelpCircle, LogOut, ChevronRight, Edit2 } from 'lucide-react';

const menuItems = [
  { icon: FileText, label: 'Invoices', description: 'Request quotes, view invoices' },
  { icon: MapPin, label: 'Shipping Locations', description: 'Warehouse addresses' },
  { icon: Settings, label: 'Settings', description: 'Account preferences' },
  { icon: HelpCircle, label: 'Support', description: 'Contact your account manager' },
];

const Profile = () => {
  // TODO: Fetch user data from API
  const userData = {
    name: 'Natnael',
    email: 'natnaeljemaneh001@gmail.com',
    phone: '+251 911 234 567',
    company: 'TechSupply Co.',
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Profile Header */}
      <div className="bg-card border-b border-border">
        <div className="p-6 pt-8 pb-6">
          <div className="flex items-center gap-4 mb-4">
            <div className="relative">
              <div className="w-20 h-20 bg-primary rounded-full flex items-center justify-center">
                <User className="w-10 h-10 text-primary-foreground" />
              </div>
              <button className="absolute bottom-0 right-0 w-6 h-6 bg-primary rounded-full flex items-center justify-center border-2 border-background">
                <Edit2 className="w-3 h-3 text-primary-foreground" />
              </button>
            </div>
            <div className="flex-1">
              <h1 className="font-display text-xl font-semibold mb-1">{userData.name}</h1>
              <p className="text-sm text-muted-foreground">natnaeljemaneh001@gmail.com</p>
            </div>
          </div>
        </div>
      </div>

      {/* Personal Information */}
      <div className="px-4 mt-6">
        <h2 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide mb-3 px-2">
          Personal Information
        </h2>
        <div className="bg-card border border-border rounded-md overflow-hidden">
          <div className="flex items-center justify-between p-4 border-b border-border">
            <div className="flex items-center gap-3">
              <Mail className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-xs text-muted-foreground mb-0.5">Email</p>
                <p className="text-sm font-medium">{userData.email}</p>
              </div>
            </div>
            <button className="p-2 hover:bg-secondary rounded-sm btn-press transition-colors">
              <Edit2 className="w-4 h-4 text-muted-foreground" />
            </button>
          </div>
          <div className="flex items-center justify-between p-4 border-b border-border">
            <div className="flex items-center gap-3">
              <Phone className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-xs text-muted-foreground mb-0.5">Phone</p>
                <p className="text-sm font-medium">{userData.phone}</p>
              </div>
            </div>
            <button className="p-2 hover:bg-secondary rounded-sm btn-press transition-colors">
              <Edit2 className="w-4 h-4 text-muted-foreground" />
            </button>
          </div>
          <div className="flex items-center justify-between p-4">
            <div className="flex items-center gap-3">
              <Building2 className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-xs text-muted-foreground mb-0.5">Company</p>
                <p className="text-sm font-medium">{userData.company}</p>
              </div>
            </div>
            <button className="p-2 hover:bg-secondary rounded-sm btn-press transition-colors">
              <Edit2 className="w-4 h-4 text-muted-foreground" />
            </button>
          </div>
        </div>
      </div>

      {/* Menu Items */}
      <div className="px-4 mt-6">
        <h2 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide mb-3 px-2">
          Account
        </h2>
        <div className="bg-card border border-border rounded-md overflow-hidden">
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
        <button className="w-full flex items-center justify-center gap-2 p-4 text-sm text-destructive hover:bg-destructive/10 border border-destructive/20 rounded-md btn-press transition-colors">
          <LogOut className="w-4 h-4" />
          Sign Out
        </button>
      </div>

      {/* App Info */}
      <div className="text-center mt-8 mb-4 text-xs text-muted-foreground">
        <p>Version 1.0.0</p>
      </div>
    </motion.div>
  );
};

export default Profile;
