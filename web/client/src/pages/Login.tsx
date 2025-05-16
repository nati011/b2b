
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import useAuthStore from "@/lib/store/useAuthStore";
import { AuthModel } from "@/lib/types";

const Login = () => {
  const {
    success,
    loading,
    error,
    login
  } = useAuthStore()
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState<AuthModel>({
    email: "",
    password: ""
  })

  useEffect(() => {
    toast(success)
  }, [success])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login(formData)
    toast.success("Successfully logged in!");
  };

  return (
    <div className="min-h-screen bg-white flex flex-col">
      <header className="fixed top-0 left-0 right-0 bg-white z-50 border-b border-gray-100">
        <nav className="container mx-auto px-4 py-4">
          <Link to="/" className="flex items-center gap-2 text-primary hover:opacity-80">
            <ArrowLeft className="w-5 h-5" />
            <span>Back to Store</span>
          </Link>
        </nav>
      </header>

      <main className="flex-1 flex items-center justify-center p-4">
        <div className="w-full max-w-md space-y-8 animate-fade-in">
          <div className="text-center">
            <h2 className="text-3xl font-semibold text-primary">Welcome back</h2>
            <p className="mt-2 text-sm text-gray-600">
              Please enter your details to sign in
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => { setFormData(prev => ({ ...prev, email: e.target.value })) }}
                placeholder="Enter your email"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={formData.password}
                onChange={(e) => { setFormData(prev => ({ ...prev, password: e.target.value })) }}
                placeholder="Enter your password"
                required
              />
            </div>

            <Button type="submit" className="w-full bg-primary text-white hover:bg-primary/90">
              Sign In
            </Button>

            <p className="text-center text-sm text-gray-600">
              Don't have an account?{" "}
              <Link to="/signup" className="text-primary hover:underline">
                Sign up
              </Link>
            </p>
          </form>
        </div>
      </main>
    </div>
  );
};

export default Login;