export default function AboutPage() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-24">
        <div className="max-w-4xl mx-auto space-y-12">
          <header className="text-center space-y-4">
            <p className="text-sm font-semibold text-primary uppercase tracking-wide">
              About us
            </p>
            <h1 className="text-4xl md:text-5xl font-bold text-gray-900">
              About Efoyetastore
            </h1>
            <p className="text-lg text-gray-600 leading-relaxed">
              Efoyetastore is an innovative e-commerce platform that connects
              customers with high-quality products from trusted partner
              suppliers.
            </p>
          </header>

          <section className="space-y-4 text-gray-700 leading-relaxed">
            <p>
              We carefully select each product and work closely with suppliers
              to maintain quality and ensure timely delivery. Some products are
              shipped directly from our partner suppliers to offer a wider
              variety and faster fulfillment.
            </p>
          </section>

          <section className="grid gap-6 md:grid-cols-2">
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8">
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">
                Our Mission
              </h2>
              <p className="text-gray-600">
                To make shopping online easy, reliable, and accessible for
                everyone in Ethiopia.
              </p>
            </div>
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8">
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">
                Our Vision
              </h2>
              <p className="text-gray-600">
                To become Ethiopia&apos;s most trusted online marketplace while
                supporting local and international suppliers.
              </p>
            </div>
          </section>

          <section className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8 space-y-4">
            <h2 className="text-2xl font-semibold text-gray-900">
              Return &amp; Refund Policy
            </h2>
            <ul className="list-disc pl-5 text-gray-600 space-y-2">
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
            <p className="text-gray-600">
              Customers can contact us via WhatsApp, email, or customer support
              for any issues.
            </p>
          </section>
        </div>
      </div>
    </main>
  );
}


