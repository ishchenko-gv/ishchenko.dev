#!/bin/bash

go run ./build
gcloud compute scp --recurse dist/* ishchenko-dev:/home/georgii/www