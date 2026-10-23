package menus

import (
	"fmt"
	"strconv"
	"strings"
)

func NewTransform(rt *Runtime) {
	fmt.Println("ops: trim, upper, replace <old> <new>")
	name, err := rt.prompt("name-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if name == "" {
		fmt.Println("name is required")
		return
	}
	raw, err := rt.prompt("op-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	op, args, err := parseTransformOp(raw)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	rt.Transforms = append(rt.Transforms, Transform{Name: name, Op: op, Args: args})
	if err := SaveTransforms(rt.dataDir(), rt.Transforms); err != nil {
		fmt.Println("error:", err)
	}
}

func ShowTransforms(rt *Runtime) {
	transforms, err := LoadTransforms(rt.dataDir())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	rt.Transforms = transforms
	fmt.Print(formatTransforms(transforms))
}

func DeleteTransform(rt *Runtime) {
	if len(rt.Transforms) == 0 {
		fmt.Println("no transforms")
		return
	}
	fmt.Print(formatTransforms(rt.Transforms))
	input, err := rt.prompt("transform-> ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Transforms) {
		fmt.Println("invalid transform")
		return
	}
	rt.Transforms = append(rt.Transforms[:n-1], rt.Transforms[n:]...)
	if err := SaveTransforms(rt.dataDir(), rt.Transforms); err != nil {
		fmt.Println("error:", err)
	}
}

func RunTransform(rt *Runtime) {
	body := ""
	if rt.State != nil {
		body = rt.State["last.body"]
	}
	if body == "" {
		fmt.Println("no body to transform")
		return
	}
	transform, ok := pickTransform(rt)
	if !ok {
		return
	}
	updated, err := applyTransform(transform, body)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if rt.State == nil {
		rt.State = map[string]string{}
	}
	rt.State["last.body"] = updated
	rt.State["last.transform"] = transform.Name
	if err := SaveState(rt.dataDir(), rt.State); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("\n%s\n", updated)
}

func pickTransform(rt *Runtime) (Transform, bool) {
	if len(rt.Transforms) == 0 {
		fmt.Println("no transforms")
		return Transform{}, false
	}
	fmt.Print(formatTransforms(rt.Transforms))
	input, err := rt.prompt("transform-> ")
	if err != nil {
		fmt.Println("error:", err)
		return Transform{}, false
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(rt.Transforms) {
		fmt.Println("invalid transform")
		return Transform{}, false
	}
	return rt.Transforms[n-1], true
}

func parseTransformOp(raw string) (string, []string, error) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return "", nil, fmt.Errorf("op is required")
	}
	switch fields[0] {
	case "trim", "upper":
		if len(fields) != 1 {
			return "", nil, fmt.Errorf("%s takes no arguments", fields[0])
		}
		return fields[0], nil, nil
	case "replace":
		if len(fields) < 3 {
			return "", nil, fmt.Errorf("replace needs old and new")
		}
		return "replace", []string{fields[1], strings.Join(fields[2:], " ")}, nil
	default:
		return "", nil, fmt.Errorf("unknown transform %q", fields[0])
	}
}

func applyTransform(transform Transform, body string) (string, error) {
	switch transform.Op {
	case "trim":
		return strings.TrimSpace(body), nil
	case "upper":
		return strings.ToUpper(body), nil
	case "replace":
		if len(transform.Args) != 2 {
			return "", fmt.Errorf("replace needs old and new")
		}
		return strings.ReplaceAll(body, transform.Args[0], transform.Args[1]), nil
	default:
		return "", fmt.Errorf("unknown transform %q", transform.Op)
	}
}
