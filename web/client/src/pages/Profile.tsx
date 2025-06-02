"use client"
import { FormEvent, useState, useCallback, useEffect } from "react";
import useAuthStore from "@/lib/store/useAuthStore";
import { Navbar } from "./Navbar";
import { Card, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { RxAvatar } from "react-icons/rx";
import { Separator } from "@/components/ui/separator";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { User } from '@/lib/types'

// @ts-ignore
export default function Settings() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const {
        success,
        retailer,
        fetchLoggedInUser
    } = useAuthStore()

    const [formData, setFormData] = useState<User>({
        id: 0,
        first_name: "",
        last_name: "",
        email: "",
        phone: "",
        username: "",
        dob: "",
        external_id: "",
    })

    const handleSubmit = async (
        event: FormEvent<HTMLFormElement>
    ): Promise<void> => {
        event.preventDefault();
        try {
            setLoading(true);
            setError("");

            setLoading(false);
        } catch (error: any) {
            setError(error.response.data.message);
            setTimeout(() => {
                setError("");
            }, 5000);
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchLoggedInUser()
        console.log(retailer)
        setFormData({
            first_name: retailer.user.first_name,
            last_name: retailer.user.last_name,
            email: retailer.user.email,
            phone: retailer.user.phone
        })
    }, [retailer]);



    return (
        <div className="mx-20 justify-center items-center">
            <div className="">
                <div className="shadow-none grid grid-cols-1 p-4">
                    <p className="font-semibold text-xl">
                        Profile Settings
                    </p>
                    <div className="my-10" />
                    <div className="w-full">
                        <form onSubmit={handleSubmit} className="space-y-6 text-left">
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
                            {/* <Button type="submit" className="w-full bg-primary text-white hover:bg-primary/90" onClick={(e)=>handleSubmit(e)}>
                                Save Changes
                            </Button> */}
                        </form>
                    </div>

                </div>
            </div>
        </div>

    );
}