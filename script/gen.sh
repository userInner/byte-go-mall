#!/bin/bash

svcName=${1}

cd app/${svcName}
cwgo client -I ../../idl --type RPC --service ${svcName} --module byte-go-mall --idl ../../idl/${svcName}.proto
cwgo server -I ../../idl --type RPC --service ${svcName} --module byte-go-mall --idl ../../idl/${svcName}.proto
