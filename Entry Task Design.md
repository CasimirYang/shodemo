[toc]



## FLOW

### Login flow

```mermaid
sequenceDiagram
		user ->> httpserver: login (http)
		httpserver ->> httpserver: convert password with md5+salt
    httpserver ->> + tcpserver : login (rpc)
    tcpserver ->> mysql : get user
    tcpserver ->> redis : cache user
    tcpserver ->> - httpserver: response userInfo
    httpserver ->> httpserver: generate jwt token
 		httpserver ->> user: response userInfo and token
```

### Get user flow

```mermaid
sequenceDiagram
		user ->> httpserver: request with token (http)
		httpserver ->> httpserver: check token
		httpserver ->> httpserver: decode token to get uid
    httpserver ->> + tcpserver : request(rpc)
    tcpserver ->> redis : get user by uid
    alt could not get user in redis
   tcpserver ->> mysql : get user by uid
    tcpserver ->> redis : cache user
	end
    tcpserver ->> - httpserver: response success
 		httpserver ->> user: response success
```



### Edit user flow

```mermaid
sequenceDiagram
		user ->> httpserver: request with token (http)
		httpserver ->> httpserver: check token
		httpserver ->> httpserver: decode token to get uid
    httpserver ->> + tcpserver : request(rpc)
    tcpserver ->> mysql : edit user by uid
    tcpserver ->> redis : delete cache
    tcpserver ->> - httpserver: response success
 		httpserver ->> user: response success
```



### Upload profile flow

```mermaid
sequenceDiagram
	  user ->> httpserver: request with token (http)
		httpserver ->> httpserver: check token
		httpserver ->> httpserver: decode token to get uid
		httpserver ->> httpserver: convert file to uri path
    httpserver ->> + tcpserver : edit flow 
    tcpserver ->> - httpserver: response success
 		httpserver ->> user: response success
```



## API

### common code

| Code |                                        |
| ---- | -------------------------------------- |
| 0    | success                                |
| 1    | unknown System error                   |
| 2    | InvalidParams                          |
| 3    | login fail (password or useName error) |
| 4    | InvalidToken                           |

### Login

**??????URL???** 

- ` /uc/login `

**???????????????**

- POST 

**Header???**

```js
Content-Type:application/json
```

**body ?????????** 

| ?????????   | ?????? | ??????   |
| :------- | :--- | :----- |
| userName | ???   | string |
| password | ???   | string |


**????????????**

``` json
{
 "code":0,
 "message":{
 "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxNjA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck",
 "userInfo":{
       "userName":"casi",
       "nickName":"Lily",
       "password":"123",
       "profile":"www.shopee.com/1.png"
    }
  }
}
```

### getUser

**??????URL???** 

- ` /uc/getUser `

**???????????????**

- GET 

**Header???**

```js
Content-Type:application/json
token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxN jA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck"
```

**????????????**

``` json
{
 "code":0,
 "message":{
 "userInfo":{
       "userName":"casi",
       "nickName":"Lily",
       "password":"123",
       "profile":"www.shopee.com/1.png"
    }
  }
}
```

### editUser

**??????URL???** 

- ` /uc/editUser `

**???????????????**

- POST 

**Header???**

```js
Content-Type:application/json
token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxN jA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck"
```

**body ?????????** 

| ?????????   | ?????? | ??????   |
| :------- | :--- | ------ |
| nickName | ???   | string |

**????????????**

``` json
{
 "code":0,
 "message":{
 "userInfo":{
       "userName":"casi",
       "nickName":"Lily",
       "password":"123",
       "profile":"www.shopee.com/1.png"
    }
  }
}
```

### Upload profile

**??????URL???** 

- ` /uc/uploadProfile `

**???????????????**

- POST 

**Header???**

```js
Content-Type: multipart/form-data
token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxN jA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck"
```

**body ?????????** 

| ?????????  | ?????? | ?????? |
| :------ | :--- | ---- |
| profile | ???   | file |

**????????????**

``` json
{
 "code":0,
 "message":{
 "userInfo":{
       "userName":"casi",
       "nickName":"Lily",
       "password":"123",
       "profile":"www.shopee.com/1.png"
    }
  }
}
```



## Table design

```sql
CREATE SCHEMA `spo_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ;

CREATE TABLE `user_base_info_tab` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(11) DEFAULT NULL,
  `nick_name` varchar(30) NOT NULL,
  `password` varchar(100) NOT NULL,
  `profile` varchar(100) NOT NULL,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_uname_password` (`user_name`,`password`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4
```



![image-20210706130118534](/Users/jinghua.yang/Library/Application Support/typora-user-images/image-20210706130118534.png)
