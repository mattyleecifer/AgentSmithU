package gui

import (
	"agentsmithu/agent"
	"agentsmithu/config"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func chat(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			// remove empty messages
			ag.Messages.Deletelines()

			// Check what to display
			if len(ag.Messages) == 1 {
				// If only system prompt, show the empty page
				render(w, chatpage, data)
			} else {
				// Display existing messages
				for i, item := range ag.Messages[1:] {
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
						<button hx-delete="/chat/edit/` + index + `" hx-swap="outerHTML" hx-target="closest .message">Delete</button>
					</div>`
					msg := message{
						Content: template.HTML(content),
					}
					if item.Role == "assistant" {
						msg.Content = template.HTML(`<div style="display: flex; width: 100%; background-color: #393939">` + content + `</div>`)
					}

					data.Messages = append(data.Messages, msg)
				}
				render(w, chatpage, data)
			}
		}

		if r.Method == http.MethodPost {
			rawtext := r.FormValue("text")
			if strings.TrimSpace(rawtext) == "" {
				render(w, "", nil)
				return
			}
			if strings.TrimSpace(rawtext) == "!" {
				ag.Setprompt()
				render(w, `<div id="message" hx-target="#main-content" hx-post="/chat/clear/" hx-trigger="load"></div>`, nil)
				return
			}
			query := agent.Message{
				Role:    agent.RoleUser,
				Content: rawtext,
			}
			ag.Messages = append(ag.Messages, query)
			// text := agent.Messages[len(agent.Messages)-1].Content

			data.Header = template.HTML(`<div id="message" class="message">`)
			data.Role = agent.RoleUser
			data.Content = rawtext
			data.Index = strconv.Itoa(len(ag.Messages) - 1)

			render(w, chatnewpage, data)
		}

		if r.Method == http.MethodPut {
			// Get agent response
			response, err := ag.Getresponse()
			if err != nil {
				fmt.Println(err)
			}

			w.Header().Set("HX-Trigger-After-Settle", `tokenupdate`)

			data.Role = agent.RoleAssistant
			data.Header = template.HTML(`<div id="message" class="message" style="background-color: #393939">`)

			data.Content = response.Content
			data.Index = strconv.Itoa(len(ag.Messages) - 1)

			render(w, chatnewpage, data)
		}
	}
}

func chatedit(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				Edittext:  ag.Messages[id].Content,
				MessageID: id,
			}
			render(w, chateditpage, data)
		}

		if r.Method == http.MethodPost {
			id, err := strconv.Atoi(query)
			if err != nil {
				fmt.Println(err)
			}
			edittext := r.FormValue("edittext")
			ag.Messages[id].Content = edittext
			newtext := `<pre style="white-space: pre-wrap; font-family: inherit;">` + edittext + `</pre>`
			render(w, newtext, nil)
		}

		if r.Method == http.MethodDelete {
			err := ag.Messages.Clearlines(query)
			if err != nil {
				fmt.Println(err)
			}
			// don't rerender page
			// r.Method = http.MethodGet
			// agent.hchat(w, r)
		}
	}
}

func chatsave(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			currentTime := time.Now()
			filename := currentTime.Format("20060102150405")
			data := struct {
				Filename string
			}{
				Filename: filename,
			}
			render(w, chatsavepage, data)
		}

		if r.Method == http.MethodPost {
			filename := r.FormValue("filename")
			config.Save(ag.Messages, "Chats", filename)
			render(w, "Chat Saved!", nil)
		}
	}
}

func chatdata(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimPrefix(r.URL.Path, "/chat/data/")
		if r.Method == http.MethodGet {
			if query == "" {
				var data struct {
					Filelist []string
				}
				filelist, err := config.GetSaveFileList("Chats")
				if err != nil {
					fmt.Println(err)
				}
				data.Filelist = filelist
				render(w, chatfilespage, data)
			} else {
				_, err := config.Load(ag, "Chats", query)
				if err != nil {
					fmt.Println(err)
				}
				r.Method = http.MethodGet
				// agent.hchat(w, r)
				w.Header().Set("HX-Redirect", "/")
			}

		}

		if r.Method == http.MethodDelete {
			err := config.Delete("Chats", query)
			if err != nil {
				fmt.Println(err)
			}
			render(w, "<p>Chat Deleted</p>", nil)
		}
	}
}

func chatclear(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newChat := []agent.Message{}
		newChat = append(newChat, ag.Messages[0])
		ag.Messages = newChat
		r.Method = http.MethodGet
		// agent.hchat(w, r)
		w.Header().Set("HX-Redirect", "/")
	}
}

func reset(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config.Reset(ag)
		r.Method = http.MethodGet
		// agent.hchat(w, r)
		w.Header().Set("HX-Redirect", "/")
	}
}

func settings(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			data := struct {
				Model     string
				Maxtokens int
				Callcost  float64
				ModelURL  string
			}{
				Model:     ag.Model,
				Maxtokens: ag.Maxtokens,
				Callcost:  config.CallCost,
				ModelURL:  ag.Modelurl,
			}
			render(w, settingspage, data)
		}
		if r.Method == http.MethodPut {
			apikey := r.FormValue("apikey")
			if apikey != "" {
				ag.Api_key = apikey
			}
			ag.Model = r.FormValue("chatmodel")
			ag.Maxtokens, _ = strconv.Atoi(r.FormValue("maxtokens"))
			config.CallCost, _ = strconv.ParseFloat(r.FormValue("callcost"), 64)
			ag.Modelurl = r.FormValue("modelurl")

			r.Method = http.MethodGet
			// agent.hchat(w, r)
			w.Header().Set("HX-Redirect", "/")
		}
	}
}

func sidebar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("HX-Trigger-After-Settle", `tokenupdate`)
		render(w, sidebarpage, nil)
	}
	if r.Method == http.MethodDelete {
		button := `<div class="sidebar" id="sidebar" style="width: 0; background-color: transparent;"><button id="floating-button" hx-get="/sidebar/" hx-target="#sidebar" hx-swap="outerHTML">Show Menu</button></div>`
		render(w, button, nil)
	}
}

func tokenupdate(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("htokenupdate")
		estcost := (float64(ag.Tokencount) / 1000) * config.CallCost
		tokencount := strconv.Itoa(ag.Tokencount)
		estcoststr := strconv.FormatFloat(estcost, 'f', 6, 64)
		render(w, "#Tokens: "+tokencount+"<br>$Est: "+estcoststr, nil)
	}
}
