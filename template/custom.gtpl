{{- $TABLE := .Table}}
{{- $D_NUM := len .Table.Data}}
{{- $F_NUM := len .Table.Headers -}}

{{- range $k,$v := .Table.Data}}
    {{- range $hk,$hv := $TABLE.Headers -}}
        {{- if not .IsVoid -}}
    {{(strs_index $v (dec .Index))}}{{if ne $hk (dec $F_NUM)}},{{end}}
        {{- end -}}
    {{- end -}} {{if ne $k (dec $D_NUM)}},{{end}}
{{- end}}
