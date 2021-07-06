[toc]



## FLOW

### Login flow

```mermaid
sequenceDiagram
		user ->> httpserver: login (http)
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

### Login

**请求URL：** 

- ` /uc/login `

**请求方式：**

- POST 

**Header：**

```js
Content-Type:application/json
```

**body 参数：** 

| 参数名   | 必选 | 类型   |
| :------- | :--- | :----- |
| userName | 是   | string |
| password | 是   | string |


**返回示例**

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

**请求URL：** 

- ` /uc/getUser `

**请求方式：**

- GET 

**Header：**

```js
Content-Type:application/json
token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxN jA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck"
```

**返回示例**

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

**请求URL：** 

- ` /uc/editUser `

**请求方式：**

- POST 

**Header：**

```js
Content-Type:application/json
token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxN jA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck"
```

**body 参数：** 

| 参数名   | 必选 | 类型   |
| :------- | :--- | ------ |
| nickName | 是   | string |

**返回示例**

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

**请求URL：** 

- ` /uc/uploadProfile `

**请求方式：**

- POST 

**Header：**

```js
Content-Type: multipart/form-data
token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inhpbmd5ZSIsInBhc3N3b3JkIjoiMTIzNDU2IiwiZXhwIjoxN jA2MzAwOTQ4LCJpc3MiOiJnaW4tYmxvZyJ9.Nv6e46XYoKfRjlgCBYnajB_CIRzZKepf09cw6KP3kck"
```

**body 参数：** 

| 参数名  | 必选 | 类型 |
| :------ | :--- | ---- |
| profile | 是   | file |

**返回示例**

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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4
```



![image-20210706130118534](/Users/jinghua.yang/Library/Application Support/typora-user-images/image-20210706130118534.png)
