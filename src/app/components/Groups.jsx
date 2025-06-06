import Image from "next/image";
import Link from "next/link";
import Styles from "../global.module.css"

export default function Groups() {
  return (
    <div className={Styles.groups}>
      <h1>Groups</h1>
      {[1, 2, 3, 4, 5].map((_, i) => (
        <div key={i}>
          <div>
            <Image src="/iconMale.png" alt="profile" width={40} height={40} />
            <h5>Group Name</h5>
          </div>
          <Link href="/join">Join</Link>
        </div>
      ))}
    </div>
  );
}
