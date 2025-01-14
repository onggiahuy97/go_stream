#!/usr/bin/env python3
import os
import sys

def print_tree(start_path, prefix=""):
    """
    Recursively prints the directory structure starting from `start_path`
    using a tree-like format.
    """
    # Get all entries in the directory, excluding hidden files/folders
    entries = [entry for entry in os.listdir(start_path) if not entry.startswith('.')]
    entries.sort()

    # Iterate over each entry to print them with the correct tree branch symbols
    for index, entry in enumerate(entries):
        # Determine the connector for the current item
        connector = "└── " if index == len(entries) - 1 else "├── "

        # Build the full path of the current item
        path = os.path.join(start_path, entry)

        # Print the current item
        print(prefix + connector + entry)

        # If the current item is a directory, recurse into it
        if os.path.isdir(path):
            # For directories, extend the prefix properly
            extension = "    " if index == len(entries) - 1 else "│   "
            print_tree(path, prefix + extension)

def main():
    # If called like './print_tree <target_dir>', use that directory; otherwise, use '.'
    start_dir = sys.argv[1] if len(sys.argv) > 1 else "."

    print("🥶 Project structure:")
    print("👉 " + os.path.abspath(start_dir))
    print_tree(start_dir)

if __name__ == "__main__":
    main()


