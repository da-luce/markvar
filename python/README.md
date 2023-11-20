# Markvar-Python

## Overview

Markvar-Python is a Python package designed to facilitate the management of text
variables within project environments in conjunction with the
[Markvar command line tool](https://github.com/da-luce/markvar). It allows users
to save key-value pairs to a .markvar file, which can then be utilized to
automatically update text and Markdown documents with script-generated text.

## Installation

Install Markvar-python using pip:

```bash
pip install markvar
```

## Usage

Adding a variable:

```python
from markvar import markvar

markvar.mark('variable_name', 'Some content for the variable')
```

Now, this variable is available to use within markdown documents:

```markdown
Text... {{var:variable_name}} ...
```

Use the Markvar command line tool to replace this variable with your desired
text. More at: [https://github.com/da-luce/markvar](https://github.com/da-luce/markvar).
