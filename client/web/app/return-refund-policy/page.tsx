"use client";

export default function ReturnRefundPolicyPage() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-24">
        <div className="max-w-4xl mx-auto space-y-10">
          <header className="text-center space-y-4">
            <p className="text-sm font-semibold text-primary uppercase tracking-wide">
              Policy
            </p>
            <h1 className="text-4xl md:text-5xl font-bold text-gray-900">
              Return & Refund Policy
            </h1>
            <p className="text-lg text-gray-600 leading-relaxed">
              Please read this policy carefully before making a purchase. By
              placing an order, you agree to the terms below.
            </p>
          </header>

          <section className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8 space-y-6">
            <h2 className="text-2xl font-semibold text-gray-900">
              Return & Refund Policy
            </h2>
            <ul className="list-disc list-inside space-y-2 text-gray-700">
              <li>
                Returns are accepted within 2 days of delivery if the product is
                damaged or defective.
              </li>
              <li>
                For items shipped from partner suppliers, returns may take
                longer.
              </li>
              <li>
                Refunds are processed once the product is returned and verified.
              </li>
            </ul>
          </section>

          <section className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8 space-y-6">
            <h2 className="text-2xl font-semibold text-gray-900">
              Contact Us
            </h2>
            <p className="text-gray-700 leading-relaxed">
              Customers can contact us via WhatsApp, email, or customer support
              for any issues.
            </p>
            <p className="text-gray-700 leading-relaxed">
              For return or refund requests, reach us via the contact options on
              our{" "}
              <a href="/contact" className="text-primary hover:underline">
                Contact
              </a>{" "}
              page. We will respond as soon as possible.
            </p>
          </section>

          <p className="text-sm text-gray-500 text-center pt-4">
            Efoyeta Store PLC reserves the right to update this policy.
            Continued use of our services after changes constitutes acceptance of
            the updated policy.
          </p>
        </div>
      </div>
    </main>
  );
}
