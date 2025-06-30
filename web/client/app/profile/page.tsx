"use client"
import { FormEvent, useEffect, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { User, UserIdentity } from '@/lib/types';
import { Button } from "@/components/ui/button";
import { PiSpinner } from "react-icons/pi";
import { useUserStore } from "@/lib/store/useAuthStore";
import { toast } from "sonner";

export default function Settings() {
    const {
        user,
        loading,
        success,
        error,
        fetchUser,
        updateProfile
    } = useUserStore()

    const [formData, setFormData] = useState<Partial<UserIdentity>>({
        first_name: "",
        last_name: "",
        email: "",
        phone: "",
        username: "",
        dob: "",
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    useEffect(()=>{
        fetchUser()
    },[])
    useEffect(()=>{
        if(user){
        setFormData(user)
    }
    },[user])

    useEffect(()=>{
        if(success){
            toast.success(success)
        }
        if(error){
            toast.error(error)
        }
    },[error,success ])

    return (
        <div className="flex flex-col gap-4 p-6 md:p-10 my-32">
            <div className="flex flex-1 items-center justify-center">
                <div className="w-full max-w-lg">
                    <form className="flex flex-col gap-6" >
                        <div className="flex flex-col items-center gap-2 text-left">
                            <h1 className="text-2xl font-bold">Profile Settings</h1>
                            <p className="text-balance text-sm text-muted-foreground">
                                Enter your details to update your account
                            </p>
                        </div>

                        <div className="grid gap-6">
                            {error && (
                                <div className="bg-red-100/[0.2] rounded-md p-2 border border-red-900">
                                    <p className="text-red-900 text-center font-medium">{error}</p>
                                </div>
                            )}

                            {/* Personal Information */}
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="first_name">First Name<span className='text-red-500'>*</span></Label>
                                    <Input
                                        id="first_name"
                                        name="first_name"
                                        placeholder="John"
                                        required
                                        value={formData.first_name}
                                        onChange={handleChange}
                                    />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="last_name">Last Name<span className='text-red-500'>*</span></Label>
                                    <Input
                                        id="last_name"
                                        name="last_name"
                                        placeholder="Doe"
                                        required
                                        value={formData.last_name}
                                        onChange={handleChange}
                                    />
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
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="phone">Phone Number<span className='text-red-500'>*</span></Label>
                                <Input
                                    id="phone"
                                    name="phone"
                                    type="tel"
                                    placeholder="+251966961629"
                                    required
                                    value={formData.phone}
                                    onChange={handleChange}
                                />
                            </div>


                            <Button onClick={()=>{updateProfile(formData)}} className="w-full" disabled={loading}>
                                {loading ? (
                                    <div className="flex gap-2">
                                        <PiSpinner className="animate-spin" />
                                        <span>Updating account...</span>
                                    </div>
                                ) : (
                                    <span>Update Profile</span>
                                )}
                            </Button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}