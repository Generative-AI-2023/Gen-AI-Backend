# HolidAI backend

This is the backend for HolidAI, an AI-based trip planner

# Code Explanation

The HolidAI backend responds to trip requests. The trip information is sent to HolidAI as a JSON file, which the backend parses and forms a guided prompt with. The backend then returns the openAI response to the prompt in a format that the front end can use.

# File writeup

## main.go

This is the main backend code in go

## Procfile

Sets commands to start binary compile by makefile.

## Makefile

Builds binary of main.go. Make sure that names in Procfile and Makefile match.

## heroku.yml

Sets language and name of binary

## go.mod

Documents modules used in main.go

## go.sum

Provides version numbers and more for modules in go.mod

## Dockerfile

Sets up, executes and deploys buildpack

## .gitignore

Specifies files to ignore in git