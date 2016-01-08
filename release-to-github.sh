#!/bin/bash

pushd $BUILD_DIR

#if [ "Xjenkins-starschema" == "X$CI_COMMITTER_USERNAME" ]; then echo "Modification by build server, no need to deploy"; exit 0; fi
#if [ "X" == "X$GITHUB_USER" ]; then echo "Need GITHUB_USER environment variable"; exit 10; fi
if [ "X" == "X$GITHUB_TOKEN" ]; then echo "Need GITHUB_TOKEN environment variable"; exit 10; fi
#if [ "X" == "X$GITHUB_EMAIL" ]; then echo "Need GITHUB_EMAIL environment variable"; exit 10; fi
if [ "X" == "X$HOME" ]; then echo "Need HOME environment variable"; exit 10; fi

#echo "Set github user..."
#git --version
#git config --global push.default simple
#git config --global user.name $GITHUB_USER
#git config --global user.email $GITHUB_EMAIL

#echo "Updating version..."
#VERSION=`npm version patch -m "Version %s [ci skip][skip ci]"`
#if [ $? -ne 0 ]; then echo "updating package version failed"; exit 10; fi
#echo $VERSION

#echo "Building javascript files"
#npm install
#node_modules/gulp/bin/gulp.js build
#if [ $? -ne 0 ]; then echo "gulp build failed"; exit 10; fi

#echo "Removing unnecessary files"
#rm -rf *.sh
#rm -rf *.yml
#rm -rf Dockerfile
#rm -rf build*
#rm -rf test
#rm -rf tasks

#echo "Creating release package..."
#rm -rf node_modules
#npm install --production
#if [ $? -ne 0 ]; then echo "installing production packages failed"; exit 10; fi
#pushd $HOME
#PCKG_NAME=$PACKAGE-$VERSION
#PCKG_FILE=$PCKG_NAME.tbz
#mkdir $PCKG_NAME
#cp -Rf $BUILD_DIR/* $PCKG_NAME
#tar cjf $BUILD_DIR/$PCKG_FILE $PCKG_NAME
#popd
#du -h $PCKG_FILE

echo "Uploading new version to Github..."
git push --force "https://$GITHUB_TOKEN@github.com/$OWNER/$PACKAGE.git" HEAD:master --tags
if [ $? -ne 0 ]; then echo "uploading new version failed"; exit 10; fi

echo "Creating Github realase..."
RELEASE_ID=`curl -H "Authorization: token $GITHUB_TOKEN" -d "{\"tag_name\": \"VERSION\"}" "https://api.github.com/repos/$OWNER/$PACKAGE/releases"| jsawk "return this.id"`
if [ $? -ne 0 ]; then echo "creating new release failed"; exit 10; fi
echo $RELEASE_ID

echo "Uploading Github realase asset..."
curl --progress-bar \
     -H "Content-Type: application/octet-stream" \
     -H "Authorization: token $GITHUB_TOKEN" \
     --retry 3 \
     --data-binary @$PCKG_FILE \
     "https://uploads.github.com/repos/$OWNER/$PACKAGE/releases/$RELEASE_ID/assets?name=$PCKG_FILE"
if [ $? -ne 0 ]; then echo "uploading release asset failed"; exit 10; fi