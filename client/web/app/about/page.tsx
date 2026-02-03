"use client";

export default function AboutPage() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 lg:py-24">
        <div className="max-w-4xl mx-auto space-y-12">
          {/* Header Section */}
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

          {/* About Efoyetastore Section */}
          <section className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8 space-y-4">
            <p className="text-gray-700 leading-relaxed">
              We carefully select each product and work closely with suppliers
              to maintain quality and ensure timely delivery. Some products are
              shipped directly from our partner suppliers to offer a wider
              variety and faster fulfillment.
            </p>
          </section>

          {/* About the Founder Section */}
          <section className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8 space-y-4">
            <h2 className="text-2xl md:text-3xl font-semibold text-gray-900 mb-4">
              About the Founder
            </h2>
            <div className="space-y-4 text-gray-700 leading-relaxed">
              <p>
                Efoyetastore General Trading PLC was founded by{" "}
                <span className="font-semibold text-gray-900">Aragaw Melak</span>
                , a dedicated entrepreneur with over 5 years of experience in
                dropshipping and e-commerce. Born in a rural part of Ethiopia,
                Aragaw started from humble beginnings and, through hard work and
                determination, is now a self-made millionaire.
              </p>
              <p>
                His mission with Efoyetastore is to bring high-quality products
                to Ethiopia and eventually reach the international market.
                Aragaw believes in trust, transparency, and customer
                satisfaction, ensuring every shopping experience is safe, smooth,
                and reliable.
              </p>
            </div>
          </section>

          {/* Mission & Vision Section */}
          <section className="grid gap-6 md:grid-cols-2">
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8">
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">
                Our Mission
              </h2>
              <p className="text-gray-600 leading-relaxed">
                To make shopping online easy, reliable, and accessible for
                everyone in Ethiopia.
              </p>
            </div>
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sm:p-8">
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">
                Our Vision
              </h2>
              <p className="text-gray-600 leading-relaxed">
                To become Ethiopia&apos;s most trusted online marketplace while
                supporting local and international suppliers.
              </p>
            </div>
          </section>
        </div>
      </div>
    </main>
  );
}


