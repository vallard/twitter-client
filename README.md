# Twitstream

A go twitter client.  Mostly inspired by   
https://github.com/matryer/goblueprints/blob/master/chapter6/twittervotes/main.go
but I wanted to turn it into a package.  So I split it into two parts.  

The main program sets it up and the twitstream package does the Twitter work. 

## Setup
You'll need to set some environment variables for this to work. 
First you need to register an application on Twitter.  Then take the keys they
give you and export them in your .bash_profile.  It will look something like this: 
```
export STRIPSTOCK_TWITTER_CONSUMER_KEY=1234...
export STRIPSTOCK_TWITTER_CONSUMER_SECRET=123...
export STRIPSTOCK_TWITTER_ACCESS_TOKEN=abcd...
export STRIPSTOCK_TWITTER_ACCESS_SECRET=abcd...
```

## Running the program.  

You can run it with: 
```
go run twitter-client.go earth wind fire water
```

The program will then stream to your heart's content.  You cancel with ctrl-c and the
program should exit gracefully.  

## Enhancements

You can enhance the Get function to stream it to something like kafka or NSQ.  
