import Styles from "./global.module.css";
import Groups from "./components/Groups";
import Friends from "./components/Friends";
import Image from 'next/image';

export default function Home() {
  return (
    <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <Groups />
      </div>

      {/* Center Content */}
      <div className={Styles.centerContent}>
        <div className={Styles.post}>
          <section className={Styles.userinfo}>
            <div>
              <Image src="/iconMale.png" alt="notification" width={25} height={25} />
              <p>userName</p>
            </div>
            <div>
              <p>time/time/time</p>
            </div>
          </section>

          <section className={Styles.content}>
            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ratione, ipsum!
          </section>

          <section className={Styles.footer}>
            <div className={Styles.action}>
              <Image src="/Like.svg" alt="like" width={20} height={20} />
              <p>10</p>
            </div>
            <div className={Styles.action}>
              <Image src="/comment.svg" alt="comment" width={20} height={20} />
              <p>5</p>
            </div>
          </section>
        </div>

        <div className={Styles.post}>
          <section className={Styles.userinfo}>
            <div>
              <p>groupeName @</p>
              <Image src="/iconMale.png" alt="notification" width={25} height={25} />
              <p>userName</p>
            </div>
            <div>
              <p>time/time/time</p>
            </div>
          </section>

          <section className={Styles.content}>
            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ratione, ipsum!
          </section>

          <section className={Styles.footer}>
            <div className={Styles.action}>
              <Image src="/Like2.svg" alt="like" width={20} height={20} />
              <p>10</p>
            </div>
            <div className={Styles.action}>
              <Image src="/comment.svg" alt="comment" width={20} height={20} />
              <p>5</p>
            </div>
          </section>
        </div>

        <div className={Styles.post}>
          <section className={Styles.userinfo}>
            <div>
              <Image src="/iconMale.png" alt="notification" width={25} height={25} />
              <p>userName</p>
            </div>
            <div>
              <p>time/time/time</p>
            </div>
          </section>

          <section className={Styles.content}>
            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ratione, ipsum!
          </section>

          <section className={Styles.footer}>
            <div className={Styles.action}>
              <Image src="/Like.svg" alt="like" width={20} height={20} />
              <p>10</p>
            </div>
            <div className={Styles.action}>
              <Image src="/comment.svg" alt="comment" width={20} height={20} />
              <p>5</p>
            </div>
          </section>
        </div>

        <div className={Styles.post}>
          <section className={Styles.userinfo}>
            <div>
              <p>groupeName @</p>
              <Image src="/iconMale.png" alt="notification" width={25} height={25} />
              <p>userName</p>
            </div>
            <div>
              <p>time/time/time</p>
            </div>
          </section>

          <section className={Styles.content}>
            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ratione, ipsum!
          </section>

          <section className={Styles.footer}>
            <div className={Styles.action}>
              <Image src="/Like.svg" alt="like" width={20} height={20} />
              <p>10</p>
            </div>
            <div className={Styles.action}>
              <Image src="/comment.svg" alt="comment" width={20} height={20} />
              <p>5</p>
            </div>
          </section>
        </div>

        <div className={Styles.post}>
          <section className={Styles.userinfo}>
            <div>
              <Image src="/iconMale.png" alt="notification" width={25} height={25} />
              <p>userName</p>
            </div>
            <div>
              <p>time/time/time</p>
            </div>
          </section>

          <section className={Styles.content}>
            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ratione, ipsum!
          </section>

          <section className={Styles.footer}>
            <div className={Styles.action}>
              <Image src="/Like2.svg" alt="like" width={20} height={20} />
              <p>10</p>
            </div>
            <div className={Styles.action}>
              <Image src="/comment.svg" alt="comment" width={20} height={20} />
              <p>5</p>
            </div>
          </section>
        </div>

        <div className={Styles.post}>
          <section className={Styles.userinfo}>
            <div>
              <p>groupeName @</p>
              <Image src="/iconMale.png" alt="notification" width={25} height={25} />
              <p>userName</p>
            </div>
            <div>
              <p>time/time/time</p>
            </div>
          </section>

          <section className={Styles.content}>
            Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ratione, ipsum!
          </section>

          <section className={Styles.footer}>
            <div className={Styles.action}>
              <Image src="/Like.svg" alt="like" width={20} height={20} />
              <p>10</p>
            </div>
            <div className={Styles.action}>
              <Image src="/comment.svg" alt="comment" width={20} height={20} />
              <p>5</p>
            </div>
          </section>
        </div>
      </div>

      {/* Right Sidebar */}
      <div className={Styles.thirdSide}>
        <Friends />
      </div>
    </div>
  );
}

