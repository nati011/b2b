import { ShieldX, Home, ArrowLeft, AlertTriangle } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function Forbidden() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-background via-background to-muted/20">
      <div className="text-center space-y-8 px-4">
        {/* Icon */}
        <div className="flex justify-center">
          <div className="flex aspect-square size-24 items-center justify-center rounded-full bg-destructive/10 border border-destructive/20">
            <ShieldX className="h-12 w-12 text-destructive" />
          </div>
        </div>

        {/* Content */}
        <div className="space-y-4 max-w-md mx-auto">
          <h1 className="text-4xl font-bold tracking-tight text-foreground">
            403
          </h1>
          <h2 className="text-2xl font-semibold text-foreground">
            Access Forbidden
          </h2>
          <p className="text-muted-foreground text-lg leading-relaxed">
            Sorry, you don't have permission to access this page. Please contact your administrator if you believe this is an error.
          </p>
        </div>



        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Button asChild variant="outline" className="gap-2">
            <Link href="/">
              <ArrowLeft className="h-4 w-4" />
              Go Back
            </Link>
          </Button>
          <Button asChild className="gap-2">
            <Link href="/">
              <Home className="h-4 w-4" />
              Go Home
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
} 