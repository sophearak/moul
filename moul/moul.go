package moul

import "fmt"

type Collection struct {
	Name     string `json:"name"`
	Src      string `json:"src"`
	Svg      string `json:"svg"`
	SrcHd    string `json:"srcHd"`
	Width    int    `json:"width"`
	WidthHd  int    `json:"widthHd"`
	Height   int    `json:"height"`
	HeightHd int    `json:"heightHd"`
}

func Build() {
	template := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title><%= siteName %></title>

	<style type="text/css">
		*,
		:before,
		:after {
			box-sizing: border-box;
		}
		body {
			margin: 0;
			font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif;
			line-height: 1.5;
		}
		.cover {
			position: relative;
			height: 500px;
		}
		.cover-wrap {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			text-align: center;
		}
		.cover-wrap picture {
			position: relative;
			display: block;
			height: 100%;
			z-index: -1;
		}
		.cover-wrap picture img {
			width: 100%;
			height: 100%;
			-o-object-fit: cover;
			object-fit: cover;
		}
		.profile {
			text-align: center;
			margin: 32px 0 16px;
			font-size: 0;
		}
		.profile a img {
			width: 120px;
			height: 120px;
			border-radius: 50%;
			border: 2px solid transparent;
			transition: border 250ms cubic-bezier(0.4, 0, 0.2, 1), box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1);
		}
		.profile a img:hover {
			box-shadow: 0 1px 2px 0 rgba(60,64,67,0.302), 0 2px 6px 2px rgba(60,64,67,0.149);
		}
		.name {
			margin: 0 0 8px;
			text-align: center;
			font-weight: 900;
			color: #111;
			line-height: 1;
		}
		.info {
			max-width: 500px;
			text-align: center;
			margin: 0 auto 32px;
			color: #444;
		}
		.collection {
			position: relative;
			margin: 10px auto;
		}
		.collection figure {
			margin: 0;
		}
		.collection figure a {
			display: block;
			font-size: 0;
			float: left;
		}
		.collection figure a img {
			position: absolute;
			background-size: cover;
		}
	</style>
</head>
<body>
    <div class="cover">
        <div class="cover-wrap">
            <picture>
                <img class="js-img" src="" alt="">
            </picture>
        </div>
    </div>
    <div class="profile">
        <a href="">
            <img class="js-img" src="" alt="">
        </a>
    </div>
    <h1 class="name"><%= siteName %></h1>
    <p class="info"></p>
	<div id="collection"><div>
	<input type="hidden" id="photo-collection">

	<script type="text/javascript">
		const $ = document.querySelector.bind(document);
		const $$ = document.querySelectorAll.bind(document);

		document.addEventListener("DOMContentLoaded", () => {
			const photoTags = $$('.js-img');
			photoTags.forEach(photo => {
				const sizes = Math.ceil(
				  (photo.getBoundingClientRect().width / window.innerWidth) * 100
				);
				photo.setAttribute('sizes', sizes + 'vw');
			});
		});
	</script>
</body>
</html>`
	fmt.Print(template)
}
