package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type filedata struct {
	binding    strings.Builder
	action     strings.Builder
	hasVariant bool
}

type ranges struct {
	binding rng
	action  struct {
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
	action, binding []string
}

var (
	keybindingPattern   = regexp.MustCompile(`^#(\w+{.*?}|{.*?}|\w+)(((\+(\w+{.*?}|{.*?}|\w+)))+)?`)
	variantPattern      = regexp.MustCompile(`{.*?}`)
	bindingRangePattern = regexp.MustCompile(`([0-9]|[a-z])-([0-9]|[a-z])`)
	actionRangePattern  = regexp.MustCompile(`(?m)^(([0-9]+)-([0-9]+))|(([a-z])-([a-z]))$`)
	numericalPattern    = regexp.MustCompile(`([0-9]+)-([0-9]+)`)
	alphabeticalPattern = regexp.MustCompile(`([a-z])-([a-z])`)
)

func parse(file string, data *[]filedata) (shell string, err error) {
	if data == nil {
		return "", errors.New("empty value was passed to parse function")
	}

	configFile, err := os.Open(file)

	if err != nil {
		return
	}

	defer configFile.Close()

	reader := bufio.NewReader(configFile)

	lineNumber := 0
	shell = "/bin/sh"
	wasKeybinding := false
	wasPrefix := false
	datum := []filedata{}
	index := 0

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
			datum = append(datum, filedata{})
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

		// ignore comments (##+)
		if strings.HasPrefix(lineStr, "##") {
			continue
		}

		// decide whether the line is a keybinding or not
		if strings.HasPrefix(lineStr, "#") {
			if isPrefix {
				log.Fatalf("a keybinding can't be that long, line %d, file %s", lineNumber, file)
				os.Exit(1)
			}
			// erase spaces for key validation
			lineStr = strings.ReplaceAll(lineStr, " ", "")

			if keybindingPattern.MatchString(lineStr) {
				if datum[index].action.Len() != 0 {
					index++
					datum = append(datum, filedata{})
				}
				replace := func(origin *string, old, new string) {
					*origin = strings.ReplaceAll(*origin, old, new)
				}

				// trim # prefix
				lineStr := lineStr[1:]
				// lowercase whole string, since xgbutil accepts any case
				lineStr = strings.ToLower(lineStr)
				// replace shorthands
				replace(&lineStr, "super", "mod4")
				replace(&lineStr, "ctrl", "control")
				replace(&lineStr, "+", "-")
				// overwrite previous prefix if needed
				if wasKeybinding {
					if datum[index].binding.Len() != 0 {
						datum[index].binding.Reset()
						fmt.Println(fmt.Sprintf("overwriting %d", lineNumber))
						fmt.Println(fmt.Sprintf("previous - \"%s\"", datum[index].binding.String()))
						fmt.Println(fmt.Sprintf("new - \"%s\"", lineStr))
					}
				}
				datum[index].binding.Write([]byte(lineStr))
				datum[index].hasVariant = len(variantPattern.FindStringIndex(lineStr)) > 0
				wasKeybinding = true
			}
		} else {
			wasKeybinding = false
			if isPrefix {
				if wasPrefix {
					datum[index].action.Write(line)
				} else {
					if datum[index].action.Len() != 0 {
						datum[index].action.Write([]byte("\n"))
					}
					datum[index].action.Write(line)
					wasPrefix = true
				}
				continue
			}

			if wasPrefix {
				datum[index].action.Write(line)
				wasPrefix = false
			} else {
				if datum[index].action.Len() != 0 {
					datum[index].action.Write([]byte("\n"))
				}
				datum[index].action.Write(line)
			}
		}
	}

	if err == io.EOF {
		err = nil
	} else {
		return
	}

	*data = nil
	for _, d := range datum {
		if d.hasVariant {
			replicated, e := replicate(d.binding.String(), d.action.String())
			if e != nil {
				err = errors.New(fmt.Sprintf("can't register %s keybinding, error (%s)", strings.TrimPrefix(d.binding.String(), "#"), e.Error()))
				return
			}
			for _, repl := range replicated {
				*data = append(*data, filedata{binding: repl.binding, action: repl.action})
			}
		} else {
			*data = append(*data, filedata{binding: d.binding, action: d.action})
		}
	}

	if len(*data) == 1 && ((*data)[0].action.String() == "" || (*data)[0].binding.String() == "") {
		err = errors.New("config file does not contain any binding/action")
		return
	}

	return
}

func replicate(binding, action string) (replicated []*filedata, err error) {
	// find all variants
	bindingVariants, actionVariants := variantPattern.FindAllString(binding, -1), variantPattern.FindAllString(action, -1)

	// make sure the amount of variants do match
	if len(bindingVariants) != len(actionVariants) {
		err = errors.New("the amount of variants in a keybinding and it's action do not match")
		return
	}

	var bindingVars, actionVars [][]string

	// extract variant members
	extract := func(from []string, where *[][]string) {
		for _, f := range from {
			*where = append(*where, strings.Split(strings.TrimSuffix(strings.TrimPrefix(f, "{"), "}"), ","))
		}
	}

	extract(bindingVariants, &bindingVars)
	extract(actionVariants, &actionVars)

	// validate the amount of variant memebers do match
	for i, b := range bindingVars {
		if len(b) != len(actionVars[i]) {
			err = errors.New("the amount of variant members in a keybinding and it's action do not match")
			return
		}
	}

	// validate and extract ranges
	var rngs []ranges
	rngs, err = extractRanges(bindingVars, actionVars)
	if err != nil {
		return
	}

	var expandedBindingRanges, expandedActionRanges []string

	// expands a range in a keybinding({1-9} -> {1},{2},{3},{...},{9})
	expandRange := func(r ranges, binding, acton string, bindings, actions *[]string) {
		// bindings
		for bIn := r.binding.start; bIn != r.binding.end+1; bIn++ {
			*bindings = append(*bindings, strings.Replace(
				binding,
				fmt.Sprintf("%s-%s", r.binding.startStr, r.binding.endStr),
				fmt.Sprintf("%c", rune(bIn)),
				1,
			))
			if r.action.skip {
				*actions = append(*actions, acton)
			}
		}
		// actions
		if !r.action.skip {
			for aIn := r.action.start; aIn != r.action.end+1; aIn++ {
				if r.action.numerical {
					*actions = append(*actions, strings.Replace(
						action,
						fmt.Sprintf("%s-%s", r.action.startStr, r.action.endStr),
						fmt.Sprintf("%d", aIn),
						1,
					))
				} else {
					*actions = append(*actions, strings.Replace(
						action,
						fmt.Sprintf("%s-%s", r.action.startStr, r.action.endStr),
						fmt.Sprintf("%c", rune(aIn)),
						1,
					))
				}
			}
		}
	}

	for len(rngs) > 0 {
		if len(expandedBindingRanges) > 0 {
			var newBindingRanges, newActionRanges []string

			if len(expandedActionRanges) != len(expandedBindingRanges) {
				err = errors.New("an unknown error occurred whilst expanding keybinding and action ranges")
			}

			for i := 0; i != len(expandedBindingRanges); i++ {
				expandRange(rngs[0], expandedBindingRanges[i], expandedActionRanges[i], &newBindingRanges, &newActionRanges)
			}

			expandedBindingRanges, expandedActionRanges = newBindingRanges, newActionRanges
		} else {
			expandRange(rngs[0], binding, action, &expandedBindingRanges, &expandedActionRanges)
		}
		rngs = rngs[1:]
	}

	if len(expandedActionRanges) != len(expandedBindingRanges) {
		err = errors.New("an unknown error occurred whilst expanding keybinding and action ranges")
	}

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

	for i, r := 0, 0; i != len(expandedBindingRanges); i++ {
		var replicatedBindings, replicatedActions []string
		vGroup := &variantGroup{}
		vGroup.action = variantPattern.FindAllString(expandedActionRanges[i], -1)
		vGroup.binding = variantPattern.FindAllString(expandedBindingRanges[i], -1)

		if !(len(vGroup.action) == len(vGroup.binding) && len(vGroup.action) > 0) {
			err = errors.New("can not extract variant groups")
			return
		}

		for len(vGroup.binding) > 0 {
			bVariantMembers := strings.Split(strings.TrimSuffix(strings.TrimPrefix(vGroup.binding[0], "{"), "}"), ",")
			aVariantMembers := strings.Split(strings.TrimSuffix(strings.TrimPrefix(vGroup.action[0], "{"), "}"), ",")
			if len(replicatedBindings) > 0 {
				var newBindingVariants, newActionVariants []string

				for _, alreadyR := range replicatedBindings {
					replicateVariant(alreadyR, vGroup.binding[0], bVariantMembers, &newBindingVariants)
				}

				for _, alreadyR := range replicatedActions {
					replicateVariant(alreadyR, vGroup.action[0], aVariantMembers, &newActionVariants)
				}

				replicatedBindings, replicatedActions = newBindingVariants, newActionVariants
			} else {
				replicateVariant(expandedBindingRanges[i], vGroup.binding[0], bVariantMembers, &replicatedBindings)
				replicateVariant(expandedActionRanges[i], vGroup.action[0], aVariantMembers, &replicatedActions)
			}
			vGroup.binding = vGroup.binding[1:]
			vGroup.action = vGroup.action[1:]
		}

		if len(replicatedBindings) != len(replicatedActions) {
			err = errors.New("replication went wrong")
			return
		}

	appender:
		for i := 0; i != len(replicatedBindings); i++ {
			replicatedBindings[i] = strings.ReplaceAll(replicatedBindings[i], "--", "-")
			if i > 0 {
				for _, aR := range replicated {

					if aR.binding.String() == replicatedBindings[i] {
						continue appender
					}
				}
			}
			replicated = append(replicated, &filedata{})
			replicated[r].binding.WriteString(replicatedBindings[i])
			replicated[r].action.WriteString(replicatedActions[i])
			r++
		}
	}

	return
}

func extractRanges(bindingVars, actionVars [][]string) (r []ranges, err error) {
	// range patterns for binding and action and range errors
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
					aVar             = actionVars[bIn][vIn]
					aRange           []string
					aRangeValidation = true
				)
				// make sure action variant is also a range (or _)
				if !actionRangePattern.MatchString(aVar) {
					if aVar == "_" {
						// in case it's _, skip the range validation
						aRangeValidation = false
					} else {
						err = errors.New("the indexes of ranges for a keybinding and it's action do not match")
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
						for _, a := range []rune(aRange[1]) {
							aStartStr += string(a)
						}
						for _, a := range []rune(aRange[2]) {
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
						err = errors.New("the ranges of a keybinding and it's action do not match")
						return
					}
				}
				r = append(r, ranges{})
				r[len(r)-1].binding.start = bStart
				r[len(r)-1].binding.startStr = bStartStr
				r[len(r)-1].binding.end = bEnd
				r[len(r)-1].binding.endStr = bEndStr

				r[len(r)-1].action.start = aStart
				r[len(r)-1].action.startStr = aStartStr
				r[len(r)-1].action.end = aEnd
				r[len(r)-1].action.endStr = aEndStr
				r[len(r)-1].action.skip = aRangeValidation == false
				r[len(r)-1].action.numerical = aNumerical
			}
		}
	}
	return
}
