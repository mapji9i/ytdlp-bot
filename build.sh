#!/bin/zsh

package=$1

if [[ -z "$package" ]]; then
  package="./pkg/main"
fi
package_split=(${(@s:/:)package})
package_name=${package_split[-1]}
	
platforms=("linux/amd64")

for platform in "${platforms[@]}"
do
	platform_split=(${(@s:/:)platform})
	GOOS=${platform_split[1]}
	GOARCH=${platform_split[2]}
  
	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi	
	mkdir -p ./build
	env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=1 CC=x86_64-unknown-linux-gnu-gcc  go build -o ./build/$output_name $package
	#env CGO_ENABLED=1 CGO_CFLAGS="-D_LARGEFILE64_SOURCE" CC="zig cc -target x86_64-linux" CXX="zig c++ -target x86_64-linux" GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
