#!/bin/bash

if [[ $# != 1 ]]; then
    echo "argument can not be less or more than 1"
    exit 400
fi

FEATURE="${1,}"

UPPERCASE_INDEX=()
for ((i = 0; i < ${#FEATURE}; i++)); do
    char="${FEATURE:$i:1}"

    if [[ "$char" == [[:upper:]] ]]; then
        UPPERCASE_INDEX+=($i)
    fi
done

PARENTFOLDER="${FEATURE,,}"
LENGTH=${#PARENTFOLDER}
COUNT=0
for INDEX in ${UPPERCASE_INDEX[@]}; do 
    PARENTFOLDER="${PARENTFOLDER:0:$INDEX+$COUNT}_${PARENTFOLDER:$INDEX+$COUNT:$LENGTH}"
    COUNT=$((COUNT + 1))
done

mkdir $PARENTFOLDER
cd $PARENTFOLDER
mkdir dtos handler mocks repository usecase

FILES=("interfaces.go" "entities.go" "dtos/request.go" "dtos/response.go" "handler/controller.go" "usecase/service.go" "repository/model.go")

for FILE in ${FILES[@]}; do
    sed -e "s/placeholder/${FEATURE}/g" -e "s/Placeholder/${FEATURE^}/g" -e "s/_blueprint/${PARENTFOLDER}/g" ../_blueprint/$FILE > $FILE
done

# setup routes
sed -e "s/placeholder/${FEATURE}/g" -e "s/Placeholder/${FEATURE^}/g" -e "s/_blueprint/${PARENTFOLDER}/g" ../_blueprint/routes.go > "../../routes/${PARENTFOLDER}.go"
