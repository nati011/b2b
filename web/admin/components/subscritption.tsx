'use client'
import { useSubscriptionStore } from "@/app/libs/store/useSubscriptionStore";
import { Button } from "./ui/button";
import { toast } from "sonner";
import { useEffect } from "react";

const formatDate = (dateString: string | undefined): string => {
    if (!dateString) return "N/A";
    const date = new Date(dateString);
    return isNaN(date.getTime()) ? "Invalid Date" : date.toLocaleDateString();
  };
  
  const handleRenewSubscription = () => {
    toast.info("Renewal functionality coming soon.");
    // TODO: Integrate 
  };

const SubscriptionMessage = () => {
  const {
    subscription,
    subscriptionError,
    subscriptionLoading,
    fetchSubscription,
  } = useSubscriptionStore();


  useEffect(() => {
    fetchSubscription();
  }, [fetchSubscription]);

    return (
        <div className="w-full h-10 bg-amber-300/50 border-b-3 border-amber-300 px-4">
                        <div className="flex justify-between items-center">
              <div className="space-y-1">
                <h4 className="text-lg font-semibold">
                  {subscription?.subscription_plan_name || "Unknown Plan"}
                </h4>
                <p className="text-sm text-muted-foreground">
                  Purchased on {formatDate(subscription?.created_date.toLocaleString())}
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
        </div>
    )
}

export default SubscriptionMessage;