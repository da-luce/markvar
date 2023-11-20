import json
import os
import fcntl

mark_var_path = None
error = False


def _create_markvar_file(root_dir: str) -> None:
    # Using 'with' ensures the file is properly closed
    with open(os.path.join(root_dir, ".markvar"), "a") as f:
        pass


def _find_mark_file(root_dir: str) -> str:
    global mark_var_path, error

    result = []
    name = ".markvar"
    for root, dirs, files in os.walk(root_dir):
        if name in files:
            result.append(os.path.join(root, name))

    if len(result) > 1:
        print("Warning: found more than one .markvar file. Please delete all but one")
        error = True
        return ""

    if len(result) == 0:
        print("Did not find a .markvar file. Creating " + os.path.join(root_dir, name))
        _create_markvar_file(root_dir)
        return os.path.join(root_dir, name)

    mark_var_path = result[0]
    return mark_var_path


def mark(var: str, data: str, root_dir: str) -> None:
    global mark_var_path, error

    if error:
        return

    if mark_var_path is None:
        mark_var_path = _find_mark_file(root_dir)
        if not mark_var_path:
            error = True
            print("Error: Unable to locate or create .markvar file.")
            return

    try:
        with open(mark_var_path, "r+") as f:
            # Attempt to acquire a lock
            fcntl.flock(f, fcntl.LOCK_EX)

            try:
                markvar_data = json.load(f)
            except json.JSONDecodeError:
                markvar_data = {}

            markvar_data[var] = data
            f.seek(0)  # Resets file pointer to the beginning
            f.truncate()  # Clears the file
            json.dump(markvar_data, f, indent=4)

            # Release the lock
            fcntl.flock(f, fcntl.LOCK_UN)

    except IOError as e:
        print(f"Error: Unable to write to .markvar file - {e}")
