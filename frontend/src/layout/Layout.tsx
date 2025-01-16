import { Outlet } from "react-router-dom";
import styles from "./Layout.module.scss";
import { MainMenu } from "@shared/components/MainMenu/MainMenu";

export default function Layout() {
  return (
    <main className={styles.layout}>
      <MainMenu />
      <Outlet />
    </main>
  );
}
