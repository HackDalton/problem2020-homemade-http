# HackDalton: Homemade HTTP (Writeup)

> Warning! There are spoilers ahead

When we connect to the server, we are immediatly prompted with an ASCII art Captcha:

```
Solve this captcha to connect:
   ___            ____
  / _ \     _    | ___|
 | | | |  _| |_  |___ \
 | |_| | |_   _|  ___) |
  \___/    |_|   |____/
```

Upon correctly solving the captcha, we recieve a message from the server:
```
Now accepting HTTP/1.1 requests.
```

HTTP, short for Hypertext Transfer Protocol, is a protocol used for data communication on the internet. You can read more about HTTP on [Wikipedia](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol), or read [the specification](https://tools.ietf.org/html/rfc7230) from the Internet Engineering Task Force (IETF). HTTP is the foundation that your browser (like Chrome or Firefox) uses to communicate with websites.

In this puzzle, you need to handwrite HTTP requests instead of using your browser. We know from the problem description that we need to access the website "problems.hackdalton.com."

We can start by making a `GET` request for `http://problems.hackdalton.com/`:
```http
GET / HTTP/1.1
Host: problems.hackdalton.com


```

Let's break apart what each line of this request does:
- `GET / HTTP/1.1`
    - `GET` specifies that we're making a [GET request](https://tools.ietf.org/html/rfc2068#section-9.3), generally used to retreive information from the server.
    - `/` specifies that we're requesting the `/` page from the server.
    - `HTTP/1.1` specifies that we're using HTTP version 1.1.
- `Host: problems.hackdalton.com`
    - This is the [Host header](https://tools.ietf.org/html/rfc2068#section-14.23). It specifies that we're requesting the page from the server at problems.hackdalton.com.

The server immediatly responds to us!
```http
HTTP/1.1 200 OK
Content-Length: 315
Accept-Ranges: bytes
Content-Type: text/html; charset=utf-8
Last-Modified: Sun, 17 May 2020 20:38:36 GMT

<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Cool Site</title>
</head>

<body>
    <h1>Hi!</h1>
    <p>Welcome to my super cool site! You can log in to continue:</p>
    <a href="/login.html">Login!</a>
</body>

</html>
```
We can also break apart the response line by line.
- `HTTP/1.1 200 OK`
    - `HTTP/1.1` again is tells us that the server is using HTTP version 1.1.
    - `200 OK` tells us that the server is responding with [status code `200`](https://tools.ietf.org/html/rfc2068#section-10.2.1), which means `OK`.
- `Content-Length: 315`
    - This is the [Content-Length header](https://tools.ietf.org/html/rfc2616#section-14.13). It says how long the response body will be so that the client reading the message will know when it can stop listening for more.
- `Accept-Ranges: bytes`
    - This is the [Accept-Ranges header](https://tools.ietf.org/html/rfc2616#section-14.5). It tells the client what types of responses the server can understand. In this example, the server can receive information in the `bytes` range (numbers 0<sub>10</sub>-255<sub>10</sub>, or 0<sub>2</sub>-11111111<sub>2</sub>)
- `Content-Type: text/html; charset=utf-8`
    - This is the [Content-Type header](https://tools.ietf.org/html/rfc2616#section-14.17). It tells the client what type of file the server is sending and how it is encoded. `text/html` means that it is a [Hypertext Markup Language (HTML)](https://en.wikipedia.org/wiki/HTML) file, and `charset=utf-8` means that the file is encoded with [UTF-8](https://en.wikipedia.org/wiki/UTF-8).
- `Last-Modified: Sun, 17 May 2020 20:38:36 GMT`
    - This is the [Last-Modified header](https://tools.ietf.org/html/rfc2616#section-14.29). It tells the client when the file being served was last changed. It is used for [caching](https://en.wikipedia.org/wiki/Cache_(computing)).

Everything below this is the response body, which in this case is an HTML file.

```html
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Cool Site</title>
</head>

<body>
    <h1>Hi!</h1>
    <p>Welcome to my super cool site! You can log in to continue:</p>
    <a href="/login.html">Login!</a>
</body>

</html>
```

We can see that there is an [`<a>` tag](https://www.w3.org/TR/2017/REC-html52-20171214/textlevel-semantics.html#the-a-element) in the html file, which represents a link. The [`href` attribute](https://www.w3.org/TR/2017/REC-html52-20171214/links.html#element-attrdef-a-href) is set to the destination of the link, which is `/login.html`. We can then make another request to find the content of that page.

```http
GET /login.html HTTP/1.1
Host: problems.hackdalton.com
```
Thes server responds:
```http
HTTP/1.1 200 OK
Content-Length: 519
Accept-Ranges: bytes
Content-Type: text/html; charset=utf-8
Last-Modified: Sun, 17 May 2020 17:23:13 GMT

<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login | My Cool Site</title>
</head>

<body>
    <h1>Log in</h1>
    <p>Use your username and password to login to the site.</p>
    <form action="/login" method="post">
        <input type="text" name="username" placeholder="Username">
        <input type="password" name="password" placeholder="Password">
        <input type="submit" value="Submit">
    </form>
</body>

</html>
```
This is the same as the last request. We can now look at the content of the HTML file. The [`<form>` element](https://www.w3.org/TR/2017/REC-html52-20171214/sec-forms.html#the-form-element) is used to display a login form to the user. We can see that the [`method` attribute](https://www.w3.org/TR/2017/REC-html52-20171214/sec-forms.html#element-attrdef-form-method) is set to `post`, which means that the form will make an [HTTP `POST` request](https://tools.ietf.org/html/rfc2616#section-9.5) to `/login`, which is set as the form's [`action` attribute](https://www.w3.org/TR/2017/REC-html52-20171214/sec-forms.html#element-attrdef-form-action).

We need to mimick the action that the form will take when filled out correctly. This can be done by making a [POST request](https://tools.ietf.org/html/rfc2616#section-9.5) to the server.

```http
POST /login HTTP/1.1
Content-Type: application/x-www-form-urlencoded; charset=utf-8
Host: problems.hackdalton.com
Content-Length: 30

username=admin&password=secret
```

We can again go through this line by line.
- `POST /login HTTP/1.1`
    - `POST` specifies that we're making a [POST request](https://tools.ietf.org/html/rfc2616#section-9.5), generally used to send updated data to the server.
- `Content-Type: application/x-www-form-urlencoded; charset=utf-8`
    - This is the Content-Type header that we previously saw in the response. This time, we're using [`application/x-www-form-urlencoded`](https://www.w3.org/TR/html401/interact/forms.html#h-17.13.4.1), which is the format use by default by HTML `<form>` elements.
- `Host: problems.hackdalton.com`
    - This is the same as our last request.
- `Content-Length: 30`
    - We have to specify the length of our response in bytes so that the server knows when we're done sending it.

Below that is the body in `application/x-www-formurlencoded` with the following values:
| key | value |
|----------|--------|
| username | admin |
| password | secret |

The server responds with our flag, this time with Content-Type text/plain, because it is just plain text.

```http
HTTP/1.1 200 OK
Content-Length: 43
Content-Type: text/plain

hackDalton{wh0_n33ds_4_br0ws3r_4nyw4y_ebxei2eYDA}
```