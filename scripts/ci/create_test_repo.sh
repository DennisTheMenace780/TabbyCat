#!/bin/bash

# Step 1: Create a directory in testdata/TestRepo2
mkdir -p testdata/TestRepo

# Step 2: Change directories into TestRepo2
cd testdata/TestRepo || exit

# Step 3: Initialize a Git repository
git init

git config --global user.email "you@example.com"
git config --global user.name "Your Name"

# Step 4: Create a file called file.txt with "hello world" in it
echo "hello world" > file.txt

# Optionally, you can add and commit the file to the Git repository:
git add file.txt && git commit --no-verify -m "Inital Commit"

count=1
while [ $count -lt 4 ]; do
    git checkout -b "branch-$count"
    count=$((count + 1))
done
