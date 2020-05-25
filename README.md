# dxhd

## daky's X11 hotkey daemon

dxhd is heavily inspired by [sxhkd](https://github.com/baskerville/sxhkd), written in Go, and has quite elegant syntax for configuration files!

## Installation

* Arch User Repository

```sh
git clone https://aur.archlinux.org/dxhd-git.git
cd dxhd-git
makepkg -si
```

or use an AUR helper like yay - `yay -S dxhd-git`

* From the source

```sh
git clone https://github.com/dakyskye/dxhd.git
cd dxhd
go build -o dxhd .
```

and copy `dxhd` executable file to somewhere in your `$PATH`

* From releases

Download the `dxhd` executable file from the latest release, from [releases page](https://github.com/dakyskye/dxhd/releases)

and copy `dxhd` executable file to somewhere in your `$PATH`

## Features

* basic keybindings (all presed together)
* variants in keybindings and their actions ({a,b,c})
* ranges in keybindings and their actions ({1-9} in a keybinding, {11-19} or whatever range of `9-1` in the action)
* support for any shell scripting language (sh, bash, ksh, zsh, python, perl etc.) given as a shebang
* support for scripting, as much as a user wishes!

## Configuration

The default config file is located at `~/.config/dxhd/dxhd.sh`, however, dxhd can use a file from any path, by passing it to `-c`:

```sh
dxhd -c /my/custom/path/to/a/config/file
```

A dxhd config file should containt a shebang (defaults to `/bin/sh`) on top of a file, which is where binding actions take action

## Syntax

```
<file>

#! shebang

## a comment
######### a comment

# key + combo
<what to do>

# key + combo + with + {1-9,0,a-z} + ranges
<what to do {with,these,ranges}>
```

### Demo

!outdated demo!

![demo gif](./dxhd_demo.gif)

### Roadmap

* [x] basic keybindings
* [ ] released keybindings ([#4](https://github.com/dakyskye/dxhd/issues/4))
* [x] ranges ([#5](https://github.com/dakyskye/dxhd/issues/5))
* ~~[x] formatting ([#6](https://github.com/dakyskye/dxhd/issues/6))~~
* [ ] daemonisation ([#3](https://github.com/dakyskye/dxhd/issues/3))

## License

Licensed under the [**MIT**](https://choosealicense.com/licenses/mit/) license.
