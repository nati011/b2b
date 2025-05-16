import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import useAuthStore from "@/lib/store/useAuthStore"
import { useEffect, useState } from "react"
import { AuthModel } from "@/lib/types"
import { GoEye, GoEyeClosed } from "react-icons/go";
import { toast } from "sonner"

interface Errors {
    employeeNumber?: string;
    password?: string;
    general?: string;
}


export function LoginForm({
    className,
    ...props
}: React.ComponentPropsWithoutRef<"form">) {
    const {
        success,
        loading,
        error,
        login
    } = useAuthStore()
    const [showPassword, setShowPassword] = useState(false);
    const [errors, setErrors] = useState<Errors>({});
    const [formData, setFormData] = useState<AuthModel>({
        email: "",
        password: ""
    })

    const handleSubmit = () => {
        login(formData)
    }
    useEffect(() => {
        toast(success)
    }, [success])
    return (
        <form className={cn("flex flex-col gap-6", className)} {...props}>
            <div className="flex flex-col items-left gap-2 text-left">
                <h1 className="text-3xl font-bold">Welcome back</h1>
                <p className="text-balance text-md text-muted-foreground">
                    Enter your email below to login to your account
                </p>
            </div>
            <div className="grid gap-6">
                <div className="grid gap-2">
                    <Label htmlFor="email">Email</Label>
                    <Input
                        id="email"
                        type="email"
                        placeholder="john.doe@example.com"
                        required value={formData.email}
                        onChange={(e) => { setFormData(prev => ({ ...prev, email: e.target.value })) }}
                        onBlur={(e) => {
                            if (!e.target.value) {
                                setErrors((prev) => ({ ...prev, password: 'Email is required' }));
                            } else {
                                setErrors((prev) => ({ ...prev, password: '' }));
                            }
                        }}

                    />
                </div>
                <div className="grid gap-2">
                    <div className="flex items-center">
                        <Label htmlFor="password">Password</Label>
                        <a
                            href="#"
                            className="ml-auto text-sm underline-offset-4 hover:underline"
                        >
                            Forgot your password?
                        </a>
                    </div>
                    <div className="relative">
                        <Input
                            type={showPassword ? 'text' : 'password'}
                            className={`block rounded-md dark:text-white border-0 py-1.5 pl-7 pr-12 w-full text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-[#003949] sm:text-sm sm:leading-6 ${errors.password ? 'ring-red-300' : ''
                                }`}
                            placeholder="Password"
                            value={formData.password}
                            onChange={(e) => { setFormData(prev => ({ ...prev, password: e.target.value })) }}
                            onBlur={(e) => {
                                if (!e.target.value) {
                                    setErrors((prev) => ({ ...prev, password: 'Password is required' }));
                                } else {
                                    setErrors((prev) => ({ ...prev, password: '' }));
                                }
                            }}
                        />
                        <button
                            type="button"
                            className="absolute inset-y-0 right-0 pr-3 flex items-center"
                            onClick={() => setShowPassword(!showPassword)}
                        >
                            {showPassword ? <GoEye /> : <GoEyeClosed />}
                        </button>
                    </div>


                </div>
                <Button className="w-full" onClick={handleSubmit}>
                    Login
                </Button>

            </div>
            <div className="text-center text-sm">
                Don&apos;t have an account?{" "}
                <a href="#" className="underline underline-offset-4">
                    Sign up
                </a>
            </div>
        </form>
    )
}
