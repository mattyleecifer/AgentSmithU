package main

// create, edit, and make function calls

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func (agent *Agent) hfunction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		type Funcdef struct {
			Name        string
			Description string
		}

		var data struct {
			Currentfunctions []Funcdef
			Savedfunctions   []string
		}

		for _, item := range agent.Functions {
			newfunc := Funcdef{
				Name:        item.Name,
				Description: item.Description,
			}
			data.Currentfunctions = append(data.Currentfunctions, newfunc)
		}

		data.Savedfunctions, _ = getsavefilelist("Functions")
		render(w, hfunctionpage, data)
	}

	if r.Method == http.MethodPost {
		newfunction := Function{
			Name:        r.FormValue("functionname"),
			Description: r.FormValue("functiondescription"),
			Parameters:  r.FormValue("edittext"),
		}

		agent.addFunction(newfunction)

		r.Method = http.MethodGet
		agent.hfunction(w, r)
	}

	query := strings.TrimPrefix(r.URL.Path, "/function/")

	if r.Method == http.MethodPatch {
		var data Function
		if query != "" {
			for _, function := range agent.Functions {
				if query == function.Name {
					data.Name = function.Name
					data.Description = function.Description
					data.Parameters = function.Parameters
					break
				}
			}
		}
		render(w, hfunctioneditpage, data)
	}

	if r.Method == http.MethodDelete {
		agent.removeFunction(query)

		// looks like this reloads the page, edit to make it not
		r.Method = http.MethodGet
		agent.hfunction(w, r)
	}
}

func (agent *Agent) hfunctiondata(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimPrefix(r.URL.Path, "/function/data/")

	if r.Method == http.MethodGet {
		functionname := query
		newfunction, err := agent.loadFunction(functionname)
		if err != nil {
			fmt.Println(err)
		}

		agent.addFunction(newfunction)
		agent.hfunction(w, r)
	}

	if r.Method == http.MethodPost {
		newfunction := Function{
			Name:        r.FormValue("functionname"),
			Description: r.FormValue("functiondescription"),
			Parameters:  r.FormValue("edittext"),
		}

		agent.savefile(newfunction, "Functions", newfunction.Name)

		// reloads page
		r.Method = http.MethodGet
		agent.hfunction(w, r)
	}

	if r.Method == http.MethodDelete {
		functionname := query
		deletefile("Functions", functionname)

		// reloads page
		agent.hfunction(w, r)
	}
}

func (agent *Agent) hfunctionrun(w http.ResponseWriter, r *http.Request) {
	rawquery := strings.TrimPrefix(r.URL.Path, "/function/run/")
	query := strings.Split(rawquery, "/")
	fmt.Println(rawquery, query)

	var function Function

	for _, f := range agent.Functions {
		if f.Name == query[0] {
			function = f
			break
		}
	}

	if function.Name == "" {
		return
	}

	response := agent.runFunction(function)

	w.Header().Set("HX-Trigger-After-Settle", `tokenupdate`)

	var data struct {
		Header   template.HTML
		Role     string
		Content  string
		Index    string
		Function string
	}
	data.Header = template.HTML(`<div id="message" class="message" style="background-color: #393939">`)
	data.Role = RoleAssistant
	data.Content = response.Content
	data.Index = strconv.Itoa(len(agent.Messages) - 1)
	render(w, hchatnewpage, data)
}
