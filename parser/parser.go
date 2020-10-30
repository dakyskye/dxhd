package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dakyskye/dxhd/logger"
	"github.com/sirupsen/logrus"
)

// event type
type EventType int8

// event definitions
const (
	EvtKeyPress EventType = iota
	EvtKeyRelease
	EvtButtonPress
	EvtButtonRelease
)

// FileData holds the data of parsed file
type FileData struct {
	OriginalBinding string
	Binding         strings.Builder
	Command         strings.Builder
	EvtType         EventType
	hasVariant      bool
}

// ranges hold the data of a keybinding and its command
type ranges struct {
	binding rng
	command struct {
		rng
		skip      bool
		numerical bool
	}
}

type rng struct {
	start, end       int
	startStr, endStr string
}

type variantGroup struct {
	binding, command []string
}

// global regular expressions, compiled once at run-time
var (
	keybindingPattern   = regexp.MustCompile(`^#(((!?@?)|@?!?)\w+{.*?}|((!?@?)|@?!?){.*?}|((!?@?)|@?!?)\w+)(((\+(((!?@?)|@?!?)\w+{.*?}|((!?@?)|@?!?){.*?}|((!?@?)|@?!?)\w+)))+)?`)
	variantPattern      = regexp.MustCompile(`{.*?}`)
	bindingRangePattern = regexp.MustCompile(`([0-9]|[a-z])-([0-9]|[a-z])`)
	commandRangePattern = regexp.MustCompile(`(?m)^(([0-9]+)-([0-9]+))|(([a-z])-([a-z]))$`)
	numericalPattern    = regexp.MustCompile(`([0-9]+)-([0-9]+)`)
	alphabeticalPattern = regexp.MustCompile(`([a-z])-([a-z])`)
	mouseBindPattern    = regexp.MustCompile(`mouse([0-9]+)`)
	xfKeyPattern        = regexp.MustCompile(`XF86\w+`)
)

func Parse(what interface{}, data *[]FileData) (shell, globals string, err error) {
	if data == nil {
		return "", "", errors.New("empty value was passed to parse function")
	}

	var reader *bufio.Reader

	switch what.(type) {
	case string:
		configFile, e := os.Open(what.(string))

		if e != nil {
			err = e
			return
		}

		defer func() {
			e := configFile.Close()
			if e != nil {
				if err == nil {
					err = e
				} else {
					logger.L().WithError(err).Debug("failed to close config file")
				}
			}
		}()

		reader = bufio.NewReader(configFile)
	case []byte:
		reader = bufio.NewReader(bytes.NewReader(what.([]byte)))
	default:
		err = errors.New("invalid type was passed to Parse function")
	}

	lineNumber := 0
	shell = "/bin/sh"
	wasKeybinding := false
	wasPrefix := false
	datum := []FileData{}
	index := 0
	globalsBuilder := new(strings.Builder)
	globalsEnded := false

	// read file line by line
	for {
		lineNumber++

		var (
			line     []byte
			isPrefix bool
		)

		line, isPrefix, err = reader.ReadLine()

		if err != nil {
			break
		}

		if index+1 != len(datum) {
			datum = append(datum, FileData{})
		}

		lineStr := string(line)

		// skip the shebang
		if lineNumber == 1 && strings.HasPrefix(lineStr, "#!") {
			shell = lineStr[2:]
			continue
		}

		// skip an empty line
		if lineStr == "" {
			continue
		}

		if !strings.HasPrefix(lineStr, "#") && !globalsEnded {
			globalsBuilder.WriteString(lineStr + "\n")
			continue
		}

		// ignore comments (##+)
		if strings.HasPrefix(lineStr, "##") {
			if !globalsEnded {
				globalsEnded = true
			}
			continue
		}

		// decide whether the line is a keybinding or not
		if strings.HasPrefix(lineStr, "#") {
			if isPrefix {
				err = fmt.Errorf("a keybinding can't be that long, line %d, file %s", lineNumber, what)
				return
			}
			if !globalsEnded {
				globalsEnded = true
			}
			// erase spaces for key validation
			lineStr = strings.ReplaceAll(lineStr, " ", "")

			if keybindingPattern.MatchString(lineStr) {
				if datum[index].Command.Len() != 0 {
					index++
					datum = append(datum, FileData{})
				}
				// trim # prefix
				lineStr := lineStr[1:]

				// overwrite previous prefix if needed
				if wasKeybinding {
					if datum[index].Binding.Len() != 0 {
						datum[index].Binding.Reset()
						logger.L().WithFields(logrus.Fields{"file": what, "line": lineNumber}).Info("overwriting keybinding")
						logger.L().WithFields(logrus.Fields{"old": datum[index].Binding.String(), "new": lineStr}).Debug("overwriting keybinding")
					}
				}

				// getEventType merges two events into one type
				getEventType := func(old, new EventType) (evt EventType) {
					switch old {
					case EvtKeyPress:
						evt = new
					case EvtKeyRelease:
						if new != EvtKeyPress {
							evt = new
						}
					case EvtButtonPress:
						if new == EvtKeyRelease || new == EvtButtonRelease {
							evt = EvtButtonRelease
						}
					case EvtButtonRelease:
						evt = EvtButtonRelease
					default:
						evt = new
					}
					return
				}

				// set to -1, in case a keybinding is a single letter
				datum[index].EvtType = -1
				for _, key := range strings.Split(lineStr, "+") {
					if len(key) > 1 {
						if strings.HasPrefix(key, "@mouse") {
							datum[index].EvtType = getEventType(datum[index].EvtType, EvtButtonRelease)
						} else if strings.HasPrefix(key, "mouse") {
							datum[index].EvtType = getEventType(datum[index].EvtType, EvtButtonPress)
						} else if strings.HasPrefix(key, "@") {
							datum[index].EvtType = getEventType(datum[index].EvtType, EvtKeyRelease)
						} else {
							datum[index].EvtType = getEventType(datum[index].EvtType, EvtKeyPress)
						}

					}
				}
				// means a keybinding was the single letter
				if datum[index].EvtType == -1 {
					datum[index].EvtType = EvtKeyPress
				}
				_, err = datum[index].Binding.WriteString(lineStr)
				if err != nil {
					return
				}
				datum[index].hasVariant = len(variantPattern.FindStringIndex(lineStr)) > 0
				wasKeybinding = true
			}
		} else {
			wasKeybinding = false
			if isPrefix {
				if wasPrefix {
					datum[index].Command.Write(line)
				} else {
					if datum[index].Command.Len() != 0 {
						datum[index].Command.Write([]byte("\n"))
					}
					datum[index].Command.Write(line)
					wasPrefix = true
				}
				continue
			}

			if wasPrefix {
				datum[index].Command.Write(line)
				wasPrefix = false
			} else {
				if datum[index].Command.Len() != 0 {
					datum[index].Command.Write([]byte("\n"))
				}
				datum[index].Command.Write(line)
			}
		}
	}

	// we are ok when we reach the end of line
	if err == io.EOF {
		err = nil
	} else {
		return
	}

	// xgb requires these shorthands to be replaced to what they are called internally
	replaceShorthands := func(data *FileData) (err error) {
		data.OriginalBinding = data.Binding.String()
		data.Binding.Reset()

		modified := data.OriginalBinding

		// extract xf86 keys if any
		matches := xfKeyPattern.FindAllString(modified, -1)
		indexes := xfKeyPattern.FindAllStringIndex(modified, -1)
		if len(matches) != len(indexes) {
			err = errors.New("can not process XF86 keys properly")
			return
		}

		// lowercase whole line
		modified = strings.ToLower(modified)

		// put XF86 keys as they were before in lineStr
		for in, index := range indexes {
			modified = strings.Replace(modified, modified[index[0]:index[1]], matches[in], 1)
		}

		modified = strings.ReplaceAll(data.OriginalBinding, "+", "-")
		modified = strings.ReplaceAll(modified, "super", "mod4")
		modified = strings.ReplaceAll(modified, "alt", "mod1")
		modified = strings.ReplaceAll(modified, "ctrl", "control")
		modified = strings.ReplaceAll(strings.ReplaceAll(modified, "@", ""), "!", "")
		// replace mouseN with N
		if data.EvtType == EvtButtonPress || data.EvtType == EvtButtonRelease {
			modified = mouseBindPattern.ReplaceAllString(modified, "$1")
		}
		_, err = data.Binding.WriteString(modified)
		return
	}

	*data = nil
	for _, d := range datum {
		// replicate a keybinding and it's command if it has variants
		if d.hasVariant {
			replicated, e := replicate(d.Binding.String(), d.Command.String())
			if e != nil {
				err = fmt.Errorf("can't register %s keybinding, error (%s)", strings.TrimPrefix(d.Binding.String(), "#"), e.Error())
				return
			}
			for _, repl := range replicated {
				repl.EvtType = d.EvtType
				err = replaceShorthands(repl)
				if err != nil {
					return
				}
				*data = append(*data, FileData{OriginalBinding: repl.OriginalBinding, Binding: repl.Binding, Command: repl.Command, EvtType: d.EvtType})
			}
		} else {
			err = replaceShorthands(&d)
			if err != nil {
				return
			}
			*data = append(*data, FileData{OriginalBinding: d.OriginalBinding, Binding: d.Binding, Command: d.Command, EvtType: d.EvtType})
		}
	}

	// means config file was empty
	if len(*data) == 1 && ((*data)[0].Command.String() == "" || (*data)[0].Binding.String() == "") {
		err = errors.New("config file does not contain any binding")
		return
	}

	// build string from globalsBuilder
	globals = globalsBuilder.String()

	return
}

// replicate replicates variants
func replicate(binding, command string) (replicated []*FileData, err error) {
	// find all the variants
	bindingVariants, commandVariants := variantPattern.FindAllString(binding, -1), variantPattern.FindAllString(command, -1)

	// make sure the amount of variants do match
	if len(bindingVariants) != len(commandVariants) {
		err = errors.New("the amount of variants in a keybinding and its command do not match")
		return
	}

	var bindingVars, commandVars [][]string

	// extract variant members
	extract := func(from []string, where *[][]string) {
		for _, f := range from {
			*where = append(*where, strings.Split(strings.TrimSuffix(strings.TrimPrefix(f, "{"), "}"), ","))
		}
	}

	extract(bindingVariants, &bindingVars)
	extract(commandVariants, &commandVars)

	// validate the amount of variant memebers do match
	for i, b := range bindingVars {
		if len(b) != len(commandVars[i]) {
			err = errors.New("the amount of variant members in a keybinding and its command do not match")
			return
		}
	}

	// validate and extract ranges
	var rngs []ranges
	rngs, err = extractRanges(bindingVars, commandVars)
	if err != nil {
		return
	}

	var expandedBindingRanges, expandedCommandRanges []string

	// expands a range in a keybinding({1-9} -> {1},{2},{3},{...},{9})
	expandRange := func(r ranges, binding, acton string, bindings, commands *[]string) {
		// bindings
		for bIn := r.binding.start; bIn != r.binding.end+1; bIn++ {
			*bindings = append(*bindings, strings.Replace(
				binding,
				fmt.Sprintf("%s-%s", r.binding.startStr, r.binding.endStr),
				fmt.Sprintf("%c", rune(bIn)),
				1,
			))
			if r.command.skip {
				*commands = append(*commands, acton)
			}
		}
		// commands
		if !r.command.skip {
			for aIn := r.command.start; aIn != r.command.end+1; aIn++ {
				if r.command.numerical {
					*commands = append(*commands, strings.Replace(
						command,
						fmt.Sprintf("%s-%s", r.command.startStr, r.command.endStr),
						fmt.Sprintf("%d", aIn),
						1,
					))
				} else {
					*commands = append(*commands, strings.Replace(
						command,
						fmt.Sprintf("%s-%s", r.command.startStr, r.command.endStr),
						fmt.Sprintf("%c", rune(aIn)),
						1,
					))
				}
			}
		}
	}

	// for as long as we have unexpanded ranges, expand them
	for len(rngs) > 0 {
		if len(expandedBindingRanges) > 0 {
			var newBindingRanges, newCommandRanges []string

			if len(expandedCommandRanges) != len(expandedBindingRanges) {
				err = errors.New("an unknown error occurred whilst expanding keybinding and command ranges")
			}

			for i := 0; i != len(expandedBindingRanges); i++ {
				expandRange(rngs[0], expandedBindingRanges[i], expandedCommandRanges[i], &newBindingRanges, &newCommandRanges)
			}

			expandedBindingRanges, expandedCommandRanges = newBindingRanges, newCommandRanges
		} else {
			expandRange(rngs[0], binding, command, &expandedBindingRanges, &expandedCommandRanges)
		}
		rngs = rngs[1:]
	}

	if len(expandedCommandRanges) != len(expandedBindingRanges) {
		err = errors.New("an unknown error occurred whilst expanding keybinding and command ranges")
	}

	// replicateVariant replaces pattern with each member of variants group
	replicateVariant := func(in, pattern string, variants []string, where *[]string) {
		for _, v := range variants {
			if v == "_" {
				*where = append(*where, strings.Replace(
					in,
					pattern,
					"",
					1,
				))
			} else {
				*where = append(*where, strings.Replace(
					in,
					pattern,
					v,
					1,
				))
			}
		}
	}

	// in case our keybinding and command had no ranges
	if len(expandedBindingRanges) == 0 {
		expandedBindingRanges = append(expandedBindingRanges, binding)
		expandedCommandRanges = append(expandedCommandRanges, command)
	}

	// do replicate every variant member
	for i, r := 0, 0; i != len(expandedBindingRanges); i++ {
		var replicatedBindings, replicatedCommands []string
		vGroup := &variantGroup{}
		vGroup.command = variantPattern.FindAllString(expandedCommandRanges[i], -1)
		vGroup.binding = variantPattern.FindAllString(expandedBindingRanges[i], -1)

		if !(len(vGroup.command) == len(vGroup.binding) && len(vGroup.command) > 0) {
			err = errors.New("can not extract variant groups")
			return
		}

		// for as long as we have binding AND command in a variant group
		for len(vGroup.binding) > 0 {
			// extract variant members
			bVariantMembers := strings.Split(strings.TrimSuffix(strings.TrimPrefix(vGroup.binding[0], "{"), "}"), ",")
			aVariantMembers := strings.Split(strings.TrimSuffix(strings.TrimPrefix(vGroup.command[0], "{"), "}"), ",")
			// if we already replicated a variant, use it
			if len(replicatedBindings) > 0 {
				var newBindingVariants, newCommandVariants []string

				for _, alreadyR := range replicatedBindings {
					replicateVariant(alreadyR, vGroup.binding[0], bVariantMembers, &newBindingVariants)
				}

				for _, alreadyR := range replicatedCommands {
					replicateVariant(alreadyR, vGroup.command[0], aVariantMembers, &newCommandVariants)
				}

				replicatedBindings, replicatedCommands = newBindingVariants, newCommandVariants
			} else {
				replicateVariant(expandedBindingRanges[i], vGroup.binding[0], bVariantMembers, &replicatedBindings)
				replicateVariant(expandedCommandRanges[i], vGroup.command[0], aVariantMembers, &replicatedCommands)
			}
			vGroup.binding = vGroup.binding[1:]
			vGroup.command = vGroup.command[1:]
		}

		if len(replicatedBindings) != len(replicatedCommands) {
			err = errors.New("replication went wrong")
			return
		}

		// append replicated bindings and command to the return result
	appender:
		for i := 0; i != len(replicatedBindings); i++ {
			// we get ++ when we replace underscore literal with nothing
			replicatedBindings[i] = strings.ReplaceAll(replicatedBindings[i], "++", "+")
			if i > 0 {
				for _, aR := range replicated {

					if aR.Binding.String() == replicatedBindings[i] {
						continue appender
					}
				}
			}
			replicated = append(replicated, &FileData{})
			_, err = replicated[r].Binding.WriteString(replicatedBindings[i])
			if err != nil {
				return
			}
			_, err = replicated[r].Command.WriteString(replicatedCommands[i])
			if err != nil {
				return
			}
			r++
		}
	}

	return
}

// extracts every range from a config file
func extractRanges(bindingVars, commandVars [][]string) (r []ranges, err error) {
	// range patterns for binding and command and range errors
	var (
		rangeParseErr   = errors.New("could not parse a range")
		invalidRangeErr = errors.New("invalid parse given")
	)
	// iterate over a binding variants, and replicate ranges
	for bIn, bVars := range bindingVars {
		// iterate over the variants of the binding variants
		for vIn, bVar := range bVars {
			// check if the variant is a range
			if bindingRangePattern.MatchString(bVar) {
				bRange := bindingRangePattern.FindStringSubmatch(bVar)
				if len(bRange) != 3 {
					err = rangeParseErr
					return
				}
				var (
					aVar             = commandVars[bIn][vIn]
					aRange           []string
					aRangeValidation = true
				)
				// make sure command variant is also a range (or _)
				if !commandRangePattern.MatchString(aVar) {
					if aVar == "_" {
						// in case it's _, skip the range validation
						aRangeValidation = false
					} else {
						err = errors.New("the indexes of ranges for a keybinding and its command do not match")
						return
					}
				}

				// int values for comparison
				// string values for assignment
				var (
					bStart, bEnd       int
					bStartStr, bEndStr string
				)

				// convert strings to runes to int
				bStart = int([]rune(bRange[1])[0])
				bEnd = int([]rune(bRange[2])[0])

				// make sure the given range is valid
				if bStart >= bEnd {
					err = invalidRangeErr
					return
				}

				// assign string values
				bStartStr = string(rune(bStart))
				bEndStr = string(rune(bEnd))

				var (
					aStart, aEnd       int
					aStartStr, aEndStr string
					aNumerical         = false
				)

				if aRangeValidation {
					if numericalPattern.MatchString(aVar) { // is it a numerical range?
						aRange = numericalPattern.FindStringSubmatch(aVar)
						if len(aRange) != 3 {
							err = rangeParseErr
							return
						}
						aStart, err = strconv.Atoi(aRange[1])
						if err != nil {
							return
						}
						aEnd, err = strconv.Atoi(aRange[2])
						if err != nil {
							return
						}
						for _, a := range aRange[1] {
							aStartStr += string(a)
						}
						for _, a := range aRange[2] {
							aEndStr += string(a)
						}
						aNumerical = true
					} else if alphabeticalPattern.MatchString(aVar) { // well, is it an alphabetical range?
						aRange = alphabeticalPattern.FindStringSubmatch(aVar)
						if len(aRange) != 3 {
							err = rangeParseErr
							return
						}
						aStart = int([]rune(aRange[1])[0])
						aEnd = int([]rune(aRange[2])[0])
						aStartStr = string(rune(aStart))
						aEndStr = string(rune(aEnd))
					} else { // it's an invalid range
						err = invalidRangeErr
						return
					}
					// make sure the ranges match
					// 1-9 compared to 11-19 is a valid range
					// 4-8 compared to a-e is a valid range also
					// b start-end        a start-end
					if (bStart - bEnd) != (aStart - aEnd) {
						err = errors.New("the ranges of a keybinding and its command do not match")
						return
					}
				}
				r = append(r, ranges{})
				r[len(r)-1].binding.start = bStart
				r[len(r)-1].binding.startStr = bStartStr
				r[len(r)-1].binding.end = bEnd
				r[len(r)-1].binding.endStr = bEndStr

				r[len(r)-1].command.start = aStart
				r[len(r)-1].command.startStr = aStartStr
				r[len(r)-1].command.end = aEnd
				r[len(r)-1].command.endStr = aEndStr
				r[len(r)-1].command.skip = !aRangeValidation
				r[len(r)-1].command.numerical = aNumerical
			}
		}
	}
	return
}
