ALTER TABLE public.orders
  ADD COLUMN IF NOT EXISTS referral_code VARCHAR(20);

COMMENT ON COLUMN public.orders.referral_code IS 'optional referral code submitted with the order';
