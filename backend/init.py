import os
import sys
from datetime import datetime, timedelta

def write_env_file(env_path):
    if os.path.exists(env_path):
        return
    with open(env_path, "w") as f:
        f.write("JWT_SECRET_KEY=#SocialNetwork@zone01!\n")

def main():
    script_dir = os.path.dirname(os.path.abspath(sys.argv[0]))
    private_dir = os.path.abspath(os.path.join(script_dir, "..", "private"))
    os.makedirs(private_dir, exist_ok=True)

    env_path = os.path.join(private_dir, ".env")
    write_env_file(env_path)

if __name__ == "__main__":
    main()
