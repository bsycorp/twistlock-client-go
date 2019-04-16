#!/bin/bash
docker login -u $DOCKERUSER -p $DOCKERPASS
docker build . -t bsycorp/twistlock-controller:$TRAVIS_BRANCH
docker push bsycorp/twistlock-controller
