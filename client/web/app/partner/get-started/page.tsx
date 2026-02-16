"use client";

import { useSearchParams } from "next/navigation";
import { useState, useMemo } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ArrowLeft, Check } from "lucide-react";
import { PiTelegramLogo } from "react-icons/pi";
import { cn } from "@/lib/utils";

type StepId = "plan" | "details" | "review";

const steps: { id: StepId; label: string }[] = [
  { id: "plan", label: "Plan" },
  { id: "details", label: "Your Details" },
  { id: "review", label: "Review" },
];

const planFromId: Record<string, { name: string; price: number; period: string }> = {
  monthly: { name: "Monthly", price: 500, period: "per month" },
  yearly: { name: "Year", price: 4500, period: "per year" },
};

export default function PartnerGetStartedPage() {
  const searchParams = useSearchParams();
  const planId = searchParams.get("plan") || "yearly";
  const plan = planFromId[planId] || planFromId.yearly;

  const [currentStep, setCurrentStep] = useState<StepId>("plan");
  const [form, setForm] = useState({
    fullName: "",
    email: "",
    phone: "",
    businessName: "",
  });

  const currentStepIndex = steps.findIndex((s) => s.id === currentStep);

  const handleNext = () => {
    const idx = steps.findIndex((s) => s.id === currentStep);
    if (idx < steps.length - 1) setCurrentStep(steps[idx + 1].id);
  };

  const handleBack = () => {
    const idx = steps.findIndex((s) => s.id === currentStep);
    if (idx > 0) setCurrentStep(steps[idx - 1].id);
  };

  const canProceedFromPlan = true;
  const canProceedFromDetails = form.fullName.trim() && form.email.trim() && form.phone.trim();
  const canProceed = useMemo(() => {
    if (currentStep === "plan") return canProceedFromPlan;
    if (currentStep === "details") return canProceedFromDetails;
    return true;
  }, [currentStep, canProceedFromPlan, canProceedFromDetails]);

  const telegramUrl = "https://t.me/efoyetastore";

  return (
    <div className="min-h-screen bg-gray-50">
      <main className="container px-4 py-8 mx-auto">
        <Card className="max-w-2xl mx-auto shadow-none px-4">
          <div className="space-y-6">
            <div className="flex items-center gap-4 pt-2">
              <Button variant="ghost" size="icon" asChild className="h-9 w-9">
                <Link href="/#pricing">
                  <ArrowLeft className="h-5 w-5" />
                </Link>
              </Button>
              <h1 className="text-2xl font-semibold">Become a Partner</h1>
            </div>

            <div className="flex items-center justify-center gap-4 pb-4">
              {steps.map((step, index) => {
                const isActive = currentStepIndex === index;
                const isCompleted = currentStepIndex > index;
                const canNavigate = isCompleted || isActive;
                return (
                  <div key={step.id} className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => canNavigate && setCurrentStep(step.id)}
                      disabled={!canNavigate}
                      className={cn(
                        "flex items-center gap-2 transition-all rounded-lg px-2 py-1",
                        canNavigate && "cursor-pointer hover:bg-accent hover:opacity-90",
                        !canNavigate && "cursor-not-allowed opacity-50"
                      )}
                    >
                      <div
                        className={cn(
                          "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-all",
                          isActive && "bg-primary text-primary-foreground ring-2 ring-primary ring-offset-2",
                          isCompleted && "bg-primary text-primary-foreground",
                          !isActive && !isCompleted && "bg-secondary text-muted-foreground"
                        )}
                      >
                        {isCompleted ? <Check className="w-4 h-4" /> : index + 1}
                      </div>
                      <span
                        className={cn(
                          "text-sm font-medium",
                          isActive ? "text-foreground" : "text-muted-foreground"
                        )}
                      >
                        {step.label}
                      </span>
                    </button>
                    {index < steps.length - 1 && (
                      <div
                        className={cn("w-12 h-0.5", isCompleted ? "bg-primary" : "bg-border")}
                        aria-hidden
                      />
                    )}
                  </div>
                );
              })}
            </div>

            <div className="bg-white rounded-lg border border-gray-100 p-6">
              {currentStep === "plan" && (
                <div className="space-y-6">
                  <h2 className="text-xl font-semibold mb-4">Selected Plan</h2>
                  <div className="rounded-lg border-2 border-primary bg-primary/5 p-6">
                    <h3 className="text-lg font-bold text-gray-900">{plan.name}</h3>
                    <p className="text-2xl font-bold text-primary mt-2">
                      {plan.price.toLocaleString()} ETB
                    </p>
                    <p className="text-sm text-gray-600 mt-1">{plan.period}</p>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    You selected this plan from our pricing section. Continue to add your details so we can get in touch.
                  </p>
                </div>
              )}

              {currentStep === "details" && (
                <div className="space-y-6">
                  <h2 className="text-xl font-semibold mb-4">Your Details</h2>
                  <div className="space-y-4">
                    <div>
                      <label htmlFor="fullName" className="text-sm font-medium text-foreground block mb-2">
                        Full name
                      </label>
                      <input
                        id="fullName"
                        type="text"
                        value={form.fullName}
                        onChange={(e) => setForm((f) => ({ ...f, fullName: e.target.value }))}
                        placeholder="Your full name"
                        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                      />
                    </div>
                    <div>
                      <label htmlFor="email" className="text-sm font-medium text-foreground block mb-2">
                        Email
                      </label>
                      <input
                        id="email"
                        type="email"
                        value={form.email}
                        onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
                        placeholder="you@example.com"
                        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                      />
                    </div>
                    <div>
                      <label htmlFor="phone" className="text-sm font-medium text-foreground block mb-2">
                        Phone
                      </label>
                      <input
                        id="phone"
                        type="tel"
                        value={form.phone}
                        onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))}
                        placeholder="+251 ..."
                        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                      />
                    </div>
                    <div>
                      <label htmlFor="businessName" className="text-sm font-medium text-foreground block mb-2">
                        Business name (optional)
                      </label>
                      <input
                        id="businessName"
                        type="text"
                        value={form.businessName}
                        onChange={(e) => setForm((f) => ({ ...f, businessName: e.target.value }))}
                        placeholder="Your business or shop name"
                        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                      />
                    </div>
                  </div>
                </div>
              )}

              {currentStep === "review" && (
                <div className="space-y-6">
                  <h2 className="text-xl font-semibold mb-4">Review</h2>
                  <div className="space-y-4 rounded-lg border border-gray-200 p-4">
                    <div>
                      <p className="text-xs text-muted-foreground">Plan</p>
                      <p className="font-medium">{plan.name} – {plan.price.toLocaleString()} ETB {plan.period}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Name</p>
                      <p className="font-medium">{form.fullName || "—"}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Email</p>
                      <p className="font-medium">{form.email || "—"}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Phone</p>
                      <p className="font-medium">{form.phone || "—"}</p>
                    </div>
                    {form.businessName && (
                      <div>
                        <p className="text-xs text-muted-foreground">Business</p>
                        <p className="font-medium">{form.businessName}</p>
                      </div>
                    )}
                  </div>
                  <div className="bg-primary/5 border border-primary/20 rounded-lg p-4 space-y-2">
                    <p className="text-sm font-medium text-foreground">Next step</p>
                    <p className="text-sm text-muted-foreground">
                      Contact us on Telegram with the details above. We will confirm your plan and guide you through the rest.
                    </p>
                  </div>
                  <Button size="lg" className="w-full" asChild>
                    <a
                      href={telegramUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center justify-center gap-2"
                    >
                      <PiTelegramLogo className="w-5 h-5" />
                      Contact us on Telegram
                    </a>
                  </Button>
                </div>
              )}
            </div>

            {currentStep !== "review" && (
              <div className="flex gap-4">
                {currentStepIndex > 0 && (
                  <Button variant="outline" onClick={handleBack} className="flex-1">
                    Back
                  </Button>
                )}
                <Button
                  className="flex-1"
                  size="lg"
                  onClick={handleNext}
                  disabled={!canProceed}
                >
                  Continue
                </Button>
              </div>
            )}
          </div>
        </Card>
      </main>
    </div>
  );
}
