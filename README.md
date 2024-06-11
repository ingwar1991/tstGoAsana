# tstGoAsana

### Run
It's possible to run program in docker container
Go to docker/ folder and run "docker compose up" command from there
It will create container in which program can run

### Program
The program itself is located at project/

### Config
Do not forget to create project/conf.yml from project/conf.yml.dist and add a valid Asana PAT token there

#### Progress
To check the state of the task at the end of live coding check the commit 4569c148ddf8f3fe4f34919bf74d61a60f24bdf4

### Scaling
As for scaling, I don't think there's anything can be done with one token. Due to pretty small limit( 50 GET requests per minute ) 
there's no sense in scaling without multiple tokens. 

But if we can get multiple tokens( from users authorizing in our Asana APP, by ex ) something can be done. 
We can run the whole program functionality in multiple threads giving different tokens to them. We can give a dediced token
to each thread or just share a list of tokens with some logic of swithcing those, that are already close to rate limit.

Unfortunatelly, Asana API doesn't provide any call for getting total number of projects/users, so it's not possible to know how many
threads are required. So, the way I would suggest to handle that is to send the program the supposed number of projects/users and
program will create threads accordingly. So, let's say we expect there to be 3000 projects and our calculations showed us, that
the most optimal number for one thread for getting & saving projects is 1000 per unit. Thus the program will create 3 threads.
Considering that we are getting 100 projects per api call( can be adjusted too ), threads will be created with offsets
respectively:
1. from 0 to 900
2. from 1000 to 1900
3. from 2000 to 2900

There's the possibility that we have guessed the number wrong. If we guessed the lower number then the last thread just would have to
process more units. If less, then the program will create threads that won't return data, just errors from Asana API
