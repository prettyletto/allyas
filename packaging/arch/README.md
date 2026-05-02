# Arch Packaging

This directory is for local package testing and a future AUR submission. It is not a deployment step.

Build and install locally with:

```sh
makepkg -si
```

The package only installs the `allyas` binary. It does not create user config, does not run `allyas init`, and does not edit shell rc files.

After installing the package, each user still runs:

```sh
allyas init
allyas install --shell zsh --manual
```

or:

```sh
allyas install --shell zsh --auto
```
