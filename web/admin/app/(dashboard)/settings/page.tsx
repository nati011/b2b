"use client";
import { useEffect, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { PiSpinner } from "react-icons/pi";
import { toast } from "sonner";
import { Card,CardContent } from "@/components/ui/card";

export default function Settings() {
//   const { user, loading, success, error, fetchUser, updateProfile } =
//     useUserStore();

  const [formData, setFormData] = useState<any>({
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    username: "",
    dob: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev: any) => ({
      ...prev,
      [name]: value,
    }));
  };

//   useEffect(() => {
//     fetchUser();
//   }, []);
//   useEffect(() => {
//     if (user) {
//       setFormData(user);
//     }
//   }, [user]);

//   useEffect(() => {
//     if (success) {
//       toast.success(success);
//     }
//     if (error) {
//       toast.error(error);
//     }
//   }, [error, success]);

  return (
    <div className="space-y-6">
          <div className="space-y-6">
        <Card className="rounded-sm border-2 border-gray-200 shadow-none">
          <CardContent className="">
            <div className="mb-4">
              <h3 className="text-xl font-bold">Profile Information</h3>
              <p className="text-sm text-muted-foreground">
                Enter your details to update your account
              </p>
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-1 gap-6">
                <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="first_name">
                    First Name<span className="text-red-500">*</span>
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
                    Last Name<span className="text-red-500">*</span>
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
                <Label htmlFor="phone">
                  Phone Number<span className="text-red-500">*</span>
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
              <Button
                // onClick={() => {
                //   updateProfile(formData);
                // }}
                // className="w-full"
                // disabled={loading}
              >
               <span>Update Profile</span>
              </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
</div>

<div className="space-y-6">
<Card className="rounded-sm border-2 border-gray-200 shadow-none">
  <CardContent className="">
    <div className="mb-4">
      <h3 className="text-xl font-bold">Subscriptions</h3>
      <p className="text-sm text-muted-foreground">
        Enter your details to update your account
      </p>
    </div>

    <div className="space-y-4">
      <div className="flex justify-between">
        <div className="space-y-2">
        <h3 className="text-lg font-bold">Annual</h3>
      <p className="text-sm text-muted-foreground">
        Feb 12 2025
      </p>
        </div>
        <Button
                // onClick={() => {
                //   updateProfile(formData);
                // }}
                // className="w-full"
                // disabled={loading}
              >
               <span>Renew Subscription</span>
              </Button>
      </div>
    </div>
  </CardContent>
</Card>
</div>
</div>
  );
}
