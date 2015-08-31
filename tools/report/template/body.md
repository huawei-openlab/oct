
The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
{{range .HostDetails}}
|{{.Object}}|{{.Distribution}}|{{.Resource}}|{{.Containers}}| {{.Command}}|
{{end}}

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
{{range .ContainerDetails}}
|{{.Class}}|{{.Distribution}}|[{{.ConfigFile}}](#{{.ConfigFile}})|
{{end}}

After running the `Command` in each OS and container, we get {{len .Collects}} result(s).

{{range .Collects}}
* [{{.Name}}](#{{.Name}})
{{end}}

###{{.TestCase.Name}}

```
{{.TestCase.Content}}
```

{{range .Files}}

###{{.Name}}

```
{{.Content}}
```
{{end}}

{{range .Collects}}

###{{.Name}}

```
{{.Content}}
```
{{end}}
