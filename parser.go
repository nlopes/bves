package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/xmlpath.v2"
)

type VBNode struct {
	Name     string  `json:"name"`
	Provider string  `json:"provider"`
	URL      string  `json:"url"`
	Size     float64 `json:"size"`
	Id       int     `json:"id"`
}

func parse_vbes_html(body []byte) (nodes []VBNode, err error) {
	xmlroot, err := xmlpath.ParseHTML(bytes.NewReader(body))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse html: %s", err))
	}
	elems := xmlpath.MustCompile("//table[@id='dataTable']/tbody//td")

	iter := elems.Iter(xmlroot)
	node := VBNode{}
	for i, j := 0, 1; iter.Next(); i++ {
		mod := i % 4
		switch {
		case mod == 0:
			node.Name = strings.TrimSpace(iter.Node().String())
		case mod == 1:
			node.Provider = strings.TrimSpace(iter.Node().String())
		case mod == 2:
			node.URL = strings.TrimSpace(iter.Node().String())
		case mod == 3:
			node.Size, err = strconv.ParseFloat(strings.TrimSpace(iter.Node().String()), 32)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Could not parse size %s: %s",
					iter.Node().String(), err))
			}
			node.Id = j
			j++
			nodes = append(nodes, node)
			node = VBNode{}
		}
	}
	return nodes, nil
}
