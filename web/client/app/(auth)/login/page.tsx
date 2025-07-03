"use client"
import { useRouter, useSearchParams } from 'next/navigation'
import { signIn } from "next-auth/react";
import { useEffect, useState } from 'react';
import { PiSpinner } from 'react-icons/pi'
import { FcGoogle } from "react-icons/fc";
import { toast } from 'sonner'

import { InitResetPassword } from '@/app/actions/auth'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"


interface Errors {
    email?: string;
    password?: string;
    general?: string;
}
export default function LoginPage() {
    const searchParams = useSearchParams()
    const [errors, setErrors] = useState<Errors>({});
    const [loading, setLoading] = useState(false)
    const callbackUrl = searchParams.get("callbackUrl") || "/"
    const error = searchParams.get("error")
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const router = useRouter()



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
                const result = await signIn("credentials", {
                    email,
                    password,
                    redirect: true,
                    callbackUrl,
                })
                if (result?.error) {
                    setErrors((prev) => ({ ...prev, general: "Invalid Credentials" }));
                    setTimeout(() => {
                        setErrors({});
                    }, 5000);
                }
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
    useEffect(() => {
        if (error) {
            setErrors((prev) => ({ ...prev, general: "Invalid Credentials" }));
        }
    }, [error])
    const handleInitResetPassword = async () => {
        try {
            const success = await InitResetPassword(email)
            toast.success(success)
            router.push("/check-your-email")
        } catch (error: any) {
            toast.error(error.message)
            // TEMP
            router.push("/check-your-email")
        }
    }


    return (
        <div className="flex flex-col gap-4 p-6 md:p-10 my-32">
            <div className="flex flex-1 items-center justify-center">
                <div className="w-full max-w-lg">

                    <form className="flex flex-col gap-6" onSubmit={handleSubmit}>
                        <div className="flex flex-col items-center gap-2 text-left">
                            <h1 className="text-2xl font-bold">Login to your account</h1>
                            <p className="text-balance text-sm text-muted-foreground">
                                Enter your credentials below to login to your account
                            </p>
                        </div>
                        <div className="grid gap-6">
                            {
                                errors.general && (
                                    <div className="bg-red-100/[0.2] rounded-md p-2 border border-red-900">
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
                                    <div
                                        onClick={() => { handleInitResetPassword() }}
                                        className="ml-auto text-sm underline-offset-4 hover:underline cursor-pointer"
                                    >
                                        <p>Forgot your password?</p>
                                    </div>
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
                            Don&apos;t have an account?{" "}
                            <a href="/signup" className="underline underline-offset-4">
                                Sign up
                            </a>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}
