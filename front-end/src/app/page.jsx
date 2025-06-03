import Image from "next/image";
import Styles from "./global.module.css";
import Link from "next/link";

export default function Home() {
  return (
    <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <div>
          <h1>Groups</h1>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Group Name</h5>
            <Link href="/join">Join</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Group Name</h5>
            <Link href="/join">Join</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Group Name</h5>
            <Link href="/join">Join</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Group Name</h5>
            <Link href="/join">Join</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Group Name</h5>
            <Link href="/join">Join</Link>
          </div>
        </div>

        <div>
          <h1>Friends</h1>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>User Name</h5>
            <p>Start of last message ...</p>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>User Name</h5>
            <p>Start of last message ...</p>
          </div>
        </div>
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
        <div>
          <h1>Requests</h1>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Username</h5>
            <Link href="/join">Reject</Link>
            <Link href="/join">Accept</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Username</h5>
            <Link href="/join">Reject</Link>
            <Link href="/join">Accept</Link>
          </div>
        </div>

        <div>
          <h1>Suggestions</h1>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Username</h5>
            <Link href="/join">Add</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Username</h5>
            <Link href="/join">Add</Link>
          </div>
          <div>
            <Image src="/login2.svg" alt="profile" width={20} height={20} />
            <h5>Username</h5>
            <Link href="/join">Add</Link>
          </div>
        </div>
      </div>
    </div>
  );
}
