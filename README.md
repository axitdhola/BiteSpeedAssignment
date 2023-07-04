# Resume Link
<pre> https://drive.google.com/file/d/1tc9FfLHI3Av_ygS-GhkAMZhVIkayBI1j/view </pre>

# How To Run Code
## Local
1. Change database credential in file src/constants
2. In terminal go to src/server and run command <pre> go run main.go </pre>
3. Hit the url http://localhost:8000/identify with POST request and pass the payload 

 <pre>{
    "email":"email12@hillvalley.edu",
    "phoneNumber": "1234567"
}</pre>