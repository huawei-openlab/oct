
The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
{{range .HostDetails}}|{{.Object}}|{{.Distribution}}|{{.Resource}}|{{.Containers}}| {{.Command}}|
{{end}}

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
{{range .ContainerDetails}}|{{.Class}}|{{.Distribution}}|[{{.ConfigFile}}](#{{.ConfigTag}})|
{{end}}

After running the `Command` in each OS and container, we get {{len .Collects}} result(s).

{{range .Collects}}
* [{{.Name}}](#{{.Tag}})
{{end}}

###{{.TestCase.Tag}}

```
{{.TestCase.Content}}
```

{{range .Files}}

###{{.Tag}}

```
{{.Content}}
```
{{end}}

{{range .Collects}}

###{{.Tag}}

```
{{.Content}}
```
{{end}}
