# dxhd

## daky's X11 hotkey daemon

This hotkey daemon is quite stable already and can be used *in production*, however, it's still WIP as it lacks some good features, like ranges and key release event support.

### Testing

* git clone the repo
* install Go programming language
* run `go build -o dxhd .`
* execute `./dxhd`

### Demo

![demo gif](./dxhd_demo.gif)

### Roadmap

* [x] basic keybindings (all pressed together)
* [ ] released keybindings (take action on key release event)
* [ ] ranges (1-9 and a-z)
* [ ] formatting (to expand ranges, i.e. `{1-9}` -> `%{F1}1-9`, which means 11-19)
* [ ] reloading as a command (dxhd reload)

## License

Licensed under the [**MIT**](https://choosealicense.com/licenses/mit/) license.
