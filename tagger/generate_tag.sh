#!/bin/bash

echo "running tag script..."

major=${MAJOR_VERSION:1}
minor=${MINOR_VERSION:13}
rel=beta

branch_name=$CI_COMMIT_REF_NAME

git checkout $branch_name

# Get the commit message of the last commit on the branch
commit_message=$(git log -1 --pretty=%B)

if [[ "$commit_message" != *tag-it* ]]; then
    echo "Commit message does not contain 'tag-it'. Skipping tag scripts."
    exit 1
fi

read -p "Enter target branch [master, qa, rls] : " rel

if [[ "$rel" != "master" && "$rel" != "qa" && "$rel" != "rls" ]]; then
    echo "Error: Invalid keyword. Please enter 'master', 'qa', or 'rls'."
    exit 1
fi

if [[ "$rel" == "master" ]]; then
    echo "here"
    rel="beta"
fi

tag_regex="v([0-9]+)\.([0-9]+)\.([0-9]+)-.*"

latest_tag=$(git tag -l --sort=-version:refname "v*-$rel-*" | head -n 1)
echo "latest tag is $latest_tag"

if [[ $latest_tag =~ $tag_regex ]]; then
    major="${BASH_REMATCH[1]}"
    minor="${BASH_REMATCH[2]}"
    patch="${BASH_REMATCH[3]}"

    echo "major is $major"
    echo "minor is $minor"
    echo "patch is $patch"

    patch=$((patch + 1))

    new_tag="v$major.$minor.$patch-$rel-$branch_name"
    echo "new tag is $new_tag"

    git tag "$new_tag"

    git push origin "$new_tag"

else 
    echo "Error: latest tag doesn't match the expected format"
fi