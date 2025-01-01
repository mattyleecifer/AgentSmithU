package main

// create, edit, and make prompt calls - this will allow users to make commandline or api promptcalls.

import (
	. "AgentSmithU/agent"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func hprompt(agent *Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			var data struct {
				Name         string
				Description  string
				Parameters   string
				Savedprompts []string
			}

			data.Name = agent.Prompt.Name
			data.Description = agent.Prompt.Description
			data.Parameters = agent.Messages[0].Content
			data.Savedprompts, _ = getsavefilelist("Prompts")

			render(w, hpromptspage, data)
		}

		if r.Method == http.MethodPost {
			newprompt := PromptDefinition{
				Name:        r.FormValue("promptname"),
				Description: r.FormValue("promptdescription"),
				Parameters:  r.FormValue("edittext"),
			}

			agent.Prompt = newprompt
			agent.Setprompt()

			w.Header().Set("HX-Redirect", "/")
			// r.Method = http.MethodGet
			// agent.hchat(w, r)
		}
	}
}

func hpromptdata(agent *Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimPrefix(r.URL.Path, "/prompt/data/")

		if r.Method == http.MethodGet {
			var data struct {
				Name         string
				Description  string
				Parameters   string
				Savedprompts []string
			}

			prompt := PromptDefinition{}

			loaddata, err := loadfile(agent, "Prompts", query)
			if err != nil {
				fmt.Println(err)
			}

			_ = json.Unmarshal(loaddata, &prompt)

			data.Name = prompt.Name
			data.Description = prompt.Description
			data.Parameters = prompt.Parameters
			data.Savedprompts, _ = getsavefilelist("Prompts")

			render(w, hpromptspage, data)
		}

		if r.Method == http.MethodPost {
			newprompt := PromptDefinition{
				Name:        r.FormValue("promptname"),
				Description: r.FormValue("promptdescription"),
				Parameters:  r.FormValue("edittext"),
			}

			savefile(newprompt, "Prompts", newprompt.Name)

			htmldata := `
		<div id="prompt-` + newprompt.Name + `" hx-swap-oob="delete"></div>
		<div id="prompt-` + newprompt.Name + `" style="display: flex;">
			<div style="text-align: left; float: left;">` + newprompt.Name + `</div>
			<div style="float: right; margin-left: auto; display: inline;">
				<button hx-target='#main-content' hx-get='/prompt/data/` + newprompt.Name + `'>Load</button>
				<button hx-target='#prompt-` + newprompt.Name + `' hx-delete='/prompt/data/` + newprompt.Name + `' hx-swap='delete' hx-confirm='Are you sure?'>Delete</button>
			</div>
		</div>`
			// this should actually pop up a new row in the saves list with the new save
			// same with functions - like chats
			render(w, htmldata, nil)
		}

		if r.Method == http.MethodDelete {
			deletefile("Prompts", query)
		}
	}
}
