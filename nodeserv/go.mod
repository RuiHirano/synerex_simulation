module nodeid-server

require (
	github.com/google/gops v0.3.5
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/synerex/synerex_alpha/nodeapi v0.0.1
	github.com/synerex/synerex_alpha/util v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/grpc v1.28.0
)

replace github.com/synerex/synerex_alpha/nodeapi => ../nodeapi

replace github.com/synerex/synerex_alpha/util => ../util

go 1.13
