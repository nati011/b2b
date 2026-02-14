import { Product, Brand } from '@/types/product';

export const products: Product[] = [
  {
    id: '1',
    name: 'Industrial Safety Helmet Pro',
    brand: 'SafeGuard',
    price: 89,
    originalPrice: 119,
    category: 'Equipment',
    images: [
      'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=800&q=80',
      'https://images.unsplash.com/photo-1581094794329-c8112a89af12?w=800&q=80',
    ],
    description: 'ANSI-certified industrial safety helmet with adjustable suspension system. Designed for maximum protection in harsh work environments.',
    details: ['ANSI Z89.1 Type I certified', 'UV-resistant ABS shell', '4-point suspension', 'Accessory slots compatible'],
    fit: 'Adjustable fit: 6.5" to 8" head circumference',
    sizes: ['Standard', 'Large'],
    colors: ['Yellow', 'White', 'Orange'],
    isNew: true,
    isTrending: true,
  },
  {
    id: '2',
    name: 'Cordless Impact Wrench 20V',
    brand: 'PowerMax',
    price: 349,
    category: 'Tools',
    images: [
      'https://images.unsplash.com/photo-1572981779307-38b8cabb2407?w=800&q=80',
      'https://images.unsplash.com/photo-1530124566582-a618bc2615dc?w=800&q=80',
    ],
    description: 'High-torque cordless impact wrench delivering 1,200 ft-lbs of max torque. Includes 2 batteries and rapid charger.',
    details: ['1,200 ft-lbs max torque', 'Brushless motor', '20V 5.0Ah battery', '3-speed settings'],
    fit: 'Weight: 6.2 lbs with battery',
    sizes: ['1/2" Drive'],
    colors: ['Red/Black'],
    isNew: true,
  },
  {
    id: '3',
    name: 'Heavy-Duty Pallet Jack',
    brand: 'LiftPro',
    price: 475,
    category: 'Equipment',
    images: [
      'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=800&q=80',
      'https://images.unsplash.com/photo-1553413077-190dd305871c?w=800&q=80',
    ],
    description: 'Professional-grade pallet jack with 5,500 lb capacity. Precision-machined hydraulics for smooth operation.',
    details: ['5,500 lb capacity', 'Fork dimensions: 48" x 27"', 'Nylon wheels', '3-position control handle'],
    fit: 'Lowered height: 2.9" | Raised height: 7.5"',
    sizes: ['Standard', 'Narrow'],
    colors: ['Yellow', 'Blue'],
    isTrending: true,
  },
  {
    id: '4',
    name: 'Steel-Toe Work Boots',
    brand: 'ToughStep',
    price: 159,
    category: 'Footwear',
    images: [
      'https://images.unsplash.com/photo-1520639888713-7851133b1ed0?w=800&q=80',
      'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&q=80',
    ],
    description: 'Waterproof steel-toe boots with slip-resistant outsole. ASTM F2413 rated for impact and compression protection.',
    details: ['ASTM F2413-18 rated', 'Waterproof leather', 'Electrical hazard safe', 'Composite shank'],
    fit: 'True to size. Wide widths available.',
    sizes: ['7', '8', '9', '10', '11', '12', '13'],
    colors: ['Brown', 'Black'],
    isNew: true,
    isTrending: true,
  },
  {
    id: '5',
    name: 'Digital Multimeter Pro',
    brand: 'TechMeasure',
    price: 189,
    category: 'Tools',
    images: [
      'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&q=80',
      'https://images.unsplash.com/photo-1517077304055-6e89abbf09b0?w=800&q=80',
    ],
    description: 'Professional-grade True RMS multimeter with Bluetooth connectivity. Auto-ranging with 10,000 count display.',
    details: ['CAT IV 600V / CAT III 1000V', 'True RMS AC/DC', 'Bluetooth data logging', 'IP54 rated'],
    fit: 'Includes test leads, temperature probe, carrying case',
    sizes: ['Standard'],
    colors: ['Yellow/Black'],
    isNew: true,
  },
  {
    id: '6',
    name: 'Industrial First Aid Kit',
    brand: 'MedReady',
    price: 149,
    category: 'Accessories',
    images: [
      'https://images.unsplash.com/photo-1603398938378-e54eab446dde?w=800&q=80',
      'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800&q=80',
    ],
    description: 'OSHA-compliant first aid kit for up to 100 workers. Wall-mountable metal case with organized compartments.',
    details: ['OSHA compliant', '446 pieces included', 'Wall mountable', 'Refillable design'],
    fit: 'Dimensions: 15" x 10" x 5"',
    sizes: ['50 Person', '100 Person'],
    colors: ['White/Red'],
    isTrending: true,
  },
  {
    id: '7',
    name: 'Warehouse Storage Rack',
    brand: 'RackMaster',
    price: 899,
    category: 'Equipment',
    images: [
      'https://images.unsplash.com/photo-1553413077-190dd305871c?w=800&q=80',
      'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=800&q=80',
    ],
    description: 'Heavy-duty pallet racking system with 10,000 lb capacity per level. Quick-assembly design with safety locks.',
    details: ['10,000 lb per level', 'Powder-coated steel', 'Adjustable beam heights', 'Seismic design'],
    fit: 'Dimensions: 96"H x 48"D x 120"W',
    sizes: ['8ft', '10ft', '12ft'],
    colors: ['Blue/Orange', 'Green/Orange'],
  },
  {
    id: '8',
    name: 'High-Visibility Safety Vest',
    brand: 'SafeGuard',
    price: 24,
    category: 'Accessories',
    images: [
      'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&q=80',
      'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=800&q=80',
    ],
    description: 'ANSI Class 3 high-visibility safety vest with reflective strips. Breathable mesh construction.',
    details: ['ANSI/ISEA 107 Class 3', '360° reflectivity', 'Breathable mesh', 'Multiple pockets'],
    fit: 'Sizes run large. Order one size down for fitted look.',
    sizes: ['S', 'M', 'L', 'XL', '2XL', '3XL'],
    colors: ['Lime Green', 'Orange'],
    isTrending: true,
  },
  {
    id: '9',
    name: 'Industrial Pressure Washer',
    brand: 'CleanForce',
    price: 1295,
    category: 'Equipment',
    images: [
      'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&q=80',
      'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=800&q=80',
    ],
    description: 'Commercial-grade pressure washer with 4,000 PSI output. CAT triplex pump with thermal relief valve.',
    details: ['4,000 PSI @ 4.0 GPM', 'CAT triplex pump', 'Honda GX390 engine', '50ft hose included'],
    fit: 'Weight: 185 lbs | Dimensions: 42"L x 24"W x 38"H',
    sizes: ['Standard'],
    colors: ['Red/Black'],
    isNew: true,
  },
  {
    id: '10',
    name: 'LED Work Light Tower',
    brand: 'BrightSite',
    price: 549,
    category: 'Equipment',
    images: [
      'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=800&q=80',
      'https://images.unsplash.com/photo-1581094794329-c8112a89af12?w=800&q=80',
    ],
    description: 'Portable LED light tower with 50,000 lumens output. Telescoping mast extends to 12 feet.',
    details: ['50,000 lumens', '12ft telescoping mast', 'Generator or battery powered', 'IP65 rated'],
    fit: 'Collapsed: 6ft | Extended: 12ft',
    sizes: ['Standard'],
    colors: ['Yellow', 'White'],
    isNew: true,
    isTrending: true,
  },
];

export const brands: Brand[] = [
  { id: '1', name: 'SafeGuard', logo: 'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=200&q=80' },
  { id: '2', name: 'PowerMax', logo: 'https://images.unsplash.com/photo-1572981779307-38b8cabb2407?w=200&q=80' },
  { id: '3', name: 'LiftPro', logo: 'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=200&q=80' },
  { id: '4', name: 'ToughStep', logo: 'https://images.unsplash.com/photo-1520639888713-7851133b1ed0?w=200&q=80' },
  { id: '5', name: 'TechMeasure', logo: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=200&q=80' },
  { id: '6', name: 'MedReady', logo: 'https://images.unsplash.com/photo-1603398938378-e54eab446dde?w=200&q=80' },
  { id: '7', name: 'RackMaster', logo: 'https://images.unsplash.com/photo-1553413077-190dd305871c?w=200&q=80' },
  { id: '8', name: 'CleanForce', logo: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=200&q=80' },
];

export const getProductById = (id: string): Product | undefined => {
  return products.find(p => p.id === id);
};

export const getNewArrivals = (): Product[] => {
  return products.filter(p => p.isNew);
};

export const getTrendingProducts = (): Product[] => {
  return products.filter(p => p.isTrending);
};

export const getProductsByCategory = (category: string): Product[] => {
  return products.filter(p => p.category === category);
};
