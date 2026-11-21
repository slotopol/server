package slot

import (
	"fmt"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

const maxSafeFloat = 1 << 53 // limit of exact integers in the significand of a float64

type Number interface {
	~uint64 | ~float64
}

// MarshalYAML implements the yaml.Marshaler interface.
func (sp StatPage) MarshalYAML() (any, error) {
	var root = &yaml.Node{
		Kind: yaml.MappingNode,
	}

	var addPair = func(key string, keyStyle yaml.Style, valNode *yaml.Node) {
		root.Content = append(root.Content,
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Style: keyStyle,
				Value: key,
			},
			valNode,
		)
	}

	addPair("N", yaml.DoubleQuotedStyle, numNode(sp.N))
	addPair("C", 0, alignedMatrixNode(sp.C))
	addPair("S", 0, alignedMatrixNode(sp.S))

	if sp.Cfg != 0 {
		addPair("Cfg", 0, numNode(sp.Cfg))
	}
	if sp.Cfs != 0 {
		addPair("Cfs", 0, numNode(sp.Cfs))
	}
	if sp.Qfs != 0 {
		addPair("Qfs", 0, numNode(sp.Qfs))
	}
	if len(sp.Cb) > 0 {
		addPair("Cb", 0, sliceNode(sp.Cb))
	}
	if len(sp.Sb) > 0 {
		addPair("Sb", 0, sliceNode(sp.Sb))
	}
	if len(sp.Cj) > 0 {
		addPair("Cj", 0, sliceNode(sp.Cj))
	}

	return root, nil
}

func formatNum[T Number](v T) string {
	switch val := any(v).(type) {
	case uint64:
		return strconv.FormatUint(val, 10)
	case float64:
		if val >= maxSafeFloat {
			return strconv.FormatFloat(val, 'e', -1, 64)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

// numWidth calculates the column width based on the largest number
// it contains, without performing a float64 -> uint64 conversion.
func numWidth[T Number](maxVal T) int {
	switch val := any(maxVal).(type) {
	case uint64:
		var w = 1
		for val >= 10 {
			val /= 10
			w++
		}
		return w
	case float64:
		if val >= maxSafeFloat {
			// For the e-format, we return the actual record length for the
			// maximum number (with a buffer to accommodate a full 17-digit
			// mantissa, in case the maximum number has trailing zeros)
			return max(len(strconv.FormatFloat(val, 'e', -1, 64)), 22)
		}
		var w = 1
		for val >= 10 {
			val /= 10
			w++
		}
		return w
	default:
		return len(fmt.Sprint(maxVal))
	}
}

func numNode[T Number](v T) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: formatNum(v),
	}
}

func sliceNode[T Number](s []T) *yaml.Node {
	node := &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FlowStyle,
	}
	for _, v := range s {
		node.Content = append(node.Content, numNode(v))
	}
	return node
}

// alignedMatrixNode constructs a table with columns aligned to the right.
func alignedMatrixNode[T Number](m [][]T) *yaml.Node {
	if len(m) == 0 {
		return &yaml.Node{
			Kind:  yaml.SequenceNode,
			Style: yaml.FlowStyle,
		}
	}

	var node = &yaml.Node{
		Kind: yaml.SequenceNode,
	}

	var maxVals []T
	for _, row := range m {
		for j, val := range row {
			if j >= len(maxVals) {
				maxVals = append(maxVals, val)
			} else if val > maxVals[j] {
				maxVals[j] = val
			}
		}
	}

	var widths = make([]int, len(maxVals))
	for j, maxV := range maxVals {
		widths[j] = numWidth(maxV)
	}

	var sb strings.Builder
	for i, row := range m {
		sb.Reset()
		sb.WriteByte('[')
		for j, val := range row {
			if j > 0 {
				sb.WriteString(", ")
			}
			fmt.Fprintf(&sb, "%*s", widths[j], formatNum(val))
		}
		sb.WriteByte(']')

		node.Content = append(node.Content, &yaml.Node{
			Kind:        yaml.ScalarNode,
			Value:       sb.String(),
			LineComment: fmt.Sprintf("sym %d", i+1),
		})
	}

	return node
}

func statPageSliceNode(pages []*StatPage) (*yaml.Node, error) {
	seq := &yaml.Node{
		Kind: yaml.SequenceNode,
	}

	for i, p := range pages {
		if p.N == 0 {
			break
		}
		if p == nil {
			seq.Content = append(seq.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!null",
				Value: "null",
			})
			continue
		}

		var raw, err = p.MarshalYAML()
		if err != nil {
			return nil, err
		}

		var elemNode, ok = raw.(*yaml.Node)
		if !ok {
			return nil, fmt.Errorf("expected *yaml.Node, got %T", raw)
		}
		if i > 0 {
			elemNode.HeadComment = fmt.Sprintf("cascade %d", i)
		}

		seq.Content = append(seq.Content, elemNode)
	}

	return seq, nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (ss StatSnap) MarshalYAML() (any, error) {
	var root = &yaml.Node{
		Kind: yaml.MappingNode,
	}

	var addPair = func(key string, keyStyle yaml.Style, valNode *yaml.Node) {
		root.Content = append(root.Content,
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Style: keyStyle,
				Value: key,
			},
			valNode,
		)
	}

	if ss.EC > 0 {
		addPair("EC", 0, numNode(ss.EC))
	}

	if ss.Q > 0 {
		addPair("Q", 0, numNode(ss.Q))
	}

	if len(ss.A) > 0 {
		var valNode, err = statPageSliceNode(ss.A)
		if err != nil {
			return nil, err
		}
		addPair("A", 0, valNode)
	}

	return root, nil
}
