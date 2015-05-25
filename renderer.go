// Copyright 2015 Rentabiliweb Europe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package be2bill

import (
	"bytes"
	"fmt"
	"html/template"
)

type Renderer interface {
	Render(params, options Options) string
}

const (
	formPath        = "/front/form/process"
	defaultEncoding = "UTF-8"

	formTemplate = `<form method="post" action="{{.Url}}"{{range $name, $value := .Attributes}} {{name $name}}="{{$value}}"{{end}}>{{template "hidden" .Hidden}}
  {{template "submit" .Submit}}
</form>`

	hiddenTemplate = `{{range $name, $value := .}}
  <input type="hidden" name="{{$name}}" value="{{$value}}" />{{end}}`

	submitTemplate = `<input type="submit"{{range $name, $value := .}} {{name $name}}="{{$value}}"{{end}} />`
)

type templateContents struct {
	Url        string
	Attributes Options
	Hidden     Options
	Submit     Options
}

type htmlRenderer struct {
	url      string
	encoding string
}

func newHTMLRenderer(url string) Renderer {
	return &htmlRenderer{
		url:      url + formPath,
		encoding: defaultEncoding,
	}
}

func (p *htmlRenderer) Render(params, htmlOptions Options) string {
	funcMap := template.FuncMap{
		"name": safeHTMLAttributeName,
	}
	formTpl := template.Must(template.New("form").Funcs(funcMap).Parse(formTemplate))
	template.Must(formTpl.New("hidden").Parse(hiddenTemplate))
	template.Must(formTpl.New("submit").Parse(submitTemplate))

	var data templateContents

	// url
	data.Url = p.url

	// form attributes
	if formOptions, ok := htmlOptions[HTMLOptionForm].(Options); ok {
		data.Attributes = formOptions
	} else {
		data.Attributes = make(Options)
	}

	// hidden input fields
	data.Hidden = make(Options)
	for name, value := range params {
		if valueMap, ok := value.(Options); ok {
			for _, k := range valueMap.sortedKeys() {
				data.Hidden[fmt.Sprintf("%s[%s]", name, k)] = valueMap[k]
			}
		} else {
			data.Hidden[name] = value
		}
	}

	// submit input attributes
	if submitOptions, ok := htmlOptions[HTMLOptionSubmit].(Options); ok {
		data.Submit = submitOptions
	} else {
		data.Submit = make(Options)
	}

	// render
	var buf bytes.Buffer
	formTpl.Execute(&buf, data)

	// return
	return buf.String()
}

func safeHTMLAttributeName(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}
