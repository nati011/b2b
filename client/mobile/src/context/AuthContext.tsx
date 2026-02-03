import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { GetCustomerByEmail } from '@/lib/api/customer';

interface UserInfo {
  id: string;
  name: string;
  email: string;
  roles?: string[];
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: UserInfo | null;
  token: string | null;
  login: (token: string, userInfo: UserInfo) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<UserInfo | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Helper function to fetch and save customer data
  const fetchAndSaveCustomerData = async (email: string) => {
    try {
      console.log('🔄 Fetching customer data for email:', email);
      const customer = await GetCustomerByEmail(email);
      
      if (customer) {
        // Save customer credentials to localStorage
        localStorage.setItem('customer_id', customer.id.toString());
        localStorage.setItem('customer_data', JSON.stringify(customer));
        console.log('✅ Customer data saved to localStorage:', {
          customerId: customer.id,
          email: customer.email,
        });
      } else {
        console.warn('⚠️ Customer not found for email:', email);
        // Don't throw error - customer might not exist yet, they can create profile later
      }
    } catch (error) {
      console.error('❌ Error fetching customer data:', error);
      // Don't throw error - allow login to proceed even if customer fetch fails
    }
  };

  // Load persisted credentials on mount
  useEffect(() => {
    const loadPersistedAuth = async () => {
      try {
        const storedToken = localStorage.getItem('auth_token');
        const storedUserId = localStorage.getItem('user_id');
        const storedEmail = localStorage.getItem('user_email');
        const storedName = localStorage.getItem('user_name');
        const storedCustomerId = localStorage.getItem('customer_id');
        const storedCustomerData = localStorage.getItem('customer_data');

        if (storedToken && storedUserId && storedEmail) {
          // Restore authentication state from localStorage
          setToken(storedToken);
          setUser({
            id: storedUserId,
            name: storedName || '',
            email: storedEmail,
            roles: [],
          });
          setIsAuthenticated(true);

          // If customer data is missing, fetch it
          if (!storedCustomerId || !storedCustomerData) {
            console.log('⚠️ Customer data missing, fetching from API...');
            await fetchAndSaveCustomerData(storedEmail);
          } else {
            // Verify customer data matches current user email
            try {
              const customerDataParsed = JSON.parse(storedCustomerData);
              if (customerDataParsed.email?.toLowerCase() !== storedEmail.toLowerCase()) {
                console.log('⚠️ Stored customer data email mismatch, fetching fresh data...');
                await fetchAndSaveCustomerData(storedEmail);
              }
            } catch (parseError) {
              console.error('Error parsing stored customer data, fetching fresh...');
              await fetchAndSaveCustomerData(storedEmail);
            }
          }
        }
      } catch (error) {
        console.error('Error loading persisted auth:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadPersistedAuth();
  }, []);

  const login = async (newToken: string, userInfo: UserInfo) => {
    // Store token and user info
    localStorage.setItem('auth_token', newToken);
    localStorage.setItem('user_id', userInfo.id);
    localStorage.setItem('user_email', userInfo.email);
    localStorage.setItem('user_name', userInfo.name);

    setToken(newToken);
    setUser(userInfo);
    setIsAuthenticated(true);

    // Fetch and save customer data after login
    await fetchAndSaveCustomerData(userInfo.email);
  };

  const logout = () => {
    // Clear all auth data
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('user_email');
    localStorage.removeItem('user_password'); // Clear password
    localStorage.removeItem('user_name');
    localStorage.removeItem('customer_id');
    localStorage.removeItem('customer_data'); // Clear customer data

    setToken(null);
    setUser(null);
    setIsAuthenticated(false);
    // Navigation will be handled by components using this context
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        user,
        token,
        login,
        logout,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

