"use client"
import Image from 'next/image'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Eye, EyeOff } from 'lucide-react'

import { signIn } from "next-auth/react";
import { useEffect, useState } from 'react';
import { PiSpinner } from 'react-icons/pi'
import { useRouter, useSearchParams } from 'next/navigation'
import { InitResetPassword, RegisterRetailer } from '@/app/actions/auth'
import { toast } from 'sonner'
import { FcGoogle } from 'react-icons/fc'
import { RegisterRequest } from '@/lib/types'

interface Errors {
    tin?: string;
    latitude?: string;
    longitude?: string;
    general_zone?: string;
    region?: string;
    woreda?: string;
    first_name?: string;
    last_name?: string;
    email?: string;
    phone?: string;
    password?: string;
    confirmPassword?: string;
    general?: string;
}

export default function SignUpPage() {
    const searchParams = useSearchParams()
    const [errors, setErrors] = useState<Errors>({});
    const [loading, setLoading] = useState(false)
    const callbackUrl = searchParams.get("callbackUrl") || "/"
    const error = searchParams.get("error")
    const router = useRouter()

    const [formData, setFormData] = useState({
        first_name: "",
        last_name: "",
        email: "",
        phone: "",
        password: "",
        confirmPassword: ""
    });

    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setErrors({});
        const newErrors: Errors = {};
        if (!formData.first_name) newErrors.first_name = "First name is required";
        if (!formData.last_name) newErrors.last_name = "Last name is required";
        if (!formData.email) newErrors.email = "Email is required";
        if (!formData.phone) newErrors.phone = "Phone is required";
        if (!formData.password) newErrors.password = "Password is required";
        if (formData.password !== formData.confirmPassword) newErrors.confirmPassword = "Passwords don't match";

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        try {
            setLoading(true);
            const response = await RegisterRetailer(formData as RegisterRequest)
            console.log(response)
            toast.success("Account created successfully!");
            router.push('/login')
        } catch (error: any) {
            setErrors(prev => ({ ...prev, general: error.message }));
            setTimeout(() => {
                setErrors({});
            }, 5000);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (error) {
            setErrors((prev) => ({ ...prev, general: "An error occurred" }));
        }
    }, [error]);

    return (
        <div className="flex flex-col gap-4 p-6 md:p-10 my-16">
            <div className="flex flex-1 items-center justify-center">
                <div className="w-full max-w-2xl">
                    <form className="flex flex-col gap-6" onSubmit={handleSubmit}>
                        <div className="flex flex-col items-center gap-2 text-left">
                            <h1 className="text-2xl font-bold">Create your account</h1>
                            <p className="text-balance text-sm text-muted-foreground">
                                Enter your details to create an account
                            </p>
                        </div>
                        <div className="grid gap-6">
                            {errors.general && (
                                <div className="bg-red-100/[0.2] rounded-md p-2 border border-red-900">
                                    <p className="text-red-900 text-center font-medium">{errors.general}</p>
                                </div>
                            )}

                            {/* Personal Information */}
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="first_name">First Name<span className='text-red-500 '>*</span></Label>
                                    <Input
                                        id="first_name"
                                        name="first_name"
                                        placeholder="John"
                                        required
                                        value={formData.first_name}
                                        onChange={handleChange}
                                    />
                                    {errors.first_name && <p className="text-red-500 text-xs">{errors.first_name}</p>}
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="last_name">Last Name<span className='text-red-500 '>*</span></Label>
                                    <Input
                                        id="last_name"
                                        name="last_name"
                                        placeholder="Doe"
                                        required
                                        value={formData.last_name}
                                        onChange={handleChange}
                                    />
                                    {errors.last_name && <p className="text-red-500 text-xs">{errors.last_name}</p>}
                                </div>
                            </div>

                            {/* Contact Information */}
                            <div className="grid gap-2">
                                <Label htmlFor="email">Email*</Label>
                                <Input
                                    id="email"
                                    name="email"
                                    type="email"
                                    placeholder="m@example.com"
                                    required
                                    value={formData.email}
                                    onChange={handleChange}
                                />
                                {errors.email && <p className="text-red-500 text-xs">{errors.email}</p>}
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="phone">Phone Number<span className='text-red-500 '>*</span></Label>
                                <Input
                                    id="phone"
                                    name="phone"
                                    type="tel"
                                    placeholder="+251966961629"
                                    required
                                    value={formData.phone}
                                    onChange={handleChange}
                                />
                                {errors.phone && <p className="text-red-500 text-xs">{errors.phone}</p>}
                            </div>



                            {/* Password Fields */}
                            <div className="grid gap-2 relative">
                                <Label htmlFor="password">Password<span className='text-red-500 '>*</span></Label>
                                <div className="relative">
                                    <Input
                                        id="password"
                                        name="password"
                                        type={showPassword ? "text" : "password"}
                                        placeholder="********"
                                        required
                                        value={formData.password}
                                        onChange={handleChange}
                                    />
                                    <button
                                        type="button"
                                        className="absolute right-3 top-1/2 transform -translate-y-1/2"
                                        onClick={() => setShowPassword(!showPassword)}
                                    >
                                        {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                    </button>
                                </div>
                                {errors.password && <p className="text-red-500 text-xs">{errors.password}</p>}
                            </div>

                            <div className="grid gap-2 relative">
                                <Label htmlFor="confirmPassword">Confirm Password<span className='text-red-500 '>*</span></Label>
                                <div className="relative">
                                    <Input
                                        id="confirmPassword"
                                        name="confirmPassword"
                                        type={showConfirmPassword ? "text" : "password"}
                                        placeholder="********"
                                        required
                                        value={formData.confirmPassword}
                                        onChange={handleChange}
                                    />
                                    <button
                                        type="button"
                                        className="absolute right-3 top-1/2 transform -translate-y-1/2"
                                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                    >
                                        {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                    </button>
                                </div>
                                {errors.confirmPassword && <p className="text-red-500 text-xs">{errors.confirmPassword}</p>}
                            </div>

                            <Button type="submit" className="w-full" disabled={loading}>
                                {loading ? (
                                    <div className="flex gap-2">
                                        <PiSpinner className="animate-spin" />
                                        <span>Creating account...</span>
                                    </div>
                                ) : (
                                    <span>Sign Up</span>
                                )}
                            </Button>
                        </div>
                        <div className="flex gap-2 items-center justify-center">
                            <div className="border h-[0.2px] w-full"></div>
                            <p>
                                Or
                            </p>
                            <div className="border h-[0.1px] w-full"></div>
                        </div>
                        <button onClick={() => signIn('google', { callbackUrl: "/dashboard" })} className="border-2 w-full font-semibold rounded p-2 flex justify-center gap-4 items-center">
                            <FcGoogle size={24} className="" />
                            <span>
                                Continue with Google
                            </span>
                        </button>
                        <div className="text-center text-sm">
                            Already have an account?{" "}
                            <a href="/login" className="underline underline-offset-4">
                                Login
                            </a>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}