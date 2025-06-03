import os
import sys
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.backends import default_backend
import datetime

def generate_self_signed_cert(cert_path, key_path):
    # Generate private key
    key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend()
    )

    # Write private key to file
    with open(key_path, "wb") as f:
        print(f.write(key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption()
            )))

    # Build certificate subject and issuer (self-signed)
    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, u"MA"),  # Morocco country code
        x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, u"Oujda Orientale"),  # Your region
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, u"Zone01Oujda talen"),  # Your org or zone name
        x509.NameAttribute(NameOID.COMMON_NAME, u"localhost"),
    ])

    cert = x509.CertificateBuilder().subject_name(
        subject
    ).issuer_name(
        issuer
    ).public_key(
        key.public_key()
    ).serial_number(
        x509.random_serial_number()
    ).not_valid_before(
        datetime.datetime.utcnow()
    ).not_valid_after(
        # Certificate valid for 1 year
        datetime.datetime.utcnow() + datetime.timedelta(days=365)
    ).add_extension(
        x509.SubjectAlternativeName([x509.DNSName(u"localhost")]),
        critical=False,
    ).sign(key, hashes.SHA256(), default_backend())

    # Write cert to file
    with open(cert_path, "wb") as f:
        print(f.write(cert.public_bytes(serialization.Encoding.PEM)))

    print(f"✅ Generated TLS cert: {cert_path}")
    print(f"✅ Generated TLS key: {key_path}")

def write_env_file(env_path):
    content = "JWT_SECRET_KEY=#SotialNetwork@zone01!\n"
    with open(env_path, "w") as f:
        f.write(content)
    print(f"✅ Created env file: {env_path}")

def main():
    script_dir = os.path.dirname(os.path.abspath(sys.argv[0]))
    private_dir = os.path.abspath(os.path.join(script_dir,"server", "private"))

    # Ensure private dir exists
    os.makedirs(private_dir, exist_ok=True)

    cert_path = os.path.join(private_dir, "cert.pem")
    key_path = os.path.join(private_dir, "key.pem")
    env_path = os.path.join(private_dir, ".env")

    generate_self_signed_cert(cert_path, key_path)
    write_env_file(env_path)

if __name__ == "__main__":
    main()
