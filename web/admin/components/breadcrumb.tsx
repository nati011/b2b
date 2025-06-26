import React from 'react';
import Link from 'next/link';
import { ChevronRight, Home } from 'lucide-react';

interface BreadcrumbItem {
  title: string;
  href: string;
}

interface BreadcrumbProps {
  page: BreadcrumbItem[];
  heading: string;
  subheading?: string;
}

const Breadcrumb: React.FC<BreadcrumbProps> = ({ page, heading, subheading }) => {
  return (
    <div className="mb-6 space-y-2 print:hidden">
      {/* Breadcrumb Navigation */}
      <nav className="flex items-center space-x-1 text-sm text-muted-foreground">
        <Link
          href="/"
          className="flex items-center hover:text-foreground transition-colors"
        >
          <Home className="h-4 w-4" />
        </Link>
        {page.map((item, index) => (
          <React.Fragment key={index}>
            <ChevronRight className="h-4 w-4" />
            <Link
              href={item.href}
              className={`hover:text-foreground transition-colors ${index === page.length - 1 ? 'text-foreground font-medium' : ''
                }`}
            >
              {item.title}
            </Link>
          </React.Fragment>
        ))}
      </nav>

      {/* Page Header */}
      <div className="space-y-1">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">
          {heading}
        </h1>
        {subheading && (
          <p className="text-muted-foreground">
            {subheading}
          </p>
        )}
      </div>
    </div>
  );
};

export default Breadcrumb;