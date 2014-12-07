go-deathbycaptcha
=================

A Death by Captcha implementation for Go
(http://www.deathbycaptcha.com/)

This is a library implemented with Web API of Death By Captcha service.

Install:
	go get github.com/antonyho/go-deathbycaptcha

Usage:
	dbc := deathbycaptcha.New()
	solvedCaptcha, err := dbc.Solve("<dbcUsername>", "<dbcPassword>", "/abs/path/to/imagefile")
