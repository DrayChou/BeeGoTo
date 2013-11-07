{{template "header.tpl"}}

<a href="{{.AuthCodeURL}}">{{.AuthCodeURL}}</a>
{{.JsonStr}}

{{template "footer.tpl"}}