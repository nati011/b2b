'use client';
import React from 'react'
import { SessionProvider as AuthSessionProvider } from 'next-auth/react'

interface SessionProviderProps {
  children: React.ReactNode
}

const SessionProvider: React.FC<SessionProviderProps> = ({ children }) => {
  return (
    <AuthSessionProvider
    
    >
      {children}
    </AuthSessionProvider>
  )
}

export default SessionProvider