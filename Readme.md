LLConf
======

[![Build Status](https://travis-ci.org/d3media/llconf.png?branch=master)](https://travis-ci.org/d3media/llconf)

## LLConf - Lisp Like configuration management ##

LLConf is a configuration management system. It is, just like Ruby and Cheff inspired
by CFEngine. It has three main goals:

* Make it feel natural to describe a machine state instead of a setup process
* Be extendible
* Keep it simple


## Config Files ##

LLConf concatenates all files in the input folder recursively. You therefore are free to spilt
your configuration up in as may files as your want.

## Promises ##

LLConf evolved arround the concept of promises that are to be kept by a
machine. Those familiar with CFEngine heard of the concept but the way
its done in LLConf is quite different. In LLConf, everything done is a promise.
There are no classes to select the promises you want to be met. Also, the
amount of special promise types is kept to minimum. Instead what LLconf does
is to help you building your own special promises by the use of named promises.

### Named Promises ###

    ( <some name> ( other promise ))

Named promises serve the purpose of naming a promise. This way, you can
reuse the promise and improve the readablilty of your setup. There use is
mandatory. Every other buildin promise has to be evaluated by testing
a named promise.

### Boolean ###

    (and (promisses...)) (or (promisses ...))

The promises "(and)" and "(or)" are essential to combine promises and by that, create a
tree of promises that have to be fullfilled. Both promises take a list of other promises.
These are evaluated in order. The (and) promise stops evaluating and fails as soon as one
promise in the list fails. The or promise stops evaluating and returns sucess as soon as one
promise is successful.

#### Execution ####

At the and of the day, running llconf boils down to executing shell commands and
interpreting there outputs. Since a promise can only be met or not, the evaluation
of a execution is a simple check if the programm execution succeded or not. To ease
reporting there are two ways to execute programms. They only differ by their name
and the bucket their outcomes are storend into for later monitoring.

#### Tests ####

     (test "command" "argument 1" "argument 2" ... "argument n")

Oftentimes you just want to check stuff and execute something else depending on the outcome.
For example if you want to check if a process is running you don't want a report on how
often the check that the process is running succeded. Ideally this should be always.

#### Changes ####

     (change "command" "argument 1" "argument 2" ... "argument n")

When the process isn't isn't running and you have to change that, you surly want that be reported
so if that change has to be done regualry, something ought to be wrong with either your machine or
the setup you are tring to implement.

#### Pipes ####

     (pipe (test) (test) (change) ... )

The standart unix toolbox is one of the most flexible tools for system adimistration out there.
This all comes down to the simple concept of pipes. LLConf also supports pipes using the pipe promise.
The standart output of every command is passed down as the standart input of the next one in the list.
For obvious reasons you can only use (change) and (test) promises inside a (pipe) promise.

#### Changing directories ####

     (indir ( change / test ))

To execute a program inside a specific directory, for example running a "git checkout" you can use the
indir promise.

## Getters ##

Two types of "Getters" in LLConf. Both allow you to retrieve a string and use it whereever you would use
a string. The first one is the simple argument getter.

### Arguments ###

You can pass arguments whenever you test a named promise.

For example

    (hello world "foo" "bar")

will invoke the named promise "hello world" with the arguments "foo" and "bar". You can use these
arguments using the argument getter [arg:n]. In this example [arg:0] will return "foo" and [arg:1] will
return "bar"

### Variables ###

Another type of getters allow you to use named variables. First, the scope of variables is similar
to the scope of arguments, so basically variables are "named promise scoped". This means that every
named promise creates a new dictionary by copying the dictionary of the surrounding promise. This way
you can overwrite variables locally without affecting evaluations globally.

To get the contents of a variable you simply write

     [var:name] (eg. [var:favorite_colour])

Variables can be set by using on of the following special promises:

     (setvar "name" "value")

This promise stores the value "value" under the name "name"

     (readvar "name (cmd))

This promise stores the standart output of the invoked command under the name "name". The command
can be a (test) a (change) or a (pipe) promise.

### Join ###

Sometimes you have to combine the contents of multible variables to one string. This is done by the
"join" "getter" which simply joins the value of every getter following the join, eg:

    (hello world
        (test "echo" [join [arg:0] " " [arg:1]]))

which will, when invoked like in the sample above, run the command "echo" with the argument "hello world".

### Template Editing ###

LLConf leverages go's template engine. It expects json as the input to the template engine.

       (template "{json}" "template-file" "output-file")

Since the template engine can handle arrays and objects you can easily can adapt the template to your
needs. LLConf is not able to edit files. It is in my oppinion very dangerous to edit a file based
on regular expressions, since you cant be really sure that the config file you are editing is

## Updateing Config Files ##

LLConf keeps the parsed promise-tree in memory and only updates it if there is new and valid input.
This way you can easily change the input files, using git or other means, and don't have to worry about
backing up the last working state.


## Samples ##

#### Keep mysql running ####

    (mysql running
        (or (process match "mysqld")
            (change "/etc/init.d/mysql" "start")))

    (process match
        (test "ps" "-C" [arg:0]))

It simply checks there is a mysqld process running and, if not, fires up the init script



#### Get THE source and keep it up to date ####


    (keep the source "git://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git" "/usr/src/linux")

    (keep the source
        (or (git repo present and uptodate)
            (change "git" "clone" [arg:0])))

    (git repo present and uptodate
        (and (git repo present)
             (git repo uptodate)))

    (git repo present
        (test "test" "-d" [join [arg:0] "/.git"]))

    (git repo uptodate
       (and (indir [arg:0] (test "git" "fetch" "--all")
        (or (indir [arg:0] (test "git" "diff" "--quiet" "HEAD@{upstream}"))
            (indir [arg:0] (change "git" "reset" "--hard" "HEAD@{upstream}")))))


this example is a little more complicated. The promise named "keep the source" checks if there
is a git clone present at the location specified as the second argument and if this repo is uptodate.
If it isn't, it clones the repository specified as the first argument.
