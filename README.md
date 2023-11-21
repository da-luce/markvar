# Markvar

## Introduction

Markvar is a versatile tool designed to integrate script-generated data into
Markdown and other text documents. It provides both a command-line interface and
a Python package, enabling seamless updates of documents with dynamically
generated content

## Features

* Automated Data Insertion: Effortlessly updates Markdown documents with fresh
  data from scripts.
* Customizable Placeholders: Uses HTML comments in Markdown or custom
  placeholders in other text formats for easy data mapping.

## Installation

### Command Line Tool

Requires Go (1.15 or higher)

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/Markvar.git
   ```

2. Navigate to the cloned directory:

   ```bash
   cd Markvar
   ```

3. Build the application using Go:

   ```bash
   go build
   ```

### Python Package

   ```bash
   pip install markvar
   ```

## Usage

Prepare your Markdown file or text with ID-based placeholders for data insertion:

```markdown
My Data:

| A           | B           |
|-------------|-------------|
| {{var:A1}}  | {{var:B1}}  |
| {{var:A2}}  | {{var:B2}}  |

The mean of column A is {{var:A_mean}}, while the mean of column B is
{{var:B_mean}}
```

Use the Python (or applicable package to generate your data):

```python
from markvar import mark

A = numpy.random.randint(low=0, high=100, size=2)
B = numpy.random.randint(low=0, high=100, size=2)

# Record variables to markvar environment
mark("A1", str(A[0]), "")
mark("A2", str(A[1]), "")
mark("B1", str(B[0]), "")
mark("B2", str(B[1]), "")

A_mean = np.mean(A)
B_mean = np.mean(B)

# Record variables to markvar environment
mark("A_mean", str(A_mean), "")
mark("B_mean", str(B_mean), "")
```

Now, use the command line tool to replace the placeholders in a new copy of
markdown file:

```bash
./markvar update -file ./example.md
```

This will generate a new file `example.md.updated` with the replaced bindings.
