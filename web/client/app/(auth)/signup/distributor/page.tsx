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
import { useRouter } from "next/navigation";
import dynamic from "next/dynamic";
import { Checkbox } from "@/components/ui/checkbox";
import ImageUpload from "@/components/ImageUpload";

const Map = dynamic(() => import("@/components/map"), { ssr: false });

const steps = [
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
];

export default function DistributorsForm() {
  const router = useRouter();
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
  });
  const [markerPosition, setMarkerPosition] = useState<[number, number]>([
    8.9934609, 38.7714897,
  ]);
  const [useCurrentLocation, setUseCurrentLocation] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);

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

  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep((prev) => prev + 1);
    }
  };

  const handleBack = () => {
    if (currentStep > 0) {
      setCurrentStep((prev) => prev - 1);
    }
  };

  const handleSubmit = async () => {
    setLoading(true);
    try {
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000));
      toast.success("Distributor registered successfully!");
      router.push("/distributors");
    } catch (e) {
      toast.error("Failed to register distributor.");
    } finally {
      setLoading(false);
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

      {/* Stepper */}
      <div className="flex items-center justify-between mb-6">
        {steps.map((step, idx) => (
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
            {idx < steps.length - 1 && (
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
            <h3 className="text-xl font-bold">{steps[currentStep].title}</h3>
            <p className="text-sm text-muted-foreground">
              {steps[currentStep].description}
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
              <ImageUpload onChange={function (value: string[]): void {
                throw new Error("Function not implemented.");
              } }/>
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
              {currentStep < steps.length - 1 ? (
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
    </div>
  );
}
