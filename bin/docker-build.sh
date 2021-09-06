#!/bin/sh

checkOk() {
	# shellcheck disable=SC2181
	if [ "$2" != 0 ]; then
		# shellcheck disable=SC2028
		echo "
        ############   build $1 fail!   ############
        "
		exit 1
	else
		# shellcheck disable=SC2028
		echo "
        ############   build $1 success!   ############
        "
	fi
}

build() {
	local v=$2
	if [ "$2" = '' ]; then
		echo "缺少版本号 如 v0.0.1, 默认版本号为 latest"
		v='latest'
	fi

	p=$(pwd)
	cd "$p"

	docker build -t docker.pkg.github.com/kainhuck/yao-proxy/"$1":"$v" -f cmd/"$1"/Dockerfile .
  docker push docker.pkg.github.com/kainhuck/yao-proxy/"$1":"$v"

	checkOk $1 $?
}

buildAll() {
	local v=$1
	if [ "$v" = '' ]; then
		echo "缺少版本号 如 v0.0.1, 默认版本号为 latest"
		v='latest'
	fi

	build local "$v"
	build remote "$v"

	# wait
	echo "------end------"
}

case $1 in
"all")
	# shellcheck disable=SC2119
	buildAll "$2"
	;;
*)
	build "$1" "$2"
	;;
esac