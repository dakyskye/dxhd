# dxhd

## daky's X11 hotkey daemon

_dxhd_ is heavily inspired by [sxhkd](https://github.com/baskerville/sxhkd),
written in Go, and has an elegant syntax for configuration files!

Thanks [JetBrains](https://jetbrains.com) for providing dxhd with free licenses.

## READ THIS FIRST

Issue [#39](https://github.com/dakyskye/dxhd/issues/39) is opened to track the
rewrite process of _dxhd_.  The reason why
[rewrite](https://github.com/dakyskye/dxhd/tree/rewrite) is required is that the
current codebase is terrible (but the app works well, does not matter for an
end-user).  There is only one known bug in the parser which I have documented
here. dxhd is being rewritten and it will be resolved!

**the bug**:

```sh
#!/bin/sh

# super + {a,b}
echo it was either {aaaaa,bbbbbbb}
echo I want to print {aaaaa,bbbbbbb}
echo I can print anything {tho, though}
```

Parser will error on this. There is a workaround you can use (for some cases):

```sh
#!/bin/sh

# super + {a,b}
what={aaaaa,bbbbbbb}
echo it was either "$what"
echo I want to print "$what"
# echo I can print anything {tho, though} <-- good luck
```

## Installation

**NOTE:** the git version, a.k.a. the master version is usually more bug-free
than the released, binary ones, since introduced bugs first get fixed in this
version.

* Manual Arch User Repository installation

```sh
git clone https://aur.archlinux.org/dxhd-git.git
# or binary version - git clone https://aur.archlinux.org/dxhd-bin.git
cd dxhd-git
# or cd dxhd-bin if you cloned binary one
makepkg -si
```

* Using an [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers)

| AUR helper                                  | Command        |
|---------------------------------------------|----------------|
| [paru](https://github.com/morganamilo/paru) | `paru -S dxhd` |
| [yay](https://github.com/Jguer/yay)         | `yay -S dxhd`  |

* From the source

```sh
git clone https://github.com/dakyskye/dxhd.git
cd dxhd
make fast
```

Copy the _dxhd_ executable file somewhere in your `$PATH`

... or alternatively run `make install`, which builds and copies the built
executable to `/usr/bin/` directory.

* From releases

Download the _dxhd_ executable file from the latest release, from
[releases page](https://github.com/dakyskye/dxhd/releases), then copy `dxhd`
executable file somewhere in your `$PATH`.

**Note:** `go get`ting _dxhd_ is possible, but not recommended.  Read more
[here](https://github.com/dakyskye/dxhd#why-is-go-getting-dxhd-not-recommended)

## Features (what's inside parentheses, are just minimal example patterns)

| Feature                                                                                                        | Description                                                              |
|----------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| key press events                                                                                               | `super + key`, where `key` is a non-modifier key                         |
| key release events                                                                                             | `super + @key` where `key` is a non-modifier key, and `@` is a specifier |
| mouse button press events                                                                                      | `mouseN` where `n` is button number                                      |
| mouse button release events                                                                                    | `@mouseN` where `n` is button number, and `@` is  a specifier            |
| variants                                                                                                       | `{a,b,c}`                                                                |
| ranges                                                                                                         | `{1-9}`, `{a-z}`, `{1-3,5-9,i-k,o-z}`                                    |
| in-place reloading                                                                                             | `dxhd -r`                                                                |
| calculating the time parsing a config file took                                                                | `dxhd -p`                                                                |
| editing config files quickly                                                                                   | `dxhd -e i3.py`                                                          |
| running as a daemon                                                                                            | `dxhd -b`                                                                |
| running interactively                                                                                          | `dxhd -i`                                                                |
| support for any shell scripting language                                                                       | sh, bash, ksh, zsh, python, perl, etc. given as a shebang                |
| support for global variable declarations in a config                                                           | -                                                                        |
| support for scripting, as much as a user wishes!                                                               | -                                                                        |
| support for running as many dxhd instances simultaneously as you want, to logically separate your keybindings  | -                                                                        |

## Configuration

The default config file is `~/.config/dxhd/dxhd.sh`, however, _dxhd_ can read a
file from any path, by passing it to `-c`:

```sh
dxhd -c /my/custom/path/to/a/config/file
```

A _dxhd_ config file should contain a shebang (defaults to `/bin/sh`) on top of
a file, which will be the shell used for executing commands.

## Syntax

\* config file *
```
#!/shebang

test=5 # a globally declared variable for each keybinding command

## a comment
######### also a comment

# modifier + keys
<what to do>

# modifier + @keys
<what to do on release event>
```

## Running

By just running `dxhd`, you only get information level logs, however, you can
set `DEBUG` environment variable, which will output more information, like what
bindings are registered, what command failed etc.

To kill every running instance of dxhd, you can use built-in `-k` flag, which
under the hood uses `pkill` command to kill instances.

## Daemonisation

~~Rather than dxhd self daemonising itself, let other programs do their job.~~

The `--background` (`-b`) flag is a simple workaround for *daemonising* _dxhd_.
It uses `/usr/sh` shell to achieve it, as Go does not allow forking a process
without executing it.

### For further help, join the developer's Discord guild

<a target="_blank" href="https://discord.gg/x5RuZCN">
	<img src="https://img.shields.io/discord/627168403005767711?color=%238577ce&label=dakycord&logo=discord&logoColor=%23FFFFFF&style=plastic" alt="dakycord">
</a>

## Examples

* [shell](https://github.com/dakyskye/dxhd/tree/master/examples/dxhd.sh)
* [python](https://github.com/dakyskye/dxhd/tree/master/examples/dxhd.py)
* [author's config](https://github.com/dakyskye/dotfiles/tree/master/dxhd/)

## License

Licensed under the [**MIT**](https://choosealicense.com/licenses/mit/) license.

## FAQ

### Why was dxhd made

Because I had (and have) 20 workspaces, and `sxhkd` did not allow me to have
`11-19` range, that was one of the main reasons I started developing dxhd

### What makes dxhd better than sxhkd

* _dxhd_ uses shebang to determine which shell to use (so you don't have to set
  an environment variable)
* _dxhd_ config file syntax matches shell, python, perl and probably some other
  languages syntax
* _dxhd_ config lets you declare global variables for each keybinding command
* _dxhd_ is great with scripting, because of it's elegant syntax.  multi line
  scripts do not need `\` at the end of line
* _dxhd_ allows you to have different range in a keybinding's command, for
  example, `1-9` in a keybinding, and `11-19` in it's command
* _dxhd_ has support for mouse bindings out of the box, no patching required!

### How do I port my sxhkd config to dxhd

It is simple enough! (I personally used Vim macros when I did it.. Vim users
will get it)
* convert any line starting with single `#` to a *dxhd comment* (so ## or more)
* put `#` before every keybinding (`super + a` to `# super + a`)
* remove spaces before lines (`  echo foo` to `echo foo`) (optional)
* remove every end-line backslash (`echo bar \` to `echo bar`)
  (probably optional, unsure)

So you'd end up with:

```sh
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

Yes! And no.  dxhd has released key events and ranges, but no chords (yet -
[wip](https://github.com/dakyskye/dxhd/issues/8))

### How do global variables inside a config file work

Everything after (if there is) the shebang before the first comment/keybinding
is collected and passed to each keybinding command

A shell example:

```sh
#!/bin/sh

INFO="$(wmctrl -m)"

## print info about my WM
# super + i
echo "Info about your WM:"
echo "$INFO"
```

A Python example:

```py
#!/usr/bin/python

foo="foo bar"

## print the value of foo variable
# super + i
print(foo)
```

### Is dxhd faster than sxhkd

They haven't benchmarked yet, so I don't know.  However, been using _dxhd_ since
the first release and haven't noticed any speed loss!

### Why is the released binary file ~~+8mb~~ ~~+6mb~~ +3mb

Because it's statically built, to make sure it will work on any supported
machine!

### Why is go getting dxhd not recommended

Whilst `go get`ting _dxhd_ should work fine, it's not recommended, because we
can't know what version of dxhd you use in case you want to open a bug report or
so. +It's not like _dxhd_ has any bug issue is not opened for already, since the
developer of dxhd himself uses _dxhd_ daily, but still.
