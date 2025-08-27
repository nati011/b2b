"use client";

import { DistributorRequest } from "@/lib/types";
import "leaflet/dist/leaflet.css";
import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { PiSpinner } from "react-icons/pi";
import { useRouter, useSearchParams } from "next/navigation";
import dynamic from "next/dynamic";
import { Checkbox } from "@/components/ui/checkbox";
import ImageUpload from "@/components/ImageUpload";
// import { RegisterDistributor } from "@/app/actions/auth";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { CheckCircle, Eye, EyeOff } from "lucide-react";
import usePlanstore from "@/lib/store/usePricingPlan";
import { Skeleton } from "@/components/ui/skeleton";
import useDistributorStore from "@/lib/store/useDistributorStore";
import usePartnerStore from "@/lib/store/usePaymentStore";

const Map = dynamic(() => import("@/components/map"), { ssr: false });

const baseSteps = [
  {
    title: "Profile Information",
    description: "Register a new distributor",
  },
  {
    title: "Business Information",
    description: "Add Business Information",
  },
  {
    title: "Location",
    description: "Set your business location",
  },
  {
    title: "Business Documents",
    description: "Set your business documents",
  },
  {
    title: "Security",
    description: "Set a password for your account",
  },
];

export default function DistributorsForm() {
  const router = useRouter();
  const { plans, loading: plansLoading, error: plansError, fetchPricingPlan } = usePlanstore();
  const { distributor_id, register, buySubscription, loading, error, success } = useDistributorStore();
  const { partners, fetchPaymentPartners, loading: partnersLoading } = usePartnerStore();
  const [showSuccess, setShowSuccess] = useState(false);
  const [selectedPartnerId, setSelectedPartnerId] = useState<number | null>(null);
  const [formData, setFormData] = useState<DistributorRequest>({
    name: "",
    tin: "",
    latitude: "",
    longitude: "",
    general_zone: "",
    region: "",
    woreda: "",
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    username: "",
    dob: "",
    external_id: "",
    licence_url: "",
    password:"",
    confirm_password:""
  });
  const [markerPosition, setMarkerPosition] = useState<[number, number]>([
    8.9934609, 38.7714897,
  ]);
  const [useCurrentLocation, setUseCurrentLocation] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [imageError, setImageError] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleReset = () => {
    setFormData({
      name: "",
      tin: "",
      latitude: "",
      longitude: "",
      general_zone: "",
      region: "",
      woreda: "",
      first_name: "",
      last_name: "",
      email: "",
      phone: "",
      username: "",
      dob: "",
      external_id: "",
      licence_url: "",
      password:"",
    confirm_password:""
    });
    setMarkerPosition([8.9934609, 38.7714897]);
    setUseCurrentLocation(false);
    setCurrentStep(0);
    toast.info("Form reset");
  };

  useEffect(() => {
    if (useCurrentLocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const newPosition: [number, number] = [
            position.coords.latitude,
            position.coords.longitude,
          ];
          setMarkerPosition(newPosition);
          setFormData((prev) => ({
            ...prev,
            latitude: newPosition[0].toString(),
            longitude: newPosition[1].toString(),
          }));
        },
        (error) => {
          console.error("Error getting location:", error);
        }
      );
    }
  }, [useCurrentLocation]);

  const searchParams = useSearchParams();
  const plan_id = searchParams.get("plan_id");
  const needsPlanSelection = !plan_id;

  useEffect(() => {
    void fetchPricingPlan();
  }, [fetchPricingPlan]);

  useEffect(() => {
    if (!needsPlanSelection) {
      void fetchPaymentPartners();
    }
  }, [needsPlanSelection, fetchPaymentPartners]);

  const formatPrice = (price: number) => `$${price.toLocaleString()}`;
  const termToPeriod = (termInMonth: number) => {
    if (termInMonth === 1) return "per month";
    if (termInMonth === 12) return "per year";
    return `for ${termInMonth} months`;
  };

  const handleSelectPlan = (id: string | number) => {
    const newSearch = new URLSearchParams(Array.from(searchParams.entries()));
    newSearch.set("plan_id", String(id));
    router.replace(`?${newSearch.toString()}`);
    setCurrentStep(0);
  };

  const handleNext = () => {
    if (currentStep < baseSteps.length - 1) {
      setCurrentStep((prev) => prev + 1);
    }
  };

  const handleBack = () => {
    if (currentStep > 0) {
      setCurrentStep((prev) => prev - 1);
    }
  };

  const handleSubmit = async () => {
    if(formData.password != formData.confirm_password){
      toast.error("Password don't match")
      return
    }
    try {
      await register(formData);
      if (error) {
        toast.error(error);
        return;
      }
      toast.success("Distributor registered successfully!");
      setShowSuccess(true);
    } catch (e: any) {
      toast.error(e?.message || "Failed to register distributor.");
    }
  };

  const handleBuySubscription = async () => {
    if (!plan_id) {
      toast.error("Please select a plan first.");
      return;
    }
    if (!selectedPartnerId) {
      toast.error("Please select a payment partner.");
      return;
    }
    try {
      await buySubscription({
        distributor_id: distributor_id,
        subscription_plan_id: Number(plan_id),
        payment_partner_id: selectedPartnerId,
      } as any);
      if (error) {
        toast.error(error);
        return;
      }
      if (success && success.startsWith("http")) {
        window.location.href = success;
        return;
      }
      toast.success(success || "Subscription initiated successfully.");
    } catch (e: any) {
      toast.error(e?.message || "Failed to buy subscription.");
    }
  };

  return (
    <div className="grid grid-cols-1 gap-4 mx-auto max-w-6xl my-10">
      <div className="mb-6">
        <p className="text-xl font-semibold text-black">Welcome Aboard!</p>
        <p className="text-md font-medium text-gray-700">
          Fill in the forms accordingly to register.
        </p>
      </div>

      {needsPlanSelection && (
        <Card className="rounded-sm border-2 border-gray-200 shadow-none">
          <CardContent className="pt-6">
            <div className="mb-4">
              <h3 className="text-xl font-bold">Select Plan</h3>
              <p className="text-sm text-muted-foreground">Choose a pricing plan</p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {plansLoading && (
                Array.from({ length: 3 }).map((_, idx) => (
                  <div key={idx} className="relative bg-background p-6 rounded-lg border border-border">
                    <div className="text-center mb-6">
                      <Skeleton className="h-6 w-32 mx-auto mb-2" />
                      <Skeleton className="h-8 w-24 mx-auto mb-1" />
                      <Skeleton className="h-4 w-28 mx-auto" />
                      <div className="mt-4 space-y-2">
                        <Skeleton className="h-3 w-56 mx-auto" />
                        <Skeleton className="h-3 w-44 mx-auto" />
                      </div>
                    </div>
                    <Skeleton className="h-9 w-full" />
                  </div>
                ))
              )}

              {!plansLoading && plans.length === 0 && (
                <div className="md:col-span-3 text-center text-muted-foreground">No pricing plans available.</div>
              )}

              {!plansLoading && plans.map((plan) => (
                <div key={(plan as any).id ?? plan.name} className={`relative bg-background p-6 rounded-lg border border-border`}>
                  <div className="text-center mb-6">
                    <h3 className="text-lg font-medium mb-2">{plan.name}</h3>
                    <div className="text-2xl font-bold text-primary mb-1">{formatPrice(plan.price)}</div>
                    <div className="text-xs text-muted-foreground">{termToPeriod((plan as any).term_in_month)}</div>
                  </div>
                  <Button className="w-full" onClick={() => handleSelectPlan((plan as any).id ?? plan.name)}>
                    Select Plan
                  </Button>
                </div>
              ))}

              {plansError && (
                <div className="md:col-span-3 text-center text-destructive">{plansError}</div>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Stepper appears after plan is selected */}
      {!needsPlanSelection && (
        <>
          <div className="flex items-center justify-between mb-6">
            {baseSteps.map((step, idx) => (
              <div key={step.title} className="flex-1 flex flex-col items-center">
                <div
                  className={`rounded-full w-8 h-8 flex items-center justify-center text-white font-bold ${
                    idx === currentStep
                      ? "bg-primary"
                      : idx < currentStep
                      ? "bg-primary"
                      : "bg-gray-300"
                  }`}
                >
                  {idx + 1}
                </div>
                <span
                  className={`text-xs mt-2 text-center ${
                    idx === currentStep
                      ? "text-cyan-700 font-semibold"
                      : "text-gray-500"
                  }`}
                >
                  {step.title}
                </span>
                {idx < baseSteps.length - 1 && (
                  <div className="w-full h-1 bg-gray-200 my-2">
                    <div
                      className={`h-1 ${
                        idx < currentStep ? "bg-primary" : "bg-gray-200"
                      }`}
                      style={{ width: "100%" }}
                    />
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Step Content */}
          <Card className="rounded-sm border-2 border-gray-200 shadow-none">
            <CardContent className="pt-6">
              <div className="mb-4">
                <h3 className="text-xl font-bold">{baseSteps[currentStep].title}</h3>
                <p className="text-sm text-muted-foreground">
                  {baseSteps[currentStep].description}
                </p>
              </div>

              {currentStep === 0 && (
                <div className="space-y-4">
                  <div className="grid grid-cols-1 gap-6">
                    <div className="space-y-4">
                      <div className="grid gap-2">
                        <Label htmlFor="FirstName">First Name</Label>
                        <Input
                          id="FirstName"
                          name="FirstName"
                          value={formData.first_name}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              first_name: e.target.value,
                            }));
                          }}
                          placeholder="Abebe"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="LastName">Last Name</Label>
                        <Input
                          id="LastName"
                          name="LastName"
                          value={formData.last_name}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              last_name: e.target.value,
                            }));
                          }}
                          placeholder="Kebede"
                          required
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Email">Email</Label>
                        <Input
                          id="Email"
                          name="Email"
                          value={formData.email}
                          type="email"
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              email: e.target.value,
                            }));
                          }}
                          placeholder="abebe.kebede@example.com"
                          required
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="phone">Phone</Label>
                        <Input
                          id="phone"
                          name="phone"
                          value={formData.phone}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              phone: e.target.value,
                            }));
                          }}
                          placeholder="+25191234566"
                          required
                        />
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {currentStep === 1 && (
                <div className="space-y-4">
                  <div className="grid grid-cols-1 gap-6">
                    <div className="space-y-4">
                      <div className="grid gap-2">
                        <Label htmlFor="Name">Name</Label>
                        <Input
                          id="Name"
                          name="Name"
                          value={formData.name}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              name: e.target.value,
                            }));
                          }}
                          placeholder="Business Name"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Tin">Tin</Label>
                        <Input
                          id="Tin"
                          name="Tin"
                          value={formData.tin}
                          type="number"
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              tin: e.target.value,
                            }));
                          }}
                          placeholder="1234567890"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="General Zone">General Zone</Label>
                        <Input
                          id="General Zone"
                          name="General Zone"
                          value={formData.general_zone}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              general_zone: e.target.value,
                            }));
                          }}
                          placeholder="Bole"
                          required
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Region">Region</Label>
                        <Input
                          id="Region"
                          name="Region"
                          value={formData.region}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              region: e.target.value,
                            }));
                          }}
                          placeholder="region-001"
                          required
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Woreda">Woreda</Label>
                        <Input
                          id="Woreda"
                          name="Woreda"
                          value={formData.woreda}
                          onChange={(e) => {
                            setFormData((prev) => ({
                              ...prev,
                              woreda: e.target.value,
                            }));
                          }}
                          placeholder="woreda-001"
                          required
                        />
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {currentStep === 2 && (
                <div className="space-y-4">
                  <div className="flex items-center space-x-2 my-4">
                    <Checkbox
                      id="useLocation"
                      checked={useCurrentLocation}
                      onCheckedChange={(checked: any) =>
                        setUseCurrentLocation(checked as boolean)
                      }
                    />
                    <label htmlFor="useLocation">Use my current location</label>
                  </div>
                  <div className="h-[400px] rounded-lg overflow-hidden">
                    <Map
                      markerPosition={markerPosition}
                      useCurrentLocation={useCurrentLocation}
                      setMarkerPosition={setMarkerPosition}
                      setFormData={setFormData}
                    />
                  </div>
                </div>
              )}

              {currentStep === 3 && (
                <div className="space-y-4">
                  <ImageUpload
                    onChange={(value: string[]) => {
                      setFormData((prev) => ({
                        ...prev,
                        licence_url: value[0] || "",
                      }));
                      setImageError("");
                    }}
                    value={formData.licence_url ? [formData.licence_url] : []}
                  />
                  {imageError && (
                    <p className="text-red-600 text-sm mt-2">{imageError}</p>
                  )}
                </div>
              )}

              {currentStep === 4 && (
                <div className="space-y-4">
                  <div className="grid gap-2 relative">
                    <Label htmlFor="password">
                      Password<span className="text-red-500 ">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        id="password"
                        name="password"
                        type={showPassword ? "text" : "password"}
                        placeholder="********"
                        required
                        value={formData.password}
                        onChange={(e)=>
                          setFormData((prev) => ({
                            ...prev,
                            password: e.target.value,
                          }))
                        }
                      />
                      <button
                        type="button"
                        className="absolute right-3 top-1/2 transform -translate-y-1/2"
                        onClick={() => setShowPassword(!showPassword)}
                      >
                        {showPassword ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </button>
                    </div>
                    
                  </div>

                  <div className="grid gap-2 relative">
                    <Label htmlFor="confirmPassword">
                      Confirm Password<span className="text-red-500 ">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        id="confirmPassword"
                        name="confirmPassword"
                        type={showConfirmPassword ? "text" : "password"}
                        placeholder="********"
                        required
                        value={formData.confirm_password}
                        onChange={ (e)=>setFormData((prev) => ({
                          ...prev,
                          confirm_password: e.target.value,
                        }))}
                      />
                      <button
                        type="button"
                        className="absolute right-3 top-1/2 transform -translate-y-1/2"
                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      >
                        {showConfirmPassword ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </button>
                    </div>
                    
                  </div>

                </div>
              )}

              {/* Stepper Navigation */}
              <div className="flex items-center justify-between mt-8">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleBack}
                  disabled={currentStep === 0}
                >
                  Back
                </Button>
                <div className="flex items-center space-x-4">
                  <Button type="button" variant="outline" onClick={handleReset}>
                    Reset
                  </Button>
                  {currentStep < (baseSteps.length - 1) ? (
                    <Button type="button" onClick={handleNext}>
                      Next
                    </Button>
                  ) : (
                    <Button type="submit" onClick={handleSubmit} disabled={loading}>
                      {loading ? (
                        <>
                          <PiSpinner className="animate-spin text-white mr-2" />
                          Loading
                        </>
                      ) : (
                        <>Register</>
                      )}
                    </Button>
                  )}
                </div>
              </div>

            </CardContent>
          </Card>
        </>
      )}

      <Dialog open={showSuccess} onOpenChange={setShowSuccess}>
        <DialogContent className="p-6 text-center">
          <CheckCircle className="w-16 h-16 text-green-600 mx-auto mb-4" />
          <h1 className="text-3xl font-bold text-green-800 mb-2">
            Thank You for Registering!
          </h1>
          <p className="text-green-700 mb-4">
            Our team will get back to you shortly after reviewing your profile information.
          </p>

          {/* Subscription payment section */}
          <div className="text-left space-y-3">
            <h3 className="text-lg font-semibold">Complete Subscription Payment</h3>
            <p className="text-sm text-muted-foreground">Select a payment partner and proceed to payment.</p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3 mt-2">
              {partnersLoading && <p className="text-sm text-gray-500">Loading partners...</p>}
              {!partnersLoading && partners
                .filter((p) => p.payment_method !== "MANUAL_PAYMENT")
                .map((p) => (
                  <button
                    key={p.id}
                    onClick={() => setSelectedPartnerId(p.id)}
                    className={`border rounded-md p-3 text-left ${selectedPartnerId === p.id ? "border-primary" : "border-gray-200"}`}
                  >
                    <div className="font-medium">{p.name}</div>
                    <div className="text-xs text-muted-foreground">{p.payment_method}</div>
                  </button>
                ))}
            </div>
            <div className="flex justify-end mt-4">
              <Button onClick={handleBuySubscription} disabled={!selectedPartnerId || loading}>
                {loading ? (
                  <>
                    <PiSpinner className="animate-spin text-white mr-2" />
                    Processing
                  </>
                ) : (
                  <>Pay Subscription</>
                )}
              </Button>
            </div>
          </div>

          {error && <p className="text-red-600 text-sm mt-2">{error}</p>}
        </DialogContent>
      </Dialog>
    </div>
  );
}
