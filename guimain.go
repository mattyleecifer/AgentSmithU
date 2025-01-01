package main

import (
	. "AgentSmithU/agent"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
)

func gui(agent *Agent) {
	http.HandleFunc("/", RequireAuth(index))
	http.HandleFunc("/auth/", hauth)
	http.HandleFunc("/chat/", RequireAuth(hchat(agent)))
	http.HandleFunc("/chat/edit/", RequireAuth(chatedit(agent)))
	http.HandleFunc("/chat/save/", RequireAuth(hchatsave(agent)))
	http.HandleFunc("/chat/data/", RequireAuth(hchatdata(agent)))
	http.HandleFunc("/chat/clear/", RequireAuth(hchatclear(agent)))
	http.HandleFunc("/chat/reset/", RequireAuth(hreset(agent)))
	http.HandleFunc("/settings/", RequireAuth(hsettings(agent)))
	http.HandleFunc("/sidebar/", RequireAuth(hsidebar))
	http.HandleFunc("/tokenupdate/", RequireAuth(htokenupdate(agent)))
	http.HandleFunc("/prompt/", RequireAuth(hprompt(agent)))
	http.HandleFunc("/prompt/data/", RequireAuth(hpromptdata(agent)))
	http.HandleFunc("/function/", RequireAuth(hfunction(agent)))
	http.HandleFunc("/function/data/", RequireAuth(hfunctiondata(agent)))
	http.Handle("/static/", http.FileServer(http.FS(hcss)))
	fmt.Println("Running GUI on http://127.0.0.1"+port, "(ctrl-click link to open)")
	log.Fatal(http.ListenAndServe(port, nil))
	// log.Fatal(http.ListenAndServeTLS(port, "certificate.crt", "private.key", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	render(w, hindexpage, nil)
}

func render(w http.ResponseWriter, html string, data any) {
	// Render the HTML template
	// fmt.Println("Rendering...")
	w.WriteHeader(http.StatusOK)
	tmpl, err := template.New(html).Parse(html)
	if err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// allowedIps = append(allowedIps, GetLocalIP(), "127.0.0.1")

		// allowedIps = append(allowedIps, "127.0.0.1")
		// fmt.Println("\nAllowed ips: ", allowedIps)
		// Get the IP address of the client
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// fmt.Println("\nConnecting IP: ", ip)
		// Check if the client's IP is in the list of allowed IP
		if allowAllIps {
			handler.ServeHTTP(w, r)
			return
		} else {
			for _, allowedIp := range allowedIps {
				if ip == allowedIp {
					// If the client's IP is in the list of allowed IPs, allow access to the proxy server
					handler.ServeHTTP(w, r)
					return
				}
			}
		}

		if authstring != "" {
			hauth(w, r)
		}

		// If the client's IP is not in the list of allowed IPs, return a 403 Forbidden error
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("HX-Redirect", `/auth/`)
	}
}

func hauth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth := r.FormValue("auth")
		if auth == authstring {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			allowedIps = append(allowedIps, ip)
		}
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		render(w, hauthpage, nil)
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
