# dxhd

## daky's X11 hotkey daemon

dxhd is heavily inspired by [sxhkd](https://github.com/baskerville/sxhkd), written in Go, and has an elegant syntax for configuration files!

## Installation

**NOTE:** the git version, a.k.a. the master version is usually more bug free than the released, binary ones, since introduced bugs first get fixed in this version.

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
make fast
```

and copy `dxhd` executable file to somewhere in your `$PATH`

or alternatively run `make install`, which builds and copies the built executable to `/usr/bin/` directory.

* From releases

Download the `dxhd` executable file from the latest release, from [releases page](https://github.com/dakyskye/dxhd/releases)

and copy `dxhd` executable file to somewhere in your `$PATH`

**Note:** `go get`ting dxhd is possible, but not recommended. Read more [here](https://github.com/dakyskye/dxhd#why-is-go-getting-dxhd-not-recommended)

## Features (what's inside parentheses, are just minimal example patterns)

* key press events (`super + key`, where `key` is a non-modifier key)
* key release events (`super + @key` where `key` is a non-modifier key, and `@` is a specifier)
* mouse button press events (`mouseN` where `n` is button number)
* mouse button release events (`@mouseN` where `n` is button number, and `@` is a specifier)
* variants (`{a,b,c}`)
* ranges (`{1-9}`, `{a-z}`, `{1-3,5-9,i-k,o-z}`)
* in-place reloading (`dxhd -r`)
* calculating the time parsing a config file took (`dxhd -p`)
* support for any shell scripting language (sh, bash, ksh, zsh, python, perl etc.) given as a shebang
* support for scripting, as much as a user wishes!

### Demo

! outdated gif !

![demo gif](./dxhd_demo.gif)

## Configuration

The default config file is located at `~/.config/dxhd/dxhd.sh`, however, dxhd can read a file from any path, by passing it to `-c`:

```sh
dxhd -c /my/custom/path/to/a/config/file
```

A dxhd config file should contain a shebang (defaults to `/bin/sh`) on top of a file, which will be the shell used for executing commands.

## Syntax

\* config file *
```
#!/shebang

## a comment
######### also a comment

# modifier + keys
<what to do>

# modifier + @keys
<what to do on release event>
```

## Running

By just running `dxhd`, you only get information level logs, however, you can set `DEBUG` environment variable, which will output more information, like what bindings are registered, what command failed etc.

To kill every running instance of dxhd, you can use built-in `-k` flag, which under the hood uses `pkill` command to kill instances.

## Daemonisation

Rather then dxhd self daemonising itself, let other programs do their job.

Use `systemd`, `runit`, `openrc` or other Linux init system to start dxhd on system startup,
or let your DE/WM start it by adding an ampersand at the end `dxhd -c path/to/config &`,
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

## FAQ

### Why was dxhd made

Because I had (and have) 20 workspaces, and `sxhkd` did not allow me to have `11-19` range,
that was one of the main reasons I started developing dxhd

### What makes dxhd better than sxhkd

* dxhd uses shebang to determine which shell to use (so you don't have to set an environment variable)
* dxhd config file syntax matches shell, python, perl and probably some other languages syntax
* dxhd is great with scripting, because of it's elegant syntax.  multi line scripts do not need `\` at the end of line
* dxhd allows you to have different range in a keybinding's action, for example, `1-9` in a keybinding, and `11-19` in it's action
* dxhd has support for mouse bindings out of the box, no patching required!

### How do I port my sxhkd config to dxhd

It is simple enough! (I personally used Vim macros when I did it.. Vim users will get it)
* convert any line starting with single `#` to a *dxhd comment* (so ## or more)
* put `#` before every keybinding (`super + a` to `# super + a`)
* remove spaces before lines (`  echo foo` to `echo foo`) (optional)
* remove every end-line backslash (`echo bar \` to `echo bar`) (probably optional, unsure)

So you'd end up with:

```
# print hello world
super + a
	echo hello \
	echo world
```

to

```sh
#!/bin/sh

## print hello world
# super + a
echo hello
echo world
```

### I use ranges, released key events and chords from sxhkd, does dxhd have them

Yes! And no.  dxhd has released key events and ranges, but no chords (yet - [wip](https://github.com/dakyskye/dxhd/issues/8))

### Is dxhd faster than sxhkd

They haven't benchmarked yet, so I don't know.
However, been using dxhd since the first release and haven't noticed any speed loss!

### Why is the released binary file ~~+8mb~~ +6mb

Because it's statically built, to make sure it will work on any supported machine!

### Why is go getting dxhd not recommended

Whilst `go get`ting dxhd should work fine, it's not recommended, because we can't know what version of dxhd you use in case you want to open a bug report or so.
+It's not like dxhd has any bug issue is not opened for already, since the developer of dxhd himself uses dxhd daily, but still.
