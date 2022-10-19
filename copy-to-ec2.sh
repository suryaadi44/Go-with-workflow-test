#!/bin/bash
scp -r -i ~/personal-aws.pem ./program/* ubuntu@13.229.209.129:/home/ubuntu/app
