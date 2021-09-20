# API Generation

#### Install NodeJS

* https://nodejs.org/en/

#### Install openapi-generator-cli  
`npm install @openapitools/openapi-generator-cli -g`
* More info:
  * https://openapi-generator.tech/docs/installation/
    
#### Generate API

`cd src/nwdaf/deploy/`

`sudo openapi-generator-cli generate -i TS29520_Nnwdaf_AnalyticsInfo.yaml -g go --skip-validate-spec -o api/Nnwdaf_AnalyticsInfo/`

`sudo openapi-generator-cli generate -i TS29520_Nnwdaf_EventsSubscription.yaml -g go --skip-validate-spec -o api/Nnwdaf_EventsSubscription/`

* Optional: Example using `openapi-generator-cli` docker version

  `sudo docker run --rm -v $(pwd):/local openapitools/openapi-generator-cli generate -i /local/TS29520_Nnwdaf_AnalyticsInfo.yaml -g go --skip-validate-spec -o /local/Nnwdaf_AnalyticsInfo/`

