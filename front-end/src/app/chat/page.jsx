"use client";
import Styles from "../global.module.css";
import LStyle from "./style.module.css";
import Groups from "../components/Groups";
import Friends from "../components/Friends";
import Image from 'next/image';
import Profiles from "./profiles";
import Discussion from "./discussion";
import Input from "./input";
console.log("Chat page loaded");
export default function Chat() {
    return (
        <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <Profiles />
      </div>

      {/* Center Content */}
      <div className={[Styles.centerContent , LStyle.chat_container].join(" ")}>
        <Discussion />
        <Input />
      </div>

      {/* Right Sidebar */}
      <div className={Styles.thirdSide}>
        <Friends />
        {/* TODO: Decide */}
      </div>
    </div>
    )
}