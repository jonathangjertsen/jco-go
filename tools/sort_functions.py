"""
Script to sort all the functions in a .go file.
Run `python sort_functions.py all` to see the generated output, then run with --mode=in-place to commit the changes to disk.
Make sure to have the file checked into git before running the script, so that you can easily revert it with git reset --hard HEAD.

ASSUMPTIONS:
- File has been go fmt'd (with --mode=in-place, this is done automatically)
- Multiline docstring not supported
- After the first function definition, the rest of the file is function definitions ONLY
"""
from argparse import ArgumentParser
from pathlib import Path
import subprocess

class SorryMultilineDocStringNotSupported(ValueError):
    pass

def get_preamble(lines) -> (int, str):
    """
    Scan to start of first function
    Returns index to start from, and the preamble as a string
    """
    entry_point = 0
    for i, line in enumerate(lines):
        if line.startswith("//"):
            if entry_point == 0:
                entry_point = i
        elif line.startswith("func"):
            if entry_point == 0:
                entry_point = i
            break
        else:
            entry_point = 0
    return entry_point, "".join(lines[:entry_point])

def extract_functions(lines) -> (str, list[str]):
    """
    Returns the preamble and a list of functions
    """
    entry_point, preamble = get_preamble(lines)
    functions = []
    current_function = []
    state = "root"

    def discard_function(require_empty=False):
        """Drops the currently buffered lines"""
        nonlocal current_function
        if require_empty:
            assert "".join(current_function).strip() == "", current_function
        current_function = []

    def set_state(new_state):
        """Updates the state machine, with special case for returning to root"""
        nonlocal state
        state = new_state
        if state == "root":
            push_function()

    def push_function():
        """Adds the current function to the list and clears it"""
        if current_function:
            functions.append("".join(current_function))
        discard_function()

    def add_line(line, new_state=None):
        """Adds a new line to the current function"""
        current_function.append(line)
        if new_state is not None:
            set_state(new_state)

    def sm_root(line):
        """Handler for the 'root' state"""
        if line.startswith("//"):
            add_line(line, "docstring")
        elif line.startswith("/*"):
            raise SorryMultilineDocStringNotSupported
        elif line.startswith("func"):
            add_line(line, "func")
        elif line.strip() == "":
            discard_function(require_empty=True)

    def sm_docstring(line):
        """Handler for the 'docstring' state"""
        if line.startswith("//"):
            add_line(line)
        elif line.startswith("func"):
            add_line(line, "func")

    def sm_func(line):
        """Handler for the 'func' state"""
        add_line(line, "root" if line.rstrip() == "}" else "func")

    # Create and execute state machine
    state_machine = {
        "root": sm_root,
        "docstring": sm_docstring,
        "func": sm_func,
    }
    for i, line in enumerate(lines[entry_point:]):
        try:
            state_machine[state](line)
        except Exception as exc:
            print(f"Exception on line #{i}")
            raise

    return preamble, functions

def function_name_and_receiver_type(function: str) -> str:
    """
    Returns the name of the function
    """
    for line in function.split("\n"):
        if line.startswith("func"):
            after_func = line.removeprefix("func ")
            only_whitespace = True
            get_end_of_receiver = False
            name_start_idx = 0
            name_end_idx = 0
            receiver_start_idx = 0
            receiver_end_idx = 0
            for i, char in enumerate(after_func):
                if get_end_of_receiver:
                    if char == ")":
                        receiver_end_idx = i
                        name_start_idx = i + 2
                        get_end_of_receiver = False
                else:
                    if char == "(":
                        if only_whitespace:
                            get_end_of_receiver = True
                            receiver_start_idx = i + 1
                        else:
                            name_end_idx = i
                            break
                    elif char != " ":
                        only_whitespace = False
            name = after_func[name_start_idx:name_end_idx]
            receiver = after_func[receiver_start_idx:receiver_end_idx]
            receiver_type = receiver.split(" ")[-1]
            return name, receiver_type
    raise ValueError(function)

def name_key(name: str) -> tuple[int]:
    if not name:
        return (0,)
    offset = 0
    first_letter = name[offset]
    while first_letter == "*":
        offset += 1
        first_letter = name[offset]
    if first_letter == '_':
        key = 1
    elif 'a' <= first_letter <= 'z':
        key = ord(first_letter) - ord('a') + 10
    elif 'A' <= first_letter <= 'Z':
        key = ord(first_letter) - ord('A') + 1000
    else:
        key = 1000000
    return (key, *name_key(name[1:]))

def function_key(function: str) -> tuple[int]:
    """
    Returns a key to ensure the following sort order:

    - anything starting with _
    - unexported functions (lowercase start)
    - exported functions (uppercase start)
    - functions with a receiver type
    """
    name, receiver_type = function_name_and_receiver_type(function)
    assert name, f"No name for function:\n\n{function}\n\n"
    return name_key(receiver_type) + name_key(name)


def recombine_functions(preamble: str, functions: list[str]) -> str:
    """
    Returns the new text given the preamble and functions
    """
    return "\n".join([preamble] + sorted(functions, key=function_key))

def process_file(args, file):
    if args.mode == "in-place":
        print(f"Start: {file}")
        subprocess.run(["go", "fmt", file])
    with open(file) as fp:
        new_text = recombine_functions(*extract_functions(fp.readlines()))
    if args.mode == "dry-run":
        print(new_text)
    elif args.mode == "in-place":
        with open(file, "w") as fp:
            fp.write(new_text)
    if args.mode == "in-place":
        subprocess.run(["go", "fmt", file])
        print(f"End: {file}")


def main():
    """
    Main script
    """
    parser = ArgumentParser()
    parser.add_argument("files", nargs="+")
    parser.add_argument("--mode", choices=("dry-run", "in-place", "quiet"), default="dry-run")
    args = parser.parse_args()

    if args.files == ["all"]:
        files = Path(".").glob("**/*.go")
    else:
        files = args.files

    for file in files:
        try:
            process_file(args, file)
        except Exception as exc:
            print(f"Exception in file {file}")
            raise

if __name__ == "__main__":
    main()
