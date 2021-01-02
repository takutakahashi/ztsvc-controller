{{ $svc := . }}
{{ range $i, $port := .Spec.Ports }}
frontend frontend-port-{{ $port.Port }}
  default_backend backend-port-{{ $port.Port }}
  bind *:{{ $port.Port }}
backend backend-port-{{ $port.Port }}
  server server-{{ $port.Port }} {{ $svc.Name }}:{{ $port.Port }}
{{ end }}
