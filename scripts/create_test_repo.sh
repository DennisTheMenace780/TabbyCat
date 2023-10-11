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
git add file.txt && git commit -m "Inital Commit"

git checkout -b 'branch-1'
git checkout -b 'branch-2'
git branch
git checkout -b 'branch-3'
git branch
