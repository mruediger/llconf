Vagrant Images
==============

Here are some [Veewee] (https://github.com/jedi4ever/veewee) recipies to build Vagrant images to test LLConf
To install Veewee, just run `gem install veewee`

In this folder run

    veewee vbox build llconf-CentOS-6.4-x86_64

and wait a little. To login after the setup use the the password "veewee" and the following command:

    ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p 7222 -l veewee 127.0.0.1

To export the VM and import into vagrant run:

    veewee vbox export llconf-CentOS-6.4-x86_64
    vagrant box add llconf-centos llconf-CentOS-6.4-x86_64.box

If you have done all this, you can use Vagrant like you are used to, e.g.:

    vagrant init llconf-centos
    vagrant up
    vagrant ssh
