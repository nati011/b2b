import SessionProvider from '@/app/sessionprovider'

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      {children}
    </SessionProvider>
  )
} 