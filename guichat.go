package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (agent *Agent) hchat(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Header   template.HTML
		Role     string
		Content  string
		Index    string
		Function template.HTML
	}

	if r.Method == http.MethodGet {
		type message struct {
			Role    string
			Content template.HTML
			Index   int
			Header  string
		}
		var data struct {
			Messages []message
		}

		messages := agent.Messages
		if len(messages) == 1 {
			render(w, hchatpage, data)
		} else {
			for i, item := range messages[1:] {
				var content string
				lines := strings.Split(item.Content, "\n")
				for _, line := range lines {
					content += line + "<br>"
				}
				index := strconv.Itoa(i + 1)
				content = `<div class="agent">` + item.Role + `</div>
					<div id="reply-` + index + `" class="content">
						<pre style="white-space: pre-wrap; font-family: inherit;">` + item.Content + `</pre>
					</div>
					<div class="editbutton">
						<button hx-get="/chat/edit/` + index + `" hx-target="#reply-` + index + `">Edit</button>
						<button hx-delete="/chat/edit/` + index + `" hx-target="closest .message">Delete</button>
					</div>`
				msg := message{
					Content: template.HTML(content),
				}
				if item.Role == "assistant" {
					msg.Content = template.HTML(`<div style="display: flex; width: 100%; background-color: #393939">` + content + `</div>`)
				}

				data.Messages = append(data.Messages, msg)
			}
			render(w, hchatpage, data)
		}
	}

	if r.Method == http.MethodPost {
		rawtext := r.FormValue("text")
		if strings.TrimSpace(rawtext) == "" {
			render(w, "", nil)
			return
		}
		if strings.TrimSpace(rawtext) == "!" {
			agent.setprompt()
			render(w, `<div id="message" hx-target="#main-content" hx-post="/chat/clear/" hx-trigger="load"></div>`, nil)
			return
		}
		query := Message{
			Role:    RoleUser,
			Content: rawtext,
		}
		agent.Messages = append(agent.Messages, query)
		// text := agent.Messages[len(agent.Messages)-1].Content

		data.Header = template.HTML(`<div id="message" class="message">`)
		data.Role = RoleUser
		data.Content = rawtext
		data.Index = strconv.Itoa(len(agent.Messages) - 1)

		render(w, hchatnewpage, data)
	}

	if r.Method == http.MethodPut {
		response, err := agent.getresponse()
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("HX-Trigger-After-Settle", `tokenupdate`)

		data.Role = RoleAssistant
		data.Header = template.HTML(`<div id="message" class="message" style="background-color: #393939">`)

		data.Content = response.Content
		data.Index = strconv.Itoa(len(agent.Messages) - 1)

		render(w, hchatnewpage, data)
	}
}

func (agent *Agent) hchatedit(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimPrefix(r.URL.Path, "/chat/edit/")
	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(query)
		if err != nil {
			fmt.Println(err)
		}
		data := struct {
			Edittext  string
			MessageID int
		}{
			Edittext:  agent.Messages[id].Content,
			MessageID: id,
		}
		render(w, hchatedit, data)
	}

	if r.Method == http.MethodPost {
		id, err := strconv.Atoi(query)
		if err != nil {
			fmt.Println(err)
		}
		edittext := r.FormValue("edittext")
		agent.Messages[id].Content = edittext
		newtext := `<pre style="white-space: pre-wrap; font-family: inherit;">` + edittext + `</pre>`
		render(w, newtext, nil)
	}

	if r.Method == http.MethodDelete {
		err := agent.deletelines(query)
		if err != nil {
			fmt.Println(err)
		}
		// r.Method = http.MethodGet
		// agent.hchat(w, r)
	}
}

func (agent *Agent) hchatsave(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		currentTime := time.Now()
		filename := currentTime.Format("20060102150405")
		data := struct {
			Filename string
		}{
			Filename: filename,
		}
		render(w, hchatsavepage, data)
	}

	if r.Method == http.MethodPost {
		filename := r.FormValue("filename")
		agent.savefile(agent.Messages, "Chats", filename)
		render(w, "Chat Saved!", nil)
	}
}

func (agent *Agent) hchatdata(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimPrefix(r.URL.Path, "/chat/data/")
	if r.Method == http.MethodGet {
		if query == "" {
			var data struct {
				Filelist []string
			}
			filelist, err := getsavefilelist("Chats")
			if err != nil {
				fmt.Println(err)
			}
			data.Filelist = filelist
			render(w, hchatfilespage, data)
		} else {
			_, err := agent.loadfile("Chats", query)
			if err != nil {
				fmt.Println(err)
			}
			r.Method = http.MethodGet
			agent.hchat(w, r)
		}

	}

	if r.Method == http.MethodDelete {
		err := deletefile("Chats", query)
		if err != nil {
			fmt.Println(err)
		}
		render(w, "<p>Chat Deleted</p>", nil)
	}
}

func (agent *Agent) hchatclear(w http.ResponseWriter, r *http.Request) {
	agent.setprompt()
	r.Method = http.MethodGet
	agent.hchat(w, r)
}

func (agent *Agent) hreset(w http.ResponseWriter, r *http.Request) {
	agent.reset()
	r.Method = http.MethodGet
	agent.hchat(w, r)
}

func (agent *Agent) hsettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := struct {
			Model     string
			Maxtokens int
			Callcost  float64
		}{
			Model:     agent.model,
			Maxtokens: maxtokens,
			Callcost:  callcost,
		}
		render(w, hsettingspage, data)
	}
	if r.Method == http.MethodPut {
		apikey := r.FormValue("apikey")
		if apikey != "" {
			agent.api_key = apikey
		}
		agent.model = r.FormValue("chatmodel")
		maxtokens, _ = strconv.Atoi(r.FormValue("maxtokens"))
		callcost, _ = strconv.ParseFloat(r.FormValue("callcost"), 64)

		r.Method = http.MethodGet
		agent.hchat(w, r)
	}
}

func hsidebar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("HX-Trigger-After-Settle", `tokenupdate`)
		render(w, hsidebarpage, nil)
	}
	if r.Method == http.MethodDelete {
		button := `<div class="sidebar" id="sidebar" style="width: 0; background-color: transparent;"><button id="floating-button" hx-get="/sidebar/" hx-target="#sidebar" hx-swap="outerHTML">Show Menu</button></div>`
		render(w, button, nil)
	}
}

func (agent *Agent) htokenupdate(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("htokenupdate")
	estcost := (float64(agent.tokencount) / 1000) * callcost
	tokencount := strconv.Itoa(agent.tokencount)
	estcoststr := strconv.FormatFloat(estcost, 'f', 6, 64)
	render(w, "#Tokens: "+tokencount+"<br>$Est: "+estcoststr, nil)
}
