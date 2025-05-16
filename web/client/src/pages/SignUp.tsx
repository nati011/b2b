
import { useState } from "react";
import { Link } from "react-router-dom";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import useAuthStore from "@/lib/store/useAuthStore";
import { Register } from "@/lib/types";

const SignUp = () => {
  const {
    success,
    loading,
    error,
    signup
  } = useAuthStore()
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState<Register>({
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    password: "",
    confirm_password: "",
  }

  )
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.password !== formData.confirm_password) {
      toast.error("Passwords don't match!");
      return;
    }
    signup(formData)
    toast.success("Account created successfully!");
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
            <h2 className="text-3xl font-semibold text-primary">Create an account</h2>
            <p className="mt-2 text-sm text-gray-600">
              Enter your details to get started
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">First Name</Label>
                <Input
                  id="firstName"
                  type="text"
                  value={formData.first_name}
                  onChange={(e) => { setFormData(prev => ({ ...prev, first_name: e.target.value })) }}
                  placeholder="John"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Last Name</Label>
                <Input
                  id="lastName"
                  type="text"
                  value={formData.last_name}
                  onChange={(e) => { setFormData(prev => ({ ...prev, last_name: e.target.value })) }}
                  placeholder="Doe"
                  required
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => { setFormData(prev => ({ ...prev, email: e.target.value })) }}
                placeholder="john@example.com"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Phone Number</Label>
              <Input
                id="phone_number"
                value={formData.phone}
                onChange={(e) => { setFormData(prev => ({ ...prev, phone: e.target.value })) }}
                placeholder="+251966961629"
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
                placeholder="Create a password"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Confirm Password</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={formData.confirm_password}
                onChange={(e) => { setFormData(prev => ({ ...prev, confirm_password: e.target.value })) }}
                placeholder="Confirm your password"
                required
              />
            </div>

            <Button type="submit" className="w-full bg-primary text-white hover:bg-primary/90">
              Sign Up
            </Button>

            <p className="text-center text-sm text-gray-600">
              Already have an account?{" "}
              <Link to="/login" className="text-primary hover:underline">
                Sign in
              </Link>
            </p>
          </form>
        </div>
      </main>
    </div>
  );
};

export default SignUp;
