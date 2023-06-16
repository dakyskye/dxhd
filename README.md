# dxhd

## daky's X11 hotkey daemon

`dxhd` is heavily inspired by [sxhkd](https://github.com/baskerville/sxhkd),
written in [Go](https://go.dev), and has an elegant syntax for configuration
files!

Thanks [JetBrains](https://jetbrains.com) for providing `dxhd` with free
licenses.

## READ THIS NOTICE FIRST

Dear repository visitor, I'm sorry I wrote a horrible codebase years ago which
I've tried to rewrite a couple of times now but never finished because of losing
interest sometimes, because of work sometimes, because of my "time for a new
side-project" mood sometimes.

However, I want to tell you the rewrite will definitely take place, because I
want to apply things I have learnt to `dxhd`. I want to write a proper parser
this time, I want to structure codebase properly this time, I want to show my
today's [Go](https://go.dev) knowledge & practice in `dxhd` codebase because
current master branch does not represent it. I usually put `dxhd` in my CV but
always hate to do it because current master codebase does not show my present
[Go](https://go.dev) skills but rather my past mistakes that taught me how not
to program.

The rewrite process will begin in 2022 summer: July and august.

## READ THIS SECOND (older notice)

~~`dxhd` is being rewritten in the [Rust](https://rust-lang.org/) programming
language. The two main collaborators, [dakyskye](https://github.com/dakyskye)
and [NotUnlikeTheWaves](https://github.com/NotUnlikeTheWaves), are working on
it.~~

~~Follow [issue #39](https://github.com/dakyskye/dxhd/issues/39) for more
information regarding the rewrite.~~

`dxhd` is going to be rewritten in [Go](https://go.dev) only to provide much
better codebase to make contributions easier as well as have a quality code so
others can reference from it.

The reason why a rewrite is required is that the current codebase is terrible
(but the app works well so it does not matter for an end-user). There is only
one known bug in the parser which I have documented here. It will be resolved,
as well as many things will be improved, after releasing the rewritten `dxhd`.

The bug:

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

The Git version (the master version) is usually more bug-free than the released,
binary ones, since introduced bugs first get fixed in this version.

### Manual Arch User Repository installation

```sh
git clone https://aur.archlinux.org/dxhd-git.git
# or binary version - git clone https://aur.archlinux.org/dxhd-bin.git
cd dxhd-git
# or cd dxhd-bin if you cloned binary one
makepkg -si
```

### Using an [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers)

| AUR helper                                  | Command        |
|---------------------------------------------|----------------|
| [paru](https://github.com/morganamilo/paru) | `paru -S dxhd` |
| [yay](https://github.com/Jguer/yay)         | `yay -S dxhd`  |

### From the source

```sh
git clone https://github.com/dakyskye/dxhd.git
cd dxhd
make fast
```

Copy the `dxhd` executable file somewhere in your `$PATH`

... or alternatively run `make install`, which builds and copies the built
executable to `/usr/local/bin/` directory.

### From releases

Download the `dxhd` executable file from the latest release, from [releases
page](https://github.com/dakyskye/dxhd/releases), then copy `dxhd` executable
file somewhere in your `$PATH`.

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
| support for running as many `dxhd` instances simultaneously as you want, to logically separate your keybindings  | -                                                                        |

## Configuration

The default config file is `~/.config/dxhd/dxhd.sh`, however, `dxhd` can read a
file from any path, by passing it the `-c` command line flag:

```sh
dxhd -c /my/custom/path/to/a/config/file
```

A `dxhd` config file should contain a shebang (defaults to `/bin/sh`) on top of
a file, which will be the shell used for executing commands.

## Syntax

config.sh
```sh
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

To kill every running instance of `dxhd`, you can use built-in `-k` flag, which
under the hood uses `pkill` command to kill instances.

## Daemonisation

~~Rather than `dxhd` self daemonising itself, let other programs do their job.~~

The `--background` (`-b`) flag is a simple workaround for daemonising `dxhd`. It
uses `/usr/sh` shell to achieve it, as [Go](https://go.dev) does not allow
forking a process without executing it.

## Examples

* [shell](https://github.com/dakyskye/dxhd/tree/master/examples/dxhd.sh)
* [python](https://github.com/dakyskye/dxhd/tree/master/examples/dxhd.py)
* [author's config](https://github.com/dakyskye/dotfiles/tree/master/dxhd/)

## License

Licensed under the [MIT](https://choosealicense.com/licenses/mit/) license.

## FAQ

### Why was `dxhd` made

Because I had ~~(and have)~~ 20 workspaces, and `sxhkd` did not allow me to
define `11-19` range easily which was one of the main reasons I started
developing `dxhd`.

### What makes `dxhd` better than `sxhkd`

* `dxhd` uses shebang to determine which shell to use (so you don't have to set
  an environment variable).
* `dxhd` configuration file syntax matches shell, Python, Perl and probably some
  other language syntaxes.
* `dxhd` config lets you declare global variables for each keybinding command.
* `dxhd` is great with scripting, because of its elegant syntax. Multi line
  scripts do not need `\` at the end of line.
* `dxhd` allows you to have different range in a keybinding's command, for
  example, `1-9` in a keybinding, and
* `11-19` in its body (command area).
* `dxhd` has support for mouse bindings out of the box, no patching required!

### How do I port my `sxhkd` config to `dxhd`

It is simple enough! (I personally used Vim macros when I did it... Vim users
will get it)
* convert any line starting with single `#` to a `dxhd` comment (so ## or
  more)
* put a `#` before every keybinding (`super + a` to `# super + a`)
* remove spaces before lines (`  echo foo` to `echo foo`) (optional)
* remove every end-line backslash (`echo bar \` to `echo bar`) (most likely
  optional, unsure)

So you'd end up from:

```sh
# print hello world
super + a
	echo hello \
	echo world
```

with

```sh
#!/bin/sh

## print hello world
# super + a
echo hello
echo world
```

### I use ranges, released key events and chords from `sxhkd`, does `dxhd` have them

Yes! And no. `dxhd` has released key events and ranges, but no chords (yet -
[WIP](https://github.com/dakyskye/dxhd/issues/8))

### How do global variables inside a config file work

Everything after (if there is) the shebang before the first comment/keybinding
is collected and passed to each keybinding's command.

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

```python
#!/usr/bin/python

foo="foo bar"

## print the value of foo variable
# super + i
print(foo)
```

### Is `dxhd` faster than `sxhkd`

They haven't benchmarked yet, so I don't know. However, I have been using `dxhd`
since the first release and haven't noticed any speed loss!

### Why is the released binary file ~~+8mb~~ ~~+6mb~~ +3mb

Because it's statically built, to make sure it will work on any supported
machine!
