import Link from 'next/link';
import Image from 'next/image';


export default function NotFound() {
    return (
        <div style={{ padding: '4rem', textAlign: 'center' }}>
            <Image src="/404.svg" alt="404" width={500} height={500}/>
            <h1>404 - Page Not Found</h1>
            <p>Sorry, we couldn’t find the page you’re looking for.</p>
            <Link href="/" style={{ color: '#1C1665' }}>
                <p >Go back home</p>
            </Link>
        </div>
    );
}
