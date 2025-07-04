"use client"
import { useRouter, useSearchParams } from 'next/navigation'
import { useEffect, useState } from 'react'
import { PiSpinner } from 'react-icons/pi'
import { FcGoogle } from "react-icons/fc"
import { toast } from 'sonner'

import { InitResetPassword } from '@/app/actions/auth'
import { useAuth } from '@/hooks/useAuth'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

interface Errors {
  email?: string
  password?: string
  general?: string
}

export default function LoginPage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const { login, isAuthenticated, isLoading: authLoading } = useAuth()
  
  const [errors, setErrors] = useState<Errors>({})
  const [loading, setLoading] = useState(false)
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  
  const callbackUrl = searchParams.get("callbackUrl") || "/"
  const error = searchParams.get("error")

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push(callbackUrl)
    }
  }, [isAuthenticated, router, callbackUrl])

  // Handle URL error parameter
  useEffect(() => {
    if (error) {
      setErrors((prev) => ({ ...prev, general: "Invalid credentials" }))
      setTimeout(() => {
        setErrors({})
      }, 5000)
    }
  }, [error])

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setErrors({})

    // Validation
    const newErrors: Errors = {}
    if (!email) {
      newErrors.email = "Email is required"
    }
    if (!password) {
      newErrors.password = "Password is required"
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors)
      return
    }

    try {
      setLoading(true)
      const result = await login(email, password, callbackUrl)
      
      if (!result.success) {
        setErrors((prev) => ({ ...prev, general: result.error || "Invalid credentials" }))
        setTimeout(() => {
          setErrors({})
        }, 5000)
      }
    } catch (error: any) {
      setErrors((prev) => ({ ...prev, general: error.message || "An unexpected error occurred" }))
      setTimeout(() => {
        setErrors({})
      }, 5000)
    } finally {
      setLoading(false)
    }
  }

  const handleInitResetPassword = async () => {
    if (!email) {
      toast.error("Please enter your email address first")
      return
    }

    try {
      const success = await InitResetPassword(email)
      toast.success(success)
      router.push("/check-your-email")
    } catch (error: any) {
      toast.error(error.message || "Failed to send reset email")
    }
  }

  // Show loading state while checking authentication
  if (authLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="flex items-center space-x-2">
          <PiSpinner className="h-6 w-6 animate-spin" />
          <span>Loading...</span>
        </div>
      </div>
    )
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

            {errors.general && (
              <div className="p-3 text-sm text-red-500 bg-red-50 border border-red-200 rounded-md">
                {errors.general}
              </div>
            )}

            <div className="flex flex-col gap-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="Enter your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className={errors.email ? "border-red-500" : ""}
                disabled={loading}
              />
              {errors.email && (
                <span className="text-sm text-red-500">{errors.email}</span>
              )}
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="Enter your password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={errors.password ? "border-red-500" : ""}
                disabled={loading}
              />
              {errors.password && (
                <span className="text-sm text-red-500">{errors.password}</span>
              )}
            </div>

            <div className="flex items-center justify-between">
              <Button
                type="button"
                variant="link"
                className="px-0 text-sm"
                onClick={handleInitResetPassword}
                disabled={loading}
              >
                Forgot your password?
              </Button>
            </div>

            <Button type="submit" disabled={loading} className="w-full">
              {loading ? (
                <div className="flex items-center space-x-2">
                  <PiSpinner className="h-4 w-4 animate-spin" />
                  <span>Signing in...</span>
                </div>
              ) : (
                "Sign in"
              )}
            </Button>

            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Or continue with
                </span>
              </div>
            </div>

            <Button type="button" variant="outline" className="w-full">
              <FcGoogle className="mr-2 h-4 w-4" />
              Sign in with Google
            </Button>
          </form>
        </div>
      </div>
    </div>
  )
}
