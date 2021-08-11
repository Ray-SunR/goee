module example

go 1.16

replace gee => ./gee

replace cache => ./cache

replace consistenthash => ./consistenthash

require (
	cache v0.0.0-00010101000000-000000000000
	gee v0.0.0-00010101000000-000000000000
)
