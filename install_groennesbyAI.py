import os
import subprocess
import sys
from pathlib import Path

HOME = str(Path.home())
BASE_DIR = f"{HOME}/Documents/Groennesby-AI"
SCRIPTS_DIR = f"{BASE_DIR}/scripts"

def banner():
    print("===========================================")
    print("      Installing GroennesbyAI CLI")
    print("===========================================")
    print("Created by: Dave Groennesby")
    print("A self‑made hybrid AI CLI system")
    print("-------------------------------------------")
    print("Starting installation...")
    print("")

def create_directories():
    os.makedirs(SCRIPTS_DIR, exist_ok=True)
    print(f"[OK] Created directories at {BASE_DIR}")

def write_agent_script():
    agent_code = '''
import os
import sys

def main():
    print("GroennesbyAI Agent ready.")
    while True:
        user = input(">>> ")

        if user.startswith("/help"):
            print("Available commands:")
            print("/newproject NAME")
            print("/exit")

        elif user.startswith("/newproject"):
            parts = user.split(" ", 1)
            if len(parts) < 2:
                print("Usage: /newproject NAME")
                continue
            name = parts[1]
            os.makedirs(name, exist_ok=True)
            print(f"Project {name} created.")

        elif user.startswith("/exit"):
            print("Goodbye Dave.")
            sys.exit()

        else:
            print("Command not recognized. Type /help.")

if __name__ == "__main__":
    main()
'''
    with open(f"{SCRIPTS_DIR}/agent.py", "w") as f:
        f.write(agent_code)

    print("[OK] Agent script installed")

def create_terminal_command():
    script = f'''#!/bin/zsh

cd "{BASE_DIR}"

echo "==========================================="
echo "        GroennesbyAI — Welcome"
echo "==========================================="
echo ""
echo "Your personal AI CLI is now active."
echo "Type /help for commands."
echo "-------------------------------------------"
echo ""

python3 scripts/agent.py
'''

    with open("/usr/local/bin/groennesby", "w") as f:
        f.write(script)

    subprocess.run(["sudo", "chmod", "+x", "/usr/local/bin/groennesby"])
    print("[OK] Terminal command 'groennesby' installed")

def finish():
    print("")
    print("===========================================")
    print(" GroennesbyAI installation complete")
    print("===========================================")
    print("Start your AI with:")
    print("   groennesby")
    print("-------------------------------------------")
    print("Enjoy your self‑made AI, Dave.")
    print("")

if __name__ == "__main__":
    banner()
    create_directories()
    write_agent_script()
    create_terminal_command()
    finish()