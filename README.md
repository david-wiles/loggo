# Loggo

A simple Go logging framework

## The problem

I find myself wanting to log errors and generic info messages in my projects. However, the logger in the standard 
library is somewhat lacking and I end up copying the same logging functions into all my projects.

This project moves these functions into their own repository for my own convenience and for the use of others.
The API design is somewhat based on the logger in Salesforce Commerce Cloud and the logging middleware is inspired by
the middleware described by Matt Silverlock in his blog, [Questionable Services](https://blog.questionable.services/).

## API

The API is very simple. Create a Loggo instance and give it a writer and log level. No other configuration required.

See examples in examples/
