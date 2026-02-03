import { useState, useEffect, useRef } from 'react';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { Mail, Lock, Loader2, Eye, EyeOff } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Login as LoginAPI, GoogleSSO } from '@/lib/api/auth';
import { useAuth } from '@/context/AuthContext';
import { toast } from 'sonner';

// Google OAuth Client ID - should be set via environment variable
// For now, using a placeholder that should be configured
const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID || '';

// Declare Google types
declare global {
  interface Window {
    google?: {
      accounts: {
        oauth2: {
          initTokenClient: (config: {
            client_id: string;
            scope: string;
            callback: (response: { access_token: string }) => void;
          }) => {
            requestAccessToken: () => void;
          };
        };
        id: {
          initialize: (config: {
            client_id: string;
            callback: (response: {
              credential: string;
            }) => void;
          }) => void;
          renderButton: (element: HTMLElement, config: {
            theme?: string;
            size?: string;
            text?: string;
            width?: number;
          }) => void;
        };
      };
    };
  }
}

const LoginPage = () => {
  const navigate = useNavigate();
  const { login, isAuthenticated } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isGoogleLoading, setIsGoogleLoading] = useState(false);
  const [errors, setErrors] = useState<{ email?: string; password?: string; general?: string }>({});
  const googleButtonRef = useRef<HTMLDivElement>(null);

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/');
    }
  }, [isAuthenticated, navigate]);

  const validateForm = (): boolean => {
    const newErrors: { email?: string; password?: string } = {};
    
    if (!email.trim()) {
      newErrors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      newErrors.email = 'Please enter a valid email address';
    }
    
    if (!password) {
      newErrors.password = 'Password is required';
    } else if (password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrors({});

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);

    try {
      const response = await LoginAPI({ email, password });
      
      // Store email and password for Basic Auth
      localStorage.setItem('user_email', email);
      localStorage.setItem('user_password', password);
      
      // Decode JWT to get user info
      const token = response.access_token;
      if (token) {
        const parts = token.split('.');
        if (parts.length === 3) {
          const payload = parts[1];
          const paddedPayload = payload + '='.repeat((4 - (payload.length % 4)) % 4);
          const decoded = JSON.parse(atob(paddedPayload.replace(/-/g, '+').replace(/_/g, '/')));
          
          // Update auth context with user info (this will also fetch and save customer data)
          await login(token, {
            id: decoded.sub || '',
            name: decoded.name || '',
            email: decoded.email || email,
            roles: decoded.realm_access?.roles || [],
          });
        }
      }
      
      toast.success('Login successful!', {
        description: 'Welcome back!',
      });

      // Navigate to home page after successful login
      navigate('/');
    } catch (error: any) {
      console.error('Login error:', error);
      setErrors({
        general: error.message || 'Invalid email or password. Please try again.',
      });
      toast.error('Login failed', {
        description: error.message || 'Invalid email or password',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleSignIn = async (accessToken: string) => {
    setIsGoogleLoading(true);
    setErrors({});

    try {
      // Fetch user info from Google
      const userInfoResponse = await fetch('https://www.googleapis.com/oauth2/v2/userinfo', {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });

      if (!userInfoResponse.ok) {
        throw new Error('Failed to fetch user information from Google');
      }

      const userInfo = await userInfoResponse.json();
      
      // Extract name parts
      const nameParts = (userInfo.name || '').split(' ');
      const firstName = userInfo.given_name || nameParts[0] || '';
      const lastName = userInfo.family_name || nameParts.slice(1).join(' ') || '';

      // Call backend SSO endpoint
      const response = await GoogleSSO({
        token: accessToken,
        first_name: firstName,
        last_name: lastName,
        email: userInfo.email,
      });
      
      // For Google SSO, store email for Basic Auth
      // Note: Google SSO doesn't provide password, but backend should handle SSO differently
      localStorage.setItem('user_email', userInfo.email);
      
      // Decode JWT to get user info
      const token = response.access_token;
      if (token) {
        const parts = token.split('.');
        if (parts.length === 3) {
          const payload = parts[1];
          const paddedPayload = payload + '='.repeat((4 - (payload.length % 4)) % 4);
          const decoded = JSON.parse(atob(paddedPayload.replace(/-/g, '+').replace(/_/g, '/')));
          
          // Update auth context with user info (this will also fetch and save customer data)
          await login(token, {
            id: decoded.sub || '',
            name: decoded.name || userInfo.name || '',
            email: decoded.email || userInfo.email || '',
            roles: decoded.realm_access?.roles || [],
          });
        }
      }

      toast.success('Login successful!', {
        description: 'Welcome back!',
      });

      // Navigate to home page after successful login
      navigate('/');
    } catch (error: any) {
      console.error('Google SSO error:', error);
      setErrors({
        general: error.message || 'Failed to sign in with Google. Please try again.',
      });
      toast.error('Google sign-in failed', {
        description: error.message || 'Failed to sign in with Google',
      });
    } finally {
      setIsGoogleLoading(false);
    }
  };

  // Initialize Google Sign-In
  useEffect(() => {
    if (!GOOGLE_CLIENT_ID) {
      console.warn('Google Client ID not configured. Google Sign-In will not be available.');
      return;
    }

    // Load Google Identity Services script
    const script = document.createElement('script');
    script.src = 'https://accounts.google.com/gsi/client';
    script.async = true;
    script.defer = true;
    
    script.onload = () => {
      if (window.google?.accounts.oauth2 && googleButtonRef.current) {
        // Use OAuth2 token flow for getting access token
        const tokenClient = window.google.accounts.oauth2.initTokenClient({
          client_id: GOOGLE_CLIENT_ID,
          scope: 'https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile openid',
          callback: (tokenResponse: { access_token: string }) => {
            if (tokenResponse.access_token) {
              handleGoogleSignIn(tokenResponse.access_token);
            }
          },
        });

        // Create a custom button that triggers the OAuth flow
        if (googleButtonRef.current) {
          const button = document.createElement('button');
          button.type = 'button';
          button.className = 'w-full h-12 flex items-center justify-center gap-3 bg-card hover:bg-secondary border border-border rounded-sm text-sm font-medium text-foreground transition-colors btn-press';
          button.innerHTML = `
            <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            <span>Sign in with Google</span>
          `;
          button.onclick = () => {
            setIsGoogleLoading(true);
            tokenClient.requestAccessToken();
          };
          googleButtonRef.current.innerHTML = '';
          googleButtonRef.current.appendChild(button);
        }
      }
    };

    // Check if script already exists
    const existingScript = document.querySelector('script[src="https://accounts.google.com/gsi/client"]');
    if (!existingScript) {
      document.head.appendChild(script);
    } else {
      // Script already loaded, trigger onload manually
      script.onload?.();
    }

    return () => {
      // Cleanup is handled by React
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="min-h-screen bg-background flex items-center justify-center p-4"
    >
      <div className="w-full max-w-md">
        {/* Logo/Header */}
        <motion.div
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.1 }}
          className="text-center mb-8"
        >
          <h1 className="font-display text-3xl font-bold mb-2">Welcome Back</h1>
          <p className="text-sm text-muted-foreground">
            Sign in to continue to your account
          </p>
        </motion.div>

        {/* Login Form */}
        <motion.form
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.2 }}
          onSubmit={handleSubmit}
          className="space-y-4"
        >
          {/* General Error */}
          {errors.general && (
            <div className="bg-destructive/10 border border-destructive/20 text-destructive text-sm p-3 rounded-sm">
              {errors.general}
            </div>
          )}

          {/* Email Field */}
          <div className="space-y-2">
            <label htmlFor="email" className="text-sm font-medium text-foreground">
              Email
            </label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => {
                  setEmail(e.target.value);
                  if (errors.email) setErrors({ ...errors, email: undefined });
                }}
                placeholder="Enter your email"
                className={`w-full h-12 pl-10 pr-4 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/20 ${
                  errors.email ? 'border border-destructive' : 'border border-border'
                }`}
                disabled={isLoading}
                autoComplete="email"
              />
            </div>
            {errors.email && (
              <p className="text-xs text-destructive">{errors.email}</p>
            )}
          </div>

          {/* Password Field */}
          <div className="space-y-2">
            <label htmlFor="password" className="text-sm font-medium text-foreground">
              Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                  if (errors.password) setErrors({ ...errors, password: undefined });
                }}
                placeholder="Enter your password"
                className={`w-full h-12 pl-10 pr-12 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/20 ${
                  errors.password ? 'border border-destructive' : 'border border-border'
                }`}
                disabled={isLoading}
                autoComplete="current-password"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                aria-label={showPassword ? 'Hide password' : 'Show password'}
              >
                {showPassword ? (
                  <EyeOff className="w-5 h-5" />
                ) : (
                  <Eye className="w-5 h-5" />
                )}
              </button>
            </div>
            {errors.password && (
              <p className="text-xs text-destructive">{errors.password}</p>
            )}
          </div>

          {/* Forgot Password Link */}
          <div className="flex justify-end">
            <button
              type="button"
              onClick={() => {
                // TODO: Implement forgot password
                toast.info('Forgot password feature coming soon');
              }}
              className="text-sm text-primary hover:underline"
              disabled={isLoading}
            >
              Forgot password?
            </button>
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            disabled={isLoading || isGoogleLoading}
            className="w-full h-12 text-base font-semibold tracking-wide btn-press rounded-sm"
            size="lg"
          >
            {isLoading ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Signing in...
              </>
            ) : (
              'Sign In'
            )}
          </Button>

          {/* Divider */}
          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border"></div>
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">Or continue with</span>
            </div>
          </div>

          {/* Google Sign-In Button */}
          {GOOGLE_CLIENT_ID ? (
            <div className="w-full">
              {isGoogleLoading ? (
                <Button
                  type="button"
                  disabled
                  variant="outline"
                  className="w-full h-12 text-base font-semibold tracking-wide rounded-sm border-border"
                  size="lg"
                >
                  <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                  Signing in with Google...
                </Button>
              ) : (
                <div ref={googleButtonRef} className="w-full"></div>
              )}
            </div>
          ) : (
            <Button
              type="button"
              disabled
              variant="outline"
              className="w-full h-12 text-base font-semibold tracking-wide rounded-sm border-border opacity-50"
              size="lg"
            >
              Google Sign-In (Not Configured)
            </Button>
          )}
        </motion.form>

        {/* Footer */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="mt-6 text-center text-sm text-muted-foreground"
        >
          <p>
            Don't have an account?{' '}
            <button
              type="button"
              onClick={() => {
                // TODO: Navigate to sign up page when available
                toast.info('Sign up feature coming soon');
              }}
              className="text-primary hover:underline font-medium"
            >
              Sign up
            </button>
          </p>
        </motion.div>
      </div>
    </motion.div>
  );
};

export default LoginPage;

