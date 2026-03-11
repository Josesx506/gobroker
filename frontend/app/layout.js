import './globals.css'
import Nav from '@/components/Nav'

export const metadata = {
  title: 'GoBroker',
  description: 'Real-time temperature monitoring',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <Nav />
        <main>{children}</main>
      </body>
    </html>
  )
}
