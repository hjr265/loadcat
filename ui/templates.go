// Copyright 2015 The Loadcat Authors. All rights reserved.

package ui

import (
	"html/template"
)

var (
	TplLayout = template.Must(template.New("layout.html").ParseFiles("ui/templates/layout.html"))

	TplBalancerList     = template.Must(template.Must(TplLayout.Clone()).ParseFiles("ui/templates/balancerList.html"))
	TplBalancerNewForm  = template.Must(template.Must(TplLayout.Clone()).ParseFiles("ui/templates/balancerNewForm.html"))
	TplBalancerView     = template.Must(template.Must(TplLayout.Clone()).ParseFiles("ui/templates/balancerView.html"))
	TplBalancerEditForm = template.Must(template.Must(TplLayout.Clone()).ParseFiles("ui/templates/balancerEditForm.html"))

	TplServerNewForm  = template.Must(template.Must(TplLayout.Clone()).ParseFiles("ui/templates/serverNewForm.html"))
	TplServerEditForm = template.Must(template.Must(TplLayout.Clone()).ParseFiles("ui/templates/serverEditForm.html"))
)
