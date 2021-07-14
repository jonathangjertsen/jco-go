"""
Generates release binaries
"""
from pathlib import Path
import shutil

HERE = Path(__file__).parent.absolute()
BIN_DIR = HERE.parent / "bin"
RELEASE_DIR = HERE.parent / "release"

if __name__ == "__main__":
    RELEASE_DIR.mkdir(parents=True, exist_ok=True)
    for path in BIN_DIR.glob("**/*"):
        if path.is_dir():
            continue
        dirname = path.parent.name
        archive_name = shutil.make_archive(
            base_name=dirname,
            format="gztar",
            root_dir=path.parent,
        )
        dest = RELEASE_DIR / Path(archive_name).name
        shutil.move(archive_name, dest)
        print(f"Path: {dest}, size: {dest.stat().st_size/(1000000):.1f}MB")
