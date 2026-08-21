#!/usr/bin/env python3
"""
install_groennesbyAI.py
~~~~~~~~~~~~~~~~~~~~~~
Bootstrap‑script dat in één stap een volledig werkende
Groennesby‑CLI (lokale AI‑assistant) aanmaakt, een virtuele
omgeving (venv) opzet, alle Python‑dependencies installeert,
een GGML‑model downloadt, en de CLI‑tool in de venv
registreren.

Gebruik:
    $ python3 install_groennesbyAI.py

Na afloop:
    $ source groennesby-cli/venv/bin/activate
    $ groennesby "maak een map demo en download een afbeelding"
"""

import os
import sys
import json
import subprocess
import urllib.request
import textwrap
from pathlib import Path

# ----------------------------------------------------------------------
# Helper‑functies
# ----------------------------------------------------------------------
def run(cmd: str, cwd: Path | None = None, check: bool = True) -> subprocess.CompletedProcess:
    """Run a shell command, stream its output, raise on error (if check=True)."""
    result = subprocess.run(
        cmd,
        cwd=str(cwd) if cwd else None,
        shell=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
    )
    if result.stdout:
        print(result.stdout)
    if check and result.returncode != 0:
        raise RuntimeError(f"Command failed (code {result.returncode}): {cmd}")
    return result

def write_file(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")
    print(f"✔️  {path}")

# ----------------------------------------------------------------------
# 1. Basis‑paden
# ----------------------------------------------------------------------
BASE_DIR   = Path.cwd() / "groennesby-cli"
PKG_DIR    = BASE_DIR / "groennesby"
VENV_DIR   = BASE_DIR / "venv"
MODEL_DIR  = Path.home() / ".groennesby" / "models"

# ----------------------------------------------------------------------
# 2. Schrijf de package‑bestanden (dezelfde code als in het vorige antwoord)
# ----------------------------------------------------------------------
write_file(PKG_DIR / "__init__.py", "")

# ---- cli.py -----------------------------------------------------------
cli_code = textwrap.dedent("""\
    import click
    from .model import GroennesbyLLM
    from .router import IntentRouter
    from .executor import SystemExecutor

    BANNER = r\"\"\"
       ____                       _
      / ___| ___ _ __   ___ _ __ | |_ ___  _ __ ___
     | |  _ / _ \\ '_ \\ / _ \\ '_ \\| __/ _ \\| '__/ _ \\
     | |_| |  __/ | | |  __/ | | | || (_) | | |  __/
      \\____|\\___|_| |_|\\___|_| |_|\\__\\___/|_|  \\___|

      COLLABS‑AI  –  “Let’s Collaborate”
    \"\"\"

    @click.command()
    @click.argument('prompt', nargs=-1)
    @click.option('--unsafe', is_flag=True, help='Sta onveilige shell‑commands toe')
    def main(prompt, unsafe):
        \"\"\"Start Groennesby‑CLI.
        Voorbeeld:
            $ groennesby "maak een map demo en download een afbeelding"
        \"\"\"
        click.echo(BANNER)

        user_input = " ".join(prompt).strip()
        if not user_input:
            user_input = click.prompt("Wat wil je doen?")

        llm = GroennesbyLLM()
        raw = llm.generate(user_input)

        intent = IntentRouter.parse(raw)

        executor = SystemExecutor(unsafe=unsafe)
        answer = executor.handle(intent)

        click.echo(answer)


    if __name__ == '__main__':
        main()
    """)
write_file(PKG_DIR / "cli.py", cli_code)

# ---- model.py ---------------------------------------------------------
model_code = textwrap.dedent("""\
    import os
    from pathlib import Path
    from llama_cpp import Llama

    class GroennesbyLLM:
        \"\"\"Singleton‑wrapper rond een GGML‑model (Phi‑3 mini als voorbeeld).\"\"\"
        _instance = None
        _model_path = Path.home() / ".groennesby" / "models" / "phi3-mini-ggml-q4_0.bin"

        def __new__(cls):
            if cls._instance is None:
                cls._instance = super().__new__(cls)
                cls._instance._load()
            return cls._instance

        def _load(self):
            if not self._model_path.is_file():
                raise FileNotFoundError(
                    f\"Model not found at {self._model_path}.\\n"
                    \"Run `groennesby-download-model` (automatisch bij setup) first.\"
                )
            self.llama = Llama(
                model_path=str(self._model_path),
                n_ctx=4096,
                n_batch=512,
                temperature=0.2,
                repeat_penalty=1.1,
                stop=[\"\\n\\n\", \"</s>\"],
            )

        def generate(self, user_text: str) -> str:
            \"\"\"Return a JSON string describing the intent.\"\"\"
            prompt = (
                "You are a helpful CLI assistant. "
                "Return ONLY a compact JSON describing the user's request. "
                "Do NOT add any explanations. "
                "If you cannot understand, return {\\\"intent\\\": \\\"unknown\\\"}.\\n"
                f"User: {user_text}\\nAssistant:"
            )
            out = self.llama(prompt, max_tokens=256, stop=[\"</s>\", \"\\n\\n\"])
            return out[\"choices\"][0][\"text\"].strip()
    """)
write_file(PKG_DIR / "model.py", model_code)

# ---- router.py --------------------------------------------------------
router_code = textwrap.dedent("""\
    import json
    from dataclasses import dataclass
    from typing import Any, Dict

    @dataclass
    class Intent:
        name: str
        args: Dict[str, Any]

    class IntentRouter:
        @staticmethod
        def parse(raw_json: str) -> Intent:
            \"\"\"Convert the model’s JSON output into an Intent object.\"\"\"
            try:
                data = json.loads(raw_json)
            except json.JSONDecodeError:
                return Intent(name="unknown", args={})
            return Intent(name=data.get("intent", "unknown"), args=data)
    """)
write_file(PKG_DIR / "router.py", router_code)

# ---- executor.py ------------------------------------------------------
executor_code = textwrap.dedent("""\
    import json
    import subprocess
    from pathlib import Path
    from .router import Intent

    # whitelist van veilige commando’s (naam → builder‑functie)
    SAFE_CMDS = {
        "list_dir":   lambda path="." : ["ls", "-la", path],
        "make_dir":   lambda path   : ["mkdir", "-p", path],
        "download":   lambda url, target_path: ["curl", "-L", url, "-o", target_path],
        "create_file":lambda path, content: ["sh", "-c", f'printf %s "{content}" > "{path}"'],
        "search_web": None,     # handled by separate module
        "run_shell": None,      # alleen met --unsafe
    }

    class SystemExecutor:
        def __init__(self, unsafe: bool = False):
            self.unsafe = unsafe
            self.state_path = Path.home() / ".groennesby" / "state.json"
            self.state = {"cwd": str(Path.cwd())}
            if self.state_path.is_file():
                try:
                    self.state.update(json.loads(self.state_path.read_text()))
                except Exception:
                    pass

        # ------------------------------------------------------------------
        def _save_state(self):
            self.state_path.parent.mkdir(parents=True, exist_ok=True)
            self.state_path.write_text(json.dumps(self.state))

        # ------------------------------------------------------------------
        def handle(self, intent: Intent) -> str:
            name, args = intent.name, intent.args

            if name == "unknown":
                return "❓ Ik begrijp niet wat je bedoelt."

            if name == "change_dir":
                target = Path(args.get("path", ".")).expanduser()
                if not target.is_dir():
                    return f"❌ Pad bestaat niet: {target}"
                self.state["cwd"] = str(target)
                self._save_state()
                return f"📁 Werkomgeving veranderd naar {target}"

            if name == "search_web":
                from .internet import WebSearcher
                return WebSearcher.search(args.get("query", ""))

            if name == "run_shell":
                if not self.unsafe:
                    return "⚠️  run_shell is uitgeschakeld. Gebruik --unsafe om het toe te staan."
                return self._run_cmd(args.get("cmd"))

            # normale whitelist‑commando's
            builder = SAFE_CMDS.get(name)
            if builder is None:
                return f"❌ Intent '{name}' wordt niet ondersteund."

            try:
                cmd = builder(**args) if isinstance(args, dict) else builder(*args)
            except Exception as e:
                return f"⚠️  Fout bij het bouwen van het commando: {e}"

            return self._run_cmd(cmd)

        # ------------------------------------------------------------------
        def _run_cmd(self, cmd):
            \"\"\"Execute a command list in the (eventual) cwd stored in state.\"\"\"
            try:
                result = subprocess.run(
                    cmd,
                    check=True,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True,
                    cwd=self.state.get("cwd"),
                )
                out = result.stdout.strip()
                return out if out else "✅ OK (geen output)"
            except subprocess.CalledProcessError as e:
                err = e.stderr.strip()
                return f"❌ fout: {err or 'exit status ' + str(e.returncode)}"
    """)
write_file(PKG_DIR / "executor.py", executor_code)

# ---- internet.py ------------------------------------------------------
internet_code = textwrap.dedent("""\
    import requests
    from bs4 import BeautifulSoup

    class WebSearcher:
        @staticmethod
        def search(query: str, max_results: int = 3) -> str:
            api = "https://api.duckduckgo.com/"
            params = {
                "q": query,
                "format": "json",
                "no_html": 1,
                "skip_disambig": 1,
            }
            try:
                r = requests.get(api, params=params, timeout=8,
                                 headers={"User-Agent": "Groennesby/0.1"})
                data = r.json()
            except Exception as exc:
                return f"❌ zoek‑fout: {exc}"

            abstract = data.get("AbstractText")
            if abstract:
                return f"📰 Samenvatting: {abstract}"

            results = data.get("Results", [])
            if not results:
                return "🔍 Geen resultaten gevonden."

            first = results[0]
            url = first.get("FirstURL", "")
            try:
                page = requests.get(url, timeout=8,
                                    headers={"User-Agent": "Groennesby/0.1"}).text
                soup = BeautifulSoup(page, "html.parser")
                txt = " ".join(soup.stripped_strings)[:500]
            except Exception:
                txt = "(kon geen HTML ophalen)"
            return f"🔗 {url}\\n\\n{txt}…"
    """)
write_file(PKG_DIR / "internet.py", internet_code)

# ----------------------------------------------------------------------
# 3. pyproject.toml (PEP‑621)
# ----------------------------------------------------------------------
pyproject = textwrap.dedent("""\
    [project]
    name = "groennesby-cli"
    version = "0.1.0"
    description = "Local AI assistant for macOS terminal – token‑free"
    readme = "README.md"
    requires-python = ">=3.10"
    license = {text = "Apache-2.0"}
    authors = [{name="JouwNaam", email="jouw@email.com"}]
    keywords = ["ai", "cli", "local-llm", "assistant"]

    dependencies = [
        "click>=8.1",
        "llama-cpp-python>=0.2.28",
        "requests>=2.31",
        "beautifulsoup4>=4.12",
    ]

    [project.scripts]
    groennesby = "groennesby.cli:main"
    """)
write_file(BASE_DIR / "pyproject.toml", pyproject)

# ----------------------------------------------------------------------
# 4. README.md (kort)
# ----------------------------------------------------------------------
readme = textwrap.dedent("""\
    # Groennesby‑CLI

    **COLLABS‑AI – “Let’s Collaborate”**

    Een volledig lokale, token‑vrije AI‑assistant die je via de terminal kunt oproepen:

    ```bash
    groennesby "maak een map demo en download een afbeelding"
    ```

    ## Installatie (een keer)

    1. Clone / download deze repo.
    2. Run the bootstrap script:

       ```bash
       python3 install_groennesbyAI.py
       ```

    3. Activate the virtual environment:

       ```bash
       source groennesby-cli/venv/bin/activate
       ```

    4. Gebruik de CLI:

       ```bash
       groennesby "help"
       ```

    ## Licentie

    Apache‑2.0
    """)
write_file(BASE_DIR / "README.md", readme)

# ----------------------------------------------------------------------
# 5. LICENSE (Apache‑2.0 – verkort)
# ----------------------------------------------------------------------
license_txt = textwrap.dedent("""\
    Apache License
    Version 2.0, January 2004
    http://www.apache.org/licenses/
    
    TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
    
    ... (verkort voor voorbeeld) ...
    """)
write_file(BASE_DIR / "LICENSE", license_txt)

# ----------------------------------------------------------------------
# 6. .gitignore (handmatig toevoegen, maar we maken er één)
# ----------------------------------------------------------------------
gitignore = textwrap.dedent("""\
    # Python artefacts
    __pycache__/
    *.py[cod]
    *.pyo
    *.pyd
    *.egg-info/
    *.dist-info/
    build/
    develop-eggs/
    .pytest_cache/
    .mypy_cache/
    .coverage
    htmlcov/
    .ruff_cache/
    
    # Virtual environment
    venv/
    
    # Model files (they live in ~/.groennesby)
    .groennesby/models/
    """)
write_file(BASE_DIR / ".gitignore", gitignore)

# ----------------------------------------------------------------------
# 7. Maak virtuele omgeving en installeer dependencies
# ----------------------------------------------------------------------
print("\n=== Stap 1: virtualenv aanmaken ===")
if not VENV_DIR.is_dir():
    run(f"python3 -m venv \"{VENV_DIR}\"")
else:
    print("✅ venv bestaat al")

PIP_BIN    = VENV_DIR / "bin" / "pip"
PYTHON_BIN = VENV_DIR / "bin" / "python"

print("\n=== Stap 2: pip, setuptools, wheel updaten ===")
run(f"\"{PIP_BIN}\" install --upgrade pip setuptools wheel")

print("\n=== Stap 3: package (editable) + dependencies installeren ===")
run(f"\"{PIP_BIN}\" install -e \"{BASE_DIR}\"")

# ----------------------------------------------------------------------
# 8. Model downloaden (Phi‑3 mini – 4‑bit ggml)
# ----------------------------------------------------------------------
MODEL_URL  = ("https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-ggml"
               "/resolve/main/phi3-mini-ggml-q4_0.bin")
MODEL_PATH = MODEL_DIR / "phi3-mini-ggml-q4_0.bin"

print("\n=== Stap 4: Model downloaden ===")
MODEL_DIR.mkdir(parents=True, exist_ok=True)

if MODEL_PATH.is_file():
    print(f"✅ Model bestaat al: {MODEL_PATH}")
else:
    print(f"⬇️  Downloaden van {MODEL_URL}")
    try:
        urllib.request.urlretrieve(MODEL_URL, MODEL_PATH)
        print("✅ Download voltooid")
    except Exception as exc:
        print(f"❌ Fout bij downloaden: {exc}")
        sys.exit(1)

# ----------------------------------------------------------------------
# 9. Eindbericht
# ----------------------------------------------------------------------
print("\n=== 🎉 Alles klaar! 🎉 ===")
print(f"- Project directory : {BASE_DIR}")
print(f"- Virtualenv       : {VENV_DIR}")
print("- Activeer de venv met:")
print(f"    source {VENV_DIR}/bin/activate")
print("- Voorbeeld‑gebruik:")
print('    groennesby "maak een map demo en download een kattenfoto"')
print("\nCommit & push je repo naar GitHub en deel het met de wereld!")
