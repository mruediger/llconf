LLConf - Lisp Like configuration management
===========================================

LLConf is a configuration management system. It is, just like Ruby and Cheff inspired
by CFEngine. It has three main goals:

* Make it feel natural to describe a machine state instead of a setup process
* Be extendible
* Keep it simple

# Sample Config #

    (done
      (or (and (is webserver) (apache ready))
          (and (is database) (mysql ready))))

    (apache ready
      (and (installed "apache") (apache configured)))

    (apache configured
      ...
      

