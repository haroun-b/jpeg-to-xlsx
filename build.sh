#!/bin/bash

dists=("linux" "windows" "darwin")
app_name="jpeg-to-xlsx"

for os in "${dists[@]}"; do
  bin_name="${app_name}"_$os

  if [ $os == "windows" ]; then
    bin_name=$bin_name.exe
  fi

  GOOS=$os GOARCH=amd64 go build -o ./bin/$bin_name
done

echo "Build Complete!"