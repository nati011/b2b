"use client";
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useSession } from 'next-auth/react';

interface UserInfo {
  id: string;
  name: string;
  email: string;
  roles?: string[];
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: UserInfo | null;
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
  const { data: session, status } = useSession();
  
  // Initialize state by checking localStorage immediately (client-side only)
  const getInitialAuthState = (): { isAuthenticated: boolean; user: UserInfo | null } => {
    if (typeof window === 'undefined') {
      return { isAuthenticated: false, user: null };
    }
    
    const token = localStorage.getItem('access_token');
    const storedEmail = localStorage.getItem('saved_email') || localStorage.getItem('user_email');
    const storedName = localStorage.getItem('user_name');
    const storedId = localStorage.getItem('user_id');
    const storedRoles = localStorage.getItem('user_roles');
    
    if (token && storedEmail) {
      let roles: string[] = [];
      try {
        if (storedRoles) {
          roles = JSON.parse(storedRoles);
        }
      } catch (e) {
        // Ignore parse errors
      }
      
      return {
        isAuthenticated: true,
        user: {
          id: storedId || '',
          name: storedName || '',
          email: storedEmail,
          roles: roles,
        },
      };
    }
    
    return { isAuthenticated: false, user: null };
  };

  const initialState = getInitialAuthState();
  const [isAuthenticated, setIsAuthenticated] = useState(initialState.isAuthenticated);
  const [user, setUser] = useState<UserInfo | null>(initialState.user);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkAuth = () => {
      // Check NextAuth session first
      if (session?.user) {
        // Type assertion for custom user properties
        const userData = session.user as any;
        setUser({
          id: userData.id || '',
          name: userData.name || session.user.name || '',
          email: session.user.email || '',
          roles: userData.roles || [],
        });
        setIsAuthenticated(true);
        setIsLoading(false);
        return;
      }

      // Check localStorage token (for basic auth) - check directly
      // This check can run even while session is loading
      if (typeof window !== 'undefined') {
        const token = localStorage.getItem('access_token');
        if (token) {
          // Try to get user info from localStorage
          const storedEmail = localStorage.getItem('saved_email') || localStorage.getItem('user_email');
          const storedName = localStorage.getItem('user_name');
          const storedId = localStorage.getItem('user_id');
          const storedRoles = localStorage.getItem('user_roles');
          
          if (storedEmail) {
            let roles: string[] = [];
            try {
              if (storedRoles) {
                roles = JSON.parse(storedRoles);
              }
            } catch (e) {
              // Ignore parse errors
            }
            
            setUser({
              id: storedId || '',
              name: storedName || '',
              email: storedEmail,
              roles: roles,
            });
            setIsAuthenticated(true);
            setIsLoading(false);
            return;
          }
        }
      }

      // Not authenticated - only set if session is not loading
      if (status !== 'loading') {
        setUser(null);
        setIsAuthenticated(false);
        setIsLoading(false);
      }
    };

    // Check auth immediately, and also when session status changes
    checkAuth();
  }, [session, status]);

  // Also listen for storage changes and window focus (when user logs in/out)
  useEffect(() => {
    if (typeof window === 'undefined') return;

    const recheckAuth = () => {
      // Re-check authentication
      if (session?.user) {
        const userData = session.user as any;
        setUser({
          id: userData.id || '',
          name: userData.name || session.user.name || '',
          email: session.user.email || '',
          roles: userData.roles || [],
        });
        setIsAuthenticated(true);
        return;
      }

      const token = localStorage.getItem('access_token');
      const storedEmail = localStorage.getItem('saved_email') || localStorage.getItem('user_email');
      
      if (token && storedEmail) {
        const storedName = localStorage.getItem('user_name');
        const storedId = localStorage.getItem('user_id');
        const storedRoles = localStorage.getItem('user_roles');
        
        let roles: string[] = [];
        try {
          if (storedRoles) {
            roles = JSON.parse(storedRoles);
          }
        } catch (e) {
          // Ignore parse errors
        }
        
        setUser({
          id: storedId || '',
          name: storedName || '',
          email: storedEmail,
          roles: roles,
        });
        setIsAuthenticated(true);
      } else {
        setUser(null);
        setIsAuthenticated(false);
      }
    };

    // Check on storage changes (other tabs)
    window.addEventListener('storage', recheckAuth);
    
    // Check when window regains focus (user might have logged in)
    window.addEventListener('focus', recheckAuth);
    
    // Listen for custom auth state change event (triggered after login)
    window.addEventListener('auth-state-changed', recheckAuth);
    
    return () => {
      window.removeEventListener('storage', recheckAuth);
      window.removeEventListener('focus', recheckAuth);
      window.removeEventListener('auth-state-changed', recheckAuth);
    };
  }, [session]);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        user,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

