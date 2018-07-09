BOX agent API

1. http services

1.1. Service messages

- uri:  /agent/msinfo
- Method：GET、POST
- Parameters

- Response

    {
      "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
      "ManagerIpPort": "192.168.199.181:5001"
    }

- Error Codes

| Codes        | Message    |
| ------------ |:----------:|
| 11           | JSON error |



1.2. app register

- uri:  /agent/registadd
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| regid           | string   | registration id, required.               |
| msg             | string   | message after encryption, required.      |
| regid           | string   | registration id, required.               |
| applyerid       | string   | applicant id, required.                  |
| captainid       | string   | direct manager id, required.             |
| applyeraccount  | string   | application result, 1-reject 2-approved. |
| msg             | string   |  |
| status          | string   |  |

- Response

    {
      "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message       |
| ------------ |:-------------:|
| 11           | JSON error    |
| 12           | leveldb error |


1.3. app register info query

- uri:  /agent/registlist
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| applyerid       | string   | applicant id, required.                  |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
        "RegistInfos": [//regist list
                {
                    "RegId": "1",
                    "ApplyerId": "2",
                    "CaptainId":"3",
                    "ApplyerAccount":"4",
                    "Msg":"wrwwqr",
                    "Status":"0"
                }
          ]
    }

- Error Codes

| Codes        | Message       |
| ------------ |:-------------:|
| 11           | JSON error    |
| 12           | leveldb error |



1.4. app approval create

- uri:  /agent/approvaladd
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash id, required.                       |
| name            | string   | name of approval.                        |
| appid           | string   | employee id.                             |
| captainid       | string   | PK id.                                   |
| flow            | string   | approval request data.                   |
| sign            | string   | signature.                               |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message       |
| ------------ |:-------------:|
| 11           | JSON error    |
| 12           | leveldb error |



1.5.app approval detail

- router:  /agent/approvaldetail
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash id, required.                       |

- Response

    {
      "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message       |
| ------------ |:-------------:|
| 11           | JSON error    |
| 13           | db error      |

 

1.6. app approval query

- router:  /agent/approvallist
- Method：GET、POST

- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash id, required(66 characters leading with 0x).|

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
        "ApprovalInfo": {
      "Hash": "0x240376c64624b0ff91f66115d2b52e6be47b69459199ebadd038d44a34fd921c",
      "Name": "测试模板一",
      "AppId": "1523758482398",
      "Flow": "{\n  \"flow_name\" : \"测试模板一\",\n  \"single_limit\" : \"500\",\n  \"approval_info\" : [\n    {\n      \"require\" : 1,\n      \"total\" : 2,\n      \"approvers\" : [\n        {\n          \"account\" : \"Ghkjbgjj3\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAum3CasPP3NMhIyAwXhmLuE6x5Ijn0lKqHxPWqD\\/IqH7DAtUhVgDHCGPJXhxU2djZhZKw2H6o1Ar+WxN6xeLfqTIE95Eii8fDJhCPO18P+Ia2LNoUjxPTPL+4sR3KJ0AP4qIwNoBVhv7t37MIa94Q+cHpOz0eRzSkMRFtbJlDmQEZSU5Y\\/aa3YRvIFx92wergyR+w1rqk4myYt1rcURckG8HEt5rzW6fSMzTeKt3VLQNimBZMANWsJQU3PETvFrDbiA4X2M3Qore\\/Vb6AzL965c5sdJdequLser0ctHhqdt6Hr7SINtj+e3jtOKQ3hUtKqAw6Gj8b6+tommjQitIxAwIDAQAB\",\n          \"app_account_id\" : \"1523767732096\"\n        },\n        {\n          \"account\" : \"Ghhjjjk3\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAkYbghJy1CInvX2MKJFHefRk60WijjDBdOF103ZR6FC2xuP7OClQXSg0kpzk78kkBztJfwz0WXSBIzRK77u+VLbhjyh+Fs8vEn\\/tPQMxvAP4LSMaphhIzAqFmioK+J4pykQEkXuUPIVJsq9X0rPGZMGy8GAd2kZtN9cWZktvmB31svKi2XTkRpa+AvsIl+0dEt4EQa8dn+QGwiXvpxAIH6GsK2HTqIytvNOsM2tC+h8335DUXGyuVOWZi4eYx\\/wTYVMVnqFjES1n1wkazGJ0mr0LF6iT9MChHiDNLWyynnjk6++V4kteuIjM0aCxiPeE0mF\\/pZNHsYnMYLmDGVv4\\/xwIDAQAB\",\n          \"app_account_id\" : \"1523770117303\"\n        }\n      ]\n    },\n    {\n      \"require\" : 1,\n      \"total\" : 1,\n      \"approvers\" : [\n        {\n          \"account\" : \"Bbbbbbb2\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAizTN6hpf6CqeXbXkw+SIIGLdDqXkj8dYDN0U7yIxA0YK9+HRjoSKyVVZ6aRe0kzVQI0MeLv7VaZibcCQuSN6F7B39XuZtpp21rKceHCqgCWRm3U9eJJa0AhA7UtEcxps7mEAr+Lxjlkd7KxxM7dOQylb6A7FZlZ4lqwbgP9R+KtkE2XSVDOlW94j5m\\/auAf86HOPTW\\/oNAAps3Yd0k\\/1DoI72+hEbVh3hy5fFpJWLtPVrtbNSPujmfXdkBmgdDJarbehsMYX\\/tkpZJO39GfD29TBj7IASL00IP3T2QDn0ZTwZOkC+mCUNLtwEqnBINavJhzsVuK1aCNfEc7g+kfB2QIDAQAB\",\n          \"app_account_id\" : \"1523766717563\"\n        }\n      ]\n    },\n    {\n      \"require\" : 1,\n      \"total\" : 1,\n      \"approvers\" : [\n        {\n          \"account\" : \"Aaaaaa\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAzlnIjGRBsyz\\/a1w\\/jfppqX4h0M+8QmFHWDQmm2DDvoipjZul7MzBAxcHUedlE\\/ci1vUhP+XyWBJ5pshp5xMQysUIau79a89Lrzf53GG6bx1wIZfLfifmNRyvGDyuJ0URjrCixPJoUrP7xJXJKsdLGHSqDlgXNWMYMJabw85rzP588Oj8w22W9VI8v88zQjOA01u1LY8tD\\/ud88dCGMMTULPG1QyQuL\\/n8PlmyWCbWUZvv7f9ZFvp3ApvtyibtbaJlKjAgkJ8idJEzVr+EZuSo9i6zuckK+MvSB3W\\/+lqNeXbE8dglK7mlYxZEMclyuq+4n8sKwUfzwOKYv2JK65g9wIDAQAB\",\n          \"app_account_id\" : \"1523766595967\"\n        }\n      ]\n    }\n  ]\n}",
      "Sign": "cUve5VLqDwL74PlZ4uHkp44BJW5InNA+xI+5yTaLdAXM7HOhVpni67ogjP9KAydwpLd7hb4Em72d+ihyoujdgObJO79nxftdBIjOLKqRMW650eP+PgWKEmEGSBCLWb1p0J5Dj4AxhEBWqaWne4qelYhYgSEQTNHgKYewENtI/lzrJFzhIC+ZN21tZZsFAMWGCUk2wOTVUw0/evVJLAJmOOC9PoAomVQD31Hdv4QtlT4C/91zQsK0dAyhUgUxoCknheITwGAfRkmyqkUZtIvLFpha89mlfb0JnrJHQRZ5GhdzTE7Qh7upk610+RKn1A6TrtrLvBIucWKrWA9HAwe2wA==",
      "Status": "0"
      },
    }

- Error Codes：

| Codes        | Message       |
| ------------ |:-------------:|
| 11           | JSON error    |
| 12           | leveldb error |



1.7. Approval flow PK create

- uri:  /agent/hashadd
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash id, required.                       |

- Response

    {
      "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message       |
| ------------ |:-------------:|
| 11           | JSON error    |
| 13           | db error      |
 

1.8. Create public key

- uri:  /agent/keystore
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| applyerid       | string   | applicant id, required.                  |
| publickey       | string   | public key, required.                    |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |


1.9. voucher operations

- uri:  /agent/operate
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| password        | string   | password/keyword, required.              |
| applyerid       | string   | applicant id, required.                  |
| type.           | string   | 0-add public key;1-create;2-publish;3-start;4-stop|

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |



1.10. Voucher status query

- uri:  /agent/status
- Method：GET、POST
- Parameters

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
        "Status": {
            "ServerStatus": 3,//0-unconnecte 1-uncreated 2-created 3-published 4-started
            "Status": 0,//error code of status
            "Total": 3,//number of private keys
            "HashCount": 0,//number of approval flows
            "Address": "0x5B3538f942bEAE24ed6987360193c77747D1d77d",//account address
            "ContractAddress": "0xB3a435c0329C95752858476E40F9e8fbeF292B23",//contract address
            "D": 2258783391,//random number
            "NodesAuthorized": [//Authorization status
                {
                    "ApplyerId": "1523340972122",
                    "Authorized": true
                },
                {
                    "ApplyerId": "1523340990911",
                    "Authorized": true
                },
                {
                    "ApplyerId": "1523340886398",
                    "Authorized": true
                }
            ],
            "KeyStroeStatus": [
                {
                    "ApplyerId": "1523340886398",
                    "ApplyerName": "Uuu"
                },
                {
                    "ApplyerId": "1523340972122",
                    "ApplyerName": "Igggbb"
                },
                {
                    "ApplyerId": "1523340990911",
                    "ApplyerName": "Pppp"
                }
            ]
        }
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |


1.11. approve flow

- uri:  /agent/allow
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash id, required(66 characters leading with 0x).|

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- 错误代码：

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |



1.12. reject flow

- uri:  /agent/disallow
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash id, required(66 characters leading with 0x).|

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |



1.13. approval of flow query

- uri:  /agent/approvallist
- Method：GET、POST
- Parameters

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
     
        {
          "Hash": "0x0c50cc12c7bb7531f2f7ac555024cf8210f2ccc2464b0af2374da35226d1e8ed",
          "Name": "模板一",
          "AppId": "12345",
          "Flow": "",
          "Sign": "",
          "Status": "1"
          },
          {
          "Hash": "hash",
          "Name": "来咯哦哦",
          "AppId": "a",
          "Flow": "",
          "Sign": "",
          "Status": "0"
          }
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |



1.14. create token

- uri:  /agent/tokenedit
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| tokenname       | string   | name of token, required.                 |
| decimals        | string   | |
| contractaddr    | string   | address of token contract                |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |

1.15. query token

- uri:  /agent/tokenlist
- Method：GET、POST
- Parameters

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 101          | illegal hash prefix |

1.16. operation on currency type

- uri:  /agent/coin
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| category        | int      | category of currency, required. 0-BTC    |
| used            | string   | in use or not. 0-disable;1-in use.       |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Code

| Codes        | Message             |
| ------------ |:-------------------:|
| 11           | parse error.        |
| 103          | invalid amount.     |

1.17. currency query

- uri:  /agent/coinlist
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| category        | int      | category of currency, required. 0-BTC    |
| used            | string   | in use or not. 0-disable;1-in use.       |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
        "CoinStatus": [//coin list
                {
                    "Name": "BTC",
                    "Category": 0
                }
          ]
    }

- Error Codes

| Codes        | Message             |
| ------------ |:-------------------:|
| 11           | parse error.        |
| 103          | invalid amount.     |


1.18. Withdraw request

- uri:  /agent/wtihdraw
- Method：GET、POST
- Parameters

| Field           | Type     | Description                              |
| ------------    |:--------:|:----------------------------------------:|
| hash            | string   | hash of approval flow.                   |
| wdhash          | string   | hash of withdrawl request                |
| category        | int      | type of currency, required. 0-BTC;1-ETH  |
| recaddress      | string   | address of receiver.                     |
| amount          | string   | amount.                                  |
| fee             | string   | fee.                                     |
| flow            | string   | request data.                            |
| sign            | string   | signature.                               |

- Response

    {
        "RspNo": "0",//0-Succeed, other codes see section "Error Codes"
    }

- Error Code

| Codes        | Message             |
| ------------ |:-------------------:|
| 11           | parse error.        |
| 103          | invalid amount.     |






BOX agent接口

一、http服务

1.1.管理端服务信息*

- router:  /agent/msinfo
- 请求方式：GET、POST
- 参数

- 返回值

    {
      "RspNo": "0",//0-成功，其他见错误代码
      "ManagerIpPort": "192.168.199.181:5001"
    }

- 错误代码

  code	message 
   11 	JSON处理失败



1.2.app注册*

- router:  /agent/registadd
- 请求方式：GET、POST
- 参数

        字段      	  类型  	     备注      
      regid     	string	 服务端申请表ID，必输 
       msg      	string	 加密后的注册信息，必输 
    applyerid   	string	   申请者，必输    
    captainid   	string	   直属上级，必输   
  applyeraccount	string	审批结果， 1拒绝 2同意
       msg      	string	             
      status    	string	             

- 返回值

    {
      "RspNo": "0",//0-成功，其他见错误代码
    }

- 错误代码

  code	  message  
   11 	 JSON处理失败  
   12 	leveldb处理失败



1.3. app注册信息查询*

- router:  /agent/registlist
- 请求方式：GET、POST
- 参数

     字段    	  类型  	 备注  
  applyerid	string	审批者id

- 
- 返回值

    {
        "RspNo": "0",//0-成功，其他见错误代码
        "RegistInfos": [//regist list
                {
                    "RegId": "1",
                    "ApplyerId": "2",
                  	"CaptainId":"3",
                  	"ApplyerAccount":"4",
                  	"Msg":"wrwwqr",
                  	"Status":"0"
                }
          ]
    }

- 错误代码：

  code	  message  
   11 	 JSON处理失败  
   12 	leveldb处理失败



1.4. app审批申请*

- router:  /agent/approvaladd
- 请求方式：GET、POST
- 参数

     字段    	  类型  	   备注    
    hash   	string	hashId，必输
    name   	string	   名称    
    appid  	string	  员工id   
  captainid	string	  私钥id   
    flow   	string	  原始数据   
    sign   	string	   签名    

- 返回值

    {
        "RspNo": "0",//0-成功，其他见错误代码
    }

- 错误代码：

  code	  message  
   11 	 JSON处理失败  
   12 	leveldb处理失败



1.5.app审批详情*

- router:  /agent/approvaldetail
- 请求方式：GET、POST
- 参数

   字段 	  类型  	   备注   
  hash	string	hash模板id

- 返回值

    {
      "RspNo": "0",//0-成功，其他见错误代码
    }

- 错误代码

  code	message 
   11 	JSON处理失败
   13 	 db处理失败 

 

1.6. app审批查询*

- router:  /agent/approvallist
- 请求方式：GET、POST

- 参数

   字段  	  类型  	        备注         
  appid	string	hash值，必输（带0x的66位字符）

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
      	"ApprovalInfo": {
      "Hash": "0x240376c64624b0ff91f66115d2b52e6be47b69459199ebadd038d44a34fd921c",
      "Name": "测试模板一",
      "AppId": "1523758482398",
      "Flow": "{\n  \"flow_name\" : \"测试模板一\",\n  \"single_limit\" : \"500\",\n  \"approval_info\" : [\n    {\n      \"require\" : 1,\n      \"total\" : 2,\n      \"approvers\" : [\n        {\n          \"account\" : \"Ghkjbgjj3\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAum3CasPP3NMhIyAwXhmLuE6x5Ijn0lKqHxPWqD\\/IqH7DAtUhVgDHCGPJXhxU2djZhZKw2H6o1Ar+WxN6xeLfqTIE95Eii8fDJhCPO18P+Ia2LNoUjxPTPL+4sR3KJ0AP4qIwNoBVhv7t37MIa94Q+cHpOz0eRzSkMRFtbJlDmQEZSU5Y\\/aa3YRvIFx92wergyR+w1rqk4myYt1rcURckG8HEt5rzW6fSMzTeKt3VLQNimBZMANWsJQU3PETvFrDbiA4X2M3Qore\\/Vb6AzL965c5sdJdequLser0ctHhqdt6Hr7SINtj+e3jtOKQ3hUtKqAw6Gj8b6+tommjQitIxAwIDAQAB\",\n          \"app_account_id\" : \"1523767732096\"\n        },\n        {\n          \"account\" : \"Ghhjjjk3\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAkYbghJy1CInvX2MKJFHefRk60WijjDBdOF103ZR6FC2xuP7OClQXSg0kpzk78kkBztJfwz0WXSBIzRK77u+VLbhjyh+Fs8vEn\\/tPQMxvAP4LSMaphhIzAqFmioK+J4pykQEkXuUPIVJsq9X0rPGZMGy8GAd2kZtN9cWZktvmB31svKi2XTkRpa+AvsIl+0dEt4EQa8dn+QGwiXvpxAIH6GsK2HTqIytvNOsM2tC+h8335DUXGyuVOWZi4eYx\\/wTYVMVnqFjES1n1wkazGJ0mr0LF6iT9MChHiDNLWyynnjk6++V4kteuIjM0aCxiPeE0mF\\/pZNHsYnMYLmDGVv4\\/xwIDAQAB\",\n          \"app_account_id\" : \"1523770117303\"\n        }\n      ]\n    },\n    {\n      \"require\" : 1,\n      \"total\" : 1,\n      \"approvers\" : [\n        {\n          \"account\" : \"Bbbbbbb2\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAizTN6hpf6CqeXbXkw+SIIGLdDqXkj8dYDN0U7yIxA0YK9+HRjoSKyVVZ6aRe0kzVQI0MeLv7VaZibcCQuSN6F7B39XuZtpp21rKceHCqgCWRm3U9eJJa0AhA7UtEcxps7mEAr+Lxjlkd7KxxM7dOQylb6A7FZlZ4lqwbgP9R+KtkE2XSVDOlW94j5m\\/auAf86HOPTW\\/oNAAps3Yd0k\\/1DoI72+hEbVh3hy5fFpJWLtPVrtbNSPujmfXdkBmgdDJarbehsMYX\\/tkpZJO39GfD29TBj7IASL00IP3T2QDn0ZTwZOkC+mCUNLtwEqnBINavJhzsVuK1aCNfEc7g+kfB2QIDAQAB\",\n          \"app_account_id\" : \"1523766717563\"\n        }\n      ]\n    },\n    {\n      \"require\" : 1,\n      \"total\" : 1,\n      \"approvers\" : [\n        {\n          \"account\" : \"Aaaaaa\",\n          \"itemType\" : 0,\n          \"pub_key\" : \"MIIBCgKCAQEAzlnIjGRBsyz\\/a1w\\/jfppqX4h0M+8QmFHWDQmm2DDvoipjZul7MzBAxcHUedlE\\/ci1vUhP+XyWBJ5pshp5xMQysUIau79a89Lrzf53GG6bx1wIZfLfifmNRyvGDyuJ0URjrCixPJoUrP7xJXJKsdLGHSqDlgXNWMYMJabw85rzP588Oj8w22W9VI8v88zQjOA01u1LY8tD\\/ud88dCGMMTULPG1QyQuL\\/n8PlmyWCbWUZvv7f9ZFvp3ApvtyibtbaJlKjAgkJ8idJEzVr+EZuSo9i6zuckK+MvSB3W\\/+lqNeXbE8dglK7mlYxZEMclyuq+4n8sKwUfzwOKYv2JK65g9wIDAQAB\",\n          \"app_account_id\" : \"1523766595967\"\n        }\n      ]\n    }\n  ]\n}",
      "Sign": "cUve5VLqDwL74PlZ4uHkp44BJW5InNA+xI+5yTaLdAXM7HOhVpni67ogjP9KAydwpLd7hb4Em72d+ihyoujdgObJO79nxftdBIjOLKqRMW650eP+PgWKEmEGSBCLWb1p0J5Dj4AxhEBWqaWne4qelYhYgSEQTNHgKYewENtI/lzrJFzhIC+ZN21tZZsFAMWGCUk2wOTVUw0/evVJLAJmOOC9PoAomVQD31Hdv4QtlT4C/91zQsK0dAyhUgUxoCknheITwGAfRkmyqkUZtIvLFpha89mlfb0JnrJHQRZ5GhdzTE7Qh7upk610+RKn1A6TrtrLvBIucWKrWA9HAwe2wA==",
      "Status": "0"
      },
    }

- 错误代码：

  code	  message  
   11 	 JSON处理失败  
   12 	leveldb处理失败



1.7.审批流私钥申请*

- router:  /agent/hashadd
- 请求方式：GET、POST
- 参数

   字段 	  类型  	   备注   
  hash	string	hash模板id

- 返回值

    {
      "RspNo": "0",//0-成功，其他见错误代码
    }

- 错误代码

  code	message 
   11 	JSON处理失败
   13 	 db处理失败 

 

1.8. 公钥添加*

- router:  /agent/keystore
- 请求方式：GET、POST
- 参数

     字段    	  类型  	   备注   
  applyerid	string	appId，必输
  publickey	string	 公钥，必输  

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message 
  101 	非法hash前缀



1.9. 签名机操作*

- router:  /agent/operate
- 请求方式：GET、POST
- 参数

     字段    	  类型  	            备注            
  password 	string	          关键句，必输          
  applyerId	string	         appId，可选         
    type   	string	0-添加公钥 1-创建 2-发布 3-启动 4-停止

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message 
  101 	非法hash前缀



1.10. 签名机状态查询*

- router:  /agent/status
- 请求方式：GET、POST
- 参数

- 返回值

    {
        "RspNo": "0",//0-成功，其他见错误代码
        "Status": {
            "ServerStatus": 3,//0-未连接 1-未创建 2-已创建 3-已发布 4-已启动
            "Status": 0,//错误码状态
            "Total": 3,//密钥数量
            "HashCount": 0,//审批流数量
            "Address": "0x5B3538f942bEAE24ed6987360193c77747D1d77d",//账户地址
            "ContractAddress": "0xB3a435c0329C95752858476E40F9e8fbeF292B23",//合约地址
            "D": 2258783391,//随机数
            "NodesAuthorized": [//授权情况
                {
                    "ApplyerId": "1523340972122",
                    "Authorized": true
                },
                {
                    "ApplyerId": "1523340990911",
                    "Authorized": true
                },
                {
                    "ApplyerId": "1523340886398",
                    "Authorized": true
                }
            ],
            "KeyStroeStatus": [
                {
                    "ApplyerId": "1523340886398",
                    "ApplyerName": "Uuu"
                },
                {
                    "ApplyerId": "1523340972122",
                    "ApplyerName": "Igggbb"
                },
                {
                    "ApplyerId": "1523340990911",
                    "ApplyerName": "Pppp"
                }
            ]
        }
    }

- 错误代码：

  code	message 
  101 	非法hash前缀



1.11. 审批流同意*

- router:  /agent/allow
- 请求方式：GET、POST
- 参数

   字段 	  类型  	        备注         
  hash	string	hash值，必输（带0x的66位字符）

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message 
  101 	非法hash前缀



1.12. 审批流拒绝*

- router:  /agent/disallow
- 请求方式：GET、POST
- 参数

   字段 	  类型  	        备注         
  hash	string	hash值，必输（带0x的66位字符）

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message 
  101 	非法hash前缀



1.13. 审批模板查询*

- router:  /agent/approvallist
- 请求方式：GET、POST
- 参数

- 返回值

    {
        "RspNo": "0",//0-成功，其他见错误代码
     
        {
          "Hash": "0x0c50cc12c7bb7531f2f7ac555024cf8210f2ccc2464b0af2374da35226d1e8ed",
          "Name": "模板一",
          "AppId": "12345",
          "Flow": "",
          "Sign": "",
          "Status": "1"
          },
          {
          "Hash": "hash",
          "Name": "来咯哦哦",
          "AppId": "a",
          "Flow": "",
          "Sign": "",
          "Status": "0"
          }
    }

- 错误代码：

  code	message 
  101 	非法hash前缀



1.14. 代币添加*

- router:  /agent/tokenedit
- 请求方式：GET、POST
- 参数

       字段     	  类型  	  备注   
   tokenname  	string	代币名称，必输
    decimals  	 int  	  精度   
  contractaddr	string	代币合约地址 
              	      	       

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message 
  101 	非法hash前缀

1.15. 代币查询*

- router:  /agent/tokenlist
- 请求方式：GET、POST
- 参数

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message 
  101 	非法hash前缀

1.16. 币种操作*

- router:  /agent/coin
- 请求方式：GET、POST
- 参数

     字段   	  类型  	      备注      
  category	 int  	币中分类，必输 0-BTC 
    used  	string	是否使用，0-禁用 1-使用

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message
   11 	 解析失败  
  103 	 非法金额  

1.17. 币种查询*

- router:  /agent/coinlist
- 请求方式：GET、POST
- 参数

     字段   	  类型  	      备注      
  category	 int  	币中分类，必输 0-BTC 
    used  	string	是否使用，0-禁用 1-使用

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
      	"CoinStatus": [//coin list
                {
                    "Name": "BTC",
                    "Category": 0
                }
          ]
    }

- 错误代码：

  code	message
   11 	 解析失败  
  103 	 非法金额  



1.18. 转账申请*

- router:  /agent/wtihdraw
- 请求方式：GET、POST
- 参数

      字段    	  类型  	        备注         
     hash   	string	    hash审批流模板号     
    wdhash  	string	       转账申请号       
   category 	 int  	币种分类，必输 0-BTC 1-ETH
  recaddress	string	       接受地址        
    amount  	string	        金额         
     fee    	string	        手续费        
     flow   	string	       原始数据        
     sign   	string	       签名数据        

- 返回值

    {
        "RspNo": "0"//0-成功，其他见错误代码
    }

- 错误代码：

  code	message
   11 	 解析失败  
  103 	 非法金额  


