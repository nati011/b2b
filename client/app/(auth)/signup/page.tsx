'use client'
import { useState } from "react";
import Link from "next/link";
import { Briefcase, Store } from "lucide-react";

export default function SignUpPage() {
  const [role, setRole] = useState<"supplier" | "customer">("customer");

  return (
    <div className="min-h-screen flex flex-col items-center mt-40">
      <h1 className="text-3xl font-semibold mb-8">Join as a supplier or customer</h1>
      <div className="flex gap-6 mb-8">
        <button
          className={`border rounded-lg px-8 py-8 flex flex-col items-center w-80 transition-all ${
            role === "supplier"
              ? "border-primary bg-gray-50 shadow-md"
              : "border-gray-300 bg-white"
          }`}
          onClick={() => setRole("supplier")}
        >
          <Briefcase className="h-8 w-8 mb-4" />
          <span className="text-lg font-medium mb-2">I'm a supplier.</span>
          <span
            className={`mt-4 h-5 w-5 rounded-full border-2 flex items-center justify-center ${
              role === "supplier" ? "border-primary" : "border-gray-300"
            }`}
          >
            {role === "supplier" && (
              <span className="h-3 w-3 bg-primary rounded-full block" />
            )}
          </span>
        </button>
        <button
          className={`border rounded-lg px-8 py-8 flex flex-col items-center w-80 transition-all ${
            role === "customer"
              ? "border-primary bg-gray-50 shadow-md"
              : "border-gray-300 bg-white"
          }`}
          onClick={() => setRole("customer")}
        >
          <Store className="h-8 w-8 mb-4" />
          <span className="text-lg font-medium mb-2">I'm a customer.</span>
          <span
            className={`mt-4 h-5 w-5 rounded-full border-2 flex items-center justify-center ${
              role === "customer" ? "border-primary" : "border-gray-300"
            }`}
          >
            {role === "customer" && (
              <span className="h-3 w-3 bg-primary rounded-full block" />
            )}
          </span>
        </button>
      </div>
      <Link
        href={
          role === "supplier"
            ? "/signup/supplier"
            : "/signup/customer"
        }
        className="bg-primary text-white px-8 py-3 rounded-md font-semibold hover:bg-primary transition mb-4"
      >
        {role === "supplier"
          ? "Register as Supplier"
          : "Register as Customer"}
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