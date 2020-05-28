# dxhd

## daky's X11 hotkey daemon

dxhd is heavily inspired by [sxhkd](https://github.com/baskerville/sxhkd), written in Go, and has quite elegant syntax for configuration files!

## Installation

* Arch User Repository

```sh
git clone https://aur.archlinux.org/dxhd-git.git
# or binary version - git clone https://aur.archlinux.org/dxhd-bin.git
cd dxhd-git
# or cd dxhd-bin if you cloned binary one
makepkg -si
```

or use an AUR helper like yay - `yay -S dxhd`

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

## Features (what's inside parentheses, are just minimal example patterns)

* key press events (`super + key`, where `key` is a non-modifier key)
* key release events (`super + @key` where `key` is a non-modifier key, and `@` is specifier)
* mouse button press events (`!mouseN` where `n` is button number, and `!` is specifier)
* mouse button release events (`!@mouseN` where `n` is button number, and `!` and `@` are specifiers)
* variants (`{a,b,c}`)
* ranges (`{1-9}`, `{a-z}`, `{1-3,5-9,i-k,o-z}`)
* in-place reloading (`dxhd -r`)
* support for any shell scripting language (sh, bash, ksh, zsh, python, perl etc.) given as a shebang
* support for scripting, as much as a user wishes!

### Demo

! outdated gif !

![demo gif](./dxhd_demo.gif)

## Configuration

The default config file is located at `~/.config/dxhd/dxhd.sh`, however, dxhd can use a file from any path, by passing it to `-c`:

```sh
dxhd -c /my/custom/path/to/a/config/file
```

A dxhd config file should containt a shebang (defaults to `/bin/sh`) on top of a file, which is where binding actions take action

## Syntax

\* config file *
```
#! shebang

## a comment
######### a comment

# key + combo
<what to do>

# key + combo + with + {1-9,0,a-z} + ranges
<what to do {with,these,ranges}>
```

## Daemonisation

Rather then dxhd self daemonising itself, let other programs do their job.

Use `systemd`, `runit`, `openrc` or other Linux init system to start dxhd on system startup,
or let your DE/WM start it by adding an ampersant at the end `dxhd -c path/to/config &`,
optionally, use `disown` keyword to make it not owned by the DE/WM process.

### For further help, join the developer's Discord guild

<a target="_blank" href="https://discord.gg/x5RuZCN">
	<img src="https://img.shields.io/discord/627168403005767711?color=%238577ce&label=dakyskye%27s%20discord%20guild&logo=discord&logoColor=%23FFFFFF&style=plastic">
</a>

## Examples

* [shell](https://github.com/dakyskye/dxhd/tree/master/examples/dxhd.sh)
* [python](https://github.com/dakyskye/dxhd/tree/master/examples/dxhd.py)
* [author's config](https://github.com/dakyskye/dotfiles/tree/master/dxhd/)

## License

Licensed under the [**MIT**](https://choosealicense.com/licenses/mit/) license.
