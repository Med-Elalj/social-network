import os
import sys
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.backends import default_backend
from datetime import datetime, timedelta

def generate_self_signed_cert(cert_path, key_path):
    if os.path.exists(cert_path) and os.path.exists(key_path):
        return

    key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend()
    )

    with open(key_path, "wb") as f:
        f.write(key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption()
        ))

    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, "MA"),
        x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, "Oujda Orientale"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, "Zone01Oujda"),
        x509.NameAttribute(NameOID.COMMON_NAME, "localhost"),
    ])

    cert = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(issuer)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(datetime.utcnow())
        .not_valid_after(datetime.utcnow() + timedelta(days=365))
        .add_extension(
            x509.SubjectAlternativeName([
                x509.DNSName("localhost"),
                # You can add more names or IPs here if needed
            ]),
            critical=False,
        )
        .sign(key, hashes.SHA256(), default_backend())
    )

    with open(cert_path, "wb") as f:
        f.write(cert.public_bytes(serialization.Encoding.PEM))

def write_env_file(env_path):
    if os.path.exists(env_path):
        return
    with open(env_path, "w") as f:
        f.write("JWT_SECRET_KEY=#SocialNetwork@zone01!\n")

def main():
    script_dir = os.path.dirname(os.path.abspath(sys.argv[0]))
    private_dir = os.path.abspath(os.path.join(script_dir, "..", "private"))
    os.makedirs(private_dir, exist_ok=True)

    cert_path = os.path.join(private_dir, "cert.pem")
    key_path = os.path.join(private_dir, "key.pem")
    env_path = os.path.join(private_dir, ".env")

    generate_self_signed_cert(cert_path, key_path)
    write_env_file(env_path)

if __name__ == "__main__":
    main()
