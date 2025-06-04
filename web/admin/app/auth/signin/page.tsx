"use client"
import Image from 'next/image'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

import { signIn } from "next-auth/react";
import { useState } from 'react';
import { PiSpinner } from 'react-icons/pi'


interface Errors {
    email?: string;
    password?: string;
    general?: string;
}
export default function LoginPage() {
    const [errors, setErrors] = useState<Errors>({});
    const [loading, setLoading] = useState(false)

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setErrors({});

        if (!email) {
            setErrors((prev) => ({ ...prev, email: "Email is required" }));
        }
        if (!password) {
            setErrors((prev) => ({ ...prev, password: "Password is required" }));
        }
        if (email && password) {
            try {
                setLoading(true)
                await signIn("credentials", {
                    email: email,
                    password: password,
                    callbackUrl: "/",
                    redirect: true,
                })
            } catch (error: any) {
                setErrors((prev) => ({ ...prev, general: error }));
                setTimeout(() => {
                    setErrors({});
                }, 5000);

            } finally {
                setLoading(false)
            }
        }
    };
    return (
        <div className="grid min-h-svh lg:grid-cols-2">
            <div className="flex flex-col gap-4 p-6 md:p-10">
                <div className="flex justify-center gap-2 md:justify-start">
                    <a href="#" className="flex items-center gap-2 font-medium">
                        <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
                            <Image src='/logo.png' width={150} height={100} alt="logo" />
                        </div>
                        Efoyeta Store
                    </a>
                </div>
                <div className="flex flex-1 items-center justify-center">
                    <div className="w-full max-w-lg">
                        <form className="flex flex-col gap-6" onSubmit={handleSubmit}>
                            <div className="flex flex-col items-center gap-2 text-left">
                                <h1 className="text-2xl font-bold">Login to your account</h1>
                                <p className="text-balance text-sm text-muted-foreground">
                                    Enter your email below to login to your account
                                </p>
                            </div>
                            <div className="grid gap-6">
                                {
                                    errors.general && (
                                        <div className="bg-red-100/[0.2] rounded-md p-2 ">
                                            <p className="text-red-900 text-center font-medium">{errors.general}</p>
                                        </div>
                                    )
                                }
                                <div className="grid gap-2">
                                    <Label htmlFor="email">Email</Label>
                                    <Input id="email" type="email" placeholder="m@example.com" required value={email} onChange={(e) => setEmail((e.target.value).toLowerCase())} onBlur={(e) => {
                                        if (!e.target.value) {
                                            setErrors((prev) => ({ ...prev, email: "Email is required" }))
                                        }
                                        else {
                                            setErrors((prev) => ({ ...prev, email: "" }))
                                        }
                                    }} />
                                </div>
                                <div className="grid gap-2">
                                    <div className="flex items-center">
                                        <Label htmlFor="password" >Password</Label>
                                        <a
                                            href="#"
                                            className="ml-auto text-sm underline-offset-4 hover:underline"
                                        >
                                            Forgot your password?
                                        </a>
                                    </div>
                                    <Input id="password" type="password" placeholder='********' required onChange={(e) => setPassword(e.target.value)} onBlur={(e) => {
                                        if (!e.target.value) {
                                            setErrors((prev) => ({ ...prev, password: "Password is required" }))
                                        }
                                        else {
                                            setErrors((prev) => ({ ...prev, password: "" }))
                                        }
                                    }} />
                                </div>
                                <Button type="submit" className="w-full" disabled={loading}>
                                    {
                                        loading ? <div className="flex gap-2">
                                            <PiSpinner className="animate-spin" />
                                            <span>Loading</span>
                                        </div> : <span>Login</span>
                                    }

                                </Button>

                            </div>
                            {/* <div className="text-center text-sm">
                Don&apos;t have an account?{" "}
                <a href="#" className="underline underline-offset-4">
                    Sign up
                </a>
            </div> */}
                        </form>
                    </div>
                </div>
            </div>
            <div className="relative hidden bg-muted lg:block">
                <img
                    src="/placeholder.svg"
                    alt="Image"
                    className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
                />
            </div>
        </div>
    )
}
