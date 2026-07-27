## GOFLY LIVE CHAT
Open-source live chat support system, built for modern customer service

​​Real-time messaging​​ - Instant connection between customers and support teams

Lightning-fast performance​​ - Powered by Golang for high-concurrency handling

### Technical Architecture​

A modern stack built for performance and scalability​
 
- Backend: `gin`, `jwt-go`, `websocket`, `go.uuid`, `gorm`, `cobra`  
- Frontend: `VueJS`, `ElementUI`  
- Database: `MySQL`  

---

### Installation & Usage  

#### 1. Set Up MySQL Database  
- Install and run MySQL (version ≥ 5.5).  
- Create a database:  
```sql
  CREATE DATABASE goflychat CHARSET utf8mb4;
 ```  
*  Configure Database Connection
   Copy `config/mysql.json.demo` to `config/mysql.json` and edit it (mysql.json is git-ignored and will not be committed):
```php
{
	"Server":"127.0.0.1",
	"Port":"3306",
	"Database":"goflychat",
	"Username":"goflychat",
	"Password":"goflychat"
}
```
* Install and Configure Golang
  Run the following commands:
```php
wget https://studygolang.com/dl/golang/go1.20.2.linux-amd64.tar.gz
tar -C /usr/local -xvf go1.20.2.linux-amd64.tar.gz
mv go1.20.2.linux-amd64.tar.gz /tmp
echo "PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source /etc/profile
go version
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```
* Download the Source Code

  Clone the repository in any directory:
```php
git clone https://github.com/taoshihan1991/goflylivechat.git
cd goflylivechat  
 ```  
* Initialize the Database
 ```php
 go run main.go install
 ```  
* Run the Application
```php
 go run main.go server
 ```
* ​​Build executable
```php
 go build -o gochat
```
* ​​Run binary​​:
```php
  Linux: ./gochat server (optional flags: -p 8082 -d)
  
  Windows: gochat.exe server (optional flags: -p 8082 -d)
```  
* Terminate the Process
```php
   killall gochat
``` 

Once running, the service listens on port 8081. Access via http://[your-ip]:8081.

For domain access, configure a reverse proxy to port 8081 to hide the port number.
### Customer Service Integration
Chat Link

http://127.0.0.1:8081/livechat?customer_id=agent

Optional parameters: `refer` (source page), `extra` (base64-encoded JSON `{"visitorName":"张三","visitorAvatar":"https://xxx.jpg"}` to customize visitor nickname/avatar), and `token` (a JWT access token issued by sub2api).

Visitor identity rules:
- **With `token`** = logged-in state. The backend verifies the token **offline**: HS256 signature + `exp`/`nbf` expiry, using the same `jwt.secret` as sub2api (configure it in `config/server.json` as `Sub2apiJwtSecret`, or env `GOFLY_SUB2API_JWT_SECRET`; find the secret in sub2api's database table `security_secrets` where `key='jwt_secret'`, or the deploy machine's `JWT_SECRET`). Identity is the `user_id` inside the token, nickname defaults to its email. Invalid/expired tokens are rejected. Note: offline verification is not aware of sub2api-side revocation (password change / logout) — a token stays valid until it expires.
- **Without `token`**: the page asks the visitor to fill in an email before chatting; the email is the visitor's identity and nickname.
- Same identity (token `user_id` or email) is always recognized as the same visitor across devices and browsers.

Popup Integration

```
<script>
    (function(global, document, scriptUrl, callback) {
        const head = document.getElementsByTagName('head')[0];
        const script = document.createElement('script');
        script.type = 'text/javascript';
        script.src = scriptUrl + "/static/js/chat-widget.js";
        script.onload = script.onreadystatechange = function () {
            if (!this.readyState || this.readyState === "loaded" || this.readyState === "complete") {
                callback(scriptUrl);
            }
        };
        head.appendChild(script);
    })(window, document, "http://127.0.0.1:8081", function(baseUrl) {
        CHAT_WIDGET.initialize({
            API_URL: baseUrl,
            AGENT_ID: "agent",
        });
    });
</script>
```
### Embed into sub2api
The chat page and popup widget follow sub2api's visual style (teal `#14b8a6 → #0d9488` gradient, rounded corners, Chinese UI). To embed it into a sub2api site, paste the snippet below into the sub2api page (e.g. global custom HTML), no sub2api code changes required:

```
<script>
    (function(global, document, scriptUrl, callback) {
        const head = document.getElementsByTagName('head')[0];
        const script = document.createElement('script');
        script.type = 'text/javascript';
        script.src = scriptUrl + "/static/js/chat-widget.js";
        script.onload = script.onreadystatechange = function () {
            if (!this.readyState || this.readyState === "loaded" || this.readyState === "complete") {
                callback(scriptUrl);
            }
        };
        head.appendChild(script);
    })(window, document, "https://chat.lemonzz.xyz", function(baseUrl) {
        CHAT_WIDGET.initialize({
            API_URL: baseUrl,
            AGENT_ID: "agent",
            AUTO_OPEN: false,
            // Pass the sub2api user's access_token for logged-in identity
            // TOKEN: "<sub2api access_token>",
        });
    });
</script>
```
### Important Notice  
The use of this project for illegal or non-compliant purposes, including but not limited to viruses, trojans, pornography, gambling, fraud, prohibited items, counterfeit products, false information, cryptocurrencies, and financial violations, is strictly prohibited.  

This project is intended solely for personal learning and testing purposes. Any commercial use or illegal activities are explicitly forbidden!!!  



### Copyright Notice
This project provides full-featured code but is intended ​​only for personal demonstration and testing​​. Commercial use is strictly prohibited.

By using this software, you agree to comply with all applicable local laws and regulations. ​​You are solely responsible for any legal consequences arising from misuse.​