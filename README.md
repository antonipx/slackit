slackit
=======
Slack File Upload Gateway

Abstract
--------
This service allows anonymously upload file attachments to specified Slack channel using `curl`. The use case is to allow customers and support engineers to post logs or config files directly from remote servers without uploading and downloading files.

Example usage
-------------
curl -F file=@somelog.txt https://slackit.snakeoil.com/customerA

curl -F file=@config.json https://slackit.snakeoil.com/engineering

Running the gateway service
---------------------------

```
docker run --rm -e APITOK="xoxp-slackit-token-123456" -e HOSTN="slackit.snakeoil.com" -e DATADIR=/data -v /hostdir:/data -v /etc/pki:/etc/pki -p 80:8080 -p 443:8443 --name slackit -d antonipx/slackit
```

- You need to pass APITOK variable containing Slack API Token
- You need to pass HOSTN variable containing FQDN of the HTTPS endpoint, used for Lets Encrypt host policy and usage msg
- You need to expose port 80 and 443 on the host lets Encrypt auto cert engine to work correctly
- You need to pass location of CA cerets so that communication with Slack API endpoint can be established, typically /etc/pki 
- Optionally you can map a host directory to be used for specified DATADIR. By default tmp inside the container will be used. Both LetsEncrypt certificate as well as file attachments are cached in this directory.
