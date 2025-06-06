// app/not-found.js
import Link from 'next/link';

export default function NotFound() {
    return (
        <div style={{ padding: '4rem', textAlign: 'center' }}>
            <h1>404 - Page Not Found</h1>
            <p>Sorry, we couldn’t find the page you’re looking for.</p>
            <Link href="/" style={{ color: '#1C1665' }}>
                <p >Go back home</p>
            </Link>
        </div>
    );
}
