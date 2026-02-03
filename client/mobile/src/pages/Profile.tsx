import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { User, Mail, Phone, Building2, MapPin, Settings, LogOut, ChevronRight, Loader2 } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { ListCustomers, type CustomerResponse } from '@/lib/api/customer';

const menuItems = [
  { icon: Settings, label: 'Settings', description: 'Account preferences' },
];

interface CustomerData {
  id?: number;
  full_name: string;
  email: string;
  phone_number?: string;
  city?: string;
  region?: string;
  woreda?: string;
  status?: string;
  is_active?: boolean;
}

const Profile = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, logout, isLoading: authLoading } = useAuth();
  const [customerData, setCustomerData] = useState<CustomerData>({
    full_name: '',
    email: '',
    phone_number: '',
    city: '',
    region: '',
    woreda: '',
    status: '',
  });
  const [isLoading, setIsLoading] = useState(true);

  // Load customer data from API using email from auth context
  useEffect(() => {
    const loadCustomerData = async () => {
      // Wait for auth context to load
      if (authLoading) {
        return;
      }

      // Check if user is authenticated
      if (!isAuthenticated || !user) {
        // Not logged in, redirect to login
        navigate('/login');
        return;
      }

      try {
        // Get user info from auth context
        const storedName = user.name || '';
        const storedEmail = user.email || '';

        // Initialize with data from login context
        const initialData: CustomerData = {
          full_name: storedName || 'User',
          email: storedEmail || '',
          phone_number: '',
          city: '',
          region: '',
          woreda: '',
          status: '',
        };

        // Try to fetch customer details from API
        if (storedEmail) {
          try {
            // List customers and find by email
            const response = await ListCustomers({ limit: 100 });
            const customer = response.items?.find(c => c.email?.toLowerCase() === storedEmail.toLowerCase());
            
            if (customer) {
              // Merge customer data with initial data
              const customerDataToStore = {
                id: customer.id,
                full_name: customer.full_name || storedName,
                email: customer.email || storedEmail,
                phone_number: customer.phone_number || '',
                city: customer.city || '',
                region: customer.region || '',
                woreda: customer.woreda || '',
                status: customer.status || '',
                is_active: customer.is_active,
              };
              
              setCustomerData(customerDataToStore);
              
              // Store in localStorage for use in checkout
              localStorage.setItem('customer_id', customer.id.toString());
              localStorage.setItem('customer_data', JSON.stringify(customerDataToStore));
            } else {
              // Customer not found in API, use login context data
              setCustomerData(initialData);
            }
          } catch (error) {
            console.error('Error fetching customer data:', error);
            // Use login context data as fallback
            setCustomerData(initialData);
          }
        } else {
          setCustomerData(initialData);
        }
      } catch (error) {
        console.error('Error loading customer data:', error);
        // Fallback to basic data from auth context
        setCustomerData({
          full_name: user?.name || 'User',
          email: user?.email || '',
          phone_number: '',
          city: '',
          region: '',
          woreda: '',
          status: '',
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadCustomerData();
  }, [navigate, isAuthenticated, user, authLoading]);

  const handleSignOut = () => {
    logout(); // Use auth context logout which handles cleanup
    navigate('/login'); // Navigate to login after logout
  };

  if (isLoading) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="page-transition pb-20 flex items-center justify-center min-h-screen"
      >
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading profile...</p>
        </div>
      </motion.div>
    );
  }

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
            <div className="w-20 h-20 bg-primary rounded-full flex items-center justify-center">
              <User className="w-10 h-10 text-primary-foreground" />
            </div>
            <div className="flex-1">
              <h1 className="font-display text-xl font-semibold mb-1">{customerData.full_name || 'User'}</h1>
              <p className="text-sm text-muted-foreground mb-1">{customerData.email || 'No email'}</p>
              {customerData.status && (
                <p className="text-xs text-muted-foreground capitalize">{customerData.status}</p>
              )}
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
          <div className="flex items-center p-4 border-b border-border">
            <div className="flex items-center gap-3">
              <Mail className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-xs text-muted-foreground mb-0.5">Email</p>
                <p className="text-sm font-medium">{customerData.email || 'Not provided'}</p>
              </div>
            </div>
          </div>
          
          {/* Customer Details Section */}
          <div className="p-4 space-y-3">
            {/* Full Name */}
            {customerData.full_name && (
              <div className="flex items-start gap-3">
                <User className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Full Name</p>
                  <p className="text-sm font-medium">{customerData.full_name}</p>
                </div>
              </div>
            )}

            {/* Email */}
            {customerData.email && (
              <div className="flex items-start gap-3">
                <Mail className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Email</p>
                  <p className="text-sm font-medium">{customerData.email}</p>
                </div>
              </div>
            )}

            {/* Phone Number */}
            {customerData.phone_number && (
              <div className="flex items-start gap-3">
                <Phone className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Phone Number</p>
                  <p className="text-sm font-medium">{customerData.phone_number}</p>
                </div>
              </div>
            )}

            {/* Woreda */}
            {customerData.woreda && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Woreda</p>
                  <p className="text-sm font-medium">{customerData.woreda}</p>
                </div>
              </div>
            )}

            {/* City */}
            {customerData.city && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">City</p>
                  <p className="text-sm font-medium">{customerData.city}</p>
                </div>
              </div>
            )}

            {/* Region */}
            {customerData.region && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Region</p>
                  <p className="text-sm font-medium">{customerData.region}</p>
                </div>
              </div>
            )}

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
            const isDisabled = item.disabled || false;
            const isUpcoming = item.upcoming || false;
            return (
              <button
                key={item.label}
                disabled={isDisabled}
                className={`w-full flex items-center justify-between p-4 transition-colors ${
                  isDisabled 
                    ? 'opacity-50 cursor-not-allowed' 
                    : 'btn-press hover:bg-secondary'
                } ${
                  index !== menuItems.length - 1 ? 'border-b border-border' : ''
                }`}
              >
                <div className="flex items-center gap-3">
                  <Icon className={`w-5 h-5 ${isDisabled ? 'text-muted-foreground/50' : 'text-muted-foreground'}`} />
                  <div className="text-left">
                    <div className="flex items-center gap-2">
                      <p className={`text-sm ${isDisabled ? 'text-muted-foreground/70' : 'font-medium'}`}>
                        {item.label}
                      </p>
                      {isUpcoming && (
                        <span className="px-2 py-0.5 text-xs font-medium bg-primary/10 text-primary rounded-full">
                          Coming Soon
                        </span>
                      )}
                    </div>
                    <p className={`text-xs ${isDisabled ? 'text-muted-foreground/50' : 'text-muted-foreground'}`}>
                      {item.description}
                    </p>
                  </div>
                </div>
                {!isDisabled && <ChevronRight className="w-4 h-4 text-muted-foreground" />}
              </button>
            );
          })}
        </div>
      </div>

      {/* Sign Out */}
      <div className="px-4 mt-6">
        <button 
          onClick={handleSignOut}
          className="w-full flex items-center justify-center gap-2 p-4 text-sm text-destructive hover:bg-destructive/10 border border-destructive/20 rounded-md btn-press transition-colors"
        >
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
