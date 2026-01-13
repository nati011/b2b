"use client";

import { MessageCircle, ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PiTelegramLogo } from "react-icons/pi";

export default function Contact() {
  const telegramUrl = "https://t.me/efoyetastore1";

  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-24">
        {/* Header Section */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          <p className="text-sm font-semibold text-primary mb-3 uppercase tracking-wide">
            Let&apos;s start a conversation
          </p>
          <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 mb-6">
            Get in touch with us
          </h1>
          <p className="text-lg text-gray-600 leading-relaxed">
            Have a project in mind or just want to learn more about how Efoyeta
            Store can elevate your business? We&apos;re here to help. Reach out
            to us on Telegram, and let&apos;s start a conversation.
          </p>
        </div>

        {/* Telegram CTA Section */}
        <div className="max-w-2xl mx-auto">
          <div className="bg-white rounded-2xl shadow-lg p-8 md:p-12 border border-gray-100">
            <div className="text-center mb-8">
              <div className="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-6">
                <PiTelegramLogo className="w-10 h-10 text-primary" />
              </div>
              <h2 className="text-3xl font-bold text-gray-900 mb-4">
                Chat with us on Telegram
              </h2>
              <p className="text-gray-600 text-lg mb-8">
                Get instant responses and personalized support from our team. 
                Click the button below to start a conversation with us on Telegram.
              </p>
            </div>

            <div className="flex flex-col items-center gap-4">
              <Button
                size="lg"
                asChild
                className="w-full md:w-auto bg-primary hover:bg-primary/90 text-white px-8 py-6 text-lg font-semibold shadow-lg hover:shadow-xl transition-all duration-300"
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

          {/* Additional Info */}
          <div className="mt-12 text-center">
            <p className="text-gray-600">
              We typically respond within a few hours during business hours.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
