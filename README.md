# KaaS - Klaus as a Service

Nexmo TTL with a twist

HTTP frontend listens on port 9090

Uses a in-memory cache with TTL to check if the number has already been "klaused"

If not a call is made from unknown number using the provided Nexmo credentials

-"Hallo, du hast bin Klaused!" ~"In the mood" by Klaus Wunderlich starts playing

Environment config

    APP_ID=<nexmo application id>
    API_KEY=<nexmo api key>
    API_SECRET=<nexmo api secret>
    NCCO_URL=https://your nexmo [ncco](https://developer.nexmo.com/voice/voice-api/ncco-reference) file

## private.key

You need to save your Nexmo private key in the folder of klaus, named as private.key
