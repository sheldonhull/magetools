#!/usr/bin/env bash
if ! command -v dotnet &> /dev/null
then
    echo "dotnet could not be found"
else
    dotnet tool install --global GitVersion.Tool && echo "✅  gitversion installed" || echo "‼ failed to install gitversion"
fi
