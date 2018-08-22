Abstract
--------
This service allows upload file attachments to specified Slack channel using just `curl`. The use case is to allow field personnel, customers or automation to post logs or config files directly from remote servers without uploading and downloading files by hand.

Example usage
-------------
curl -F file=@logfile.txt https://slackit.snakeoil.com/custchannel

curl -F file=@config.json https://slackit.snakeoil.com/engineering


Running the gateway service
---------------------------

```
docker run --rm -e APITOK="xoxp-slackit-token-123456" -e HOSTN="slackit.snakeoil.com" -v /hostdir:/tmp -v /etc/pki:/etc/pki -p 80:8080 -p 443:8443 --name slackit -d antonipx/slackit
```

- You need to pass APITOK variable containing Slack Legacy API Token.
- You need to pass HOSTN variable containing FQDN of remote URL endpoint, used for Lets Encrypt host policy and usage msg.
- You need to expose both ports 80 and 443 on the host for Encrypt auto cert engine to work correctly. 
- You need to map location of CA cerets so that communication with Slack API endpoint can be established, typically /etc/pki.
- You need to map an external directory to be used as /tmp in the container. Make sure it has correct permissions.

