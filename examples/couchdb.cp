<html>
<body>
<h1>Couchdb example</h1>
	<p>This sample installs couchdb</p>
<pre>
(done 
  (couch)
)

(couchdb 
  (and 
    (installed "couchdb") 
    (configured_couchdb)
  )
)

(configured_couchdb
 (local.ini certificate.cert certificate.key)
)


(local.ini
  (exec "
        cat >> /etc/couchdb/local.ini << EOF
	[httpd]
	port = 5984
	[log]
	level = debug
        [ssl]
        cert_file = /etc/pki/tls/certs/certificate.cert
        key_file = /etc/pki/tls/certs/certificate.key
	[admins]
	admin = mysecretpassword
        EOF
        "
  )
)

(certificate.cert
  (exec "
        cat >> /etc/couchdb/local.ini << EOF
-----BEGIN CERTIFICATE----- 
MIIE+zCCA+OgAwIBAgIDAj8hMA0GCSqGSIb3DQEBBQUAMEAxCzAJBgNVBAYTAlVT
MRcwFQYDVQQKEw5HZW9UcnVzdCwgSW5jLjEYMBYGA1UEAxMPR2VvVHJ1c3QgU1NM
IENBMB4XDTEzMDUxNjA5NDc0MFoXDTE0MDUxOTA2MjcxM1owgZsxKTAnBgNVBAUT
IFhxMG9TcUQ2ZTJScGZWNG1xZkpKTWFST3I2MHZURFQtMQswCQYDVQQGEwJERTEQ
ZWRpYSBBRzETMBEGA1UECxMKVGVjaG5vbG9neTETMBEGA1UEAwwKKi5kM3N2Lm5l
dDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN+/gSu+2w45wckB0o4z
8Q3nGvqu1/hHIQ5JXfqymr8x9ac+X9fnCGqIqHwa2LMKc+PxCk7w2JZiSM6uPXhS
yVBR7QEjAe2twJ1mvl6lfzpC6ghhk3mttK+oT+9NVrFFhbsXTKKRT7bvBTJRfrfE
YIu4fz0fIb5fs/uIIMRnt2K+3wSxk5Go3reTJ6YxW8R7R0Q0u5i1+XKITiZR1JNf
s2pXfFp/DMnwXSItEv4+mxqP5Z5sc6QRbYgAV48nHuPdM6WMD7MjfqQJmZXbvZjI
HwYDVR0RBBgwFoIKKi5kM3N2Lm5ldIIIZDNzdi5uZXQwPQYDVR0fBDYwNDAyoDCg
cnVzdC5jb20wMwYIKwYBBQUHMAKGJ2h0dHA6Ly9ndHNzbC1haWEuZ2VvdHJ1c3Qu
7ksCAwEAAaOCAaAwggGcMB8GA1UdIwQYMBaAFEJ5VBthzVUrPmPVPEhX9Z/7Rc5K
MA4GA1UdDwEB/wQEAwIEsDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIw
LoYsaHR0cDovL2d0c3NsLWNybC5nZW90cnVzdC5jb20vY3Jscy9ndHNzbC5jcmww
HQYDVR0OBBYEFM4F4bamyZFBKI081/6qYcublONyMAwGA1UdEwEB/wQCMAAwbwYI
KwYBBQUHAQEEYzBhMCoGCCsGAQUFBzABhh5odHRwOi8vZ3Rzc2wtb2NzcC5nZW90
Y29tL2d0c3NsLmNydDBMBgNVHSAERTBDMEEGCmCGSAGG+EUBBzYwMzAxBggrBgEF
1p6MZpB9QANAYi+/F0Ypy9IUdoNIi4pr2qknL2SdCXcmg20GAf5GI8Gxst44fBV4
BQcCARYlaHR0cDovL3d3dy5nZW90cnVzdC5jb20vcmVzb3VyY2VzL2NwczANBgkq
hkiG9w0BAQUFAAOCAQEAEbOU8Cz5wqZ7RGCFfcDsalhKdtZSrHSsdPqdm7mO+lT0
vgbWPysiLEa3VwPZQNcBnvT5jdDldR7Q3mE9zzm8tCoEXw5LblSBM43AOcquusbh
tKpWRuPlYy/7FEbGZOUs0wa+ON7/Xi4rplZx6d8UYw8NuDTGFQ+XV3QV/Cc/iVyX
bel/ZU6biGbjyror5mOetLA1PYXpV1ET8/80Xe4YPYVxh62xF5ZH6pBc4if7aWQi
aVEtZhnyG7hEdOsz44xughwe10ZpaRqAMa5o5yGaOKFnD/GcreFfEeuRpnCBBwhh
bs8YMMgbo4vy6KqC/KLQ3J8i8HQOqWu6/4rLBDELfg==
-----END CERTIFICATE-----
        "
  )
)

(certificate.key
  (exec "
        cat >> /etc/couchdb/local.ini << EOF
-----BEGIN CERTIFICATE----- 
MIIE+zCCA+OgAwIBAgIDAj8hMA0GCSqGSIb3DQEBBQUAMEAxCzAJBgNVBAYTAlVT
MRcwFQYDVQQKEw5HZW9UcnVzdCwgSW5jLjEYMBYGA1UEAxMPR2VvVHJ1c3QgU1NM
IENBMB4XDTEzMDUxNjA5NDc0MFoXDTE0MDUxOTA2MjcxM1owgZsxKTAnBgNVBAUT
ZWRpYSBBRzETMBEGA1UECxMKVGVjaG5vbG9neTETMBEGA1UEAwwKKi5kM3N2Lm5l
8Q3nGvqu1/hHIQ5JXfqymr8x9ac+X9fnCGqIqHwa2LMKc+PxCk7w2JZiSM6uPXhS
yVBR7QEjAe2twJ1mvl6lfzpC6ghhk3mttK+oT+9NVrFFhbsXTKKRT7bvBTJRfrfE
YIu4fz0fIb5fs/uIIMRnt2K+3wSxk5Go3reTJ6YxW8R7R0Q0u5i1+XKITiZR1JNf
s2pXfFp/DMnwXSItEv4+mxqP5Z5sc6QRbYgAV48nHuPdM6WMD7MjfqQJmZXbvZjI
IFhxMG9TcUQ2ZTJScGZWNG1xZkpKTWFST3I2MHZURFQtMQswCQYDVQQGEwJERTEQ
dDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN+/gSu+2w45wckB0o4z
LoYsaHR0cDovL2d0c3NsLWNybC5nZW90cnVzdC5jb20vY3Jscy9ndHNzbC5jcmww
HwYDVR0RBBgwFoIKKi5kM3N2Lm5ldIIIZDNzdi5uZXQwPQYDVR0fBDYwNDAyoDCg
cnVzdC5jb20wMwYIKwYBBQUHMAKGJ2h0dHA6Ly9ndHNzbC1haWEuZ2VvdHJ1c3Qu
7ksCAwEAAaOCAaAwggGcMB8GA1UdIwQYMBaAFEJ5VBthzVUrPmPVPEhX9Z/7Rc5K
MA4GA1UdDwEB/wQEAwIEsDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIw
Y29tL2d0c3NsLmNydDBMBgNVHSAERTBDMEEGCmCGSAGG+EUBBzYwMzAxBggrBgEF
HQYDVR0OBBYEFM4F4bamyZFBKI081/6qYcublONyMAwGA1UdEwEB/wQCMAAwbwYI
KwYBBQUHAQEEYzBhMCoGCCsGAQUFBzABhh5odHRwOi8vZ3Rzc2wtb2NzcC5nZW90
1p6MZpB9QANAYi+/F0Ypy9IUdoNIi4pr2qknL2SdCXcmg20GAf5GI8Gxst44fBV4
BQcCARYlaHR0cDovL3d3dy5nZW90cnVzdC5jb20vcmVzb3VyY2VzL2NwczANBgkq
hkiG9w0BAQUFAAOCAQEAEbOU8Cz5wqZ7RGCFfcDsalhKdtZSrHSsdPqdm7mO+lT0
tKpWRuPlYy/7FEbGZOUs0wa+ON7/Xi4rplZx6d8UYw8NuDTGFQ+XV3QV/Cc/iVyX
bel/ZU6biGbjyror5mOetLA1PYXpV1ET8/80Xe4YPYVxh62xF5ZH6pBc4if7aWQi
aVEtZhnyG7hEdOsz44xughwe10ZpaRqAMa5o5yGaOKFnD/GcreFfEeuRpnCBBwhh
vgbWPysiLEa3VwPZQNcBnvT5jdDldR7Q3mE9zzm8tCoEXw5LblSBM43AOcquusbh
bs8YMMgbo4vy6KqC/KLQ3J8i8HQOqWu6/4rLBDELfg==
-----END CERTIFICATE-----
        "
  )
)

(installed
 (or
  (and (debian) (apt_installed [arg 0]))
  (and (fedora) (yum_installed [arg 0]))
 )
)

(apt_installed (exec "/usr/bin/test -x /usr/bin/apt"))
(yum_installed (exec "/usr/bin/test -x /usr/bin/yum"))
(debian (exec "/usr/bin/grep -q Debian /etc/system-release"))
(fedora (exec "/usr/bin/grep -q Fedora /etc/system-release"))

</pre>
</body>
</html>
