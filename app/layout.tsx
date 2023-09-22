import type { Metadata } from 'next'
export const metadata: Metadata = {
  title: "Schooldown",
  description: "School Countdown",
}
export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <link rel="favicon" href="/favicon.ico" />
      <body>{children}</body>
    </html>
  )
}
