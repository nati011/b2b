"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { useState } from "react";
import { PiSpinner } from "react-icons/pi";
import { useParams, useSearchParams } from "next/navigation";
import { ResetPassword } from "@/app/actions/auth";
import { toast } from "sonner";

interface Errors {
  password?: string;
  confirmPassword?: string;
  general?: string;
}
export default function LoginPage() {
  const routeParam = useParams<{ id: string }>();
  const [errors, setErrors] = useState<Errors>({});
  const [loading, setLoading] = useState(false);
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const handleResetPassword = async () => {
    setLoading(true);
    try {
      const response = await ResetPassword(routeParam.id, password);
      toast.success(response);
    } catch (error: any) {
      console.log(error);
      toast.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col gap-4 p-6 md:p-10 my-32">
      <div className="flex flex-1 items-center justify-center">
        <div className="w-full max-w-lg">
          <div className="flex flex-col gap-6">
            <div className="flex flex-col items-center gap-2 text-left">
              <h1 className="text-2xl font-bold">Reset Password</h1>
              {/* <p className="text-balance text-sm text-muted-foreground">
                Lorem Ipsum and so on
              </p> */}
            </div>
            <div className="grid gap-6">
              {errors.general && (
                <div className="bg-red-100/[0.2] rounded-md p-2 border border-red-900">
                  <p className="text-red-900 text-center font-medium">
                    {errors.general}
                  </p>
                </div>
              )}
              <div className="grid gap-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="********"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  onBlur={(e) => {
                    if (!e.target.value) {
                      setErrors((prev) => ({
                        ...prev,
                        password: "Password is required",
                      }));
                    } else {
                      setErrors((prev) => ({ ...prev, password: "" }));
                    }
                  }}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="confirm-password">Confirm Password</Label>
                <Input
                  className={`${errors.confirmPassword ? "ring-red-300" : ""}`}
                  id="confirm-password"
                  type="password"
                  placeholder="********"
                  required
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  value={confirmPassword}
                  onBlur={(e) => {
                    if (!e.target.value) {
                      setErrors((prev) => ({
                        ...prev,
                        confirmPassword: "Confirm Password is required",
                      }));
                    } else if (e.target.value != password) {
                      setErrors((prev) => ({
                        ...prev,
                        confirmPassword: "Passwords don't match",
                      }));
                    } else {
                      setErrors((prev) => ({ ...prev, confirmPassword: "" }));
                    }
                  }}
                />
              </div>
              {errors.confirmPassword && (
                <span className="text-red-500 font-semibold text-sm">
                  {errors.confirmPassword}
                </span>
              )}
              <Button
                onClick={() => {
                  handleResetPassword();
                }}
                className="w-full"
                disabled={loading}
              >
                {loading ? (
                  <div className="flex gap-2">
                    <PiSpinner className="animate-spin" />
                    <span>Loading</span>
                  </div>
                ) : (
                  <span>Reset Password</span>
                )}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
