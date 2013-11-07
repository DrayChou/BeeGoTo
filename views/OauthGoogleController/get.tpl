{{template "header.tpl"}}

<a href="{{.AuthCodeURL}}">{{.AuthCodeURL}}</a>

<form action="" method="post">
  <p>code: <input type="text" name="code" id="code" /></p>
  <input type="submit" value="Submit" />
</form>

{{template "footer.tpl"}}