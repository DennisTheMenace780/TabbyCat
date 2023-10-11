#!/bin/bash

# Step 1: Create a directory in testdata/TestRepo2
mkdir -p testdata/TestRepo

# Step 2: Change directories into TestRepo2
cd testdata/TestRepo || exit

# Step 3: Initialize a Git repository
git init

# Step 4: Create a file called file.txt with "hello world" in it
echo "hello world" > file.txt

# Optionally, you can add and commit the file to the Git repository:
git add file.txt
git commit -m "Initial commit"
git checkout -b 'branch-1'
git checkout -b 'branch-2'
git checkout -b 'branch-3'
git checkout master
git branch
