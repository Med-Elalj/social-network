//go:build !useproxy

package handlers

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Go Backend</title>
			<style>
				:root {
					--bg-dark: #1b1a17;
					--bark: #2c2a27;
					--accent-yellow: #f9c80e;
					--text-gray: #ccc;
					--text-light: #eee;
				}
				body {
					background-color: var(--bg-dark);
					color: var(--text-light);
					font-family: 'Inter', "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
				}
				.container {
					background-color: var(--bark);
					border: 2px solid var(--accent-yellow);
					padding: 2rem 3rem;
					text-align: center;
					box-shadow: 0 0 15px rgba(249, 200, 14, 0.2);
					border-radius: 2px;
				}
				h1 {
					color: var(--accent-yellow);
					font-size: 2rem;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>⚙️ The Server is Running</h1>
			</div>
		</body>
		</html>
	`)
}
