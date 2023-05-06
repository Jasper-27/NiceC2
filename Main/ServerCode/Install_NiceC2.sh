#!/bin/bash


echo THANKS FOR CHOOSING NICEC2 - JASPER 
                                             
## Install the server code
cd Server
sh install.sh

## Get back 
cd ..

## Install the Console code 
cd Console
sh install.sh