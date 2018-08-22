docker run --rm -e APITOK="xoxp-slackit-token-123456" -e HOSTN="slackit.snakeoil.com" -v /hostdir:/tmp -v /etc/pki:/etc/pki -p 80:8080 -p 443:8443 --name slackit -d antonipx/slackit 
