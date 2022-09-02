package main

import "os"

func buildStyle() {
	data := `@font-face {
			font-family: 'proxima-nova';
			src: url('proximanova-regular.woff') format('woff'),
				 url('proximanova-regular.woff2') format('woff2');
			font-weight: normal;
			font-style: normal;
		}
		
		@font-face {
			font-family: 'proxima-nova';
			src: url('proximanova-italic.woff') format('woff'),
				 url('proximanova-italic.woff2') format('woff2');
			font-weight: normal;
			font-style: italic;
		}
		
		@font-face {
			font-family: 'proxima-nova';
			src: url('proximanova-bold.woff') format('woff'),
				 url('proximanova-bold.woff2') format('woff2');
			font-weight: bold;
			font-style: normal;
		}
		
		@font-face {
			font-family: 'proxima-nova';
			src: url('proximanova-boldit.woff') format('woff'),
				 url('proximanova-boldit.woff2') format('woff2');
			font-weight: bold;
			font-style: italic;
		}
		
		@font-face {
			font-family: 'proxima-nova';
			src: url('proximanova-light.woff') format('woff'),
				 url('proximanova-light.woff2') format('woff2');
			font-weight: 300;
		}
		
		@font-face {
			font-family: 'myriad-pro';
			src: url('mp.woff') format('woff'),
				 url('mp.woff2') format('woff2');
			font-weight: normal;
			font-style: normal;
		}
		
		@font-face {
			font-family: 'myriad-pro';
			src: url('mp-italic.woff') format('woff'),
				 url('mp-italic.woff2') format('woff2');
			font-weight: normal;
			font-style: italic;
		}
		
		@font-face {
			font-family: 'myriad-pro';
			src: url('mp-bold.woff') format('woff'),
				 url('mp-bold.woff2') format('woff2');
			font-weight: bold;
			font-style: normal;
		}
		
		@font-face {
			font-family: 'myriad-pro';
			src: url('mp-boldit.woff') format('woff');
			font-weight: bold;
			font-style: italic;
		}
		
		@font-face {
			font-family: 'myriad-pro';
			src: url('mp-semibold.woff') format('woff'),
				 url('mp-semibold.woff2') format('woff2');
			font-weight: 600;
			font-style: normal;
		}
		
		body {
			background: #000;
			color: #fff;
			font-family: "proxima-nova", "Helvetica Neue", Helvetica, Arial, "Lucida Grande", "Lucida Sans", Tahoma, sans-serif; 
		}
		
		#container {
			width: 800px;
			height: 480px;
			margin: 0px auto;
			border: 1px solid #000;
			position: relative;
		}
		#container.last_played {
			background: #111111;
		}
		
		h1 {
			font-family: 'myriad-pro','Gill Sans','Gill Sans MT',Calibri,'Lucida Grande','Lucida Sans Unicode','Lucida Sans',Tahoma,sans-serif;
			font-size: 18px;
			margin-top: 0px;
			margin-bottom: 10px;
			color: #8D9194;
			font-weight: normal;
			text-transform: uppercase;
		}
		
		.artwork {
			background-image: url("<?php echo $trackInfo['albumArt']; ?>");
			background-position: bottom;
			background-repeat: no-repeat;
			background-size: 120%;
			opacity: .5;
			filter: blur(80px);
			-webkit-filter: blur(80px);
			height: 100%;
			width: 100%;
			overflow: hidden;
		}
		.art_image {
			width: 400px;
			border-radius: 6px;
			position: relative;
			z-index: 2;
			margin-top: 30px;
			box-shadow: 0 0 20px 4px rgba(0, 0, 0, 0.7);
		}
		
		#main {
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
		}
		
		.track {
			font-weight: bold;
			font-size: 32px;
			margin-top: 4px;
			margin-bottom: 0px;
			overflow: hidden;
			display: -webkit-box;
			-webkit-line-clamp: 1;
			-webkit-box-orient: vertical;
			text-align: center;
		}
		
		.artist {
			font-size: 28px;
			margin-bottom: 4px;
			color: #F4F4F4;
			text-align: center;
		}
		
		.album {
			font-size: 20px;
			color: #8D9194;
			font-weight: 100;
			text-align: center;
			font-style: italic;
			overflow: hidden;
			display: -webkit-box;
			-webkit-line-clamp: 1;
			-webkit-box-orient: vertical;
		}
		
		.number {
			font-weight: 100;
			margin-top: 7px;
			font-size: 13px;
			text-transform: uppercase;
			color: #8D9194;
			text-align: center;
		}
		
		.user {
			font-family: 'myriad-pro','Gill Sans','Gill Sans MT',Calibri,'Lucida Grande','Lucida Sans Unicode','Lucida Sans',Tahoma,sans-serif;
			font-weight: bold;
			font-size: 32px;
			margin: 13px 0px 0px 20px;
		}
		
		.list {
			margin-left: 20px;
			position: relative;
			font-family: 'myriad-pro','Gill Sans','Gill Sans MT',Calibri,'Lucida Grande','Lucida Sans Unicode','Lucida Sans',Tahoma,sans-serif;
			color: #D0D0D0;
		}
		
		.stat_box {
			background-color: #1C1C1E;
			border-radius: 8px;
			padding: 9px 10px 10px 14px;
			margin: 11px 0px 20px 20px;
			float: left;
			width: calc(50% - 54px);
			height: 140px;
			color: #D0D0D0;
			text-transform: uppercase;
			font-size: 13px;
			font-weight: bold;
			display:flex;
			justify-content:center;
			align-items:center;
			position: relative;
		}
		.stat_box:last-of-type {
			margin-right: 20px;
		}
		.stat_box .track {
			font-size: 24px;
			text-transform: none;
			color: #fff;
			overflow: hidden;
			display: -webkit-box;
			-webkit-line-clamp: 2;
			-webkit-box-orient: vertical;
			margin-top: 0px;
		}
		.stat_box .scrobbles {
			margin-top: 0px;
			font-size: 45px;
			color: #fff;
			text-align: center;
		}
		.stat_box .artist {
			margin-top: 8px;
			text-transform: none;
			font-size: 16px;
			color: #8D9194;
		}
		.stat_box .header {
			position: absolute;
			top: 0;
			left: 0;
			margin: 10px 0px 0px 15px;
		}
		
		.albums {
			margin-top: 20px;
			margin-left: 53px;
		}
		.top_albums {
			width: 130px;
			margin-left: 16px;
			margin-bottom: 10px;
			border-radius: 6px;
			box-shadow: 0 0 10px 4px rgba(0, 0, 0, 0.7);
		}`

	os.Create("/media/scratch/display/style.css")
	f, _ := os.OpenFile("/media/scratch/display/style.css", os.O_WRONLY, 0777)
	defer f.Close()

	f.WriteString(data)
}

func buildCssNorm() {
	data := `/*! normalize.css v8.0.1 | MIT License | github.com/necolas/normalize.css */

	/* Document
	   ========================================================================== */
	
	/**
	 * 1. Correct the line height in all browsers.
	 * 2. Prevent adjustments of font size after orientation changes in iOS.
	 */
	
	 html {
		line-height: 1.15; /* 1 */
		-webkit-text-size-adjust: 100%; /* 2 */
	  }
	  
	  /* Sections
		 ========================================================================== */
	  
	  /**
	   * Remove the margin in all browsers.
	   */
	  
	  body {
		margin: 0;
	  }
	  
	  /**
	   * Render the main element consistently in IE.
	   */
	  
	  main {
		display: block;
	  }
	  
	  /**
	   * Correct the font size and margin on h1 elements within section and
	   * article contexts in Chrome, Firefox, and Safari
	   */
	  
	  h1 {
		font-size: 2em;
		margin: 0.67em 0;
	  }
	  
	  /* Grouping content
		 ========================================================================== */
	  
	  /**
	   * 1. Add the correct box sizing in Firefox.
	   * 2. Show the overflow in Edge and IE.
	   */
	  
	  hr {
		box-sizing: content-box; /* 1 */
		height: 0; /* 1 */
		overflow: visible; /* 2 */
	  }
	  
	  /**
	   * 1. Correct the inheritance and scaling of font size in all browsers.
	   * 2. Correct the odd em font sizing in all browsers.
	   */
	  
	  pre {
		font-family: monospace, monospace; /* 1 */
		font-size: 1em; /* 2 */
	  }
	  
	  /* Text-level semantics
		 ========================================================================== */
	  
	  /**
	   * Remove the gray background on active links in IE 10.
	   */
	  
	  a {
		background-color: transparent;
	  }
	  
	  /**
	   * 1. Remove the bottom border in Chrome 57-
	   * 2. Add the correct text decoration in Chrome, Edge, IE, Opera, and Safari.
	   */
	  
	  abbr[title] {
		border-bottom: none; /* 1 */
		text-decoration: underline; /* 2 */
		text-decoration: underline dotted; /* 2 */
	  }
	  
	  /**
	   * Add the correct font weight in Chrome, Edge, and Safari.
	   */
	  
	  b,
	  strong {
		font-weight: bolder;
	  }
	  
	  /**
	   * 1. Correct the inheritance and scaling of font size in all browsers.
	   * 2. Correct the odd em font sizing in all browsers.
	   */
	  
	  code,
	  kbd,
	  samp {
		font-family: monospace, monospace; /* 1 */
		font-size: 1em; /* 2 */
	  }
	  
	  /**
	   * Add the correct font size in all browsers.
	   */
	  
	  small {
		font-size: 80%;
	  }
	  
	  /**
	   * Prevent sub and sup elements from affecting the line height in
	   * all browsers.
	   */
	  
	  sub,
	  sup {
		font-size: 75%;
		line-height: 0;
		position: relative;
		vertical-align: baseline;
	  }
	  
	  sub {
		bottom: -0.25em;
	  }
	  
	  sup {
		top: -0.5em;
	  }
	  
	  /* Embedded content
		 ========================================================================== */
	  
	  /**
	   * Remove the border on images inside links in IE 10.
	   */
	  
	  img {
		border-style: none;
	  }
	  
	  /* Forms
		 ========================================================================== */
	  
	  /**
	   * 1. Change the font styles in all browsers.
	   * 2. Remove the margin in Firefox and Safari.
	   */
	  
	  button,
	  input,
	  optgroup,
	  select,
	  textarea {
		font-family: inherit; /* 1 */
		font-size: 100%; /* 1 */
		line-height: 1.15; /* 1 */
		margin: 0; /* 2 */
	  }
	  
	  /**
	   * Show the overflow in IE.
	   * 1. Show the overflow in Edge.
	   */
	  
	  button,
	  input { /* 1 */
		overflow: visible;
	  }
	  
	  /**
	   * Remove the inheritance of text transform in Edge, Firefox, and IE.
	   * 1. Remove the inheritance of text transform in Firefox.
	   */
	  
	  button,
	  select { /* 1 */
		text-transform: none;
	  }
	  
	  /**
	   * Correct the inability to style clickable types in iOS and Safari.
	   */
	  
	  button,
	  [type="button"],
	  [type="reset"],
	  [type="submit"] {
		-webkit-appearance: button;
	  }
	  
	  /**
	   * Remove the inner border and padding in Firefox.
	   */
	  
	  button::-moz-focus-inner,
	  [type="button"]::-moz-focus-inner,
	  [type="reset"]::-moz-focus-inner,
	  [type="submit"]::-moz-focus-inner {
		border-style: none;
		padding: 0;
	  }
	  
	  /**
	   * Restore the focus styles unset by the previous rule.
	   */
	  
	  button:-moz-focusring,
	  [type="button"]:-moz-focusring,
	  [type="reset"]:-moz-focusring,
	  [type="submit"]:-moz-focusring {
		outline: 1px dotted ButtonText;
	  }
	  
	  /**
	   * Correct the padding in Firefox.
	   */
	  
	  fieldset {
		padding: 0.35em 0.75em 0.625em;
	  }
	  
	  /**
	   * 1. Correct the text wrapping in Edge and IE.
	   * 2. Correct the color inheritance from fieldset elements in IE.
	   * 3. Remove the padding so developers are not caught out when they zero out
	   *    fieldset elements in all browsers.
	   */
	  
	  legend {
		box-sizing: border-box; /* 1 */
		color: inherit; /* 2 */
		display: table; /* 1 */
		max-width: 100%; /* 1 */
		padding: 0; /* 3 */
		white-space: normal; /* 1 */
	  }
	  
	  /**
	   * Add the correct vertical alignment in Chrome, Firefox, and Opera.
	   */
	  
	  progress {
		vertical-align: baseline;
	  }
	  
	  /**
	   * Remove the default vertical scrollbar in IE 10+.
	   */
	  
	  textarea {
		overflow: auto;
	  }
	  
	  /**
	   * 1. Add the correct box sizing in IE 10.
	   * 2. Remove the padding in IE 10.
	   */
	  
	  [type="checkbox"],
	  [type="radio"] {
		box-sizing: border-box; /* 1 */
		padding: 0; /* 2 */
	  }
	  
	  /**
	   * Correct the cursor style of increment and decrement buttons in Chrome.
	   */
	  
	  [type="number"]::-webkit-inner-spin-button,
	  [type="number"]::-webkit-outer-spin-button {
		height: auto;
	  }
	  
	  /**
	   * 1. Correct the odd appearance in Chrome and Safari.
	   * 2. Correct the outline style in Safari.
	   */
	  
	  [type="search"] {
		-webkit-appearance: textfield; /* 1 */
		outline-offset: -2px; /* 2 */
	  }
	  
	  /**
	   * Remove the inner padding in Chrome and Safari on macOS.
	   */
	  
	  [type="search"]::-webkit-search-decoration {
		-webkit-appearance: none;
	  }
	  
	  /**
	   * 1. Correct the inability to style clickable types in iOS and Safari.
	   * 2. Change font properties to inherit in Safari.
	   */
	  
	  ::-webkit-file-upload-button {
		-webkit-appearance: button; /* 1 */
		font: inherit; /* 2 */
	  }
	  
	  /* Interactive
		 ========================================================================== */
	  
	  /*
	   * Add the correct display in Edge, IE 10+, and Firefox.
	   */
	  
	  details {
		display: block;
	  }
	  
	  /*
	   * Add the correct display in all browsers.
	   */
	  
	  summary {
		display: list-item;
	  }
	  
	  /* Misc
		 ========================================================================== */
	  
	  /**
	   * Add the correct display in IE 10+.
	   */
	  
	  template {
		display: none;
	  }
	  
	  /**
	   * Add the correct display in IE 10.
	   */
	  
	  [hidden] {
		display: none;
	  }`

	os.Create("/media/scratch/display/normalize.css")
	f, _ := os.OpenFile("/media/scratch/display/normalize.css", os.O_WRONLY, 0777)
	defer f.Close()

	f.WriteString(data)
}
