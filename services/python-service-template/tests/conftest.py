import sys
from pathlib import Path

# Ensure app package is importable when running from repo root.
sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
