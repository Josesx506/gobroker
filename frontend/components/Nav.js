'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import styles from './Nav.module.css'

export default function Nav() {
  const pathname = usePathname()

  return (
    <nav className={styles.nav}>
      <span className={styles.brand}>GoBroker</span>
      <div className={styles.links}>
        <Link href="/charts" className={`${styles.link} ${pathname === '/charts' ? styles.active : ''}`}>Charts</Link>
        <Link href="/locations" className={`${styles.link} ${pathname === '/locations' ? styles.active : ''}`}>Locations</Link>
      </div>
    </nav>
  )
}
