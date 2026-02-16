"use client";

import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PiTelegramLogo } from "react-icons/pi";
import { FaWhatsapp } from "react-icons/fa";
import Link from "next/link";

export default function ContactSupplier() {
  const telegramUrl = "https://t.me/efoyetastore";
  const whatsappUrl = "https://wa.me/251937639608";

  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-24">
        {/* Header Section */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          <p className="text-sm font-semibold text-primary mb-3 uppercase tracking-wide">
            Become a Supplier
          </p>
          <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 mb-6">
            Get in touch with us
          </h1>
          <p className="text-lg text-gray-600 leading-relaxed">
            Interested in becoming an Efoyeta Store supplier? We&apos;re here
            to help. Reach out on Telegram or WhatsApp and we&apos;ll get back
            to you with next steps.
          </p>
        </div>

        {/* Contact Options Section */}
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 lg:gap-8">
            {/* Telegram CTA Section */}
            <div className="bg-white rounded-2xl shadow-lg p-8 md:p-12 border border-gray-100">
              <div className="text-center mb-8">
                <div className="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-6">
                  <PiTelegramLogo className="w-10 h-10 text-primary" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900 mb-4">
                  Chat with us on Telegram
                </h2>
                <p className="text-gray-600 text-lg mb-8">
                  Reach out on Telegram to discuss supplier opportunities.
                  We&apos;ll respond with details and next steps.
                </p>
              </div>

              <div className="flex flex-col items-center gap-4">
                <Button
                  size="lg"
                  asChild
                  className="w-full bg-primary hover:bg-primary/90 text-white px-8 py-6 text-lg font-semibold shadow-lg hover:shadow-xl transition-all duration-300"
                >
                  <a
                    href={telegramUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center justify-center gap-3"
                  >
                    <PiTelegramLogo className="w-6 h-6" />
                    <span>Open Telegram</span>
                    <ArrowRight className="w-5 h-5" />
                  </a>
                </Button>
                <p className="text-sm text-gray-500">
                  Don&apos;t have Telegram?{" "}
                  <a
                    href="https://telegram.org/apps"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline"
                  >
                    Download it here
                  </a>
                </p>
              </div>
            </div>

            {/* WhatsApp CTA Section */}
            <div className="bg-white rounded-2xl shadow-lg p-8 md:p-12 border border-gray-100">
              <div className="text-center mb-8">
                <div className="w-20 h-20 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-6">
                  <FaWhatsapp className="w-10 h-10 text-green-600" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900 mb-4">
                  Chat with us on WhatsApp
                </h2>
                <p className="text-gray-600 text-lg mb-8">
                  Connect with us on WhatsApp to inquire about becoming a
                  supplier. We&apos;ll get back to you shortly.
                </p>
              </div>

              <div className="flex flex-col items-center gap-4">
                <Button
                  size="lg"
                  asChild
                  className="w-full bg-green-600 hover:bg-green-700 text-white px-8 py-6 text-lg font-semibold shadow-lg hover:shadow-xl transition-all duration-300"
                >
                  <a
                    href={whatsappUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center justify-center gap-3"
                  >
                    <FaWhatsapp className="w-6 h-6" />
                    <span>Open WhatsApp</span>
                    <ArrowRight className="w-5 h-5" />
                  </a>
                </Button>
                <p className="text-sm text-gray-500">
                  Don&apos;t have WhatsApp?{" "}
                  <a
                    href="https://www.whatsapp.com/download"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-green-600 hover:underline"
                  >
                    Download it here
                  </a>
                </p>
              </div>
            </div>
          </div>

          {/* Additional Info */}
          <div className="mt-12 text-center space-y-2">
            <p className="text-gray-600">
              Supplier inquiries are typically answered within a few hours during business hours.
            </p>
            <p className="text-sm text-gray-500">
              <Link href="/contact" className="text-primary hover:underline">
                General contact
              </Link>
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
