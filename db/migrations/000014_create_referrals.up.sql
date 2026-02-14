CREATE TABLE IF NOT EXISTS public.affiliates (
  id SERIAL PRIMARY KEY,
  affiliate_id INT, -- Optional: ID if affiliate is also a customer
  email VARCHAR(255),
  phone_number VARCHAR(50),
  full_name VARCHAR(255) NOT NULL,
  status VARCHAR(50),
  social_media_handle VARCHAR(255),
  platform VARCHAR(100), -- Primary social media platform
  is_active BOOLEAN DEFAULT FALSE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (affiliate_id) REFERENCES public.customers(id) ON DELETE SET NULL
);

COMMENT ON TABLE public.affiliates IS 'stores affiliate/influencer profiles (can be customers or non-customers)';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_affiliates_email ON public.affiliates(email) WHERE is_deleted = FALSE AND email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_affiliates_phone_number ON public.affiliates(phone_number) WHERE is_deleted = FALSE AND phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_affiliates_status ON public.affiliates(status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_affiliates_is_active ON public.affiliates(is_active) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_affiliates_affiliate_id ON public.affiliates(affiliate_id) WHERE is_deleted = FALSE AND affiliate_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS public.referral_codes (
  id SERIAL PRIMARY KEY,
  affiliate_id INT NOT NULL,
  code VARCHAR(20) NOT NULL UNIQUE,
  is_custom BOOLEAN DEFAULT FALSE,
  status VARCHAR(50),
  expires_at TIMESTAMP,
  is_active BOOLEAN DEFAULT FALSE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (affiliate_id) REFERENCES public.affiliates(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.referral_codes IS 'stores unique referral codes/links for affiliates';
CREATE INDEX IF NOT EXISTS idx_referral_codes_code ON public.referral_codes(code);
CREATE INDEX IF NOT EXISTS idx_referral_codes_affiliate_id ON public.referral_codes(affiliate_id) WHERE is_deleted = FALSE;

-- Additional indexes for referral_codes
CREATE INDEX IF NOT EXISTS idx_referral_codes_status ON public.referral_codes(status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_referral_codes_is_active ON public.referral_codes(is_active) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_referral_codes_expires_at ON public.referral_codes(expires_at) WHERE is_deleted = FALSE AND expires_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS public.referral_relationships (
  id SERIAL PRIMARY KEY,
  affiliate_id INT NOT NULL,
  customer_id INT NOT NULL,
  referral_code_id INT NOT NULL,
  source VARCHAR(100), -- Social media platform (Instagram, TikTok, etc.)
  campaign VARCHAR(255), -- Campaign name/ID
  utm_source VARCHAR(255),
  utm_medium VARCHAR(255),
  utm_campaign VARCHAR(255),
  utm_content VARCHAR(255),
  referral_link VARCHAR(500),
  ip_address VARCHAR(45),
  user_agent VARCHAR(500),
  is_active BOOLEAN DEFAULT TRUE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (affiliate_id) REFERENCES public.affiliates(id) ON DELETE CASCADE,
  FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE,
  FOREIGN KEY (referral_code_id) REFERENCES public.referral_codes(id) ON DELETE CASCADE,
  UNIQUE(customer_id) -- One referral relationship per customer
);

COMMENT ON TABLE public.referral_relationships IS 'tracks which affiliate referred which new customer';
CREATE INDEX IF NOT EXISTS idx_referral_relationships_affiliate_id ON public.referral_relationships(affiliate_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_referral_relationships_customer_id ON public.referral_relationships(customer_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_referral_relationships_referral_code_id ON public.referral_relationships(referral_code_id) WHERE is_deleted = FALSE;

-- Additional indexes for referral_relationships
CREATE INDEX IF NOT EXISTS idx_referral_relationships_is_active ON public.referral_relationships(is_active) WHERE is_deleted = FALSE;

CREATE TABLE IF NOT EXISTS public.commissions (
  id SERIAL PRIMARY KEY,
  affiliate_id INT NOT NULL,
  referral_relationship_id INT NOT NULL,
  order_id INT NOT NULL,
  customer_id INT NOT NULL,
  amount DECIMAL(10,2), -- Commission amount
  commission_rate DECIMAL(5,2), -- Percentage (e.g., 10.00 for 10%)
  status VARCHAR(50),
  order_total DECIMAL(10,2), -- Snapshot of order total at time of commission calculation
  notes TEXT,
  paid_at TIMESTAMP,
  reversed_at TIMESTAMP,
  is_active BOOLEAN DEFAULT TRUE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (affiliate_id) REFERENCES public.affiliates(id) ON DELETE CASCADE,
  FOREIGN KEY (referral_relationship_id) REFERENCES public.referral_relationships(id) ON DELETE CASCADE,
  FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
  FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.commissions IS 'stores affiliate commissions from referred customer orders';
CREATE INDEX IF NOT EXISTS idx_commissions_affiliate_id ON public.commissions(affiliate_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_commissions_order_id ON public.commissions(order_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_commissions_referral_relationship_id ON public.commissions(referral_relationship_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_commissions_customer_id ON public.commissions(customer_id) WHERE is_deleted = FALSE;

-- Additional indexes for commissions
CREATE INDEX IF NOT EXISTS idx_commissions_status ON public.commissions(status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_commissions_is_active ON public.commissions(is_active) WHERE is_deleted = FALSE;

