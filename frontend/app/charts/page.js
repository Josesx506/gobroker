import TemperatureChart from '@/components/TemperatureChart'
import styles from './page.module.css'

const LOCATIONS = ['alpha', 'beta', 'charlie', 'delta', 'echo']

export default function ChartsPage() {
  return (
    <div className={styles.page}>
      <h1 className={styles.heading}>Temperature Overview</h1>
      <div className={styles.grid}>
        {LOCATIONS.map((loc) => (
          <TemperatureChart key={loc} locationId={loc} />
        ))}
      </div>
    </div>
  )
}
