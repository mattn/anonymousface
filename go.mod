module anonymousface

go 1.15

require (
	github.com/esimov/pigo v1.4.5
	github.com/mattn/anonymousface/statik v0.0.0-00010101000000-000000000000
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/rakyll/statik v0.1.7
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410
)

replace github.com/mattn/anonymousface/statik => ./statik
