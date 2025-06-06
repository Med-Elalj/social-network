import Image from "next/image";
import Styles from "./global.module.css";
import Groups from "./components/Groups";
import Friends from "./components/Friends";

export default function Home() {
  return (
    <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <Groups />
      </div>

      {/* Center Content */}
      <div className={Styles.centerContent}>
        <div>search</div>
        <div>postuser</div>
        <div>post from group</div>
        <div>postuser</div>
        <div>post from group</div>
        <div>postuser</div>
        <div>post from group</div>
      </div>

      {/* Right Sidebar */}
      <div className={Styles.thirdSide}>
        <Friends/>
      </div>
    </div>
  );
}
