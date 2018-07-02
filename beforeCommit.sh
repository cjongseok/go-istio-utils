#!/bin/bash

# Run tests

#dep ensure
go build .
for dir in $(ls); do
    if [ -d $dir ] && [[ $dir != "vendor" ]]; then
        pkg=./$dir
        echo "build $dir"
        go build $pkg
    fi
done

go test .
for dir in $(ls); do
    if [ -d $dir ] && [[ $dir != "vendor" ]]; then
        pkg=./$dir
        echo "test $dir"
        go test $pkg
    fi
done

#go tool vet .
for dir in $(ls); do
    if [ -d $dir ] && [[ $dir != "vendor" ]]; then
        pkg=./$dir
        echo "vet $dir"
        go tool vet $pkg
    fi
done

golint .
for dir in $(ls); do
    if [ -d $dir ] && [[ $dir != "vendor" ]]; then
        pkg=./$dir
        echo "lint $dir"
        golint $pkg
    fi
done

go fmt .
for dir in $(ls); do
    if [ -d $dir ] && [[ $dir != "vendor" ]]; then
        pkg=./$dir
        echo "fmt $dir"
        go fmt $pkg
    fi
done
rm -rf ./vendor
