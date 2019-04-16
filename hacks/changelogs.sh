#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CHANGELOG_PATH='changelogs/unreleased'
UNRELEASED=$(ls -t ${CHANGELOG_PATH})
echo -e "Generating CHANGELOG markdown from ${CHANGELOG_PATH}\n"
for entry in $UNRELEASED
do
    IFS=$'-' read -ra pruser <<<"$entry"
    contents=$(cat ${CHANGELOG_PATH}/${entry})
    echo "  * ${contents} (#${pruser[0]}, @${pruser[1]})"
done
echo -e "\nCopy and paste the list above in to the appropriate CHANGELOG file."
echo "Be sure to run: git rm ${CHANGELOG_PATH}/*"
