import Header from '@/components/Header';
import StartInterviewButton from "@/components/startInterviewButton";
import styles from '@/styles/index.module.css';

export default function Home() {
  return (
    <div className={styles.page}>
      <Header />
      <div className={styles.hero}>
        <span className={styles.cursor} aria-hidden="true" />
        <h1 className={styles.headline}>
          Think in code.<br />Interview under pressure.
        </h1>
        <p className={styles.subline}>
          A live AI coach works through problems with you — simulating a real technical interview.
        </p>
        <StartInterviewButton />
      </div>
    </div>
  );
}
