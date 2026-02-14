'use client'
import { useState } from "react";
import Link from "next/link";
import { Briefcase, Store } from "lucide-react";

export default function SignUpPage() {
  const [role, setRole] = useState<"distributor" | "retailer">("retailer");

  return (
    <div className="min-h-screen flex flex-col items-center mt-40">
      <h1 className="text-3xl font-semibold mb-8">Join as a distributor or retailer</h1>
      <div className="flex gap-6 mb-8">
        <button
          className={`border rounded-lg px-8 py-8 flex flex-col items-center w-80 transition-all ${
            role === "distributor"
              ? "border-primary bg-gray-50 shadow-md"
              : "border-gray-300 bg-white"
          }`}
          onClick={() => setRole("distributor")}
        >
          <Briefcase className="h-8 w-8 mb-4" />
          <span className="text-lg font-medium mb-2">I'm a distributor.</span>
          <span
            className={`mt-4 h-5 w-5 rounded-full border-2 flex items-center justify-center ${
              role === "distributor" ? "border-primary" : "border-gray-300"
            }`}
          >
            {role === "distributor" && (
              <span className="h-3 w-3 bg-primary rounded-full block" />
            )}
          </span>
        </button>
        <button
          className={`border rounded-lg px-8 py-8 flex flex-col items-center w-80 transition-all ${
            role === "retailer"
              ? "border-primary bg-gray-50 shadow-md"
              : "border-gray-300 bg-white"
          }`}
          onClick={() => setRole("retailer")}
        >
          <Store className="h-8 w-8 mb-4" />
          <span className="text-lg font-medium mb-2">I'm a retailer.</span>
          <span
            className={`mt-4 h-5 w-5 rounded-full border-2 flex items-center justify-center ${
              role === "retailer" ? "border-primary" : "border-gray-300"
            }`}
          >
            {role === "retailer" && (
              <span className="h-3 w-3 bg-primary rounded-full block" />
            )}
          </span>
        </button>
      </div>
      <Link
        href={
          role === "distributor"
            ? "/signup/distributor"
            : "/signup/retailer"
        }
        className="bg-primary text-white px-8 py-3 rounded-md font-semibold hover:bg-primary transition mb-4"
      >
        {role === "distributor"
          ? "Register as Distributor"
          : "Register as Retailer"}
      </Link>
      <div>
        Already have an account?{" "}
        <Link href="/login" className="text-primary underline">
          Log In
        </Link>
      </div>
    </div>
  );
}