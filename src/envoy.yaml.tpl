{{ $svc := . }}
static_resources:
  listeners:
{{- range $i, $port := .Spec.Ports }}
{{- if eq $port.Protocol "UDP" }}
  - name: listener_{{ $i }}
    reuse_port: true
    address:
      socket_address:
        protocol: UDP
        address: 0.0.0.0
        port_value: {{ $port.Port }}
    udp_listener_config:
      downstream_socket_config:
        max_rx_datagram_size: 9000
    listener_filters:
    - name: envoy.filters.udp_listener.udp_proxy
      typed_config:
        '@type': type.googleapis.com/envoy.extensions.filters.udp.udp_proxy.v3.UdpProxyConfig
        stat_prefix: service
        cluster: service_udp_{{ $i }}
        upstream_socket_config:
          max_rx_datagram_size: 9000
{{- end }}
{{- end }}
  clusters:
{{- range $i, $port := .Spec.Ports }}
{{- if eq $port.Protocol "UDP" }}
  - name: service_udp_{{ $i }}
    connect_timeout: 0.25s
    type: STATIC
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: service_udp_{{ $i }}
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: {{ $svc.Name }}
                port_value: {{ $port.Port }}
{{- end }}
{{- end }}