To curl against this API:
curl http://localhost:8080/watch?domain=yourdumbasswebsite.com

Output has 3 flavors, 

SOA has changed: 
SOA for domain yourdumbasswebsite.com has changed to ns1.example.com. hostmaster.example.com. 2 3600 1800 604800 1800

SOA Not found or dns lookup error:

HTTP/1.1 500 Internal Server Error
Content-Type: text/plain

No SOA record found:

HTTP/1.1 500 Internal Server Error
Content-Type: text/plain

No SOA record found for domain example.com

The api continuously checks and outputs every 10 minutes.  
If the SOA doesn't change the API will not return any value at all
If you choose to obsess over returning an indication that the soa hasn't changed add this to the end of the loop:

fmt.Fprintf(w, "SOA for domain %s is still %s", domain, soa)
