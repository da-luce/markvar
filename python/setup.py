from setuptools import setup, find_packages

setup(
    name="markvar",
    version="0.1",
    packages=find_packages(),
    description="Python package to be used in conjunction with the Markvar tool",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/da-luce/markvar",
    author="Dalton Luce",
    author_email="your.email@example.com",
    license="MIT",
    install_requires=[
        # Any dependencies, if required
    ],
    zip_safe=False,
)
