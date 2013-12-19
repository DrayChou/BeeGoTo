{{template "header.tpl"}}

{{if .AuthCodeURL}}

	<a href="{{.AuthCodeURL}}">{{.AuthCodeURL}}</a>

{{else}}

	{{with .JsonStr}}
		{{range .}}
			{{.User.Screen_name}}&nbsp;&nbsp;{{.Title}}&nbsp;&nbsp;{{.Text}}&nbsp;&nbsp;{{.Created_at}}<br/>
		{{end}}
	{{end}}
	
	{{if .dbu.Screen_name}}
		<h1>{{.dbu.Screen_name}}</h1>
	{{else}}
		<h1>{{.dbu.Name}}</h1>
	{{end}}
	
	{{.dbu.Signature}}<br/>
	{{.dbu.Desc}}<br/>
	<a href="{{.dbu.Alt}}"><img src="{{.dbu.Large_avatar}}" /></a><br/>

{{end}}

{{template "footer.tpl"}}