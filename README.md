[![Coverage Status](https://coveralls.io/repos/github/strangnet/go-webauthn-login/badge.svg?branch=feature/build-and-test)](https://coveralls.io/github/strangnet/go-webauthn-login?branch=feature/build-and-test) 
[![Build Status](https://travis-ci.org/strangnet/go-webauthn-login.svg?branch=develop)](https://travis-ci.org/strangnet/go-webauthn-login)

# go-webauthn-login

A simple login backend implementation for webauthn.

## Endpoints

- http://localhost servers the index.html page with a test form that handles registering and login
- `GET /api/register/begin/{username}`
- `POST /api/register/finish/{username}`
- `GET /api/login/begin/{username}`
- `POST /api/login/finish/{username}`

The current implementation is based on the work of [hbolimovsky](https://github.com/hbolimovsky) in his [webauthn-example](https://github.com/hbolimovsky/webauthn-example).
