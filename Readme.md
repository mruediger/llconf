LLConf
======

## LLConf - Lisp Like configuration management ##

LLConf is a configuration management system. It is, just like Ruby and Cheff inspired
by CFEngine. It has three main goals:

* Make it feel natural to describe a machine state instead of a setup process
* Be extendible
* Keep it simple

## Sample Config ##

    (done
      (or (and (is webserver) (apache ready))
          (and (is database) (couch ready))))

    (apache ready
      (and (installed "apache") (apache configured)))

    (couch ready
      (and (installed "couchdb") (couchdb configured)))

    (couchdb configured
      (and (couch local.ini) (couch certificate.cert) (couch certificate.key)))

    (couch local.ini
      (exec "sh" "-c"
        "cat >> /etc/couchdb/local.ini << EOF
	        [httpd]
	        port = 5984

            [log]
         	level = debug
            [ssl]
            cert_file = /etc/pki/tls/certs/certificate.cert
            key_file = /etc/pki/tls/certs/certificate.key

            [admins]
	        admin = mysecretpassword
         EOF"
      )
    )
      
    ...
