mf-connector
===


Media Fusion Connector that will be part of new Media Fusion Container

## Open Items
===

Need to define feedback receiver to post events back to vTS core

In the event of RabbitMQ connection termination, does the Connector automatically retry or do we report and await new request

Need to figure out how certificates will be accessed


## To Do
===

Define and implement proper response wrapper similar to following

```
// sample JSON envelope
{
  "status": {
    "code": 200,
    "message": 'Success'
  },
  "response": {
     ...results...
  }
}
```
Provide the ability to specify network interfaces to bind to

## Command line
===
	#connector_darwin_amd64 8989 443
	arg1 - http listen port
	arg2 - https listen port
	
Currently the program will try and bind to all interfaces. Future versions will provide the ability to specify network interfaces to bind to.


## API
===
### Get Status

#### Request

```
// Request get status
http(s)://host:port/api/v1/status/
```
#### Response

```
{
	"status": {
		"code": 200,
		"message": "all is well"
	},
	"response": {
		"sessionState": "SessionStateInit"
	}
}
```
### Request URI

#### Request
```
http(s)://host:port/api/v1/redirectUri

{
  "scheme" : "<https>"
  "ipaddress" : "<host:port>"
  "path" : "<token>"
}
```

#### Response

```
{
	"status": {
		"code": 200,
		"message": "all is well"
	},
	"response": {
		"uri": "https://idbroker.webex.com/idb/oauth2/v1/authorize?response_type=token&client_id=C71e2f13edd03a6307b9591f529345a90447d83814b6db35c26c18fc81044da2e&redirect_uri=https%3A%2F%2Fhercules.ladidadi.org%2Ffuse_redirect&scope=Identity%3ASCIM%20Identity%3AOrganization%20squared-fusion-mgmt%3Amanagement&state=YRYJUBUjJRvMCeiDy7tdD8twfHX6QPmi95Bvm-kSC6ZRNktJDasThgbP_3QpsuWgxW8_actNuO-A3ib9xh-_TMgCBgHVR06w099XWPDNhA0G_b3xVhpK9loHGUrBSPhEwHxPLocrRWsoN8lk5MxwB5U1BxhJ1lteEhfj9P_WL_OXOigunppaOK1ErWqxxxXeJE59Yp5igWmI4E9gCfMDgBuB9nwnfqOHy1IGB5UiJ-KbBVLj8cJgVgcpeSDtDa8Fej_kNO1FW34F2eQiIXajDLtNOBFwCvhm-Q6Co45poL2-jljOYIwHaNSpNaQqkWHtVaTA4kEyp_6wYPYZafqf2A=="
	}
}
```
###Get machine account

#### Request
```
http(s)://host:port/api/v1/getMachineAccount

{
  "accessToken" : "<encrypted session id>"
  "sessionId" : "<encrypted sessionId>"
}
```

#### Response
```
{
	"status": {
		"code": 200,
		"message": "all is well"
	},
	"response": {
		"username" : ""
		"password" : ""   
		"location" :  ""   
		"organization_id" : ""
		"account_id" : ""
	}
}
    
```

###Logon and Register

#### Request
```
http(s)://host:port/api/v1/logonRegister

{
  "Username"   		: "fusion-mgmnt-562df21a-a871-4cd6-933e-91da23e42b30",
  "Password"  		: "aaBB12$c6d81439-4fc1-40c1-a009-2470610d6f11",        
  "Organization_id" : "baab1ece-498c-452b-aea8-1a727413c818",
  "Account_id"    	: "9aa88e3f-4405-43b0-bc53-9246fe7f4ac0",
  "SerialNumber"	: "1234"
}
```

#### Response
```
{
	"status": {
		"code": 200,
		"message": "Logon - Register Successful"
	},
	"response": null
}
```

