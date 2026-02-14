"use client";

import { SupplierRequest } from "@/lib/types";
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
// import { RegisterSupplier } from "@/app/actions/auth";
import { ClientOnlyDialog } from "@/components/ClientOnlyDialog";
import { CheckCircle, Eye, EyeOff } from "lucide-react";
import usePlanstore from "@/lib/store/usePricingPlan";
import { Skeleton } from "@/components/ui/skeleton";
import useSupplierStore from "@/lib/store/useSupplierStore";

const Map = dynamic(() => import("@/components/map"), { ssr: false });

const baseSteps = [
  {
    title: "Profile Information",
    description: "Register a new supplier",
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

export default function SuppliersForm() {
  const router = useRouter();
  const { plans, loading: plansLoading, error: plansError, fetchPricingPlan } = usePlanstore();
  const { supplier_id, register, loading, error, success } = useSupplierStore();
  const [showSuccess, setShowSuccess] = useState(false);
  const [formData, setFormData] = useState<SupplierRequest>({
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
  const [errors, setErrors] = useState<Record<string, string | undefined>>({});

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
    setErrors({});
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
          // Clear location errors when location is set
          setErrors(prev => {
            const newErrors = { ...prev };
            delete newErrors.latitude;
            delete newErrors.longitude;
            return newErrors;
          });
        },
        (error) => {
          console.error("Error getting location:", error);
          toast.error("Failed to get your location. Please select manually on the map.");
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

  // Watch for successful registration - only show dialog if registration truly succeeded
  useEffect(() => {
    // Only show success dialog if:
    // 1. success is truthy (not null)
    // 2. supplier_id exists
    // 3. no error occurred
    // 4. dialog is not already showing
    if (success && supplier_id && !error && !showSuccess) {
      setShowSuccess(true);
    }
    // If there's an error, make sure dialog is closed
    if (error && showSuccess) {
      setShowSuccess(false);
    }
  }, [success, supplier_id, error, showSuccess]);

  const formatPrice = (price: number) => `${price.toLocaleString()} ETB`;
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

  // Validation functions matching backend requirements
  const validateEmail = (email: string): boolean => {
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailPattern.test(email);
  };

  const validatePhone = (phone: string): boolean => {
    const phonePattern = /^\+\d{12}$/;
    return phonePattern.test(phone);
  };

  const validateTin = (tin: string): boolean => {
    const tinPattern = /^\d{10}$/;
    return tinPattern.test(tin);
  };

  const validateLatitude = (lat: string): boolean => {
    const num = parseFloat(lat);
    return !isNaN(num) && num >= -90 && num <= 90;
  };

  const validateLongitude = (lng: string): boolean => {
    const num = parseFloat(lng);
    return !isNaN(num) && num >= -180 && num <= 180;
  };

  // Real-time validation functions for individual fields
  const validateField = (fieldName: string, value: string, compareValue?: string): string | null => {
    switch (fieldName) {
      case "name":
        if (!value.trim()) return "Business name is required";
        return null;
      case "first_name":
        if (!value.trim()) return "First name is required";
        return null;
      case "last_name":
        if (!value.trim()) return "Last name is required";
        return null;
      case "email":
        if (!value.trim()) return "Email is required";
        if (!validateEmail(value)) return "Please enter a valid email address";
        return null;
      case "phone":
        if (!value.trim()) return "Phone number is required";
        if (!validatePhone(value)) return "Phone must be in format +251912345678 (12 digits after +)";
        return null;
      case "tin":
        if (!value.trim()) return "TIN is required";
        if (!validateTin(value)) return "TIN must be exactly 10 digits";
        return null;
      case "general_zone":
        if (!value.trim()) return "General zone is required";
        return null;
      case "region":
        if (!value.trim()) return "Region is required";
        return null;
      case "woreda":
        if (!value.trim()) return "Woreda is required";
        return null;
      case "latitude":
        if (!value.trim()) return "Latitude is required";
        if (!validateLatitude(value)) return "Please select a valid location on the map";
        return null;
      case "longitude":
        if (!value.trim()) return "Longitude is required";
        if (!validateLongitude(value)) return "Please select a valid location on the map";
        return null;
      case "password":
        if (!value.trim()) return "Password is required";
        if (value.length < 6) return "Password must be at least 6 characters";
        return null;
      case "confirm_password":
        if (!value.trim()) return "Please confirm your password";
        if (compareValue !== undefined && value !== compareValue) return "Passwords do not match";
        return null;
      default:
        return null;
    }
  };

  const validateStep = (step: number): boolean => {
    const newErrors: Record<string, string> = {};

    if (step === 0) {
      // Profile Information
      if (!formData.first_name.trim()) {
        newErrors.first_name = "First name is required";
      }
      if (!formData.last_name.trim()) {
        newErrors.last_name = "Last name is required";
      }
      if (!formData.email.trim()) {
        newErrors.email = "Email is required";
      } else if (!validateEmail(formData.email)) {
        newErrors.email = "Please enter a valid email address";
      }
      if (!formData.phone.trim()) {
        newErrors.phone = "Phone number is required";
      } else if (!validatePhone(formData.phone)) {
        newErrors.phone = "Phone must be in format +251912345678 (12 digits after +)";
      }
    } else if (step === 1) {
      // Business Information
      if (!formData.name.trim()) {
        newErrors.name = "Business name is required";
      }
      if (!formData.tin.trim()) {
        newErrors.tin = "TIN is required";
      } else if (!validateTin(formData.tin)) {
        newErrors.tin = "TIN must be exactly 10 digits";
      }
      if (!formData.general_zone.trim()) {
        newErrors.general_zone = "General zone is required";
      }
      if (!formData.region.trim()) {
        newErrors.region = "Region is required";
      }
      if (!formData.woreda.trim()) {
        newErrors.woreda = "Woreda is required";
      }
    } else if (step === 2) {
      // Location - must have valid coordinates
      if (!formData.latitude.trim() || !formData.longitude.trim()) {
        newErrors.latitude = "Please select a location on the map or use your current location";
        newErrors.longitude = "Please select a location on the map or use your current location";
      } else {
        if (!validateLatitude(formData.latitude)) {
          newErrors.latitude = "Please select a valid location on the map";
        }
        if (!validateLongitude(formData.longitude)) {
          newErrors.longitude = "Please select a valid location on the map";
        }
      }
    } else if (step === 4) {
      // Security (Password)
      if (!formData.password) {
        newErrors.password = "Password is required";
      } else if (formData.password.length < 6) {
        newErrors.password = "Password must be at least 6 characters";
      }
      if (!formData.confirm_password) {
        newErrors.confirm_password = "Please confirm your password";
      } else if (formData.password !== formData.confirm_password) {
        newErrors.confirm_password = "Passwords do not match";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = (e?: React.MouseEvent) => {
    // Prevent default form submission if called from form
    if (e) {
      e.preventDefault();
    }
    
    // Validate current step before proceeding
    if (!validateStep(currentStep)) {
      toast.error("Please fill in all required fields correctly before proceeding");
      return;
    }
    
    if (currentStep < baseSteps.length - 1) {
      setCurrentStep((prev) => prev + 1);
      // Clear errors when moving to next step
      setErrors({});
    }
  };

  const handleBack = () => {
    if (currentStep > 0) {
      setCurrentStep((prev) => prev - 1);
    }
  };

  const handleStepClick = (stepIndex: number) => {
    // Allow navigation to any step
    if (stepIndex !== currentStep && stepIndex >= 0 && stepIndex < baseSteps.length) {
      // Navigate to the selected step
      setCurrentStep(stepIndex);
      // Clear errors for better UX when navigating
      setErrors({});
    }
  };

  const handleSubmit = async () => {
    // Validate all steps before submission
    const allStepsValid = [0, 1, 2, 4].every(step => {
      const stepErrors: Record<string, string> = {};
      
      if (step === 0) {
        if (!formData.first_name.trim()) stepErrors.first_name = "First name is required";
        if (!formData.last_name.trim()) stepErrors.last_name = "Last name is required";
        if (!formData.email.trim()) {
          stepErrors.email = "Email is required";
        } else if (!validateEmail(formData.email)) {
          stepErrors.email = "Please enter a valid email address";
        }
        if (!formData.phone.trim()) {
          stepErrors.phone = "Phone number is required";
        } else if (!validatePhone(formData.phone)) {
          stepErrors.phone = "Phone must be in format +251912345678 (12 digits after +)";
        }
      } else if (step === 1) {
        if (!formData.name.trim()) stepErrors.name = "Business name is required";
        if (!formData.tin.trim()) {
          stepErrors.tin = "TIN is required";
        } else if (!validateTin(formData.tin)) {
          stepErrors.tin = "TIN must be exactly 10 digits";
        }
        if (!formData.general_zone.trim()) stepErrors.general_zone = "General zone is required";
        if (!formData.region.trim()) stepErrors.region = "Region is required";
        if (!formData.woreda.trim()) stepErrors.woreda = "Woreda is required";
      } else if (step === 2) {
        if (!formData.latitude.trim()) {
          stepErrors.latitude = "Latitude is required";
        } else if (!validateLatitude(formData.latitude)) {
          stepErrors.latitude = "Please select a valid location on the map";
        }
        if (!formData.longitude.trim()) {
          stepErrors.longitude = "Longitude is required";
        } else if (!validateLongitude(formData.longitude)) {
          stepErrors.longitude = "Please select a valid location on the map";
        }
      } else if (step === 4) {
        if (!formData.password) {
          stepErrors.password = "Password is required";
        } else if (formData.password.length < 6) {
          stepErrors.password = "Password must be at least 6 characters";
        }
        if (!formData.confirm_password) {
          stepErrors.confirm_password = "Please confirm your password";
        } else if (formData.password !== formData.confirm_password) {
          stepErrors.confirm_password = "Passwords do not match";
        }
      }
      
      if (Object.keys(stepErrors).length > 0) {
        setErrors(prev => ({ ...prev, ...stepErrors }));
        return false;
      }
      return true;
    });

    if (!allStepsValid) {
      toast.error("Please fix all errors before submitting");
      // Navigate to first step with errors
      if (errors.first_name || errors.last_name || errors.email || errors.phone) {
        setCurrentStep(0);
      } else if (errors.name || errors.tin || errors.general_zone || errors.region || errors.woreda) {
        setCurrentStep(1);
      } else if (errors.latitude || errors.longitude) {
        setCurrentStep(2);
      } else if (errors.password || errors.confirm_password) {
        setCurrentStep(4);
      }
      return;
    }

    try {
      await register(formData);
      // Check store state directly since Zustand updates are synchronous
      // but React may not have re-rendered yet
      const storeState = useSupplierStore.getState();
      
      // Only show success dialog if registration truly succeeded
      // The useEffect will handle setting showSuccess based on success, supplier_id, and error state
      if (storeState.error || !storeState.success || !storeState.supplier_id) {
        // Ensure dialog is closed if there's an error or missing data
        setShowSuccess(false);
        return;
      }
      // If we reach here, registration succeeded - useEffect will handle showing the dialog
    } catch (error) {
      // Additional error handling - ensure dialog is closed on any error
      setShowSuccess(false);
      // Error is already handled in the store and shown via toast
    }
  };

  return (
    <div className="grid grid-cols-1 gap-4 mx-auto max-w-6xl my-6 sm:my-10 px-4 sm:px-6 lg:px-8 w-full">
      <div className="mb-6">
        <p className="text-xl font-semibold text-black">Welcome Aboard!</p>
        <p className="text-md font-medium text-gray-700">
          Fill in the forms accordingly to register.
        </p>
      </div>

      {needsPlanSelection && (
        <Card className="rounded-sm border-2 border-gray-200 shadow-none">
          <CardContent className="pt-6 px-4 sm:px-6">
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
          <div className="flex items-center justify-between mb-6 overflow-x-auto">
            {baseSteps.map((step, idx) => (
              <div 
                key={step.title} 
                className="flex-1 flex flex-col items-center cursor-pointer select-none"
                onClick={() => handleStepClick(idx)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    handleStepClick(idx);
                  }
                }}
                tabIndex={0}
                role="button"
                aria-label={`Go to step ${idx + 1}: ${step.title}`}
              >
                <div
                  className={`rounded-full w-8 h-8 flex items-center justify-center text-white font-bold transition-all hover:scale-110 ${
                    idx === currentStep
                      ? "bg-primary ring-2 ring-primary ring-offset-2"
                      : idx < currentStep
                      ? "bg-primary hover:bg-primary/90"
                      : "bg-gray-300 hover:bg-gray-400"
                  }`}
                >
                  {idx + 1}
                </div>
                <span
                  className={`text-xs mt-2 text-center transition-colors ${
                    idx === currentStep
                      ? "text-cyan-700 font-semibold"
                      : "text-gray-500 hover:text-gray-700"
                  }`}
                >
                  {step.title}
                </span>
                {idx < baseSteps.length - 1 && (
                  <div className="w-full h-1 bg-gray-200 my-2 pointer-events-none">
                    <div
                      className={`h-1 transition-colors ${
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
            <CardContent className="pt-6 px-4 sm:px-6">
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
                        <Label htmlFor="FirstName">First Name <span className="text-red-500">*</span></Label>
                        <Input
                          id="FirstName"
                          name="FirstName"
                          value={formData.first_name}
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              first_name: value,
                            }));
                            const error = validateField("first_name", value);
                            setErrors(prev => ({
                              ...prev,
                              first_name: error || undefined,
                            }));
                          }}
                          placeholder="Abebe"
                          required
                          className={errors.first_name ? "border-red-500" : ""}
                        />
                        {errors.first_name && (
                          <p className="text-sm text-red-500">{errors.first_name}</p>
                        )}
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="LastName">Last Name <span className="text-red-500">*</span></Label>
                        <Input
                          id="LastName"
                          name="LastName"
                          value={formData.last_name}
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              last_name: value,
                            }));
                            const error = validateField("last_name", value);
                            setErrors(prev => ({
                              ...prev,
                              last_name: error || undefined,
                            }));
                          }}
                          placeholder="Kebede"
                          required
                          className={errors.last_name ? "border-red-500" : ""}
                        />
                        {errors.last_name && (
                          <p className="text-sm text-red-500">{errors.last_name}</p>
                        )}
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Email">Email <span className="text-red-500">*</span></Label>
                        <Input
                          id="Email"
                          name="Email"
                          value={formData.email}
                          type="email"
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              email: value,
                            }));
                            const error = validateField("email", value);
                            setErrors(prev => ({
                              ...prev,
                              email: error || undefined,
                            }));
                          }}
                          placeholder="abebe.kebede@example.com"
                          required
                          className={errors.email ? "border-red-500" : ""}
                        />
                        {errors.email && (
                          <p className="text-sm text-red-500">{errors.email}</p>
                        )}
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="phone">Phone <span className="text-red-500">*</span></Label>
                        <Input
                          id="phone"
                          name="phone"
                          value={formData.phone}
                          onChange={(e) => {
                            let value = e.target.value;
                            // Ensure it starts with +
                            if (value && !value.startsWith('+')) {
                              value = '+' + value.replace(/[^0-9]/g, '');
                            } else {
                              // Only allow + and digits
                              value = '+' + value.replace(/[^0-9]/g, '').slice(0, 12);
                            }
                            setFormData((prev) => ({
                              ...prev,
                              phone: value,
                            }));
                            const error = validateField("phone", value);
                            setErrors(prev => ({
                              ...prev,
                              phone: error || undefined,
                            }));
                          }}
                          placeholder="+251912345678"
                          maxLength={13}
                          required
                          className={errors.phone ? "border-red-500" : ""}
                        />
                        {errors.phone && (
                          <p className="text-sm text-red-500">{errors.phone}</p>
                        )}
                        <p className="text-xs text-muted-foreground">Format: +251912345678 (12 digits after +)</p>
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
                        <Label htmlFor="Name">Name <span className="text-red-500">*</span></Label>
                        <Input
                          id="Name"
                          name="Name"
                          value={formData.name}
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              name: value,
                            }));
                            const error = validateField("name", value);
                            setErrors(prev => ({
                              ...prev,
                              name: error || undefined,
                            }));
                          }}
                          placeholder="Business Name"
                          required
                          className={errors.name ? "border-red-500" : ""}
                        />
                        {errors.name && (
                          <p className="text-sm text-red-500">{errors.name}</p>
                        )}
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Tin">TIN <span className="text-red-500">*</span></Label>
                        <Input
                          id="Tin"
                          name="Tin"
                          value={formData.tin}
                          onChange={(e) => {
                            // Only allow digits
                            const value = e.target.value.replace(/\D/g, '');
                            setFormData((prev) => ({
                              ...prev,
                              tin: value,
                            }));
                            const error = validateField("tin", value);
                            setErrors(prev => ({
                              ...prev,
                              tin: error || undefined,
                            }));
                          }}
                          placeholder="1234567890"
                          maxLength={10}
                          required
                          className={errors.tin ? "border-red-500" : ""}
                        />
                        {errors.tin && (
                          <p className="text-sm text-red-500">{errors.tin}</p>
                        )}
                        <p className="text-xs text-muted-foreground">Must be exactly 10 digits</p>
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="General Zone">General Zone <span className="text-red-500">*</span></Label>
                        <Input
                          id="General Zone"
                          name="General Zone"
                          value={formData.general_zone}
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              general_zone: value,
                            }));
                            const error = validateField("general_zone", value);
                            setErrors(prev => ({
                              ...prev,
                              general_zone: error || undefined,
                            }));
                          }}
                          placeholder="Bole"
                          required
                          className={errors.general_zone ? "border-red-500" : ""}
                        />
                        {errors.general_zone && (
                          <p className="text-sm text-red-500">{errors.general_zone}</p>
                        )}
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Region">Region <span className="text-red-500">*</span></Label>
                        <Input
                          id="Region"
                          name="Region"
                          value={formData.region}
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              region: value,
                            }));
                            const error = validateField("region", value);
                            setErrors(prev => ({
                              ...prev,
                              region: error || undefined,
                            }));
                          }}
                          placeholder="Addis Ababa"
                          required
                          className={errors.region ? "border-red-500" : ""}
                        />
                        {errors.region && (
                          <p className="text-sm text-red-500">{errors.region}</p>
                        )}
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="Woreda">Woreda <span className="text-red-500">*</span></Label>
                        <Input
                          id="Woreda"
                          name="Woreda"
                          value={formData.woreda}
                          onChange={(e) => {
                            const value = e.target.value;
                            setFormData((prev) => ({
                              ...prev,
                              woreda: value,
                            }));
                            const error = validateField("woreda", value);
                            setErrors(prev => ({
                              ...prev,
                              woreda: error || undefined,
                            }));
                          }}
                          placeholder="Bole Sub-city"
                          required
                          className={errors.woreda ? "border-red-500" : ""}
                        />
                        {errors.woreda && (
                          <p className="text-sm text-red-500">{errors.woreda}</p>
                        )}
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
                      setFormData={(callback) => {
                        setFormData(callback);
                        // Clear location errors when location is updated
                        if (errors.latitude || errors.longitude) {
                          setErrors(prev => {
                            const newErrors = { ...prev };
                            delete newErrors.latitude;
                            delete newErrors.longitude;
                            return newErrors;
                          });
                        }
                      }}
                    />
                  </div>
                  {(errors.latitude || errors.longitude) && (
                    <div className="text-sm text-red-500">
                      {errors.latitude || errors.longitude}
                    </div>
                  )}
                  {formData.latitude && formData.longitude && !errors.latitude && !errors.longitude && (
                    <div className="text-sm text-emerald-600">
                      Location selected: {parseFloat(formData.latitude).toFixed(6)}, {parseFloat(formData.longitude).toFixed(6)}
                    </div>
                  )}
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
                        onChange={(e) => {
                          const value = e.target.value;
                          setFormData((prev) => ({
                            ...prev,
                            password: value,
                          }));
                          const error = validateField("password", value);
                          setErrors(prev => ({
                            ...prev,
                            password: error || undefined,
                          }));
                          // Also validate confirm_password if it has a value
                          if (formData.confirm_password) {
                            const confirmError = validateField("confirm_password", formData.confirm_password, value);
                            setErrors(prev => ({
                              ...prev,
                              confirm_password: confirmError || undefined,
                            }));
                          }
                        }}
                        className={errors.password ? "border-red-500" : ""}
                        minLength={6}
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
                    {errors.password && (
                      <p className="text-sm text-red-500">{errors.password}</p>
                    )}
                    {formData.password && !errors.password && formData.password.length >= 6 && (
                      <p className="text-sm text-emerald-600">Password is valid</p>
                    )}
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
                        onChange={(e) => {
                          const value = e.target.value;
                          setFormData((prev) => ({
                            ...prev,
                            confirm_password: value,
                          }));
                          const error = validateField("confirm_password", value, formData.password);
                          setErrors(prev => ({
                            ...prev,
                            confirm_password: error || undefined,
                          }));
                        }}
                        className={errors.confirm_password ? "border-red-500" : ""}
                        minLength={6}
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
                    {errors.confirm_password && (
                      <p className="text-sm text-red-500">{errors.confirm_password}</p>
                    )}
                    {formData.confirm_password && !errors.confirm_password && formData.password === formData.confirm_password && (
                      <p className="text-sm text-emerald-600">Passwords match</p>
                    )}
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
                    <Button 
                      type="button" 
                      onClick={handleNext}
                      disabled={loading}
                    >
                      Next
                    </Button>
                  ) : (
                    <Button 
                      type="button" 
                      onClick={(e) => {
                        e.preventDefault();
                        handleSubmit();
                      }} 
                      disabled={loading}
                    >
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

      {showSuccess && (
      <ClientOnlyDialog open={showSuccess} onOpenChange={setShowSuccess}>
        <CheckCircle className="w-16 h-16 text-emerald-600 mx-auto mb-4" />
        <h1 className="text-3xl font-bold text-emerald-800 mb-2">
          Thank You for Registering!
        </h1>
        <p className="text-green-700 mb-4">
          Our team will get back to you shortly after reviewing your profile information.
        </p>
        <div className="flex justify-center mt-4">
          <Button onClick={() => router.push("/login")}>
            Continue to Login
          </Button>
        </div>
        {error && <p className="text-red-600 text-sm mt-2">{error}</p>}
      </ClientOnlyDialog>
      )}
    </div>
  );
}
