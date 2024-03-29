package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
)

func (agent *Agent) gui() {
	http.HandleFunc("/", RequireAuth(index))
	http.HandleFunc("/auth/", hauth)
	http.HandleFunc("/chat/", RequireAuth(agent.hchat))
	http.HandleFunc("/chat/edit/", RequireAuth(agent.hchatedit))
	http.HandleFunc("/chat/save/", RequireAuth(agent.hchatsave))
	http.HandleFunc("/chat/data/", RequireAuth(agent.hchatdata))
	http.HandleFunc("/chat/clear/", RequireAuth(agent.hchatclear))
	http.HandleFunc("/chat/reset/", RequireAuth(agent.hreset))
	http.HandleFunc("/settings/", RequireAuth(agent.hsettings))
	http.HandleFunc("/sidebar/", RequireAuth(hsidebar))
	http.HandleFunc("/tokenupdate/", RequireAuth(agent.htokenupdate))
	http.HandleFunc("/prompt/", RequireAuth(agent.hprompt))
	http.HandleFunc("/prompt/data/", RequireAuth(agent.hpromptdata))

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
