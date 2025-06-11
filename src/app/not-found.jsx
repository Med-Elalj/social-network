import Link from 'next/link';
import Image from 'next/image';


export default function NotFound() {
    return (
        <div style={{ padding: '4rem', textAlign: 'center', marginTop: '60px' }}>
            <Image src="/404.svg" alt="404" width={500} height={500} />
            <h1 style={{ color: '#8D6B0D' }}>404 - Page Not Found</h1>
            <p style={{ color: '#e0e0e0' }}>Sorry, we couldn’t find the page you’re looking for.</p>
            <Link href="/" style={{ color: '#8D6B0D' }}>
                <p >Go back home</p>
            </Link>
        </div>
    );
}
