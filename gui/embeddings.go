package gui

// Convert pages to strings otherwise the render() function won't work for both generated and templated html

import "embed"

//go:embed templates/index.html
var indexpage string

//go:embed templates/chat.html
var chatpage string

//go:embed templates/chatnew.html
var chatnewpage string

//go:embed templates/chatedit.html
var chateditpage string

//go:embed templates/chatsave.html
var chatsavepage string

//go:embed templates/chatfiles.html
var chatfilespage string

//go:embed templates/settings.html
var settingspage string

//go:embed templates/sidebar.html
var sidebarpage string

//go:embed templates/prompt.html
var promptspage string

//go:embed templates/function.html
var functionpage string

//go:embed templates/functionedit.html
var functioneditpage string

//go:embed templates/auth.html
var authpage string

//go:embed static
var css embed.FS
