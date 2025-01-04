package gui

// create, edit, and make function calls

import (
	"AgentSmithU/agent"
	"AgentSmithU/config"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func hfunction(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			type Funcdef struct {
				Name        string
				Description string
			}

			var data struct {
				Currentfunctions []Funcdef
				Savedfunctions   []string
			}

			for _, item := range ag.Functions {
				newfunc := Funcdef{
					Name:        item.Name,
					Description: item.Description,
				}
				data.Currentfunctions = append(data.Currentfunctions, newfunc)
			}

			data.Savedfunctions, _ = config.GetSaveFileList("Functions")
			render(w, hfunctionpage, data)
		}

		if r.Method == http.MethodPost {
			newfunction := agent.Function{
				Name:        r.FormValue("functionname"),
				Description: r.FormValue("functiondescription"),
				Parameters:  r.FormValue("edittext"),
			}

			ag.AddFunction(newfunction)

			w.Header().Set("HX-Redirect", "/function/")
			// r.Method = http.MethodGet
			// agent.hfunction(w, r)
		}

		query := strings.TrimPrefix(r.URL.Path, "/function/")

		if r.Method == http.MethodPatch {
			var data agent.Function
			if query != "" {
				for _, function := range ag.Functions {
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
			ag.RemoveFunction(query)

			// looks like this reloads the page, edit to make it not
			r.Method = http.MethodGet
			w.Header().Set("HX-Redirect", "/function/")
			// agent.hfunction(w, r)
		}
	}
}

func hfunctiondata(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimPrefix(r.URL.Path, "/function/data/")

		if r.Method == http.MethodGet {
			var newfunction agent.Function
			functionname := query
			filedata, err := config.Load(ag, "Functions", functionname)
			if err != nil {
				fmt.Println(err)
			}
			err = json.Unmarshal(filedata, &newfunction)
			if err != nil {
				// if function doesn't exist then don't do anything
				// This should return an error - similar to prompts
				return
			}

			ag.AddFunction(newfunction)
			// agent.hfunction(w, r)
			w.Header().Set("HX-Redirect", "/function/")
		}

		if r.Method == http.MethodPost {
			newfunction := agent.Function{
				Name:        r.FormValue("functionname"),
				Description: r.FormValue("functiondescription"),
				Parameters:  r.FormValue("edittext"),
			}

			config.Save(newfunction, "Functions", newfunction.Name)

			// should add to page rather than reload like in prompts
			// reloads page
			r.Method = http.MethodGet
			w.Header().Set("HX-Redirect", "/function/")
			// agent.hfunction(w, r)
		}

		if r.Method == http.MethodDelete {
			functionname := query
			config.Delete("Functions", functionname)

			// reloads page
			w.Header().Set("HX-Redirect", "/function/")
			// agent.hfunction(w, r)
		}
	}
}

func hfunctionrun(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawquery := strings.TrimPrefix(r.URL.Path, "/function/run/")
		query := strings.Split(rawquery, "/")
		fmt.Println(rawquery, query)

		var function agent.Function

		for _, f := range ag.Functions {
			if f.Name == query[0] {
				function = f
				break
			}
		}

		if function.Name == "" {
			return
		}

		response := ag.RunFunction(function)

		w.Header().Set("HX-Trigger-After-Settle", `tokenupdate`)

		var data struct {
			Header   template.HTML
			Role     string
			Content  string
			Index    string
			Function string
		}
		data.Header = template.HTML(`<div id="message" class="message" style="background-color: #393939">`)
		data.Role = agent.RoleAssistant
		data.Content = response.Content
		data.Index = strconv.Itoa(len(ag.Messages) - 1)
		render(w, hchatnewpage, data)
	}
}
