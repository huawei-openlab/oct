## {{.Name}}
[Test Case](#testcase) {{.Description}}

```
Owner: {{.Owner}}
License: {{.License}}
Group: {{.Group}}
```

The case has {{len .Deploys}} operation system(s):

{{range .Deploys}}
'{{.Object}}' has {{len .Containers}} container(s) deployed.
{{end}}

