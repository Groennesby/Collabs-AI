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
