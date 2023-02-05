module shell/main

go 1.19

replace shell/config => ../config

require (
	shell/builtins v0.0.0-00010101000000-000000000000
	shell/config v0.0.0-00010101000000-000000000000
)

replace shell/builtins => ../builtins
