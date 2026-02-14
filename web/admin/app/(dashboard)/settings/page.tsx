"use client";

import { useEffect, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { PiSpinner } from "react-icons/pi";
import { toast } from "sonner";
import { Card, CardContent } from "@/components/ui/card";
import { useUserStore } from "@/app/libs/store/useAuthStore";
import { useSubscriptionStore } from "@/app/libs/store/useSubscriptionStore";
import { useSession } from "next-auth/react";

interface ProfileFormData {
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  username: string;
  dob: string;
}

export default function Settings() {
  const { data: session } = useSession()
  const roles: string[] = (session as any)?.user?.roles || []
  const isAdminUser = roles.includes('superadmin') || roles.includes('admin')
  const {
    user,
    loading: profileLoading,
    success: profileSuccess,
    error: profileError,
    fetchUser,
    updateProfile,
  } = useUserStore();

  const {
    subscription,
    subscriptionError,
    subscriptionLoading,
    fetchSubscription,
  } = useSubscriptionStore();

  const [formData, setFormData] = useState<ProfileFormData>({
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    username: "",
    dob: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  useEffect(() => {
    fetchUser();
    if (!isAdminUser) {
      fetchSubscription();
    }
  }, [fetchUser, fetchSubscription, isAdminUser]);

  useEffect(() => {
    if (user) {
      setFormData({
        first_name: user.first_name || "",
        last_name: user.last_name || "",
        email: user.email || "",
        phone: user.phone || "",
        username: user.username || "",
        dob: user.dob || "",
      });
    }
  }, [user]);

  useEffect(() => {
    if (profileSuccess) toast.success(profileSuccess);
    if (profileError) toast.error(profileError);
    if (subscriptionError) toast.error(subscriptionError);
  }, [profileSuccess, profileError, subscriptionError]);

  const handleUpdateProfile = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfile(formData);
  };

  const handleRenewSubscription = () => {
    toast.info("Renewal functionality coming soon.");
    // TODO: Integrate 
  };

  const formatDate = (dateString: string | undefined): string => {
    if (!dateString) return "N/A";
    const date = new Date(dateString);
    return isNaN(date.getTime()) ? "Invalid Date" : date.toLocaleDateString();
  };

  return (
    <div className="space-y-6">
      {/* Profile Section */}
      <Card className="rounded-sm border-2 border-gray-200 shadow-none">
        <CardContent className="pt-6">
          <div className="mb-6">
            <h3 className="text-xl font-bold">Profile Information</h3>
            <p className="text-sm text-muted-foreground">
              Update your personal details below.
            </p>
          </div>

          <form onSubmit={handleUpdateProfile} className="space-y-6">
            <div className="grid grid-cols-1 gap-6">
              <div className="space-y-4">
                {/* Name Fields */}
                <div className="grid grid-cols-2 gap-4">
                  <div className="grid gap-2">
                    <Label htmlFor="first_name">
                      First Name <span className="text-red-500">*</span>
                    </Label>
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
                    <Label htmlFor="last_name">
                      Last Name <span className="text-red-500">*</span>
                    </Label>
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

                {/* Email */}
                <div className="grid gap-2">
                  <Label htmlFor="email">
                    Email <span className="text-red-500">*</span>
                  </Label>
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

                {/* Phone */}
                <div className="grid gap-2">
                  <Label htmlFor="phone">
                    Phone Number <span className="text-red-500">*</span>
                  </Label>
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

                {/* Submit Button */}
                <Button
                  type="submit"
                  className="w-full"
                  disabled={profileLoading}
                >
                  {profileLoading ? (
                    <>
                      <PiSpinner className="mr-2 h-4 w-4 animate-spin" />
                      Updating...
                    </>
                  ) : (
                    "Update Profile"
                  )}
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Subscription Section */}
      {!isAdminUser && (
      <Card className="rounded-sm border-2 border-gray-200 shadow-none">
        <CardContent className="pt-6">
          <div className="mb-6">
            <h3 className="text-xl font-bold">Subscription</h3>
            <p className="text-sm text-muted-foreground">
              Manage your current plan and renewal options.
            </p>
          </div>

          {subscriptionLoading ? (
            <div className="flex items-center justify-center py-8">
              <PiSpinner className="h-6 w-6 animate-spin text-muted-foreground" />
              <span className="ml-2 text-sm text-muted-foreground">Loading subscription...</span>
            </div>
          ) : subscription ? (
            <div className="flex justify-between items-center">
              <div className="space-y-1">
                <h4 className="text-lg font-semibold">
                  {subscription?.subscription_plan_name || "Unknown Plan"}
                </h4>
                <p className="text-sm text-muted-foreground">
                  Purchased on {formatDate(subscription?.created_date?.toLocaleString())}
                </p>
              </div>

              {subscription?.status !== "active" ? (
                <Button
                  size="sm"
                  onClick={handleRenewSubscription}
                  disabled={subscriptionLoading || false}
                >
                  Renew Subscription
                </Button>
              ) : (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800 border border-emerald-200">
                  Active
                </span>
              )}
            </div>
          ) : (
            <div className="text-center py-6 text-sm text-muted-foreground">
              No active subscription found.
            </div>
          )}
        </CardContent>
      </Card>
      )}
    </div>
  );
}